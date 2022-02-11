package cmd

import (
	"github.com/latifrons/lbserver/core"
	"github.com/latifrons/lbserver/tools"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"syscall"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Start a LbServer instance",
	Long:  `Start a LbServer instance`,
	Run: func(cmd *cobra.Command, args []string) {
		logrus.Info("LbServer Starting")
		folderConfigs := tools.EnsureFolders()
		readConfig(folderConfigs.Config)
		readPrivate(folderConfigs.Private)

		formatter := new(logrus.TextFormatter)
		formatter.TimestampFormat = "01-02 15:04:05.000000"
		formatter.FullTimestamp = true
		formatter.ForceColors = true

		logrus.SetFormatter(formatter)
		logrus.StandardLogger().SetOutput(os.Stdout)

		// init logs and other facilities before the core starts

		core := &core.Node{
			FolderConfig: folderConfigs,
			//	DataFolder: folderConfigs.Data,
		}
		core.Setup()
		core.Start()

		// prevent sudden stop. Do your clean up here
		var gracefulStop = make(chan os.Signal)

		signal.Notify(gracefulStop, syscall.SIGTERM)
		signal.Notify(gracefulStop, syscall.SIGINT)

		func() {
			sig := <-gracefulStop
			logrus.Infof("caught sig: %+v", sig)
			logrus.Info("Exiting... Please do no kill me")
			core.Stop()
			os.Exit(0)
		}()

	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	rootCmd.PersistentFlags().Int("rpc-port", 8080, "RPC port")

	_ = viper.BindPFlag("rpc.port", rootCmd.PersistentFlags().Lookup("rpc-port"))

}
