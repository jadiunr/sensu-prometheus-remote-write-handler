package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	corev2 "github.com/sensu/sensu-go/api/core/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStatusMetrics(t *testing.T) {
	assert := assert.New(t)
	event := corev2.FixtureEvent("entity1", "check1")
	event.Metrics = nil

	var apiStub = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		expectedBody := `check1,sensu_entity_name=entity1 status=0`
		assert.Contains(string(body), expectedBody)
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`{"ok": true}`))
		require.NoError(t, err)
	}))

	plugin.Endpoint = apiStub.URL
	plugin.IncludeEventStatus = true
}
