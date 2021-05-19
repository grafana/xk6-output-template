/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2016 Load Impact
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package template

import (
	"encoding/json"
	"fmt"
	"strings"
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
		Address:      null.StringFrom("template"),
		PushInterval: types.NullDurationFrom(1 * time.Second),
	}
}

func (c Config) Apply(cfg Config) Config {
	if cfg.Address.Valid {
		c.Address = cfg.Address
	}
	if cfg.PushInterval.Valid {
		c.PushInterval = cfg.PushInterval
	}
	return c
}

// ParseArg takes an arg string and converts it to a config
func ParseArg(arg string) (Config, error) {
	c := Config{}

	if !strings.Contains(arg, "=") {
		c.Address = null.StringFrom(arg)
		return c, nil
	}

	pairs := strings.Split(arg, ",")
	for _, pair := range pairs {
		r := strings.SplitN(pair, "=", 2)
		if len(r) != 2 {
			return c, fmt.Errorf("couldn't parse %q as argument for TEMPLATE output", arg)
		}
		switch r[0] {
		case "address":
			err := c.Address.UnmarshalText([]byte(r[1]))
			if err != nil {
				return c, err
			}
		case "push_interval":
			err := c.PushInterval.UnmarshalText([]byte(r[1]))
			if err != nil {
				return c, err
			}
		default:
			return c, fmt.Errorf("unknown key %q as argument for TEMPLATE output", r[0])
		}
	}

	return c, nil
}

// GetConsolidatedConfig combines {default config values + JSON config +
// environment vars + arg config values}, and returns the final result.
func GetConsolidatedConfig(jsonRawConf json.RawMessage, env map[string]string, arg string) (Config, error) {
	result := NewConfig()
	if jsonRawConf != nil {
		jsonConf := Config{}
		if err := json.Unmarshal(jsonRawConf, &jsonConf); err != nil {
			return result, err
		}
		result = result.Apply(jsonConf)
	}

	envConfig := Config{}
	for k, v := range env {
		switch k {
		case "K6_TEMPLATE_PUSH_INTERVAL":
			err := envConfig.PushInterval.UnmarshalText([]byte(v))
			if err != nil {
				return result, err
			}
		case "K6_TEMPLATE_ADDRESS":
			envConfig.Address = null.NewString(v, true)
		}
	}
	result = result.Apply(envConfig)

	if arg != "" {
		urlConf, err := ParseArg(arg)
		if err != nil {
			return result, err
		}
		result = result.Apply(urlConf)
	}

	return result, nil
}
