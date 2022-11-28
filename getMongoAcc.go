package main

import (
	"context"
	"encoding/json"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func getMongoAcc(username string) jsonAccount {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	eh(err)
	defer func() {
		err := client.Disconnect(context.TODO())
		eh(err)
	}()
	coll := client.Database("mongolang-js").Collection("users")

	var result bson.M
	err = coll.FindOne(context.TODO(), bson.D{{Key: "doc", Value: username}}).Decode(&result)
	if err == mongo.ErrNoDocuments {
		fmt.Println("No document was account with the name " + username)
		return jsonAccount{Doc: "missing account account"}
	}
	eh(err)
	jsonData, err := json.MarshalIndent(result, "", "    ")
	eh(err)
	var jsondata jsonAccount

	err = json.Unmarshal([]byte(jsonData), &jsondata)
	eh(err)

	return jsondata
}
