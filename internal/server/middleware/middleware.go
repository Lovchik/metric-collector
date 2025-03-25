package middleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"io"
	"time"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		bodyBytes, _ := io.ReadAll(c.Request.Body)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		responseBody := &bytes.Buffer{}
		writer := &responseWriter{ResponseWriter: c.Writer, body: responseBody}
		c.Writer = writer

		c.Next()

		duration := time.Since(startTime)

		log.Printf("[Request] Method: %s | URI: %s | Duration: %v", c.Request.Method, c.Request.RequestURI, duration)
		log.Printf("[Request] Body: %s", string(bodyBytes))

		log.Printf("[Response] Status: %d | Content-Length: %d", c.Writer.Status(), responseBody.Len())
		log.Printf("[Response] Body: %s", responseBody.String())
	}
}

type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b) // Записываем в буфер
	return w.ResponseWriter.Write(b)
}
