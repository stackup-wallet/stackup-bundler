// Package jsonrpc implements Gin middleware for handling JSON-RPC requests via HTTP.
package jsonrpc

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/stackup-wallet/stackup-bundler/pkg/errors"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	optionalTypePrefix = "optional_"
)

func formatConversionErrMsg(i int, call *reflect.Value) string {
	s, _ := strings.CutPrefix(call.Type().In(i).Name(), optionalTypePrefix)
	return fmt.Sprintf("Param [%d] can't be converted to %s", i, s)
}

func jsonrpcError(c *gin.Context, code int, message string, data any, id any) {
	c.JSON(http.StatusOK, gin.H{
		"jsonrpc": "2.0",
		"error": gin.H{
			"code":    code,
			"message": message,
			"data":    data,
		},
		"id": id,
	})
	c.Abort()
}

// parseRequestId checks if the JSON-RPC request contains an id field that is either NULL, Number, or String.
func parseRequestId(data map[string]any) (any, bool) {
	id, ok := data["id"]
	_, isFloat64 := id.(float64)
	_, isStr := id.(string)

	if ok && (id == nil || isFloat64 || isStr) {
		return id, true
	}
	return nil, false
}

// hasOptionalInput checks if the API method has defined an optional final input:
//  1. The input must start with the "optional_" prefix in its name.
//  2. The input must be of kind Map.
func hasOptionalInput(numIn int, call *reflect.Value) bool {
	return numIn > 0 &&
		strings.HasPrefix(call.Type().In(numIn-1).Name(), optionalTypePrefix) &&
		call.Type().In(numIn-1).Kind() == reflect.Map
}

// hasValidParamLength checks if the number of parameters in the request is correct:
//  1. Ok if the number of params equals number of method inputs.
//  2. Ok if optional input is defined and number of params is one less the number of method inputs.
func hasValidParamLength(numParams, numIn int, hasOptional bool) bool {
	return numParams == numIn || (hasOptional && numParams == numIn-1)
}

// isOptionalParamUndefined checks if the optional input has been left unset in the request.
func isOptionalParamUndefined(numParams, numIn int, hasOptional bool) bool {
	return hasOptional && numParams == numIn-1
}

