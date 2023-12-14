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

func RpcMock(mocks MethodMocks) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req mockReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			panic(err)
		}
		mock, ok := mocks[req.Method]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			if _, err := w.Write([]byte(fmt.Sprintf("method not in mocks: %s", req.Method))); err != nil {
				panic(err)
			}
			return
		}

		res := &mockRes{
			JsonRpc: req.JsonRpc,
			ID:      req.ID,
			Result:  mock,
		}
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(res); err != nil {
			panic(err)
		}
	}))
}
