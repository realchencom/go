package rs

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Datasource struct {
	Dsn string `mapstructure:"dsn"`
}

func (ds *Datasource) GetDB() (*gorm.DB, error) {
	if err := Con.Sub("datasource").Unmarshal(&ds); err != nil {
		panic(fmt.Errorf("unmarshal datasource conf failed, err:%s \n", err))
	}
	DB, err := gorm.Open(mysql.Open(ds.Dsn), &gorm.Config{})
	return DB, err
}
