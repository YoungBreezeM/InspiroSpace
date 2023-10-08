package dao

import (
	"context"
	"time"
)

type ChatContent struct {
	OpenId  string `bson:"openId" json:"openId"`
	Content string `bson:"content" json:"content"`
}

func (d *Dao) CreateChat(c *ChatContent) (err error) {

	collection := d.Mongo.Database("freebox").Collection("chatgtp")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	//
	_, err = collection.InsertOne(ctx, c)
	return
}
