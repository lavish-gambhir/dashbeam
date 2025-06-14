package middleware

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lavish-gambhir/dashbeam/pkg/apperr"
	"github.com/lavish-gambhir/dashbeam/pkg/utils"
	"github.com/lavish-gambhir/dashbeam/shared/config"
	sharedcontext "github.com/lavish-gambhir/dashbeam/shared/context"
	"github.com/lavish-gambhir/dashbeam/shared/models"
)

// AuthMiddleware provides JWT validation middleware for mobile app tokens
type AuthMiddleware struct {
	authConfig config.AuthConfig
	logger     *slog.Logger
}

// NewAuthMiddleware creates a new JWT validation middleware
func NewAuthMiddleware(authConfig config.AuthConfig, logger *slog.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		authConfig: authConfig,
		logger:     logger.With("middleware", "auth"),
	}
}

func (am *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		reqID, _ := sharedcontext.GetRequestID(ctx)
		logger := am.logger.With("fn", "RequireAuth").With("requestID", reqID)

		// Extract Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			logger.Warn("missing authorization header",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path))
			utils.WriteJSONError(w, apperr.New(apperr.Unauthorized, "authorization header required"), http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			logger.Warn("invalid authorization header format",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path))
			utils.WriteJSONError(w, apperr.New(apperr.Unauthorized, "bearer token required"), http.StatusUnauthorized)
			return
		}

		userContext, err := am.validateMobileJWT(tokenString)
		if err != nil {
			logger.Warn("JWT validation failed",
				slog.Any("error", err),
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path))
			utils.WriteJSONError(w, apperr.New(apperr.InvalidToken, "invalid or expired token"), http.StatusUnauthorized)
			return
		}

		ctxWithUser := sharedcontext.WithUserContext(ctx, userContext)
		r = r.WithContext(ctxWithUser)

		logger.Debug("authenticated request",
			slog.String("userID", userContext.UserID.String()),
			slog.String("role", userContext.Role),
			slog.String("schoolID", userContext.SchoolID.String()),
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path))

		next.ServeHTTP(w, r)
	})
}

func (am *AuthMiddleware) validateMobileJWT(tokenString string) (*models.UserContext, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, apperr.Newf(apperr.InvalidToken, "unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(am.authConfig.JWTSecretKey), nil
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
