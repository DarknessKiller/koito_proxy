package bootstrap

import (
	"log/slog"
	"net/http"
	"time"

	"koito_proxy/internal/config"
	db "koito_proxy/internal/database"
	"koito_proxy/internal/model"
	"koito_proxy/internal/repository"
	"koito_proxy/internal/rules"
	"koito_proxy/internal/service"

	"gorm.io/gorm"
)

type Bootstrap struct {
	Config         *config.Config
	DB             *gorm.DB
	RuleRepository *repository.BaseRepository[model.Rule]
	RuleEngine     *rules.RuleEngine
	RuleService    *service.RuleService
	KoitoService   *service.KoitoService
	HttpClient     *http.Client
}

func New() (*Bootstrap, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.InitiateDatabase(cfg)
	if err != nil {
		slog.Error("failed to open SQLite database", "error", err)
		return nil, err
	}

	ruleRepository := repository.NewBaseRepository(sqlDB)

	engine := rules.NewRuleEngine()

	httpclient := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        250,
			MaxIdleConnsPerHost: 40,
			IdleConnTimeout:     90 * time.Second,
		},
	}

	ruleService := service.NewRuleService(ruleRepository, engine)
	koitoService := service.NewKoitoService(ruleService, cfg, httpclient)

	return &Bootstrap{
		Config:         cfg,
		DB:             sqlDB,
		RuleRepository: ruleRepository,
		RuleEngine:     engine,
		RuleService:    ruleService,
		KoitoService:   koitoService,
		HttpClient:     httpclient,
	}, nil
}
