package rules_test

import (
	"context"
	"database/sql"
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"koito_proxy/internal/model"
	"koito_proxy/internal/rules"
)

// MockRuleRepository
type MockRuleRepository struct {
	CreateFunc  func(ctx context.Context, rule *model.Rule) error
	GetByIDFunc func(ctx context.Context, id string) (*model.Rule, error)
	GetAllFunc  func(ctx context.Context) ([]model.Rule, error)
	UpdateFunc  func(ctx context.Context, id string, rule *model.Rule) error
	DeleteFunc  func(ctx context.Context, id string) error
}

func (m *MockRuleRepository) Create(ctx context.Context, rule *model.Rule) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, rule)
	}
	return nil
}

func (m *MockRuleRepository) GetByID(ctx context.Context, id string) (*model.Rule, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, errors.New("mock GetByIDFunc not implemented")
}

func (m *MockRuleRepository) GetAll(ctx context.Context) ([]model.Rule, error) {
	if m.GetAllFunc != nil {
		return m.GetAllFunc(ctx)
	}
	return nil, errors.New("mock GetAllFunc not implemented")
}

func (m *MockRuleRepository) Update(ctx context.Context, id string, rule *model.Rule) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, id, rule)
	}
	return errors.New("mock UpdateFunc not implemented")
}

func (m *MockRuleRepository) Delete(ctx context.Context, id string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return errors.New("mock DeleteFunc not implemented")
}

