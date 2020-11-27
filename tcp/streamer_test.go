package tcp

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewStreamerCreator(t *testing.T) {
	tt := []struct {
		name    string
		enabled bool

		address func(string) string

		afterCreate func(t *testing.T, s *Streamer)

		wantErr bool
	}{
		{
			name:    "pass",
			enabled: true,

			address: func(srv string) string {
				return srv
			},

			afterCreate: func(t *testing.T, s *Streamer) {
				assert.NotNil(t, s)

				if err := s.Close(); err != nil {
					t.Fatal(err)
				}
			},
		},
		{
			name:    "create error",
			enabled: true,

			address: func(_ string) string {
				return "256.256.256.256"
			},

			afterCreate: func(t *testing.T, s *Streamer) {
				assert.Nil(t, s)
			},

			wantErr: true,
		},
	}

	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			if !test.enabled {
				t.SkipNow()
			}

			srv := httptest.NewServer(http.NewServeMux())
			defer srv.Close()

			s, err := newStreamer(test.address(srv.Listener.Addr().String()))
			if (err != nil) != test.wantErr {
				t.Error(err)
				t.FailNow()
			}

			test.afterCreate(t, s)
		})
	}
}

type MockConn struct {
	mock.Mock
}

func (c *MockConn) Write(p []byte) (n int, err error) {
	args := c.Called(p)

	return args.Get(0).(int), args.Error(1)
}

func (c *MockConn) Close() (err error) {
	return c.Called().Error(0)
}

func TestStreamer_Write(t *testing.T) {
	tt := []struct {
		name    string
		enabled bool

		conn        *MockConn
		writeInput  []interface{}
		writeOutput []interface{}

		wantErr bool

		expectedN int
	}{
		{
			name:    "pass",
			enabled: true,

			conn: new(MockConn),
			writeInput: []interface{}{
				[]byte(`{}`),
			},
			writeOutput: []interface{}{
				len([]byte(`{}`)),
				(error)(nil),
			},

			expectedN: len([]byte(`{}`)),
		},
		{
			name:    "write error",
			enabled: true,

			conn: new(MockConn),
			writeInput: []interface{}{
				[]byte(`{}`),
			},
			writeOutput: []interface{}{
				0,
				errors.New("write error"),
			},

			wantErr: true,
		},
	}

	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			if !test.enabled {
				t.SkipNow()
			}

			test.conn.
				On("Write", test.writeInput...).
				Return(test.writeOutput...)

			s := &Streamer{conn: test.conn}
			actualN, err := s.Write([]byte(`{}`))
			if (err != nil) != test.wantErr {
				t.Error(err)
				t.FailNow()
			}

			assert.Equal(t, test.expectedN, actualN)
		})
	}
}

func TestStreamer_Close(t *testing.T) {
	tt := []struct {
		name    string
		enabled bool

		conn        *MockConn
		closeOutput []interface{}

		wantErr bool
	}{
		{
			name:    "pass",
			enabled: true,

			conn: new(MockConn),
			closeOutput: []interface{}{
				(error)(nil),
			},
		},
		{
			name:    "close error",
			enabled: true,

			conn: new(MockConn),
			closeOutput: []interface{}{
				errors.New("close error"),
			},

			wantErr: true,
		},
	}

	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			if !test.enabled {
				t.SkipNow()
			}

			test.conn.
				On("Close").
				Return(test.closeOutput...)

			s := &Streamer{conn: test.conn}
			if err := s.Close(); (err != nil) != test.wantErr {
				t.Error(err)
				t.FailNow()
			}
		})
	}
}
