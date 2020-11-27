package cli

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDuration_MarshalJSON(t *testing.T) {
	tt := []struct {
		name    string
		enabled bool

		d Duration

		wantErr  bool
		expected []byte
	}{
		{
			name:    "pass",
			enabled: true,

			d: Duration(time.Second),

			expected: []byte(`"1s"`),
		},
	}

	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			if !test.enabled {
				t.SkipNow()
			}

			actual, err := test.d.MarshalJSON()
			if (err != nil) != test.wantErr {
				t.Error(err)
				t.FailNow()
			}

			assert.Equal(t, test.expected, actual)
		})
	}
}

func TestDuration_UnmarshalJSON(t *testing.T) {
	tt := []struct {
		name    string
		enabled bool

		in []byte

		expected Duration

		wantErr bool
	}{
		{
			name:    "unmarshal from number",
			enabled: true,

			in: []byte(strconv.FormatInt(int64(time.Second), 10)),

			expected: Duration(time.Second),
		},
		{
			name:    "unmarshal from string",
			enabled: true,

			in: []byte(`"1s"`),

			expected: Duration(time.Second),
		},
		{
			name:    "unmarshal error",
			enabled: true,

			wantErr: true,
		},
		{
			name:    "parse duration error",
			enabled: true,

			in: []byte(`"1"`),

			wantErr: true,
		},
		{
			name:    "invalid duration",
			enabled: true,

			in: []byte{0x7B, 0x7D}, // {}

			wantErr: true,
		},
	}

	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			if !test.enabled {
				t.SkipNow()
			}

			var actual Duration
			if err := actual.UnmarshalJSON(test.in); (err != nil) != test.wantErr {
				t.Error(err)
				t.FailNow()
			}

			assert.Equal(t, test.expected, actual)
		})
	}
}

func TestInput_Validate(t *testing.T) {
	tt := []struct {
		name    string
		enabled bool

		in Input

		wantErr  bool
		expected error
	}{
		{
			name:    "pass",
			enabled: true,

			in: Input{
				URL:    "http://127.0.0.1:8080/index.html",
				Output: "127.0.0.1:5000",
			},
		},
		{
			name:    "output as hostname-port",
			enabled: true,

			in: Input{
				URL:    "http://127.0.0.1:8080/index.html",
				Output: "mydomain.zone",
			},
		},
		{
			name:    "output as host-port",
			enabled: true,

			in: Input{
				URL:    "http://127.0.0.1:8080/index.html",
				Output: "127.0.0.1:5000",
			},
		},
		{
			name:    "empty url",
			enabled: true,

			in: Input{},

			wantErr:  true,
			expected: ErrInvalidURL,
		},
		{
			name:    "empty output",
			enabled: true,

			in: Input{
				URL: "http://127.0.0.1:8080/index.html",
			},
		},
		{
			name:    "invalid host-port output",
			enabled: true,

			in: Input{
				URL:    "http://127.0.0.1:8080/index.html",
				Output: "256.789.320.752:8135135368",
			},

			wantErr:  true,
			expected: ErrInvalidOutput,
		},
		{
			name:    "invalid hostname-port output",
			enabled: true,

			in: Input{
				URL:    "http://127.0.0.1:8080/index.html",
				Output: "gsfdsfdfd%@#fdfaf",
			},

			wantErr:  true,
			expected: ErrInvalidOutput,
		},
	}

	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			if !test.enabled {
				t.SkipNow()
			}

			actual := test.in.Validate()
			if (actual != nil) != test.wantErr {
				t.Error(actual)
				t.FailNow()
			}

			assert.Equal(t, test.expected, actual)
		})
	}
}
