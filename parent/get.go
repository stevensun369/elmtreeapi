package parent

import (
	"backend-go/db"
	"backend-go/models"
	"backend-go/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

// @desc   Get a list of the parent's students
// @route  GET /api/parent/students
// @access Private
func getStudents(c *fiber.Ctx) error {
  var parentID string
  utils.GetLocals(c.Locals("parentID"), &parentID)

  var studentIDList []string
  utils.GetLocals(c.Locals("studentIDList"), &studentIDList)

  students, err := db.GetStudents(bson.M{
    "studentID": bson.M{"$in": studentIDList},
  }, db.GradeSort)

  utils.CheckError(c, err)

  tokenString, err := utils.ParentGenerateToken(parentID, studentIDList)
  utils.CheckError(c, err)

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

// @desc   Get marks
// @route  GET /api/parent/marks/:studentID/:subjectID
// @access Private
func getMarks(c *fiber.Ctx) error {
  studentID := c.Params("studentID")
  subjectID := c.Params("subjectID")

  marks, err := db.GetMarks(bson.M{
    "studentID": studentID,
    "subject.subjectID": subjectID,
  }, db.DateSort)

  utils.CheckError(c, err)

  if len(marks) == 0 {
    marks = []models.Mark {}
  }

  return c.JSON(marks)
}

// @desc   Get truancies
// @route  GET /api/parent/truancy/:studentID/:subjectID
// @access Private
func getTruancies(c *fiber.Ctx) error {
  studentID := c.Params("studentID")
  subjectID := c.Params("subjectID")

  truancies, err := db.GetTruancies(bson.M{
    "studentID": studentID,
    "subject.subjectID": subjectID,
  }, db.DateSort)
  utils.CheckError(c, err)

  if len(truancies) == 0 {
    truancies = []models.Truancy {}
  }

  return c.JSON(truancies)
}

// @desc   Get average marks of all students
// @route  GET /api/parent/average
// @access Private
func getAverageMarks(c *fiber.Ctx) error {
  var studentIDList []string
  utils.GetLocals(c.Locals("studentIDList"), &studentIDList)

  averageMarks, err := db.GetAverageMarks(bson.M{
    "studentID": bson.M{"$in": studentIDList},
  }, db.TermSort)
  utils.CheckError(c, err)

  if len(averageMarks) == 0 {
    averageMarks = []models.AverageMark {}
  }

  return c.JSON(averageMarks)
}

// @desc   Get term marks of all students
// @route  GET /api/parent/term
// @access Private
func getTermMarks(c *fiber.Ctx) error {
  var studentIDList []string
  utils.GetLocals(c.Locals("studentIDList"), &studentIDList)

  termMarks, err := db.GetTermMarks(bson.M{
    "studentID": bson.M{"$in": studentIDList},
  }, db.TermSort)
  utils.CheckError(c, err)
  
  if len(termMarks) == 0 {
    termMarks = []models.TermMark {}
  }
  
  return c.JSON(termMarks)
}
