package cmd

import (
	"fmt"
	"github.com/latifrons/commongo/files"
	"github.com/latifrons/commongo/format"
	"github.com/latifrons/commongo/utilfuncs"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"path"
	"path/filepath"
)

// readConfig will respect {configdir}/config.toml first.
// If not found, get config from online source {configurl}
// {configdir}/injected.toml is the config issued by bootstrap server.
// finally merge env config so that any config can be override by env variables.
// Importance order:
// 1, ENV
// 2, injected.toml
// 3, config.toml or online toml if config.toml is not found
func readConfig(configFolder string) {
	configPath := path.Join(configFolder, "config.toml")

	if files.FileExists(configPath) {
		mergeLocalConfig(configPath)
	}

	// load injected config from ogbootstrap if any
	injectedPath := path.Join(configFolder, "injected.toml")
	if files.FileExists(injectedPath) {
		log.Info("merging local config file")
		mergeLocalConfig(injectedPath)
	}

	mergeEnvConfig()
}

func mergeEnvConfig() {
	// env override
	viper.SetEnvPrefix("INJ")
	viper.AutomaticEnv()
}

func readPrivate(privateFolder string) {
	configPath := path.Join(privateFolder, "private.toml")
	if files.FileExists(configPath) {
		mergeLocalConfig(configPath)
	}
}

//func writeConfig() {
//	configPath := files.FixPrefixPath(viper.GetString("rootdir"), path.Join(ConfigDir, "config_dump.toml"))
//	err := viper.WriteConfigAs(configPath)
//	utilfuncs.PanicIfError(err, "dump config")
//}

func mergeLocalConfig(configPath string) {
	absPath, err := filepath.Abs(configPath)
	utilfuncs.PanicIfError(err, fmt.Sprintf("Error on parsing config file path: %s", absPath))

	file, err := os.Open(absPath)
	utilfuncs.PanicIfError(err, fmt.Sprintf("Error on opening config file: %s", absPath))
	defer file.Close()

	viper.SetConfigType("toml")
	err = viper.MergeConfig(file)
	utilfuncs.PanicIfError(err, fmt.Sprintf("Error on reading config file: %s", absPath))
	return
}

func dumpConfig() {
	// print running config in console.
	b, err := format.PrettyJson(viper.AllSettings())
	utilfuncs.PanicIfError(err, "dump json")
	fmt.Println(b)
}
