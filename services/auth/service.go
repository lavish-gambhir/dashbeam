package auth

import (
	"log/slog"
	"net/http"

	"github.com/lavish-gambhir/dashbeam/services/auth/repository"
	"github.com/lavish-gambhir/dashbeam/shared/config"
)

type Service interface {
	RegisterRoutes(mux *http.ServeMux, prefix string)
}

type service struct {
	dashboardRepo repository.DashboardRepository
	authConfig    config.AuthConfig
	logger        *slog.Logger
}

func New(dashboardRepo repository.DashboardRepository, authConfig config.AuthConfig, logger *slog.Logger) Service {
	return &service{
		dashboardRepo: dashboardRepo,
		authConfig:    authConfig,
		logger:        logger,
	}
}

func (s *service) RegisterRoutes(parentmux *http.ServeMux, prefix string) {
	h := &handler{
		dashboardRepo: s.dashboardRepo,
		authConfig:    s.authConfig,
		logger:        s.logger,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/validate", h.handleValidateJWT)
	mux.HandleFunc("/login", h.handleDashboardLogin)
	mux.HandleFunc("/logout", h.handleDashboardLogout)
	mux.HandleFunc("/me", h.handleGetCurrentUser)
	parentmux.Handle(prefix+"/", http.StripPrefix(prefix, mux))
}
