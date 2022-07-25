package cmd

import (
	"github.com/latifrons/commongo/utilfuncs"
	"github.com/latifrons/lbserver/core"
	"github.com/latifrons/lbserver/folder"
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
	Short: "Start a Atom8Server instance",
	Long:  `Start a Atom8Server instance`,
	Run: func(cmd *cobra.Command, args []string) {
		logrus.Info("Atom8Server Starting")
		folderConfigs := folder.EnsureFolders()
		readConfig(folderConfigs.Config)
		readPrivate(folderConfigs.Private)
		dumpConfig()

		formatter := new(logrus.TextFormatter)
		formatter.TimestampFormat = "01-02 15:04:05.000000"
		formatter.FullTimestamp = true
		formatter.ForceColors = true

		lvl, err := logrus.ParseLevel(tools.ViperMustGetString("debug.log_level"))
		utilfuncs.PanicIfError(err, "log level")
		logrus.SetLevel(lvl)
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
