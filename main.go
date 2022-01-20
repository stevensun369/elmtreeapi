package main

import (
	// fiber
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"log"

	// internal database package
	"backend-go/db"

	// models (just for now)

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

  teacher.Endpoints(app)
  student.Endpoints(app)
  parent.Endpoints(app)

  app.Get("/", func (c *fiber.Ctx) error {
    return c.SendString("api is running")
  })

  log.Fatal(app.Listen(":9999"))
}