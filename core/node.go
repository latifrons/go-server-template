package core

import (
	"github.com/atom-eight/tmt-backend/dbgorm"
	"github.com/atom-eight/tmt-backend/folder"
	"github.com/atom-eight/tmt-backend/rpc"
	"github.com/atom-eight/tmt-backend/rpc/controllers"
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

	rpcServer := &rpc.RpcServer{
		C: &controllers.RpcController{
			FolderConfig:               n.FolderConfig,
			ReturnDetailedErrorMessage: viper.GetBool("debug.return_detailed_error"),
			DbOperator:                 dbOperator,
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
