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
	"errors"
	"fmt"
	"strings"

	"gopkg.in/guregu/null.v3"
)

// split argline taht looks like key1=value1,key2=value to a map[key]value with optional defaultKey for when arg is just
// `value` to make a map {defaultKey:value}
func parseArgLine(argline, defaultKey string) (map[string]string, error) {
	// this can be more complex and fix values with ',' but is probably way too much work
	if argline == "" {
		return map[string]string{}, nil
	}
	args := strings.Split(argline, ",")
	if len(args) == 1 && !strings.Contains(args[0], "=") {
		return map[string]string{defaultKey: argline}, nil
	}
	result := make(map[string]string, len(args))
	for _, arg := range args {
		k := strings.SplitN(arg, "=", 2)
		if len(k) != 2 {
			return nil, fmt.Errorf("bad argline value `%s`, expected `key=value`", arg)
		}
		result[k[0]] = k[1]

	}
	return result, nil
}

// ConfigHelper contains some useful methods for quickly parse configuration
// values. The helper functions lookup values in the provided LookupHelpers and will go through all of them so the
// latest value will be what will be used, but all of them will be parsed so if parsing one is problematic it will error
// out
type ConfigHelper struct {
	helpers []LookupHelper
	errors  []error
}

type LookupHelper interface {
	Lookup(string) (string, bool)
	HandleVarError(envVarName string, err error) error
}

// NewHelper returns a fully initialized ConfigHelper
func NewHelper(helpers ...LookupHelper) *ConfigHelper {
	return &ConfigHelper{helpers: helpers}
}

// GetNullString tries to read the requested string value from the lookup helpers.
// The namesPerLookup are in the same order as the lookups were in the NewHelper and can support multiple names with
// the first name winning, but the still going through all lookups
// TODO maybe the last name should win as well?
func (ch *ConfigHelper) GetNullString(result *null.String, namesPerLookup ...[]string) {
	ch.GetWithCallback(func(val string) error {
		*result = null.StringFrom(val)
		return nil
	}, namesPerLookup...)
}

// GetWithCallback doesn't save the parsed value directly, instead it calls
// the supplied callback function and handles any error it may return.
// The namesPerLookup are in the same order as the lookups were in the NewHelper and can support multiple names with
// the first name winning, but the still going through all lookups
// TODO maybe the last name should win as well?
func (ch *ConfigHelper) GetWithCallback(setter func(string) error, namesPerLookup ...[]string) {
	if len(ch.helpers) != len(namesPerLookup) {
		panic("bad") // TODO better
	}
	for i, helper := range ch.helpers {
		for _, name := range namesPerLookup[i] {
			if val, ok := helper.Lookup(name); ok {
				err := setter(val)
				if err != nil {
					ch.errors = append(ch.errors, helper.HandleVarError(name, err))
				}
				continue
			}
		}
	}
}

var _ LookupHelper = MapLookupHelper{}

type MapLookupHelper struct {
	m           map[string]string
	errTemplate string
}

func NewMapLookupHelper(m map[string]string, errTemplate string) MapLookupHelper {
	return MapLookupHelper{
		m: m, errTemplate: errTemplate,
	}
}

func (m MapLookupHelper) Lookup(name string) (string, bool) {
	s, ok := m.m[name]
	return s, ok
}

func (m MapLookupHelper) HandleVarError(name string, err error) error {
	return fmt.Errorf(m.errTemplate, name, err)
}

// GetErrors returns any errors that may have accumulated while parsing values.
func (ch *ConfigHelper) GetErrors() []error {
	return ch.errors
}

// GetConcatenatedError returns a single error out of any potentially
// accumulated errors, or nil if there were none.
func (ch *ConfigHelper) GetConcatenatedError(separator string) error {
	if len(ch.errors) == 0 {
		return nil
	}
	errStrings := make([]string, len(ch.errors))
	for i, e := range ch.errors {
		errStrings[i] = e.Error()
	}
	return errors.New(strings.Join(errStrings, separator))
}
