package rpc

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

const ShutdownTimeoutSeconds = 5

type RpcServer struct {
	C    *RpcController
	Port string

	router *gin.Engine
	server *http.Server
}

func (srv *RpcServer) Start() {
	srv.router = srv.C.NewRouter()

	srv.server = &http.Server{
		Addr:    ":" + srv.Port,
		Handler: srv.router,
	}

	logrus.Infof("listening Http on %s", srv.Port)
	go func() {
		// service connections
		if err := srv.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.WithError(err).Fatalf("berror in Http server")
		}
	}()
}

func (srv *RpcServer) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), ShutdownTimeoutSeconds*time.Second)
	defer cancel()
	if err := srv.server.Shutdown(ctx); err != nil {
		logrus.WithError(err).Error("berror while shutting down the Http server")
	}
	logrus.Infof("http server Stopped")
}

func (srv *RpcServer) Name() string {
	return fmt.Sprintf("rpcServer at port %s", srv.Port)
}

func (srv *RpcServer) InitDefault() {

}

var ginLogFormatter = func(param gin.LogFormatterParams) string {
	if logrus.GetLevel() < logrus.TraceLevel {
		return ""
	}
	var statusColor, methodColor, resetColor string
	if param.IsOutputColor() {
		statusColor = param.StatusCodeColor()
		methodColor = param.MethodColor()
		resetColor = param.ResetColor()
	}

	if param.Latency > time.Minute {
		// Truncate in a golang < 1.8 safe way
		param.Latency = param.Latency - param.Latency%time.Second
	}

	logEntry := fmt.Sprintf("GIN %v %s %3d %s %13v  %15s %s %-7s %s %s %s",
		param.TimeStamp.Format("2006/01/02 - 15:04:05"),
		statusColor, param.StatusCode, resetColor,
		param.Latency,
		param.ClientIP,
		methodColor, param.Method, resetColor,
		param.Path,
		param.ErrorMessage,
	)
	//logrus.Tracef("gin log %v ", logEntry)
	return logEntry
	//return ""
}

var FullDateFormatPattern = "2 Jan 2006 15:04:05"
var ShortDateFormatPattern = "2 Jan 2006"

func BreakerWrapper(c *gin.Context) {
	name := c.Request.Method + "-" + c.Request.RequestURI
	hystrix.Do(name, func() error {
		c.Next()

		statusCode := c.Writer.Status()

		if statusCode >= http.StatusInternalServerError {
			str := fmt.Sprintf("status code %d", statusCode)
			return errors.New(str)
		}

		return nil
	}, func(e error) error {
		if e == hystrix.ErrCircuitOpen {
			c.String(http.StatusAccepted, "Please try again later")
		}

		return e
	})
}

func RequestLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var buf bytes.Buffer
		tee := io.TeeReader(c.Request.Body, &buf)
		body, _ := ioutil.ReadAll(tee)
		c.Request.Body = ioutil.NopCloser(&buf)
		logrus.WithField("uri", c.Request.RequestURI).Trace("Received request")
		logrus.Trace(c.Request.Header)
		logrus.Trace(string(body))
		c.Next()
	}
}

func tryParseIntDefault(v string, d int) int {
	c, err := strconv.Atoi(v)
	if err != nil {
		return d
	}
	return c
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func ResponseLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw
		c.Next()
		fmt.Println("Response body: " + blw.body.String())
	}
}
