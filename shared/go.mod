module github.com/lavish-gambhir/dashbeam/shared

go 1.23.3

require (
	github.com/jackc/pgx/v5 v5.7.5
	github.com/lavish-gambhir/dashbeam/pkg/apperr v0.0.0
	github.com/matoous/go-nanoid/v2 v2.1.0
	go.opentelemetry.io/otel/trace v1.36.0
)

require (
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	go.opentelemetry.io/otel v1.36.0 // indirect
	golang.org/x/crypto v0.37.0 // indirect
	golang.org/x/sync v0.13.0 // indirect
	golang.org/x/text v0.24.0 // indirect
)

// Local workspace replacements
replace github.com/lavish-gambhir/dashbeam/pkg/apperr => ../pkg/apperr
