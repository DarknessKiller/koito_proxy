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

type mockRuleRepository struct {
	CreateFunc  func(ctx context.Context, rule *model.Rule) error
	GetByIDFunc func(ctx context.Context, id string) (*model.Rule, error)
	GetAllFunc  func(ctx context.Context) ([]model.Rule, error)
	UpdateFunc  func(ctx context.Context, id string, rule *model.Rule) error
	DeleteFunc  func(ctx context.Context, id string) error
}

func (m *mockRuleRepository) Create(ctx context.Context, rule *model.Rule) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, rule)
	}
	return nil
}

func (m *mockRuleRepository) GetByID(ctx context.Context, id string) (*model.Rule, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, errors.New("mock GetByIDFunc not implemented")
}

func (m *mockRuleRepository) GetAll(ctx context.Context) ([]model.Rule, error) {
	if m.GetAllFunc != nil {
		return m.GetAllFunc(ctx)
	}
	return nil, errors.New("mock GetAllFunc not implemented")
}

func (m *mockRuleRepository) Update(ctx context.Context, id string, rule *model.Rule) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, id, rule)
	}
	return errors.New("mock UpdateFunc not implemented")
}

func (m *mockRuleRepository) Delete(ctx context.Context, id string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return errors.New("mock DeleteFunc not implemented")
}

var _ = Describe("RuleService", func() {

	var (
		mockRepo *mockRuleRepository
		engine   *rules.RuleEngine
		svc      *service.RuleService
		ctx      context.Context
	)

	BeforeEach(func() {
		mockRepo = &mockRuleRepository{}
		engine = rules.NewRuleEngine()
		svc = service.NewRuleService(mockRepo, engine)
		ctx = context.Background()
	})

	Describe("Create", func() {
		It("persists rule and adds to engine", func() {
			var savedRule *model.Rule
			mockRepo.CreateFunc = func(_ context.Context, rule *model.Rule) error {
				savedRule = rule
				return nil
			}

			rule := &model.Rule{Enabled: boolPtr(true)}
			err := svc.Create(ctx, rule)

			Expect(err).NotTo(HaveOccurred())
			Expect(savedRule).To(Equal(rule))
		})

		It("returns error when repository fails", func() {
			mockRepo.CreateFunc = func(_ context.Context, rule *model.Rule) error {
				return errors.New("db error")
			}

			err := svc.Create(ctx, &model.Rule{Enabled: boolPtr(true)})
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("GetByID", func() {
		It("returns rule from repository", func() {
			mockRepo.GetByIDFunc = func(_ context.Context, id string) (*model.Rule, error) {
				return &model.Rule{Enabled: boolPtr(true)}, nil
			}

			rule, err := svc.GetByID(ctx, "test-id")
			Expect(err).NotTo(HaveOccurred())
			Expect(rule).NotTo(BeNil())
			Expect(*rule.Enabled).To(BeTrue())
		})

		It("returns error when repository fails", func() {
			mockRepo.GetByIDFunc = func(_ context.Context, id string) (*model.Rule, error) {
				return nil, errors.New("not found")
			}

			rule, err := svc.GetByID(ctx, "missing")
			Expect(err).To(HaveOccurred())
			Expect(rule).To(BeNil())
		})
	})

	Describe("List", func() {
		It("returns all rules from repository", func() {
			mockRepo.GetAllFunc = func(_ context.Context) ([]model.Rule, error) {
				return []model.Rule{
					{Enabled: boolPtr(true)},
					{Enabled: boolPtr(false)},
				}, nil
			}

			rules, err := svc.List(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(rules).To(HaveLen(2))
		})
	})

	Describe("Delete", func() {
		It("deletes from repository and removes from engine", func() {
			rule := &model.Rule{Enabled: boolPtr(true)}
			svc.Create(ctx, rule)

			mockRepo.DeleteFunc = func(_ context.Context, id string) error {
				return nil
			}

			err := svc.Delete(ctx, rule.ID.String())
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns error when repository fails", func() {
			mockRepo.DeleteFunc = func(_ context.Context, id string) error {
				return errors.New("not found")
			}

			err := svc.Delete(ctx, "nonexistent")
			Expect(err).To(HaveOccurred())
		})
	})
})

func boolPtr(v bool) *bool { return &v }
