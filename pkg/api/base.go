package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/wwqdrh/gokit/logger"
	"go.uber.org/zap"
)

type H map[string]interface{}

type OkResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

type ErrResponse struct {
	Code   int    `json:"code"`
	Msg    string `json:"msg"`
	Detail string `json:"detail"`
}

type StatusCode int

const (
	ServerOK StatusCode = iota
	ServerError
	ClientParamInvalid
)

func (c StatusCode) Code() int {
	switch c {
	case ServerOK:
		return 200
	case ServerError:
		return 500
	case ClientParamInvalid:
		return 400
	default:
		return 400
	}
}

func (c StatusCode) String() string {
	switch c {
	case ServerOK:
		return "ok"
	case ServerError:
		return "err"
	case ClientParamInvalid:
		return "invalid param"
	default:
		return "unknown"
	}
}

func ParseJSON(w http.ResponseWriter, r *http.Request, req interface{}) bool {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return false
	}

	if err := json.Unmarshal(body, &req); err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return false
	}

	return true
}

func EchoJSON(w http.ResponseWriter, code StatusCode, data interface{}) {
	w.WriteHeader(code.Code())
	if body, err := json.Marshal(OkResponse{
		Code: int(code),
		Msg:  code.String(),
		Data: data,
	}); err != nil {
		logger.DefaultLogger.Warn("json unmarashal err", zap.Error(err))

		w.Write([]byte("err"))
	} else {
		w.Write(body)
	}
}

func EchoError(w http.ResponseWriter, code StatusCode, err error) {
	w.WriteHeader(code.Code())
	if body, err := json.Marshal(ErrResponse{
		Code:   int(code),
		Msg:    code.String(),
		Detail: err.Error(),
	}); err != nil {
		logger.DefaultLogger.Warn("json unmarashal err", zap.Error(err))

		w.Write([]byte("err"))
	} else {
		w.Write(body)
	}
}
