package main

import (
	// fiber

	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	// internal database package
	"backend-go/db"

	// models (just for now)
	"backend-go/models"

	// routes
	"backend-go/parent"
	"backend-go/student"
	"backend-go/teacher"
)

func main() {
  app := fiber.New()

  app.Use(cors.New(cors.Config{
    AllowOrigins: "*", 
  }))

  db.InitDatabase()

  teacher.Routes(app)
  student.Routes(app)
  parent.Routes(app)

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

  app.Get("/", func (c *fiber.Ctx) error {
    return c.SendString("api is running")
  })

  log.Fatal(app.Listen(":9999"))
}