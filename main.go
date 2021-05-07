package main

import (
	// fiber

	"context"
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"

	// internal database package
	"backend-go/db"

	// models (just for now)
	"backend-go/models"
)

func main() {
  app := fiber.New()

  db.InitDatabase()

  app.Post("/", func (c * fiber.Ctx) error {
    var subject models.Subject
    json.Unmarshal([]byte(c.Body()), &subject)

    collection, err := db.GetCollection("subjects")
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