module github.com/lavish-gambhir/dashbeam/services/auth

go 1.23.3

require github.com/lavish-gambhir/dashbeam/shared v0.0.0

require (
	github.com/golang-jwt/jwt/v5 v5.2.2 // indirect
	github.com/google/uuid v1.6.0 // indirect
)

// Local workspace replacements
replace github.com/lavish-gambhir/dashbeam/shared => ../../shared
