package student

import (
	"backend-go/db"
	"backend-go/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func update(c *fiber.Ctx) error {
  var studentID string
  utils.GetLocals(c.Locals("studentID"), &studentID)

  student, err := db.GetStudentByID(studentID)
  utils.CheckError(c, err)
  
  tokenString, err := utils.StudentGenerateToken(student.StudentID, student.Grade, student.SubjectList)
  utils.CheckError(c, err)

  return c.JSON(bson.M{
    "subjectList": student.SubjectList,
    "grade": student.Grade,
    "token": tokenString,
  })
}