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
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v3"

	"go.k6.io/k6/lib/types"
)

func TestConfigParseArg(t *testing.T) {
	c, err := ParseArg("example.com")
	assert.Nil(t, err)
	assert.Equal(t, null.StringFrom("example.com"), c.Address)
	assert.False(t, c.PushInterval.Valid)
	assert.EqualValues(t, 0, c.PushInterval.Duration)

	c, err = ParseArg("address=example.com,push_interval=2s")
	assert.Nil(t, err)
	assert.Equal(t, null.StringFrom("example.com"), c.Address)
	assert.True(t, c.PushInterval.Valid)
	assert.EqualValues(t, time.Second*2, c.PushInterval.Duration)
}

func TestConsolidatedConfig(t *testing.T) {
	t.Parallel()
	// TODO: add more cases
	testCases := map[string]struct {
		jsonRaw json.RawMessage
		env     map[string]string
		arg     string
		config  Config
		err     string
	}{
		"default": {
			config: Config{
				Address:      null.StringFrom("template"),
				PushInterval: types.NullDurationFrom(1 * time.Second),
			},
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			// hacks around env not actually being taken into account
			os.Clearenv()
			defer os.Clearenv()
			for k, v := range testCase.env {
				require.NoError(t, os.Setenv(k, v))
			}

			config, err := GetConsolidatedConfig(testCase.jsonRaw, testCase.env, testCase.arg)
			if testCase.err != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), testCase.err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, testCase.config, config)
		})
	}
}
