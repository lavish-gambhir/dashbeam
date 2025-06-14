module github.com/lavish-gambhir/dashbeam/services/ingestion

go 1.23.3

require (
	github.com/google/uuid v1.6.0
	github.com/lavish-gambhir/dashbeam/pkg/apperr v0.0.0
	github.com/lavish-gambhir/dashbeam/pkg/utils v0.0.0-20250614071328-e3be77b9160d
	github.com/lavish-gambhir/dashbeam/shared v0.0.0
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/fsnotify/fsnotify v1.8.0 // indirect
	github.com/go-viper/mapstructure/v2 v2.2.1 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/pelletier/go-toml/v2 v2.2.3 // indirect
	github.com/redis/go-redis/v9 v9.10.0 // indirect
	github.com/sagikazarmark/locafero v0.7.0 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	github.com/spf13/afero v1.12.0 // indirect
	github.com/spf13/cast v1.7.1 // indirect
	github.com/spf13/pflag v1.0.6 // indirect
	github.com/spf13/viper v1.20.1 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.9.0 // indirect
	golang.org/x/sys v0.32.0 // indirect
	golang.org/x/text v0.24.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

// Local workspace replacements
replace github.com/lavish-gambhir/dashbeam/shared => ../../shared

replace github.com/lavish-gambhir/dashbeam/pkg => ../../pkg

replace github.com/lavish-gambhir/dashbeam/pkg/apperr => ../../pkg/apperr
