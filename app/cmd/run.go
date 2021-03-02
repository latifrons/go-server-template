package cmd

import (
	"github.com/atom-eight/tmt-backend/core"
	"github.com/atom-eight/tmt-backend/folder"
	"github.com/latifrons/commongo/mylog"
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
	Short: "Start a Atom8Server instance",
	Long:  `Start a Atom8Server instance`,
	Run: func(cmd *cobra.Command, args []string) {
		logrus.Info("Atom8Server Starting")
		folderConfigs := folder.EnsureFolders()
		readConfig(folderConfigs.Config)
		readPrivate(folderConfigs.Private)

		mylog.LogInit(mylog.LogLevel(viper.GetString("log.level")))

		mylog.InitLogger(logrus.StandardLogger(), mylog.LogConfig{
			MaxSize:    10,
			MaxBackups: 100,
			MaxAgeDays: 90,
			Compress:   true,
			LogDir:     folderConfigs.Log,
			OutputFile: "atom8",
		})

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
