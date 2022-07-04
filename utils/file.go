package utils

import (
	"fmt"
	"github.com/spf13/viper"
)

type config struct {
	*viper.Viper
}

var (
	Con config
)

func init() {
	Con = config{viper.GetViper()}
	Con.load()
}
func (c *config) load() {
	//加载./conf/app.yml.json配置文件
	c.SetConfigName("app.yml") // name of config file (without extension)
	c.SetConfigType("yaml")
	c.AddConfigPath("./conf")
	c.AddConfigPath("../conf") // call multiple times to add many search paths
	c.AddConfigPath(".")       // optionally look for config in the working directory
	//未找到配置文件，抛出恐慌异常
	if err := c.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// 配置文件未找到错误；如果需要可以忽略
			panic(err)
		} else {
			// 配置文件被找到，但产生了另外的错误
			panic(fmt.Errorf("Read Configuration fatal error in file: %s \n", err))
		}
	}
}