// handleRequest includes the core logic for parsing individual JSON-RPC requests and returning its id,
// result, and success flag.
func handleRequest(api interface{}, c *gin.Context, data map[string]any) (id any, result any, success bool) {
	id, ok := parseRequestId(data)
	if !ok {
		jsonrpcError(c, -32600, "Invalid Request", "No or invalid 'id' in request", nil)
		return id, nil, false
	}

	if data["jsonrpc"] != "2.0" {
		jsonrpcError(c, -32600, "Invalid Request", "Version of jsonrpc is not 2.0", &id)
		return id, nil, false
	}

	method, ok := data["method"].(string)
	if !ok {
		jsonrpcError(c, -32600, "Invalid Request", "No or invalid 'method' in request", &id)
		return id, nil, false
	}

	params, ok := data["params"].([]interface{})
	if !ok {
		jsonrpcError(c, -32602, "Invalid params", "No or invalid 'params' in request", &id)
		return id, nil, false
	}

	call := reflect.ValueOf(api).MethodByName(cases.Title(language.Und, cases.NoLower).String(method))
	if !call.IsValid() {
		jsonrpcError(c, -32601, "Method not found", "Method not found", &id)
		return id, nil, false
	}

	numIn := call.Type().NumIn()
	numParams := len(params)
	hasOptional := hasOptionalInput(numIn, &call)
	if !hasValidParamLength(numParams, numIn, hasOptional) {
		jsonrpcError(c, -32602, "Invalid params", "Invalid number of params", &id)
		return id, nil, false
	}
	if isOptionalParamUndefined(numParams, numIn, hasOptional) {
		params = append(params, map[string]any{})
		numParams++
	}

	args := make([]reflect.Value, numParams)
	for i, arg := range params {
		switch call.Type().In(i).Kind() {
		case reflect.Float32:
			val, ok := arg.(float32)
			if !ok {
				jsonrpcError(
					c,
					-32602,
					"Invalid params",
					formatConversionErrMsg(i, &call),
					&id,
				)
				return id, nil, false
			}
			args[i] = reflect.ValueOf(val)

		case reflect.Float64:
			val, ok := arg.(float64)
			if !ok {
				jsonrpcError(
					c,
					-32602,
					"Invalid params",
					formatConversionErrMsg(i, &call),
					&id,
				)
				return id, nil, false
			}
			args[i] = reflect.ValueOf(val)

		case reflect.Int:
			val, ok := arg.(int)
			if !ok {
				var fval float64
				fval, ok = arg.(float64)
				if ok {
					val = int(fval)
				}
			}

			if !ok {
				jsonrpcError(
					c,
					-32602,
					"Invalid params",
					formatConversionErrMsg(i, &call),
					&id,
				)
				return id, nil, false
			}
			args[i] = reflect.ValueOf(val)

		case reflect.Int8:
			val, ok := arg.(int8)
			if !ok {
				var fval float64
				fval, ok = arg.(float64)
				if ok {
					val = int8(fval)
				}
			}
			if !ok {
				jsonrpcError(
					c,
					-32602,
					"Invalid params",
					formatConversionErrMsg(i, &call),
					&id,
				)
				return id, nil, false
			}
			args[i] = reflect.ValueOf(val)

		case reflect.Int16:
			val, ok := arg.(int16)
			if !ok {
				var fval float64
				fval, ok = arg.(float64)
				if ok {
					val = int16(fval)
				}
			}
			if !ok {
				jsonrpcError(
					c,
					-32602,
					"Invalid params",
					formatConversionErrMsg(i, &call),
					&id,
				)
				return id, nil, false
			}
			args[i] = reflect.ValueOf(val)

		case reflect.Int32:
			val, ok := arg.(int32)
			if !ok {
				var fval float64
				fval, ok = arg.(float64)
				if ok {
					val = int32(fval)
				}
			}
			if !ok {
				jsonrpcError(
					c,
					-32602,
					"Invalid params",
					formatConversionErrMsg(i, &call),
					&id,
				)
				return id, nil, false
			}
			args[i] = reflect.ValueOf(val)

		case reflect.Int64:
			val, ok := arg.(int64)
			if !ok {
				var fval float64
				fval, ok = arg.(float64)
				if ok {
					val = int64(fval)
				}
			}
			if !ok {
				jsonrpcError(
					c,
					-32602,
					"Invalid params",
					formatConversionErrMsg(i, &call),
					&id,
				)
				return id, nil, false
			}
			args[i] = reflect.ValueOf(val)

		case reflect.Interface:
			args[i] = reflect.ValueOf(arg)

		case reflect.Map:
			val, ok := arg.(map[string]any)
			if !ok {
				jsonrpcError(
					c,
					-32602,
					"Invalid params",
					formatConversionErrMsg(i, &call),
					&id,
				)
				return id, nil, false
			}
			args[i] = reflect.ValueOf(val)

		case reflect.Slice:
			val, ok := arg.([]interface{})
			if !ok {
				jsonrpcError(
					c,
					-32602,
					"Invalid params",
					formatConversionErrMsg(i, &call),
					&id,
				)
				return id, nil, false
			}
			args[i] = reflect.ValueOf(val)

		case reflect.String:
			val, ok := arg.(string)
			if !ok {
				jsonrpcError(
					c,
					-32602,
					"Invalid params",
					formatConversionErrMsg(i, &call),
					&id,
				)
				return id, nil, false
			}
			args[i] = reflect.ValueOf(val)

		case reflect.Uint:
			val, ok := arg.(uint)
			if !ok {
				var fval float64
				fval, ok = arg.(float64)
				if ok {
					val = uint(fval)
				}
			}
			if !ok {
				jsonrpcError(
					c,
					-32602,
					"Invalid params",
					formatConversionErrMsg(i, &call),
					&id,
				)
				return id, nil, false
			}
			args[i] = reflect.ValueOf(val)

		case reflect.Uint8:
			val, ok := arg.(uint8)
			if !ok {
				var fval float64
				fval, ok = arg.(float64)
				if ok {
					val = uint8(fval)
				}
			}
			if !ok {
				jsonrpcError(
					c,
					-32602,
					"Invalid params",
					formatConversionErrMsg(i, &call),
					&id,
				)
				return id, nil, false
			}
			args[i] = reflect.ValueOf(val)

		case reflect.Uint16:
			val, ok := arg.(uint16)
			if !ok {
				var fval float64
				fval, ok = arg.(float64)
				if ok {
					val = uint16(fval)
				}
			}
			if !ok {
				jsonrpcError(
					c,
					-32602,
					"Invalid params",
					formatConversionErrMsg(i, &call),
					&id,
				)
				return id, nil, false
			}
			args[i] = reflect.ValueOf(val)

		case reflect.Uint32:
			val, ok := arg.(uint32)
			if !ok {
				var fval float64
				fval, ok = arg.(float64)
				if ok {
					val = uint32(fval)
				}
			}
			if !ok {
				jsonrpcError(
					c,
					-32602,
					"Invalid params",
					formatConversionErrMsg(i, &call),
					&id,
				)
				return id, nil, false
			}
			args[i] = reflect.ValueOf(val)

		case reflect.Uint64:
			val, ok := arg.(uint64)
			if !ok {
				var fval float64
				fval, ok = arg.(float64)
				if ok {
					val = uint64(fval)
				}
			}
			if !ok {
				jsonrpcError(
					c,
					-32602,
					"Invalid params",
					formatConversionErrMsg(i, &call),
					&id,
				)
				return id, nil, false
			}
			args[i] = reflect.ValueOf(val)

		default:
			if !ok {
				jsonrpcError(c, -32603, "Internal error", "Invalid method definition", &id)
				return id, nil, false
			}
		}
	}

	c.Set("json-rpc-request", data)
	value := call.Call(args)
	if err, ok := value[len(value)-1].Interface().(error); ok && err != nil {
		rpcErr, ok := err.(*errors.RPCError)

		if ok {
			jsonrpcError(c, rpcErr.Code(), rpcErr.Error(), rpcErr.Data(), &id)
		} else {
			jsonrpcError(c, -32601, err.Error(), err.Error(), &id)
		}
		return id, nil, false
	} else if len(value) > 0 {
		return id, value[0].Interface(), true
	} else {
		return id, nil, true
	}
}

