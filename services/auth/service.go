package auth

type Service interface {
	ValidateJWT(token string) (string, error)
}
