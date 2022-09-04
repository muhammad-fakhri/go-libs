package constant

import "errors"

type Environment string

const (
	EnvTest       Environment = "test"
	EnvStaging    Environment = "uat"
	EnvProduction Environment = "production"
)

func (e Environment) Validate() error {
	switch e {
	case EnvTest, EnvStaging, EnvProduction:
		return nil
	}
	return errors.New("invalid_environment")
}
