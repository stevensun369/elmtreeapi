package utils

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func Error(c *fiber.Ctx, err error) error  {
  return c.Status(500).SendString(fmt.Sprintf("%v", err))
} 

func CheckError(c *fiber.Ctx, err error) error {
  if err != nil {
    return Error(c, err)
  } else {
    return nil
  }
}

func MessageError(c *fiber.Ctx, message string) error {
  return c.Status(401).JSON(bson.M{
    "message": message,
  })
}

func CheckMessageError(c *fiber.Ctx, err error, message string) error {
  if err != nil {
    return MessageError(c, message)
  } else {
    return nil
  }
}