var _ = Describe("RuleEngine.Apply", func() {

	It("do not allow overwrite rule when there only one criteria like track name", func() {
		mockRepo := &MockRuleRepository{
			GetAllFunc: func(ctx context.Context) ([]model.Rule, error) {
				return []model.Rule{
					{
						MatchTrackName: sql.NullString{
							String: "Old Track",
							Valid:  true,
						},
						ReplaceTrackName: sql.NullString{
							String: "New Track",
							Valid:  true,
						},
						Enabled: new(true),
					},
				}, nil
			},
		}

		engine := rules.NewRuleEngine()
		rules, err := mockRepo.GetAll(context.Background())
		Expect(err).NotTo(HaveOccurred())
		engine.UpdateRules(rules)

		meta := &model.ListenBrainzTrackMetaData{
			TrackName:  "Old Track",
			ArtistName: "Artist A",
		}

		engine.Apply(meta)

		Expect(meta.TrackName).To(Equal("Old Track"))
	})

	It("does not overwrite when only artist matches", func() {
		mockRepo := &MockRuleRepository{
			GetAllFunc: func(ctx context.Context) ([]model.Rule, error) {
				return []model.Rule{
					{
						MatchArtistName: sql.NullString{
							String: "Artist A",
							Valid:  true,
						},
						ReplaceArtistName: sql.NullString{
							String: "Artist A1",
							Valid:  true,
						},
						Enabled: new(true),
					},
				}, nil
			},
		}

		engine := rules.NewRuleEngine()
		rules, err := mockRepo.GetAll(context.Background())
		Expect(err).NotTo(HaveOccurred())
		engine.UpdateRules(rules)

		meta := &model.ListenBrainzTrackMetaData{
			TrackName:  "Old Track",
			ArtistName: "Artist A",
		}

		engine.Apply(meta)

		Expect(meta.TrackName).To(Equal("Old Track"))
		Expect(meta.ArtistName).To(Equal("Artist A"))
	})

	It("overwrites track name only when both track and artist match", func() {
		mockRepo := &MockRuleRepository{
			GetAllFunc: func(ctx context.Context) ([]model.Rule, error) {
				return []model.Rule{
					{
						MatchTrackName: sql.NullString{
							String: "Old Track",
							Valid:  true,
						},
						MatchArtistName: sql.NullString{
							String: "Artist A",
							Valid:  true,
						},
						ReplaceTrackName: sql.NullString{
							String: "New Track",
							Valid:  true,
						},
						Enabled: new(true),
					},
				}, nil
			},
		}

		engine := rules.NewRuleEngine()
		rules, err := mockRepo.GetAll(context.Background())
		Expect(err).NotTo(HaveOccurred())
		engine.UpdateRules(rules)

		meta := &model.ListenBrainzTrackMetaData{
			TrackName:  "Old Track",
			ArtistName: "Artist A",
		}

		engine.Apply(meta)

		Expect(meta.TrackName).To(Equal("New Track"))
	})

	It("overwrites track name only when track and release match", func() {
		mockRepo := &MockRuleRepository{
			GetAllFunc: func(ctx context.Context) ([]model.Rule, error) {
				return []model.Rule{
					{
						MatchTrackName: sql.NullString{
							String: "Old Track",
							Valid:  true,
						},
						MatchReleaseName: sql.NullString{
							String: "Album B",
							Valid:  true,
						},
						ReplaceTrackName: sql.NullString{
							String: "New Track",
							Valid:  true,
						},
						Enabled: new(true),
					},
				}, nil
			},
		}

		engine := rules.NewRuleEngine()
		rules, err := mockRepo.GetAll(context.Background())
		Expect(err).NotTo(HaveOccurred())
		engine.UpdateRules(rules)

		meta := &model.ListenBrainzTrackMetaData{
			TrackName:   "Old Track",
			ReleaseName: "Album A",
		}

		engine.Apply(meta)

		Expect(meta.TrackName).To(Equal("Old Track"))
		Expect(meta.ReleaseName).To(Equal("Album A"))
	})

	It("overwrites artist list when track, release and artists match", func() {
		mockRepo := &MockRuleRepository{
			GetAllFunc: func(ctx context.Context) ([]model.Rule, error) {
				return []model.Rule{
					{
						MatchTrackName: sql.NullString{
							String: "Track",
							Valid:  true,
						},
						MatchArtistNames:   []string{"Artist A", "Artist B"},
						ReplaceArtistNames: []string{"Artist A", "Artist B1"},
						MatchReleaseName: sql.NullString{
							String: "Album",
							Valid:  true,
						},
						Enabled: new(true),
					},
				}, nil
			},
		}

		engine := rules.NewRuleEngine()
		rules, err := mockRepo.GetAll(context.Background())
		Expect(err).NotTo(HaveOccurred())
		engine.UpdateRules(rules)

		meta := &model.ListenBrainzTrackMetaData{
			TrackName:  "Track",
			ArtistName: "Artist A",
			AdditionalInfo: model.ListenBrainzAdditionalInfo{
				ArtistNames: []string{"Artist A", "Artist B"},
			},
			ReleaseName: "Album",
		}

		engine.Apply(meta)

		Expect(meta.ArtistName).To(Equal("Artist A"))
		Expect(meta.AdditionalInfo.ArtistNames).To(Equal([]string{"Artist A", "Artist B1"}))
	})

	It("overwrite when artist & artists matches", func() {
		mockRepo := &MockRuleRepository{
			GetAllFunc: func(ctx context.Context) ([]model.Rule, error) {
				return []model.Rule{
					{
						MatchArtistName: sql.NullString{
							String: "Artist A",
							Valid:  true,
						},
						ReplaceArtistName: sql.NullString{
							String: "Artist A1",
							Valid:  true,
						},
						MatchArtistNames:   []string{"Artist A"},
						ReplaceArtistNames: []string{"Artist A1"},
						Enabled:            new(true),
					},
				}, nil
			},
		}

		engine := rules.NewRuleEngine()
		rules, err := mockRepo.GetAll(context.Background())
		Expect(err).NotTo(HaveOccurred())
		engine.UpdateRules(rules)

		meta := &model.ListenBrainzTrackMetaData{
			TrackName:  "Old Track",
			ArtistName: "Artist A",
			AdditionalInfo: model.ListenBrainzAdditionalInfo{
				ArtistNames: []string{"Artist A"},
			},
		}

		engine.Apply(meta)

		Expect(meta.TrackName).To(Equal("Old Track"))
		Expect(meta.ArtistName).To(Equal("Artist A1"))
	})
})
