/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2021 Load Impact
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
