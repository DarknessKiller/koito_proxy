package service

import (
	"context"
	"fmt"
	"koito_proxy/internal/model"
	"koito_proxy/internal/repository"
	"koito_proxy/internal/rules"

	"github.com/segmentio/ksuid"
)

type RuleService struct {
	repo   repository.Repository[model.Rule]
	engine *rules.RuleEngine
}

func NewRuleService(repo repository.Repository[model.Rule], engine *rules.RuleEngine) *RuleService {
	return &RuleService{repo: repo, engine: engine}
}

func (s *RuleService) Create(ctx context.Context, rule *model.Rule) error {
	if err := s.repo.Create(ctx, rule); err != nil {
		return fmt.Errorf("create rule: %w", err)
	}
	s.engine.Add(*rule)
	return nil
}

func (s *RuleService) GetByID(ctx context.Context, id string) (*model.Rule, error) {
	rule, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get rule: %w", err)
	}
	return rule, nil
}

func (s *RuleService) List(ctx context.Context) ([]model.Rule, error) {
	rules, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("list rules: %w", err)
	}
	return rules, nil
}

func (s *RuleService) Update(ctx context.Context, id string, rule *model.Rule) error {
	parsedID, err := ksuid.Parse(id)
	if err != nil {
		return fmt.Errorf("parse rule id: %w", err)
	}

	if err := s.repo.Update(ctx, id, rule); err != nil {
		return fmt.Errorf("update rule: %w", err)
	}

	s.engine.Remove(parsedID)
	s.engine.Add(*rule)
	return nil
}

func (s *RuleService) Delete(ctx context.Context, id string) error {
	parsedID, err := ksuid.Parse(id)
	if err != nil {
		return fmt.Errorf("parse rule id: %w", err)
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete rule: %w", err)
	}

	s.engine.Remove(parsedID)
	return nil
}
