package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAddAndGetList(t *testing.T) {
	memCashe := NewMemoryStore()
	memCashe.add("3", "4")

	assert.NotEmpty(t, memCashe.getList())
	assert.Equal(t, memCashe.metrics["3"], "4")
}

func TestChangeAndGetList(t *testing.T) {
	memCashe := NewMemoryStore()
	memCashe.add("3", "4")
	memCashe.change("3", "6.123")

	assert.NotEmpty(t, memCashe.getList())
	assert.Equal(t, len(memCashe.getList()), 1)
	assert.Equal(t, memCashe.metrics["3"], "6.123")
}

func testRequest(t *testing.T, ts *httptest.Server, method,
	path string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}

func TestRouter(t *testing.T) {
	ts := httptest.NewServer(MetricRouter())
	defer ts.Close()

	var testTable = []struct {
		url    string
		want   string
		status int
	}{
		{"/update/gauge/metric1/14", "", http.StatusOK},
		{"/update/fasf/metric2/14.51", "", http.StatusBadRequest},
		{"/update/gauge/metric3/fsdfs", "", http.StatusBadRequest},
		{"/update/gauge/metric4/4.141", "", http.StatusOK},
		{"/update/counter/metric6/4", "", http.StatusOK},
		{"/update/counter/metric5/4.14", "", http.StatusBadRequest},
	}
	for _, v := range testTable {
		resp, get := testRequest(t, ts, "POST", v.url)
		assert.Equal(t, v.status, resp.StatusCode)
		assert.Equal(t, v.want, get)
	}
	resp, get := testRequest(t, ts, "GET", "/update/123/metric5")
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	assert.Equal(t, "404 page not found\n", get)
}
