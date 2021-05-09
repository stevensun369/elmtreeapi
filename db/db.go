package db

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"os"

	"github.com/joho/godotenv"
)

var Client *mongo.Client

func InitDatabase() {
  var err error
  _ = godotenv.Load(".env")
  if err != nil {
    log.Fatal(err)
  }
  Client, err = mongo.Connect(context.Background(), options.Client().ApplyURI(os.Getenv("MONGO_URI")))

  if err != nil {
    log.Fatal(err)
  }
}

func GetCollection(collectionName string) (*mongo.Collection, error) {
  collection := Client.Database("elmtree").Collection(collectionName)

  return collection, nil
}
