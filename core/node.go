package core

import (
	"github.com/latifrons/lbserver/folder"
	"github.com/latifrons/lbserver/rpc"
	"github.com/latifrons/lbserver/tools"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Node struct {
	FolderConfig folder.FolderConfig
	components   []Component
}

func (n *Node) Setup() {
	allowOrigins, err := tools.ReadStringArrayFromFile(viper.GetString("rpc.allow_origins"))
	if err != nil {
		logrus.WithError(err).Fatal("failed to read cors file")
	}

	rpcServer := &rpc.RpcServer{
		C: &rpc.RpcController{
			FolderConfig:               n.FolderConfig,
			ReturnDetailedErrorMessage: viper.GetBool("debug.return_detailed_error"),
			Mode:                       viper.GetString("common.mode"),
			AllowOrigins:               allowOrigins,
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
