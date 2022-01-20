package parent

import (
	"backend-go/db"
	"backend-go/models"
	"backend-go/utils"
	"context"
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

// @desc   Add student to student list
// @route  PUT /api/parent/students
// @access Private
func parentAddStudent(c *fiber.Ctx) error {
  var body map[string]string
  json.Unmarshal(c.Body(), &body)

  var parentID string
  utils.GetLocals(c.Locals("parentID"), &parentID)

  // get student and make sure it is there
  student, _ := db.GetStudentByIDAndCNP(body["studentCNP"], body["studentID"])
  if (student.StudentID == "") {
    return utils.MessageError(c, "Nu există niciun elev cu datele introduse.")
  }

  // new studentIDList
  var newStudentIDList []string
  utils.GetLocals(c.Locals("studentIDList"), &newStudentIDList)
  newStudentIDList = append(newStudentIDList, body["studentID"])

  // updating the parent
  var parent models.Parent
  if err := db.Parents.FindOneAndUpdate(context.Background(), 
    bson.M{"parentID": parentID}, 
    bson.M{"$set": bson.M{"studentIDList": newStudentIDList}},
    ).Decode(&parent); err != nil {
    return utils.Error(c, err)
  }
  parent.StudentIDList = append(parent.StudentIDList, body["studentID"])

  students, err := db.GetStudents(bson.M{
    "studentID": bson.M{"$in": newStudentIDList},
  }, db.GradeSort)
  if err != nil {
    return utils.Error(c, err)
  }

  tokenString, err := utils.ParentGenerateToken(parent.ParentID, parent.StudentIDList)
  if err != nil {
    return utils.Error(c, err)
  }

  if len(students) == 0 {
    students = []models.Student {}
  }

  return c.JSON(
    bson.M{
      "students": students,
      "token": tokenString,
    },
  )
}
