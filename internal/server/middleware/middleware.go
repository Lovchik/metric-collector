package middleware

import (
	"bytes"
	"compress/gzip"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strings"
	"time"
)

type gzipWriter struct {
	gin.ResponseWriter
	writer *gzip.Writer
}

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			log.Error(err)
			return
		}
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
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func GzipMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ow := c.Writer
		contentType := c.ContentType()

		if strings.Contains(c.GetHeader("Accept-Encoding"), "gzip") &&
			(contentType == "application/json" || contentType == "text/html") {
			gw := gzip.NewWriter(ow)
			c.Writer = &gzipWriter{ResponseWriter: ow, writer: gw}
			defer gw.Close()
			c.Header("Content-Encoding", "gzip")
		}

		if strings.Contains(c.GetHeader("Content-Encoding"), "gzip") {
			gr, err := gzip.NewReader(c.Request.Body)
			if err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			defer gr.Close()
			c.Request.Body = io.NopCloser(gr)
		}

		c.Next()
	}
}

func (g *gzipWriter) Write(data []byte) (int, error) {
	return g.writer.Write(data)
}
