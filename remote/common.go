package remote

import (
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
)

func GetRestyClient(rpcLog bool) *resty.Client {
	c := resty.New()
	c.Debug = rpcLog
	c.SetLogger(logrus.StandardLogger())
	return c
}

func GetRestyClientNoDebug() *resty.Client {
	c := resty.New()
	c.SetLogger(logrus.StandardLogger())
	c.DisableWarn = true
	return c
}
