package fishbase

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

type MockClient struct {
	response *http.Response
	err      error
}

func (c MockClient) Do(req *http.Request) (*http.Response, error) {
	// If we've been asked to return something specific, return it
	if c.response != nil || c.err != nil {
		return c.response, c.err
	}

	return &http.Response{}, nil
}

func TestGetHeartbeat(t *testing.T) {
	mc := MockClient{
		response: &http.Response{Body: io.NopCloser(bytes.NewReader([]byte("foo")))},
		err:      nil,
	}

	t.Run("Given a request is made to GetHeartbeat", func(t *testing.T) {
		t.Run("When an error is returned from the HTTP request", func(t *testing.T) {
			t.Run("Then an error should be returned", func(t *testing.T) {
				mc.err = errors.New("some error")
				c := New(mc)

				_, err := c.GetHeartbeat()

				assert.EqualError(t, err, "unable to query fishbase: some error")
			})
		})
		t.Run("When a 2xx error code is returned", func(t *testing.T) {
			t.Run("Then Operational should be retured", func(t *testing.T) {
				mc.err = nil
				mc.response.StatusCode = 200
				c := New(mc)

				hb, err := c.GetHeartbeat()
				assert.NoError(t, err)

				assert.Equal(t, HeartbeatStatusOperational, hb.Status)
			})
		})
		t.Run("When a 3xx error code is returned", func(t *testing.T) {
			t.Run("Then Operational should be returned", func(t *testing.T) {
				mc.err = nil
				mc.response.StatusCode = 300
				c := New(mc)

				hb, err := c.GetHeartbeat()
				assert.NoError(t, err)

				assert.Equal(t, HeartbeatStatusOperational, hb.Status)
			})
		})
		t.Run("When a 4xx error code is returned", func(t *testing.T) {
			t.Run("Then Down should be returned", func(t *testing.T) {
				mc.err = nil
				mc.response.StatusCode = 400
				c := New(mc)

				hb, err := c.GetHeartbeat()
				assert.NoError(t, err)

				assert.Equal(t, HeartbeatStatusDown, hb.Status)
			})
		})
		t.Run("When a 5xx error code is returned", func(t *testing.T) {
			t.Run("Then Down should be returned", func(t *testing.T) {
				mc.err = nil
				mc.response.StatusCode = 500
				c := New(mc)

				hb, err := c.GetHeartbeat()
				assert.NoError(t, err)

				assert.Equal(t, HeartbeatStatusDown, hb.Status)
			})
		})
	})
}
