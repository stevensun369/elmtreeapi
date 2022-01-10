package utils

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func GetLocals(locals interface{}, result interface{}) {
  json.Unmarshal([]byte(fmt.Sprintf("%v", locals)), &result)
} 

func SetLocals(c *fiber.Ctx, name string,  data interface{}) {
	bytes, _ := json.Marshal(data)
	json := string(bytes)
	c.Locals(name, json)
}