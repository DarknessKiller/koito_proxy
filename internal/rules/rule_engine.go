package rules

import (
	"database/sql"
	"koito_proxy/internal/model"
	"slices"
	"sync/atomic"

	"github.com/segmentio/ksuid"
)

type MatchCriteria struct {
	TrackName      sql.NullString
	ArtistName     sql.NullString
	ReleaseName    sql.NullString
	ArtistNames    []string
	RecordingMBID  sql.NullString
	DurationBucket int32
}

type Replacement struct {
	TrackName   sql.NullString
	ArtistName  sql.NullString
	ReleaseName sql.NullString
	ArtistNames []string
}

type CompiledRule struct {
	ID ksuid.KSUID

	Match   MatchCriteria
	Replace Replacement

	Priority int
	Valid    bool
}

type engineState struct {
	rules []CompiledRule
}

type RuleEngine struct {
	state atomic.Pointer[engineState]
}

func NewRuleEngine() *RuleEngine {
	engine := &RuleEngine{}

	engine.state.Store(&engineState{
		rules: nil,
	})

	return engine
}

func (e *RuleEngine) UpdateRules(rules []model.Rule) {
	compiled := make([]CompiledRule, 0, len(rules))

	for _, rule := range rules {
		if !rule.Enabled {
			continue
		}

		compiledRule := compileRule(rule)

		if compiledRule.Valid {
			compiled = append(compiled, compiledRule)
		}
	}

	e.state.Store(&engineState{
		rules: compiled,
	})
}

func (e *RuleEngine) Add(rule model.Rule) {
	if !rule.Enabled {
		return
	}

	compiledRule := compileRule(rule)

	if !compiledRule.Valid {
		return
	}

	oldState := e.state.Load()

	var existing []CompiledRule
	if oldState != nil {
		existing = oldState.rules
	}

	next := make([]CompiledRule, 0, len(existing)+1)
	next = append(next, existing...)
	next = append(next, compiledRule)

	e.state.Store(&engineState{
		rules: next,
	})
}

func compileRule(rule model.Rule) CompiledRule {
	compiled := CompiledRule{
		ID: rule.ID,
		Match: MatchCriteria{
			TrackName:      rule.MatchTrackName,
			ArtistName:     rule.MatchArtistName,
			ReleaseName:    rule.MatchReleaseName,
			ArtistNames:    slices.Clone(rule.MatchArtistNames),
			RecordingMBID:  rule.MatchMBID,
			DurationBucket: rule.MatchDurationBucket,
		},
		Replace: Replacement{
			TrackName:   rule.ReplaceTrackName,
			ArtistName:  rule.ReplaceArtistName,
			ReleaseName: rule.ReplaceReleaseName,
			ArtistNames: slices.Clone(rule.ReplaceArtistNames),
		},
	}

	compiled.Priority = calculatePriority(compiled)
	compiled.Valid = compiled.Priority > 0

	return compiled
}

func calculatePriority(rule CompiledRule) int {
	matchFieldCount := 0

	if rule.Match.TrackName.Valid {
		matchFieldCount++
	}

	if rule.Match.ArtistName.Valid {
		matchFieldCount++
	}

	if len(rule.Match.ArtistNames) > 0 {
		matchFieldCount++
	}

	if rule.Match.ReleaseName.Valid {
		matchFieldCount++
	}

	if rule.Match.RecordingMBID.Valid {
		matchFieldCount++
	}

	if rule.Match.DurationBucket != 0 {
		matchFieldCount++
	}

	if matchFieldCount < 2 {
		return 0
	}

	if rule.Match.RecordingMBID.Valid {
		return 1000
	}

	score := 0

	if rule.Match.TrackName.Valid {
		score += 30
	}

	if rule.Match.ArtistName.Valid ||
		len(rule.Match.ArtistNames) > 0 {
		score += 20
	}

	if rule.Match.ReleaseName.Valid {
		score += 25
	}

	if rule.Match.DurationBucket != 0 {
		score += 10
	}

	if rule.Match.TrackName.Valid &&
		rule.Match.ReleaseName.Valid &&
		(rule.Match.ArtistName.Valid ||
			len(rule.Match.ArtistNames) > 0) {
		score += 100
	}

	return score
}

func buildIdentity(
	metadata *model.ListenBrainzTrackMetaData,
) (recordingMBID string, durationBucket int32) {

	recordingMBID = metadata.MBIDMapping.RecordingMBID

	durationMs := metadata.AdditionalInfo.DurationMs
	if durationMs == 0 {
		durationMs = metadata.AdditionalInfo.Duration * 1000
	}

	if durationMs > 0 {
		durationBucket = int32(durationMs / 5000)
	}

	return recordingMBID, durationBucket
}

func (rule CompiledRule) Matches(
	metadata *model.ListenBrainzTrackMetaData,
	recordingMBID string,
	durationBucket int32,
) bool {

	if rule.Match.RecordingMBID.Valid {
		return rule.Match.RecordingMBID.String == recordingMBID
	}

	if rule.Match.TrackName.Valid &&
		rule.Match.TrackName.String != metadata.TrackName {
		return false
	}

	if rule.Match.ArtistName.Valid &&
		rule.Match.ArtistName.String != metadata.ArtistName {
		return false
	}

	if len(rule.Match.ArtistNames) > 0 {
		if len(metadata.AdditionalInfo.ArtistNames) !=
			len(rule.Match.ArtistNames) {
			return false
		}

		for i := range rule.Match.ArtistNames {
			if metadata.AdditionalInfo.ArtistNames[i] !=
				rule.Match.ArtistNames[i] {
				return false
			}
		}
	}

	if rule.Match.ReleaseName.Valid &&
		rule.Match.ReleaseName.String != metadata.ReleaseName {
		return false
	}

	if rule.Match.DurationBucket != 0 &&
		rule.Match.DurationBucket != durationBucket {
		return false
	}

	return true
}

func (rule CompiledRule) Apply(
	metadata *model.ListenBrainzTrackMetaData,
) {

	if rule.Replace.TrackName.Valid {
		metadata.TrackName = rule.Replace.TrackName.String
	}

	if rule.Replace.ArtistName.Valid {
		metadata.ArtistName = rule.Replace.ArtistName.String
	}

	if len(rule.Replace.ArtistNames) > 0 {
		metadata.AdditionalInfo.ArtistNames =
			slices.Clone(rule.Replace.ArtistNames)
	}

	if rule.Replace.ReleaseName.Valid {
		metadata.ReleaseName = rule.Replace.ReleaseName.String
	}
}

func (e *RuleEngine) Apply(
	metadata *model.ListenBrainzTrackMetaData,
) {
	state := e.state.Load()
	if state == nil {
		return
	}

	recordingMBID, durationBucket := buildIdentity(metadata)

	bestPriority := -1
	var matchingRules []CompiledRule

	for _, rule := range state.rules {
		if !rule.Matches(
			metadata,
			recordingMBID,
			durationBucket,
		) {
			continue
		}

		if rule.Priority > bestPriority {
			bestPriority = rule.Priority
			matchingRules = matchingRules[:0]
		}

		if rule.Priority == bestPriority {
			matchingRules = append(matchingRules, rule)
		}
	}

	for _, rule := range matchingRules {
		rule.Apply(metadata)
	}
}
