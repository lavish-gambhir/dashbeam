module github.com/lavish-gambhir/dashbeam/services/auth

go 1.23.3

require (
	github.com/golang-jwt/jwt/v5 v5.2.1
	github.com/lavish-gambhir/dashbeam/shared v0.0.0
	github.com/lavish-gambhir/dashbeam/pkg/apperr v0.0.0
)

// Local workspace replacements
replace github.com/lavish-gambhir/dashbeam/shared => ../../shared

replace github.com/lavish-gambhir/dashbeam/pkg/apperr => ../../pkg/apperr
