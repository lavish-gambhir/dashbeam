package auth

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/lavish-gambhir/dashbeam/pkg/apperr"
	"github.com/lavish-gambhir/dashbeam/pkg/utils"
	"github.com/lavish-gambhir/dashbeam/services/auth/repository"
	"github.com/lavish-gambhir/dashbeam/shared/config"
	sharedcontext "github.com/lavish-gambhir/dashbeam/shared/context"
	"github.com/lavish-gambhir/dashbeam/shared/models"
)

type handler struct {
	dashboardRepo repository.DashboardRepository
	authConfig    config.AuthConfig
	logger        *slog.Logger
}

func NewHandler(dashboardRepo repository.DashboardRepository, authConfig config.AuthConfig, logger *slog.Logger) *handler {
	log := logger.With("handler", "auth.handler")
	return &handler{
		dashboardRepo: dashboardRepo,
		authConfig:    authConfig,
		logger:        log,
	}
}

func (h *handler) handleValidateJWT(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	reqID, _ := sharedcontext.GetRequestID(ctx)
	logger := h.logger.With("fn", "handleValidateJWT").With("requestID", reqID)

	if r.Method != http.MethodGet {
		utils.WriteJSONError(w, apperr.New(apperr.BadRequest, "method not allowed"), http.StatusMethodNotAllowed)
		return
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		utils.WriteJSONError(w, apperr.New(apperr.BadRequest, "token is required"), http.StatusBadRequest)
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == authHeader {
		logger.Warn("invalid authorization header format",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path))
		utils.WriteJSONError(w, apperr.New(apperr.Unauthorized, "bearer token required"), http.StatusUnauthorized)
		return
	}

	userContext, err := h.validateMobileJWT(token)
	if err != nil {
		logger.Warn("JWT validation failed", slog.Any("error", err))
		utils.WriteJSONSuccess(w, ValidateJWTResponse{
			Valid: false,
			Error: err.Error(),
		})
		return
	}

	utils.WriteJSONSuccess(w, ValidateJWTResponse{
		Valid:       true,
		UserContext: userContext,
	})
}

func (h *handler) handleDashboardLogin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	reqID, _ := sharedcontext.GetRequestID(ctx)
	logger := h.logger.With("fn", "handleDashboardLogin").With("requestID", reqID)

	if r.Method != http.MethodPost {
		utils.WriteJSONError(w, apperr.New(apperr.BadRequest, "method not allowed"), http.StatusMethodNotAllowed)
		return
	}

	var req DashboardLoginRequest
	if err := utils.FromJson(r.Body, &req); err != nil {
		utils.WriteJSONError(w, apperr.Wrap(err, apperr.BadRequest, "invalid request body"), http.StatusBadRequest)
		return
	}

	if req.Username == "" || req.Password == "" {
		utils.WriteJSONError(w, apperr.New(apperr.BadRequest, "username and password are required"), http.StatusBadRequest)
		return
	}

	// Get user by username
	user, err := h.dashboardRepo.GetUserByUsername(ctx, req.Username)
	if err != nil {
		logger.Error("failed to get user by username", slog.Any("error", err), slog.String("username", req.Username))
		utils.WriteJSONError(w, apperr.New(apperr.InvalidCredentials, "invalid credentials"), http.StatusUnauthorized)
		return
	}

	if !user.IsActive {
		logger.Warn("inactive user attempted login", slog.String("username", req.Username))
		utils.WriteJSONError(w, apperr.New(apperr.InvalidCredentials, "account is inactive"), http.StatusUnauthorized)
		return
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		logger.Warn("invalid password attempt", slog.String("username", req.Username))
		utils.WriteJSONError(w, apperr.New(apperr.InvalidCredentials, "invalid credentials"), http.StatusUnauthorized)
		return
	}

	// Create session
	sessionID := uuid.New()
	expiresAt := time.Now().UTC().Add(h.authConfig.AccessTokenExpiry)
	session := &models.DashboardSession{
		ID:        sessionID,
		UserID:    user.ID,
		Username:  user.Username,
		FullName:  user.FullName,
		Email:     user.Email,
		CreatedAt: time.Now().UTC(),
		ExpiresAt: expiresAt,
	}

	accessToken, err := h.generateDashboardJWT(user, sessionID, h.authConfig.AccessTokenExpiry)
	if err != nil {
		logger.Error("failed to generate access token", slog.Any("error", err))
		utils.WriteJSONError(w, apperr.Wrap(err, apperr.Internal, "failed to generate access token"), http.StatusInternalServerError)
		return
	}

	if err := h.dashboardRepo.UpdateLastLogin(ctx, user.ID.String()); err != nil {
		logger.Warn("failed to update last login", slog.Any("error", err))
		// Don't fail the login for this
	}

	logger.Info("successful dashboard login", slog.String("username", user.Username))
	w.Header().Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	utils.WriteJSONSuccess(w, DashboardLoginResponse{
		Success:   true,
		User:      user,
		Session:   session,
		ExpiresAt: expiresAt,
	})
}

func (h *handler) handleDashboardLogout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteJSONError(w, apperr.New(apperr.BadRequest, "method not allowed"), http.StatusMethodNotAllowed)
		return
	}

	// TODO: maintain a blacklist of invalidated tokens
	utils.WriteJSONSuccess(w, LogoutResponse{
		Success: true,
		Message: "logged out successfully",
	})
}

