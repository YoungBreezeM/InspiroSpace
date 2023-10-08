package boot

import (
	"easygin/api"
	"easygin/internal/conf"
	"easygin/internal/dao"
	"easygin/internal/log"
	"easygin/internal/mid"
	"easygin/internal/service"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type App struct {
	conf *conf.Config
	e    *gin.Engine
}

func BootStrap(confPath string) (app *App, err error) {

	//init config
	c, err := conf.InitConfig(confPath)
	if err != nil {
		return
	}

	//init log
	log.NewLog()

	//init mysql
	// mysql, err := dao.InitMysql(c.MysqlConfig)
	// if err != nil {
	// 	return
	// }

	//init redis
	r, err := dao.InitRedis(c.RedisConfig)
	if err != nil {
		return
	}

	//init mongo
	mgo, err := dao.InitMongo(c.MongoConfig)
	if err != nil {
		return
	}

	//init dao
	d := dao.InitDao(&gorm.DB{}, r, mgo)

	//init service
	s := service.InitService(d, c)

	//init mid
	m := mid.NewMid(s)

	//init api
	e := gin.Default()
	api.InitApi(s, e, c, m)

	//init app
	app = &App{e: e, conf: c}

	return
}

func (a *App) Run() {
	port := fmt.Sprintf(":%d", a.conf.AppConfig.Port)
	a.e.Run(port)
}
