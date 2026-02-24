package middleware

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type Body struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type responseWriter struct {
	http.ResponseWriter
	buf        bytes.Buffer
	statusCode int
}

func (w *responseWriter) WriteHeader(code int) {
	w.statusCode = code
}

func (w *responseWriter) Write(b []byte) (int, error) {
	return w.buf.Write(b)
}

// WrapResponse 将 HTTP 响应包装为 {code, msg, data} 格式
func WrapResponse(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next(rw, r)

		var body Body
		if rw.statusCode >= 400 {
			body = Body{
				Code: rw.statusCode,
				Msg:  rw.buf.String(),
			}
		} else {
			var data interface{}
			if err := json.Unmarshal(rw.buf.Bytes(), &data); err != nil {
				data = rw.buf.String()
			}
			body = Body{
				Code: 200,
				Msg:  "ok",
				Data: data,
			}
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(rw.statusCode)
		json.NewEncoder(w).Encode(body)
	}
}