// Controller returns a custom Gin middleware that handles incoming JSON-RPC requests via HTTP. It maps the
// RPC method name to struct methods on the given api. For example, if the RPC request has the method field
// set to "namespace_methodName" then the controller will make a call to api.Namespace_methodName with the
// params spread as arguments.
//
// If request is valid it will also set the data on the Gin context with the key "json-rpc-request".
//
// NOTE: For batched requests in the current version, "json-rpc-request" on the Gin context contains only the
// last request in the array.
func Controller(api interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method != "POST" {
			jsonrpcError(c, -32700, "Parse error", "POST method excepted", nil)
			return
		}

		if c.Request.Body == nil {
			jsonrpcError(c, -32700, "Parse error", "No POST data", nil)
			return
		}

		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			jsonrpcError(c, -32700, "Parse error", "Error while reading request body", nil)
			return
		}

		data := make(map[string]any)
		err = json.Unmarshal(body, &data)
		if err != nil {
			var batch []map[string]any
			err = json.Unmarshal(body, &batch)
			if err != nil {
				jsonrpcError(c, -32700, "Parse error", "Error parsing json request", nil)
				return
			}

			var result []gin.H
			for _, data := range batch {
				id, res, success := handleRequest(api, c, data)
				if !success {
					return
				}

				result = append(result, gin.H{
					"jsonrpc": "2.0",
					"id":      id,
					"result":  res,
				})
			}
			c.JSON(http.StatusOK, result)
		} else if id, res, success := handleRequest(api, c, data); success {
			c.JSON(http.StatusOK, gin.H{
				"jsonrpc": "2.0",
				"id":      id,
				"result":  res,
			})
		}
	}
}
