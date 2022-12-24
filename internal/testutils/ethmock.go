package testutils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
)

type mockReq struct {
	JsonRpc string  `json:"jsonrpc"`
	ID      float64 `json:"id"`
	Method  string  `json:"method"`
}

type mockRes struct {
	JsonRpc string  `json:"jsonrpc"`
	ID      float64 `json:"id"`
	Result  any     `json:"result"`
}

type MethodMocks map[string]any

// EthMock returns a httptest.Server for mocking the return value of a JSON-RPC method call to an Ethereum node.
func EthMock(mocks MethodMocks) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req mockReq
		json.NewDecoder(r.Body).Decode(&req)
		mock, ok := mocks[req.Method]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("method not in mocks: %s", req.Method)))
			return
		}

		res := &mockRes{
			JsonRpc: req.JsonRpc,
			ID:      req.ID,
			Result:  mock,
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(res)
	}))
}
