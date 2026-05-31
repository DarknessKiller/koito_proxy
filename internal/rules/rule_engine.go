package rules

import (
	"errors"
	"koito_proxy/internal/model"
	"sync"
	"sync/atomic"
)

var (
	ErrArtistReplacementRequiresTrackMatch = errors.New("artist replacement requires track match")
)

type engineState struct {
	artistIndex  map[string][]Rule
	trackIndex   map[string][]Rule
	releaseIndex map[string][]Rule
	global       []Rule
}

type RuleEngine struct {
	state   atomic.Pointer[engineState]
	writeMu sync.Mutex
}

func NewRuleEngine(rules []Rule) *RuleEngine {
	e := &RuleEngine{}

	s := &engineState{
		artistIndex:  map[string][]Rule{},
		trackIndex:   map[string][]Rule{},
		releaseIndex: map[string][]Rule{},
	}

	for _, r := range rules {
		if err := ValidateRule(r); err != nil {
			continue
		}

		addToState(s, r)
	}

	e.state.Store(s)

	return e
}

func (e *RuleEngine) Apply(metadata *model.ListenBrainzTrackMetaData) {
	s := e.state.Load()
	if s == nil {
		return
	}

	var candidates []Rule

	if rs, ok := s.artistIndex[metadata.ArtistName]; ok {
		candidates = append(candidates, rs...)
	}

	if rs, ok := s.trackIndex[metadata.TrackName]; ok {
		candidates = append(candidates, rs...)
	}

	if rs, ok := s.releaseIndex[metadata.ReleaseName]; ok {
		candidates = append(candidates, rs...)
	}

	candidates = append(candidates, s.global...)

	for _, r := range candidates {
		if !match(r, metadata) {
			continue
		}

		if r.ReplaceTrackName.Valid {
			metadata.TrackName = r.ReplaceTrackName.String
		}

		if r.ReplaceArtistName.Valid {
			metadata.ArtistName = r.ReplaceArtistName.String
			metadata.AdditionalInfo.ArtistNames = []string{r.ReplaceArtistName.String}
		}

		if r.ReplaceReleaseName.Valid {
			metadata.ReleaseName = r.ReplaceReleaseName.String
		}
	}
}

func (e *RuleEngine) Add(r Rule) error {
	if err := ValidateRule(r); err != nil {
		return err
	}

	e.writeMu.Lock()
	defer e.writeMu.Unlock()

	current := e.state.Load()

	next := cloneState(current)

	addToState(next, r)

	e.state.Store(next)
	return nil
}

func (e *RuleEngine) Delete(id int64) {
	e.writeMu.Lock()
	defer e.writeMu.Unlock()

	current := e.state.Load()

	next := &engineState{
		artistIndex:  make(map[string][]Rule, len(current.artistIndex)),
		trackIndex:   make(map[string][]Rule, len(current.trackIndex)),
		releaseIndex: make(map[string][]Rule, len(current.releaseIndex)),
	}

	for k, rules := range current.artistIndex {
		filtered := keepRulesExcept(rules, id)

		if len(filtered) > 0 {
			next.artistIndex[k] = filtered
		}
	}

	for k, rules := range current.trackIndex {
		filtered := keepRulesExcept(rules, id)

		if len(filtered) > 0 {
			next.trackIndex[k] = filtered
		}
	}

	for k, rules := range current.releaseIndex {
		filtered := keepRulesExcept(rules, id)

		if len(filtered) > 0 {
			next.releaseIndex[k] = filtered
		}
	}

	if filtered := keepRulesExcept(current.global, id); len(filtered) > 0 {
		next.global = filtered
	}

	e.state.Store(next)
}

func keepRulesExcept(rules []Rule, id int64) []Rule {
	filtered := make([]Rule, 0, len(rules))

	for _, r := range rules {
		if r.ID == id {
			continue
		}

		filtered = append(filtered, r)
	}

	return filtered
}

func addToState(s *engineState, r Rule) {
	if !r.MatchArtistName.Valid &&
		!r.MatchTrackName.Valid &&
		!r.MatchReleaseName.Valid {

		s.global = append(s.global, r)
		return
	}

	if r.MatchArtistName.Valid {
		key := r.MatchArtistName.String
		s.artistIndex[key] = append(s.artistIndex[key], r)
	}

	if r.MatchTrackName.Valid {
		key := r.MatchTrackName.String
		s.trackIndex[key] = append(s.trackIndex[key], r)
	}

	if r.MatchReleaseName.Valid {
		key := r.MatchReleaseName.String
		s.releaseIndex[key] = append(s.releaseIndex[key], r)
	}
}

func cloneState(old *engineState) *engineState {
	s := &engineState{
		artistIndex:  make(map[string][]Rule, len(old.artistIndex)),
		trackIndex:   make(map[string][]Rule, len(old.trackIndex)),
		releaseIndex: make(map[string][]Rule, len(old.releaseIndex)),
		global:       append([]Rule(nil), old.global...),
	}

	for k, v := range old.artistIndex {
		s.artistIndex[k] = append([]Rule(nil), v...)
	}

	for k, v := range old.trackIndex {
		s.trackIndex[k] = append([]Rule(nil), v...)
	}

	for k, v := range old.releaseIndex {
		s.releaseIndex[k] = append([]Rule(nil), v...)
	}

	return s
}

func match(r Rule, m *model.ListenBrainzTrackMetaData) bool {
	if r.MatchArtistName.Valid && r.MatchArtistName.String != m.ArtistName {
		return false
	}
	if r.MatchTrackName.Valid && r.MatchTrackName.String != m.TrackName {
		return false
	}
	if r.MatchReleaseName.Valid && r.MatchReleaseName.String != m.ReleaseName {
		return false
	}
	return r.MatchArtistName.Valid || r.MatchTrackName.Valid || r.MatchReleaseName.Valid
}

func ValidateRule(r Rule) error {
	if r.ReplaceArtistName.Valid {
		if !r.MatchTrackName.Valid {
			return ErrArtistReplacementRequiresTrackMatch
		}
	}

	return nil
}
