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
