package student

import (
	"backend-go/db"
	"backend-go/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

// @desc   update student
// @route  GET /api/parent/update
// @access Private
func update(c *fiber.Ctx) error {
  var studentID string
  utils.GetLocals(c.Locals("studentID"), &studentID)

  student, err := db.GetStudentByID(studentID)
  if err != nil {
    utils.Error(c, err)
  }
  
  tokenString, err := utils.StudentGenerateToken(student.StudentID, student.Grade, student.SubjectList)
  if err != nil {
    utils.Error(c, err)
  }

  return c.JSON(bson.M{
    "subjectList": student.SubjectList,
    "grade": student.Grade,
    "token": tokenString,
  })
}