module github.com/lavish-gambhir/dashbeam

go 1.23.3

require (
	github.com/golang-jwt/jwt/v5 v5.2.2
	github.com/golang-migrate/migrate/v4 v4.18.3
	github.com/google/uuid v1.6.0
	github.com/jackc/pgx/v5 v5.7.5
	github.com/joho/godotenv v1.5.1
	github.com/lavish-gambhir/dashbeam/pkg/logger v0.0.0-00010101000000-000000000000
	github.com/lavish-gambhir/dashbeam/services/analytics v0.0.0-00010101000000-000000000000
	github.com/lavish-gambhir/dashbeam/services/auth v0.0.0-00010101000000-000000000000
	github.com/lavish-gambhir/dashbeam/services/ingestion v0.0.0-00010101000000-000000000000
	github.com/lavish-gambhir/dashbeam/shared v0.0.0
)

require (
	github.com/ClickHouse/ch-go v0.66.0 // indirect
	github.com/ClickHouse/clickhouse-go/v2 v2.36.0 // indirect
	github.com/andybalholm/brotli v1.1.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/fsnotify/fsnotify v1.8.0 // indirect
	github.com/go-faster/city v1.0.1 // indirect
	github.com/go-faster/errors v0.7.1 // indirect
	github.com/go-viper/mapstructure/v2 v2.2.1 // indirect
	github.com/google/go-github/v39 v39.2.0 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/lavish-gambhir/dashbeam/pkg/apperr v0.0.0 // indirect
	github.com/lavish-gambhir/dashbeam/pkg/utils v0.0.0-20250614162017-202e225a4254 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/matoous/go-nanoid/v2 v2.1.0 // indirect
	github.com/paulmach/orb v0.11.1 // indirect
	github.com/pelletier/go-toml/v2 v2.2.3 // indirect
	github.com/pierrec/lz4/v4 v4.1.22 // indirect
	github.com/redis/go-redis/v9 v9.10.0 // indirect
	github.com/sagikazarmark/locafero v0.7.0 // indirect
	github.com/segmentio/asm v1.2.0 // indirect
	github.com/shopspring/decimal v1.4.0 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	github.com/spf13/afero v1.12.0 // indirect
	github.com/spf13/cast v1.7.1 // indirect
	github.com/spf13/pflag v1.0.6 // indirect
	github.com/spf13/viper v1.20.1 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	go.opentelemetry.io/otel v1.36.0 // indirect
	go.opentelemetry.io/otel/trace v1.36.0 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/crypto v0.38.0 // indirect
	golang.org/x/oauth2 v0.25.0 // indirect
	golang.org/x/sync v0.14.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.25.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/lavish-gambhir/dashbeam/pkg/apperr => ./pkg/apperr

replace github.com/lavish-gambhir/dashbeam/pkg/logger => ./pkg/logger

replace github.com/lavish-gambhir/dashbeam/pkg/utils => ./pkg/utils

replace github.com/lavish-gambhir/dashbeam/shared => ./shared

replace github.com/lavish-gambhir/dashbeam/services/auth => ./services/auth

replace github.com/lavish-gambhir/dashbeam/services/quiz => ./services/quiz

replace github.com/lavish-gambhir/dashbeam/services/ingestion => ./services/ingestion

replace github.com/lavish-gambhir/dashbeam/services/analytics => ./services/analytics

replace github.com/lavish-gambhir/dashbeam/services/reporting => ./services/reporting
