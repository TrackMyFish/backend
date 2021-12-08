package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	testCases := []struct {
		desc        string
		conf        Config
		expectedErr string
	}{
		{
			desc:        "Empty host should return error",
			conf:        Config{},
			expectedErr: "host not defined",
		},
		{
			desc:        "Empty port should return error",
			conf:        Config{Host: "host"},
			expectedErr: "port not defined",
		},
		{
			desc:        "Empty username should return error",
			conf:        Config{Host: "host", Port: "1111"},
			expectedErr: "user not defined",
		},
		{
			desc:        "Empty password should return error",
			conf:        Config{Host: "host", Port: "1111", Username: "user"},
			expectedErr: "password not defined",
		},
		{
			desc:        "Empty database should return error",
			conf:        Config{Host: "host", Port: "1111", Username: "user", Password: "password"},
			expectedErr: "database not defined",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			_, err := New(tC.conf)
			assert.EqualError(t, err, tC.expectedErr)
		})
	}
}
