package dao

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type AppConf struct {
	Name   string `bson:"name" json:"name"`
	Domain string `bson:"domain" json:"domain"`
}

func (d *Dao) GetAppConf(name string) (a *AppConf, err error) {

	a = &AppConf{}
	collection := d.Mongo.Database("freeaibox").Collection("apps_conf")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	//
	c := collection.FindOne(ctx, bson.M{"name": name})
	err = c.Decode(a)
	return

}
