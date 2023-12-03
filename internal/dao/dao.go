package dao

import (
	"context"
	"easygin/internal/conf"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Dao struct {
	Mysql *gorm.DB
	Redis *redis.Client
	Mongo *mongo.Client
}

func InitDao(db *gorm.DB, r *redis.Client, mgo *mongo.Client) *Dao {
	return &Dao{
		Mysql: db,
		Redis: r,
		Mongo: mgo,
	}
}

func InitMysql(c *conf.MysqlConfig) (db *gorm.DB, err error) {
	db, err = gorm.Open("mysql", c.Dst)
	return
}

func InitRedis(c *conf.RedisConfig) (r *redis.Client, err error) {
	r = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", c.Host, c.Port),
		DB:   c.DB,
	})

	if _, err = r.Ping(context.Background()).Result(); err != nil {
		err = errors.WithStack(err)
	}

	return
}

func InitMongo(c *conf.MongoConfig) (mgo *mongo.Client, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mgo, err = mongo.Connect(ctx, options.Client().ApplyURI(c.Addr))
	if err != nil {
		return
	}

	err = mgo.Ping(ctx, readpref.Primary())
	return
}
