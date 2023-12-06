package testutils

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/metachris/flashbotsrpc"
)

func BadBuilderRpcMock() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		res := &flashbotsrpc.RelayErrorResponse{
			Error: "Mock upstream builder error",
		}
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(res); err != nil {
			panic(err)
		}
	}))
}