func (h *handler) handleGetCurrentUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	reqID, _ := sharedcontext.GetRequestID(ctx)
	logger := h.logger.With("fn", "handleGetCurrentUser").With("requestID", reqID)

	if r.Method != http.MethodGet {
		utils.WriteJSONError(w, apperr.New(apperr.BadRequest, "method not allowed"), http.StatusMethodNotAllowed)
		return
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		utils.WriteJSONError(w, apperr.New(apperr.Unauthorized, "authorization header required"), http.StatusUnauthorized)
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		utils.WriteJSONError(w, apperr.New(apperr.Unauthorized, "bearer token required"), http.StatusUnauthorized)
		return
	}

	claims, err := h.validateDashboardJWT(tokenString)
	if err != nil {
		logger.Warn("invalid dashboard JWT", slog.Any("error", err))
		utils.WriteJSONError(w, apperr.New(apperr.InvalidToken, "invalid token"), http.StatusUnauthorized)
		return
	}

	user, err := h.dashboardRepo.GetUserByUsername(ctx, claims.Username)
	if err != nil {
		logger.Error("failed to get user", slog.Any("error", err))
		utils.WriteJSONError(w, apperr.New(apperr.UserNotFound, "user not found"), http.StatusNotFound)
		return
	}

	session := &models.DashboardSession{
		ID:        claims.SessionID,
		UserID:    user.ID,
		Username:  user.Username,
		FullName:  user.FullName,
		Email:     user.Email,
		CreatedAt: claims.IssuedAt,
		ExpiresAt: claims.ExpiresAt,
	}

	utils.WriteJSONSuccess(w, CurrentUserResponse{
		User:    user,
		Session: session,
	})
}

func (h *handler) validateMobileJWT(tokenString string) (*models.UserContext, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(h.authConfig.JWTSecretKey), nil
	})

	if err != nil {
		return nil, apperr.Wrap(err, apperr.InvalidToken, "failed to parse token")
	}

	if !token.Valid {
		return nil, apperr.New(apperr.InvalidToken, "token is invalid")
	}

	claims, ok := token.Claims.(*models.JWTClaims)
	if !ok {
		return nil, apperr.New(apperr.InvalidToken, "invalid token claims")
	}

	// Create user context from claims
	userContext := &models.UserContext{
		UserID:      claims.UserID,
		Email:       claims.Email,
		Name:        claims.Name,
		Role:        claims.Role,
		SchoolID:    claims.SchoolID,
		ClassroomID: claims.ClassroomID,
		AppType:     claims.AppType,
		IssuedAt:    claims.IssuedAt.Time,
		ExpiresAt:   claims.ExpiresAt.Time,
	}

	if !userContext.IsValid() {
		return nil, apperr.New(apperr.TokenExpired, "token has expired or is invalid")
	}

	return userContext, nil
}

func (h *handler) validateDashboardJWT(tokenString string) (*dashboardClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(h.authConfig.JWTSecretKey), nil
	})

	if err != nil {
		return nil, apperr.Wrap(err, apperr.InvalidToken, "failed to parse token")
	}

	if !token.Valid {
		return nil, apperr.New(apperr.InvalidToken, "token is invalid")
	}

	mapClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, apperr.New(apperr.InvalidToken, "invalid token claims")
	}

	claims := &dashboardClaims{}

	if userID, ok := mapClaims["user_id"].(string); ok {
		if parsed, err := uuid.Parse(userID); err == nil {
			claims.UserID = parsed
		}
	}

	if username, ok := mapClaims["username"].(string); ok {
		claims.Username = username
	}

	if fullName, ok := mapClaims["full_name"].(string); ok {
		claims.FullName = fullName
	}

	if email, ok := mapClaims["email"].(string); ok {
		claims.Email = email
	}

	if sessionID, ok := mapClaims["session_id"].(string); ok {
		if parsed, err := uuid.Parse(sessionID); err == nil {
			claims.SessionID = parsed
		}
	}

	if iat, ok := mapClaims["iat"].(float64); ok {
		claims.IssuedAt = time.Unix(int64(iat), 0)
	}

	if exp, ok := mapClaims["exp"].(float64); ok {
		claims.ExpiresAt = time.Unix(int64(exp), 0)
	}

	if nbf, ok := mapClaims["nbf"].(float64); ok {
		claims.NotBefore = time.Unix(int64(nbf), 0)
	}

	if iss, ok := mapClaims["iss"].(string); ok {
		claims.Issuer = iss
	}

	if sub, ok := mapClaims["sub"].(string); ok {
		claims.Subject = sub
	}

	return claims, nil
}

func (h *handler) generateDashboardJWT(user *models.DashboardUser, sessionID uuid.UUID, expiry time.Duration) (string, error) {
	now := time.Now().UTC()

	claims := jwt.MapClaims{
		"user_id":    user.ID.String(),
		"username":   user.Username,
		"full_name":  user.FullName,
		"email":      user.Email,
		"session_id": sessionID.String(),
		"iat":        now.Unix(),
		"exp":        now.Add(expiry).Unix(),
		"nbf":        now.Unix(),
		"iss":        "dashbeam-auth",
		"sub":        user.ID.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(h.authConfig.JWTSecretKey))
}
