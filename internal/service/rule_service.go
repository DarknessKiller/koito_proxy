package service

import (
	"context"

	"koito_proxy/internal/model"
	"koito_proxy/internal/repository"
	"koito_proxy/internal/rules"

	"github.com/segmentio/ksuid"
)

type RuleService interface {
	GetAll(ctx context.Context) ([]model.Rule, error)
	GetByID(ctx context.Context, id string) (*model.Rule, error)
	Create(ctx context.Context, rule *model.Rule) error
	Update(ctx context.Context, id string, updatedRule *model.Rule) error
	Delete(ctx context.Context, id string) error
}

type ruleService struct {
	repo   repository.Repository[model.Rule]
	engine *rules.RuleEngine
}

func NewRuleService(repo repository.Repository[model.Rule], engine *rules.RuleEngine) RuleService {
	return &ruleService{
		repo:   repo,
		engine: engine,
	}
}

func (s *ruleService) GetAll(ctx context.Context) ([]model.Rule, error) {
	return s.repo.GetAll(ctx)
}

func (s *ruleService) GetByID(ctx context.Context, id string) (*model.Rule, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ruleService) Create(ctx context.Context, rule *model.Rule) error {
	if err := s.repo.Create(ctx, rule); err != nil {
		return err
	}
	if s.engine != nil {
		s.engine.Add(*rule)
	}
	return nil
}

func (s *ruleService) Update(ctx context.Context, id string, updatedRule *model.Rule) error {
	if err := s.repo.Update(ctx, id, updatedRule); err != nil {
		return err
	}

	if s.engine != nil {
		parsedID, err := ksuid.Parse(id)
		if err == nil {
			s.engine.Remove(parsedID)
		}
		s.engine.Add(*updatedRule)
	}
	return nil
}

func (s *ruleService) Delete(ctx context.Context, id string) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	if s.engine != nil {
		parsedID, err := ksuid.Parse(id)
		if err == nil {
			s.engine.Remove(parsedID)
		}
	}
	return nil
}
