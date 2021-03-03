package core

import (
	"github.com/atom-eight/tmt-backend/dbgorm"
	"github.com/atom-eight/tmt-backend/folder"
	"github.com/atom-eight/tmt-backend/oss"
	"github.com/atom-eight/tmt-backend/rpc"
	"github.com/atom-eight/tmt-backend/rpc/controllers"
	"github.com/atom-eight/tmt-backend/two_factor"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Node struct {
	FolderConfig folder.FolderConfig
	components   []Component
}

func (n *Node) Setup() {
	dbOperator := &dbgorm.DbOperator{
		//Source: fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		//	viper.GetString("mysql.username"), viper.GetString("mysql.password"),
		//	viper.GetString("mysql.url"), viper.GetString("mysql.schema")),
		Source: viper.GetString("sqlite.path"),
	}
	dbOperator.InitDefault()

	twoFactorValidator := &two_factor.TwoFactorValidator{
		Address:               viper.GetString("redis.address"),
		Password:              viper.GetString("redis.password"),
		DBId:                  viper.GetInt("redis.db_id"),
		CodeExpirationSeconds: viper.GetInt("redis.code_expiration_seconds"),
	}
	twoFactorValidator.InitDefault()

	fileUploader := &oss.FileUploader{
		Endpoint: endpoints.UsEast2RegionID,
	}
	fileUploader.InitDefault()

	rpcServer := &rpc.RpcServer{
		C: &controllers.RpcController{
			FolderConfig:               n.FolderConfig,
			ReturnDetailedErrorMessage: viper.GetBool("debug.return_detailed_error"),
			DbOperator:                 dbOperator,
			FileUploader:               fileUploader,
			S3Bucket:                   "file.894569.site",
			MaxUploadFileSize:          100 * 1024,
		},
		Port: viper.GetString("rpc.port"),
	}
	rpcServer.InitDefault()

	n.components = append(n.components, rpcServer)
}

func (n *Node) Start() {
	for _, component := range n.components {
		logrus.Infof("Starting %s", component.Name())
		component.Start()
		logrus.Infof("Started: %s", component.Name())

	}
	//n.heightEventChan <- 10943851
	logrus.Info("Node Started")

}

func (n *Node) Stop() {
	for i := len(n.components) - 1; i >= 0; i-- {
		comp := n.components[i]
		logrus.Infof("Stopping %s", comp.Name())
		comp.Stop()
		logrus.Infof("Stopped: %s", comp.Name())
	}
	logrus.Info("Node Stopped")
}
