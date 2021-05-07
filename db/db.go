package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func InitDatabase() {
  var err error
  Client, err = mongo.NewClient(options.Client().ApplyURI("mongodb+srv://stevensun:StevenS369@dev.wljsz.mongodb.net/elmtree?retryWrites=true&w=majority"))

  if err != nil {
    log.Fatal(err)
  }

  ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
  defer cancel()

  err = Client.Connect(ctx)

  if err != nil {
    log.Fatal(err)
  }

  fmt.Printf("The database has been connected")

  defer Client.Disconnect(ctx)
}


