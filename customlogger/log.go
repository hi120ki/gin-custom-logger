package customlogger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"time"

	"github.com/gin-gonic/gin"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w bodyLogWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

func CustomLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Request
		var buf bytes.Buffer
		tee := io.TeeReader(c.Request.Body, &buf)
		body, _ := ioutil.ReadAll(tee)
		c.Request.Body = ioutil.NopCloser(&buf)
		// Response
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		c.Next()

		var request interface{}
		if err := json.Unmarshal(body, &request); err != nil {
			request = string(body)
		}

		var response interface{}
		if err := json.Unmarshal(blw.body.Bytes(), &response); err != nil {
			response = blw.body.String()
		}

		logInterface := map[string]interface{}{
			"name":     "gin-custom-log",
			"host":     c.ClientIP(),
			"time":     time.Now().UnixNano(),
			"status":   c.Writer.Status(),
			"uri":      c.Request.URL.Path,
			"method":   c.Request.Method,
			"header":   c.Request.Header,
			"request":  request,
			"response": response,
		}
		logText, _ := json.Marshal(logInterface)
		fmt.Println(string(logText))
	}
}
