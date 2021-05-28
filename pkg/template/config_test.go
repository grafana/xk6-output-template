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
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v3"

	"go.k6.io/k6/lib/types"
)

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
			config: NewConfig(),
		},
		"custom argline - default argument": {
			arg: "something",
			config: Config{
				Address:      null.StringFrom("something"),
				PushInterval: types.NewNullDuration(1*time.Second, false),
			},
		},
		"custom argline - keyed": {
			arg: "address=something,push_interval=4s",
			config: Config{
				Address:      null.StringFrom("something"),
				PushInterval: types.NullDurationFrom(4 * time.Second),
			},
		},

		"precedence": {
			env: map[string]string{"K6_TEMPLATE_ADDRESS": "else", "K6_TEMPLATE_PUSH_INTERVAL": "4ms"},
			arg: "address=something",
			config: Config{
				Address:      null.StringFrom("something"),
				PushInterval: types.NullDurationFrom(4 * time.Millisecond),
			},
		},

		"early error": {
			env: map[string]string{"K6_TEMPLATE_ADDRESS": "else", "K6_TEMPLATE_PUSH_INTERVAL": "4something"},
			arg: "address=something",
			config: Config{
				Address:      null.StringFrom("something"),
				PushInterval: types.NewNullDuration(1*time.Second, false),
			},
			err: `error parsing environment variable 'K6_TEMPLATE_PUSH_INTERVAL': time: unknown unit "something" in duration "4something"`,
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
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
