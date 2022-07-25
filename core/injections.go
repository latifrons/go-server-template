package core

import (
	"fmt"
	"github.com/golobby/container/v3"
	"github.com/latifrons/lbserver/consts"
	"github.com/latifrons/lbserver/db"
	"github.com/latifrons/lbserver/debug"
	"github.com/latifrons/lbserver/tools"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

func buildDependencies() (err error) {
	singletons := []interface{}{
		//func() *cache.RedisCache {
		//	redisCache := &cache.RedisCache{
		//		Address:  tools.ViperMustGetString("redis.address"),
		//		Password: viper.GetString("redis.password"),
		//		Db:       tools.ViperMustGetInt("redis.db"),
		//	}
		//	redisCache.Init()
		//	return redisCache
		//},
		func() *debug.Flags {
			debugFlags := &debug.Flags{
				ReturnDetailError: viper.GetBool("debug.return_detail_error"),
				DbLog:             viper.GetBool("debug.db_log"),
				RequestLog:        viper.GetBool("debug.request_log"),
				ResponseLog:       viper.GetBool("debug.response_log"),
				RpcLog:            viper.GetBool("debug.rpc_log"),
				GinDebug:          viper.GetBool("debug.gin_debug"),
				SkipVerifyAdmin:   viper.GetBool("debug.skip_verify_admin"),
				SkipEmail:         viper.GetBool("debug.skip_email"),
				LogLevel:          viper.GetString("debug.log_level"),
				Swagger:           viper.GetBool("debug.swagger"),
			}
			return debugFlags
		},

		func() *consts.UrlConfig {
			//myselfUrl := tools.ViperMustGetString("tx.myself")
			urlConfig := &consts.UrlConfig{
				//DtmServer:                               dtmServerUrl,
				//TokenRegistryFetchNftsByNftCollectionId: tokenRegistryUrl + "/v1/public/fetch_nfts_by_collection_id",
			}
			return urlConfig
		},
		func() *gorm.DB {
			if tools.ViperMustGetString("mysql.user") == "" {
				logrus.Fatal("empty secret")
			}
			dbOperator := &db.DbOperator{
				Source: fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=true&loc=Local",
					tools.ViperMustGetString("mysql.user"), tools.ViperMustGetString("mysql.password"),
					tools.ViperMustGetString("mysql.url"), tools.ViperMustGetString("mysql.schema")),
				DbLog: tools.ViperMustGetBool("debug.db_log"),
			}
			dbOperator.InitDefault()
			return dbOperator.Db
		},
	}

	for _, singleton := range singletons {
		err = container.Singleton(singleton)
		if err != nil {
			return
		}
	}
	return

}
