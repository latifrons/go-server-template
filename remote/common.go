package remote

import (
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func GetRestyClient() *resty.Client {
	c := resty.New()
	c.Debug = viper.GetString("common.mode") == "debug"
	c.SetLogger(logrus.StandardLogger())
	return c
}

func GetRestyClientNoDebug() *resty.Client {
	c := resty.New()
	c.SetLogger(logrus.StandardLogger())
	c.DisableWarn = true
	return c
}
