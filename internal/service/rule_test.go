package service_test

import (
	"context"
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"koito_proxy/internal/model"
	"koito_proxy/internal/rules"
	"koito_proxy/internal/service"
)

type mockRepository struct {
	CreateFunc  func(ctx context.Context, rule *model.Rule) error
	GetByIDFunc func(ctx context.Context, id string) (*model.Rule, error)
	GetAllFunc  func(ctx context.Context) ([]model.Rule, error)
	UpdateFunc  func(ctx context.Context, id string, rule *model.Rule) error
	DeleteFunc  func(ctx context.Context, id string) error
}

func (m *mockRepository) Create(ctx context.Context, rule *model.Rule) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, rule)
	}
	return nil
}

func (m *mockRepository) GetByID(ctx context.Context, id string) (*model.Rule, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, errors.New("mock GetByIDFunc not implemented")
}

func (m *mockRepository) GetAll(ctx context.Context) ([]model.Rule, error) {
	if m.GetAllFunc != nil {
		return m.GetAllFunc(ctx)
	}
	return nil, errors.New("mock GetAllFunc not implemented")
}

func (m *mockRepository) Update(ctx context.Context, id string, rule *model.Rule) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, id, rule)
	}
	return errors.New("mock UpdateFunc not implemented")
}

func (m *mockRepository) Delete(ctx context.Context, id string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return errors.New("mock DeleteFunc not implemented")
}

var _ = Describe("RuleService", func() {

	var (
		mockRuleRepository *mockRepository
		engine             *rules.RuleEngine
		ruleService        *service.RuleService
		ctx                context.Context
	)

	BeforeEach(func() {
		mockRuleRepository = &mockRepository{}
		engine = rules.NewRuleEngine()
		ruleService = service.NewRuleService(mockRuleRepository, engine)
		ctx = context.Background()
	})

	Describe("Create", func() {
		It("persists rule and adds to engine", func() {
			var savedRule *model.Rule
			mockRuleRepository.CreateFunc = func(_ context.Context, rule *model.Rule) error {
				savedRule = rule
				return nil
			}

			rule := &model.Rule{Enabled: boolPtr(true)}
			err := ruleService.Create(ctx, rule)

			Expect(err).NotTo(HaveOccurred())
			Expect(savedRule).To(Equal(rule))
		})

		It("returns error when repository fails", func() {
			mockRuleRepository.CreateFunc = func(_ context.Context, rule *model.Rule) error {
				return errors.New("db error")
			}

			err := ruleService.Create(ctx, &model.Rule{Enabled: boolPtr(true)})
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("GetByID", func() {
		It("returns rule from repository", func() {
			mockRuleRepository.GetByIDFunc = func(_ context.Context, id string) (*model.Rule, error) {
				return &model.Rule{Enabled: boolPtr(true)}, nil
			}

			rule, err := ruleService.GetByID(ctx, "test-id")
			Expect(err).NotTo(HaveOccurred())
			Expect(rule).NotTo(BeNil())
			Expect(*rule.Enabled).To(BeTrue())
		})

		It("returns error when repository fails", func() {
			mockRuleRepository.GetByIDFunc = func(_ context.Context, id string) (*model.Rule, error) {
				return nil, errors.New("not found")
			}

			rule, err := ruleService.GetByID(ctx, "missing")
			Expect(err).To(HaveOccurred())
			Expect(rule).To(BeNil())
		})
	})

	Describe("List", func() {
		It("returns all rules from repository", func() {
			mockRuleRepository.GetAllFunc = func(_ context.Context) ([]model.Rule, error) {
				return []model.Rule{
					{Enabled: boolPtr(true)},
					{Enabled: boolPtr(false)},
				}, nil
			}

			rules, err := ruleService.List(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(rules).To(HaveLen(2))
		})
	})

	Describe("Delete", func() {
		It("deletes from repository and removes from engine", func() {
			rule := &model.Rule{Enabled: boolPtr(true)}
			ruleService.Create(ctx, rule)

			mockRuleRepository.DeleteFunc = func(_ context.Context, id string) error {
				return nil
			}

			err := ruleService.Delete(ctx, rule.ID.String())
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns error when repository fails", func() {
			mockRuleRepository.DeleteFunc = func(_ context.Context, id string) error {
				return errors.New("not found")
			}

			err := ruleService.Delete(ctx, "nonexistent")
			Expect(err).To(HaveOccurred())
		})
	})
})

func boolPtr(v bool) *bool { return &v }
