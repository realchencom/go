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
		e := fmt.Errorf("fatal error config file: %w", err)
		fmt.Println(e.Error())
		panic(e)
	}
}
