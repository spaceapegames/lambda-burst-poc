package loadtest

import (
	"github.com/spaceapegames/lambda-burst/api"
	"github.com/spaceapegames/lambda-burst/client"
	"github.com/tj/assert"
	"net/http/httptest"
	"testing"
)

func TestLoadTest_Run(t *testing.T) {

	/*
		Simulate 2 running servers.
		One will represent the ALB (srv1) and the other Lambda (srv2).
		srv1 is able to redirect to srv2
	*/

	srv2 := httptest.NewServer(api.NewServer(false, 8081, "", 0, false).Router)
	defer srv2.Close()

	srv := httptest.NewServer(api.NewServer(false, 8080, srv2.URL, 1000, false).Router)
	defer srv.Close()
	lt := NewLoadTest(5, 5, 5, client.NewClient(srv.URL))

	results, err := lt.Run()
	assert.NoError(t, err)
	_ = results
	assert.Equal(t, int64(125), results.Count)
}
