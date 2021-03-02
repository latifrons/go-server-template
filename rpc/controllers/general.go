package controllers

// swaggo
// Run swag init -g general.go

// @title Atom8 Server API
// @version 1.0
// @description Atom8 server API
// @host localhost:8080
// @BasePath /api/v1
import (
	"errors"
	"fmt"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/atom-eight/tmt-backend/dbgorm"
	"github.com/atom-eight/tmt-backend/folder"
	"github.com/atom-eight/tmt-backend/oss"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var FullDateFormatPattern = "2 Jan 2006 15:04:05"
var ShortDateFormatPattern = "2 Jan 2006"

type Platform string

const Platform_PC Platform = "PC"
const Platform_M Platform = "MO"

const DefaultCacheControl = "public; max-age=86400"

type RpcController struct {
	FolderConfig               folder.FolderConfig
	ReturnDetailedErrorMessage bool
	DbOperator                 *dbgorm.DbOperator
	FileUploader               *oss.FileUploader
	S3Bucket                   string
	MaxUploadFileSize          int64
}

func (rpc *RpcController) NewRouter() *gin.Engine {
	router := gin.New()
	if logrus.GetLevel() > logrus.DebugLevel {
		logger := gin.LoggerWithConfig(gin.LoggerConfig{
			Formatter: ginLogFormatter,
			Output:    logrus.StandardLogger().Out,
			SkipPaths: []string{"/"},
		})
		router.Use(logger)
	}
	router.Use(gin.RecoveryWithWriter(logrus.StandardLogger().Out))
	router.Use(Cors())

	rpc.addRouter(router)
	router.Use(BreakerWrapper)

	return router
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET,HEAD,OPTIONS,POST,PUT")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "*")

		//c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		//c.Writer.Header().Set("Access-Control-Expose-Headers", "Access-Control-Allow-Origin")
		//c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		defer func() {
			if err := recover(); err != nil {
				logrus.Error("gin")
			}
		}()

		c.Next()
	}
}

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
	logrus.Tracef("gin log %v ", logEntry)
	//return  logEntry
	return ""
}

func (rpc *RpcController) addRouter(router *gin.Engine) *gin.Engine {
	//router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/api/health", rpc.Health)

	vDebug := router.Group("/api/debug")
	{
		vDebug.GET("ua", rpc.DebugUA)
		vDebug.GET("panic", rpc.Panic)
	}
	v1 := router.Group("/api/v1")
	{
		v1.POST("file", rpc.UploadFile)
		v1.DELETE("file/:key", rpc.DeleteFile)
	}
	return router
}

func cors(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
}

func Response(c *gin.Context, status int, code int, msg string, data interface{}) {
	c.JSON(status, GeneralResponse{
		Code: code,
		Msg:  msg,
		Data: data,
	})
}

func ResponsePaging(c *gin.Context, status int, code int, msg string, pagingResult dbgorm.PagingResult, data interface{}) {
	c.JSON(status, PagingResponse{
		Code:   code,
		Msg:    msg,
		Limit:  pagingResult.Limit,
		Total:  pagingResult.Total,
		Offset: pagingResult.Offset,
		Data:   data,
	})
}

func (rpc *RpcController) ResponseError(c *gin.Context, err error) bool {
	if err == nil {
		return false
	}
	if err == gorm.ErrRecordNotFound {
		Response(c, http.StatusNotFound, 1, "record not found", nil)
		return true
	}
	logrus.WithError(err).Warn("internal error")
	if rpc.ReturnDetailedErrorMessage {
		Response(c, http.StatusInternalServerError, 2, err.Error(), nil)
	} else {
		Response(c, http.StatusInternalServerError, 2, "Internal error. Check your input or wait some time to retry", nil)
	}
	return true
}

func (rpc *RpcController) extractPagingQuery(c *gin.Context) dbgorm.PagingParams {
	return dbgorm.PagingParams{
		Offset:    tryParseIntDefault(c.DefaultQuery("offset", "0"), 0),
		Limit:     tryParseIntDefault(c.DefaultQuery("limit", "10"), 10),
		NeedTotal: true,
	}
}

func (rpc *RpcController) ResponseEmptyQuery(c *gin.Context, value string) bool {
	if value == "" {
		Response(c, http.StatusBadRequest, 1, "param missing", nil)
		return true
	}
	return false
}

func (rpc *RpcController) ToStringArray(query string) (arr []string, err error) {
	return strings.Split(query, "$"), nil
}

func (rpc *RpcController) UploadFile(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if rpc.ResponseError(c, err) {
		return
	}

	size := fileHeader.Size
	if size > rpc.MaxUploadFileSize {
		Response(c, http.StatusNotAcceptable, 1, "Exceed max file size", nil)
		return
	}

	fileName := fileHeader.Filename
	f, err := fileHeader.Open()
	if rpc.ResponseError(c, err) {
		return
	}

	result, err := rpc.FileUploader.Upload(f, rpc.S3Bucket, fileName, DefaultCacheControl)
	if rpc.ResponseError(c, err) {
		return
	}
	logrus.WithField("path", result.Location).Info("file uploaded")
	Response(c, http.StatusOK, 1, "", result.Location)
}

func (rpc *RpcController) DeleteFile(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		return
	}

	result, err := rpc.FileUploader.Remove(rpc.S3Bucket, key)
	if rpc.ResponseError(c, err) {
		return
	}
	logrus.Info("file deleted")
	Response(c, http.StatusOK, 1, "", result.VersionId)
}

func tryParseIntDefault(v string, d int) int {
	c, err := strconv.Atoi(v)
	if err != nil {
		return d
	}
	return c
}
