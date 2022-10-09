package jsonrpc

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/iancoleman/strcase"
)

func jsonrpcError(c *gin.Context, code int, message string, data string, id *float64) {
	c.JSON(http.StatusOK, gin.H{
		"result":  nil,
		"jsonrpc": "2.0",
		"error": gin.H{
			"code":    code,
			"message": message,
			"data":    data,
		},
		"id": id,
	})
}

func HandleRequest(c *gin.Context, api interface{}) {
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

	data := make(map[string]interface{})
	err = json.Unmarshal(body, &data)
	if err != nil {
		jsonrpcError(c, -32700, "Parse error", "Error parsing json request", nil)
		return
	}

	id, ok := data["id"].(float64)
	if !ok {
		jsonrpcError(c, -32600, "Invalid Request", "No or invalid 'id' in request", nil)
		return
	}

	if data["jsonrpc"] != "2.0" {
		jsonrpcError(c, -32600, "Invalid Request", "Version of jsonrpc is not 2.0", &id)
		return
	}

	method, ok := data["method"].(string)
	if !ok || method != strcase.ToSnake(method) {
		jsonrpcError(c, -32600, "Invalid Request", "No or invalid 'method' in request", &id)
		return
	}

	params, ok := data["params"].([]interface{})
	if !ok {
		jsonrpcError(c, -32602, "Invalid params", "No or invalid 'params' in request", &id)
		return
	}

	call := reflect.ValueOf(api).MethodByName(strcase.ToCamel(method))
	if !call.IsValid() {
		jsonrpcError(c, -32601, "Method not found", "Method not found", &id)
		return
	}

	if call.Type().NumIn() != len(params) {
		jsonrpcError(c, -32602, "Invalid params", "Invalid number of params", &id)
		return
	}

	args := make([]reflect.Value, len(params))
	for i, arg := range params {
		switch call.Type().In(i).Kind() {
		case reflect.Float32:
			val, ok := arg.(float32)
			if !ok {
				jsonrpcError(c, -32602, "Invalid params", fmt.Sprintf("Param [%d] can't be converted to %v", i, call.Type().In(i).String()), &id)
				return
			}
			args[i] = reflect.ValueOf(val)

		case reflect.Float64:
			val, ok := arg.(float64)
			if !ok {
				jsonrpcError(c, -32602, "Invalid params", fmt.Sprintf("Param [%d] can't be converted to %v", i, call.Type().In(i).String()), &id)
				return
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
				jsonrpcError(c, -32602, "Invalid params", fmt.Sprintf("Param [%d] can't be converted to %v", i, call.Type().In(i).String()), &id)
				return
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
				jsonrpcError(c, -32602, "Invalid params", fmt.Sprintf("Param [%d] can't be converted to %v", i, call.Type().In(i).String()), &id)
				return
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
				jsonrpcError(c, -32602, "Invalid params", fmt.Sprintf("Param [%d] can't be converted to %v", i, call.Type().In(i).String()), &id)
				return
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
				jsonrpcError(c, -32602, "Invalid params", fmt.Sprintf("Param [%d] can't be converted to %v", i, call.Type().In(i).String()), &id)
				return
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
				jsonrpcError(c, -32602, "Invalid params", fmt.Sprintf("Param [%d] can't be converted to %v", i, call.Type().In(i).String()), &id)
				return
			}
			args[i] = reflect.ValueOf(val)

		case reflect.Interface:
			args[i] = reflect.ValueOf(arg)

		case reflect.Map:
			val, ok := arg.(map[interface{}]interface{})
			if !ok {
				jsonrpcError(c, -32602, "Invalid params", fmt.Sprintf("Param [%d] can't be converted to %v", i, call.Type().In(i).String()), &id)
				return
			}
			args[i] = reflect.ValueOf(val)

		case reflect.Slice:
			val, ok := arg.([]interface{})
			if !ok {
				jsonrpcError(c, -32602, "Invalid params", fmt.Sprintf("Param [%d] can't be converted to %v", i, call.Type().In(i).String()), &id)
				return
			}
			args[i] = reflect.ValueOf(val)

		case reflect.String:
			val, ok := arg.(string)
			if !ok {
				jsonrpcError(c, -32602, "Invalid params", fmt.Sprintf("Param [%d] can't be converted to %v", i, call.Type().In(i).String()), &id)
				return
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
				jsonrpcError(c, -32602, "Invalid params", fmt.Sprintf("Param [%d] can't be converted to %v", i, call.Type().In(i).String()), &id)
				return
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
				jsonrpcError(c, -32602, "Invalid params", fmt.Sprintf("Param [%d] can't be converted to %v", i, call.Type().In(i).String()), &id)
				return
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
				jsonrpcError(c, -32602, "Invalid params", fmt.Sprintf("Param [%d] can't be converted to %v", i, call.Type().In(i).String()), &id)
				return
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
				jsonrpcError(c, -32602, "Invalid params", fmt.Sprintf("Param [%d] can't be converted to %v", i, call.Type().In(i).String()), &id)
				return
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
				jsonrpcError(c, -32602, "Invalid params", fmt.Sprintf("Param [%d] can't be converted to %v", i, call.Type().In(i).String()), &id)
				return
			}
			args[i] = reflect.ValueOf(val)

		default:
			if !ok {
				jsonrpcError(c, -32603, "Internal error", "Invalid method definition", &id)
				return
			}
		}

	}

	result := call.Call(args)
	if len(result) > 0 {
		c.JSON(http.StatusOK, gin.H{
			"result":  result[0].Interface(),
			"jsonrpc": "2.0",
			"id":      id,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"result":  nil,
			"jsonrpc": "2.0",
			"id":      id,
		})
	}
}
