package template

import (
	"encoding/json"
	"time"

	"gopkg.in/guregu/null.v3"

	"go.k6.io/k6/lib/types"
)

// Config is the config for the kafka collector
type Config struct {
	Address      null.String        `json:"address" envconfig:"K6_TEMPLATE_ADDRESS"`
	PushInterval types.NullDuration `json:"push_interval" envconfig:"K6_TEMPLATE_PUSH_INTERVAL"`
}

// NewConfig creates a new Config instance with default values for some fields.
func NewConfig() Config {
	return Config{
		Address:      null.NewString("template", false),
		PushInterval: types.NewNullDuration(1*time.Second, false),
	}
}

// GetConsolidatedConfig combines {default config values + JSON config +
// environment vars + arg config values}, and returns the final result.
func GetConsolidatedConfig(_ json.RawMessage, env map[string]string, argline string) (Config, error) {
	result := NewConfig()
	arglines, err := parseArgLine(argline, "address")
	if err != nil {
		return result, err
	}
	ch := NewHelper(
		NewMapLookupHelper(env, "error parsing environment variable '%s': %w"),
		NewMapLookupHelper(arglines, "error parsing argument line key '%s': %w"),
	)
	ch.GetNullString(&result.Address, []string{"K6_TEMPLATE_ADDRESS"}, []string{"address"})

	// TODO this can also be a helper  but is here for illustrative purposes
	ch.GetWithCallback(func(val string) error {
		pushInterval, err := time.ParseDuration(val)
		if err != nil {
			return err
		}
		result.PushInterval = types.NewNullDuration(pushInterval, true)

		return nil
	}, []string{"K6_TEMPLATE_PUSH_INTERVAL"}, []string{"push_interval"})

	return result, ch.GetConcatenatedError("\n")
}
