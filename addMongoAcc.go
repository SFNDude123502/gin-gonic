package main

import (
	"context"
	"encoding/base64"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func addMongoAcc(username string, password string, auth int) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	coll := client.Database("mongolang-js").Collection("users")

	newPass := base64.StdEncoding.EncodeToString([]byte(password))

	login_daata := bson.D{{Key: "password", Value: newPass}, {Key: "auth", Value: auth}}

	coll.InsertOne(context.TODO(), bson.D{{Key: "doc", Value: username}, {Key: "login_data", Value: login_daata}})
}
