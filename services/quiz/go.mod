module github.com/lavish-gambhir/dashbeam/services/quiz

go 1.23.3

require (
	github.com/google/uuid v1.6.0
	github.com/gorilla/websocket v1.5.3
	github.com/lavish-gambhir/dashbeam/shared v0.0.0
	github.com/lavish-gambhir/dashbeam/pkg/apperr v0.0.0
)

// Local workspace replacements
replace github.com/lavish-gambhir/dashbeam/shared => ../../shared

replace github.com/lavish-gambhir/dashbeam/pkg/apperr => ../../pkg/apperr
