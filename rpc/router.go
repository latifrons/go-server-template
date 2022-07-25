package rpc

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/swag/example/basic/docs"
	"time"
)

func (rpc *RpcController) NewRouter() *gin.Engine {
	if rpc.Flags.GinDebug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	if rpc.Flags.RequestLog {
		if logrus.GetLevel() >= logrus.DebugLevel {
			logger := gin.LoggerWithConfig(gin.LoggerConfig{
				SkipPaths: []string{"/", "/health"},
			})
			router.Use(logger)
			router.Use(RequestLoggerMiddleware())
			router.Use(ResponseLoggerMiddleware())

		} else {
			logger := gin.LoggerWithConfig(gin.LoggerConfig{
				SkipPaths: []string{"/", "/health"},
			})
			router.Use(logger)
		}
	}

	router.Use(gin.RecoveryWithWriter(logrus.StandardLogger().Out))
	router.Use(cors.New(cors.Config{
		AllowOrigins:     rpc.AllowOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
		AllowHeaders:     []string{"Accept", "User-Agent", "Origin", "Content-Length", "Content-Type", "Authorization", "X-Requested-With"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	rpc.addRouter(router)
	router.Use(BreakerWrapper)

	return router
}

func (rpc *RpcController) addRouter(router *gin.Engine) *gin.Engine {
	docs.SwaggerInfo.BasePath = "/v1"
	if rpc.Flags.Swagger {
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}

	router.GET("/health", rpc.Health)

	return router
}
