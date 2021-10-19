package student

import (
	"backend-go/models"
	"backend-go/utils"
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)


func update(c *fiber.Ctx) error {
  studentIDLocals := fmt.Sprintf("%v", c.Locals("studentID"))
  var studentID string
  json.Unmarshal([]byte(studentIDLocals), &studentID)

  student, err := models.GetStudentById(studentID)
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }
  
  tokenString, err := utils.StudentGenerateToken(student.StudentID, student.Grade, student.SubjectList)
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  return c.JSON(bson.M{
    "subjectList": student.SubjectList,
    "grade": student.Grade,
    "token": tokenString,
  })
}