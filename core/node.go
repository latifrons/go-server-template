package core

import (
	"github.com/latifrons/lbserver/cache"
	"github.com/latifrons/lbserver/folder"
	"github.com/latifrons/lbserver/rpc"
	"github.com/latifrons/lbserver/tools"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/yitter/idgenerator-go/idgen"
)

type Node struct {
	FolderConfig folder.FolderConfig
	components   []Component
}

func (n *Node) Setup() {
	err := buildDependencies()
	if err != nil {
		logrus.WithError(err).Fatal("failed to inject")
	}

	var redisCache *cache.RedisCache
	if tools.ViperMustGetString("common.id_gen") != "local" {
		intCmd := redisCache.Rdb.Incr(tools.GetContext(10), "snowflake_worker")
		if intCmd.Err() != nil {
			logrus.WithError(intCmd.Err()).Fatal("failed to incr snowflake worker")
		}
		idgen.SetIdGenerator(idgen.NewIdGeneratorOptions(uint16(intCmd.Val())))
	}

	allowOrigins, err := tools.ReadStringArrayFromFile(viper.GetString("rpc.allow_origins"))
	if err != nil {
		logrus.WithError(err).Fatal("failed to read cors file")
	}

	rpcServer := &rpc.RpcServer{
		C: &rpc.RpcController{
			AllowOrigins: allowOrigins,
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
