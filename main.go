package main

import (
	// fiber

	"context"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	// internal database package

	// models (just for now)
	"backend-go/models"
)

var Client *mongo.Client

func mongoConnect() (*mongo.Client, error) {
  client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb+srv://stevensun:StevenS369@dev.wljsz.mongodb.net/elmtree?retryWrites=true&w=majority"))

  if err != nil {
    log.Fatal(err)
  }

  // err = client.Ping(context.Background(), readpref.Primary())

  // if err != nil {
  //   log.Fatal(err)
  // }

  return client, nil
}

func getCollection(dbName string, collectionName string) (*mongo.Collection, error) {

  collection := Client.Database(dbName).Collection(collectionName)

  return collection, nil
}

func main() {
  app := fiber.New()

  var err error
  Client, err = mongoConnect()
  if err != nil {
    log.Fatal(err)
  }

  app.Get("/", func (c * fiber.Ctx) error {
    subject := models.Subject{
      SubjectID: "000000000",
      Name: "testing",
      Grade: models.Grade{
        GradeID: "000000000",
        SchoolID: "000000000",
        YearFrom: 2020,
        YearTo: 2021,
        GradeNumber: 9,
        GradeLetter: "A",
      },
    }

    collection, err := getCollection("elmtree", "subjects")
    if err != nil {
      c.Status(500).SendString(fmt.Sprintf("%v", err))
    }

    insertResult, err := collection.InsertOne(context.Background(), subject)
    if err != nil {
      c.Status(500).SendString(fmt.Sprintf("%v", err))
    }   

    return c.JSON(insertResult)
  })

  app.Listen(":1000")
}