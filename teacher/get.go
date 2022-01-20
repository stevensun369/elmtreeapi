package teacher

import (
	"github.com/gofiber/fiber/v2"

	// internal packages
	"backend-go/db"
	"backend-go/models"
	"backend-go/utils"

	// std

	// mongodb
	"go.mongodb.org/mongo-driver/bson"
)

// @desc    Get all students that a non-homeroom teacher can interact with
// @route   GET /api/teacher/students
// @access  Private
func getStudents(c *fiber.Ctx) error {
  var subjectList []models.Subject
  utils.GetLocals(c.Locals("subjectList"), &subjectList)

  var subjectIDList []string
  for _, subject := range subjectList {
    subjectIDList = append(subjectIDList, subject.SubjectID)
  }

  students, err := db.GetStudents(bson.M{
    "subjectList.subjectID":  bson.M{"$in": subjectIDList},
  }, db.EmptySort)
  if err != nil {
    return utils.Error(c, err)
  }

  if len(students) == 0 {
    students = []models.Student {}
  }

  return c.JSON(students)
}

// @desc    Get marks
// @route   GET /api/teacher/marks/:subjectID/:studentID
// @access  Private
func getMarks(c *fiber.Ctx) error {
  subjectID := c.Params("subjectID")
  studentID := c.Params("studentID")

  marks, err := db.GetMarks(bson.M{
    "subject.subjectID": subjectID,
   "studentID": studentID,
  }, db.DateSort)
  if err != nil {
    return utils.Error(c, err)
  }

  if len(marks) == 0 {
    marks = []models.Mark {}
  }

  return c.JSON(marks)
}

// @desc    Get truancies
// @route   GET /api/teacher/truancies/:subjectID/:studentID
// @access  Private
func getTruancies(c *fiber.Ctx) error {
  subjectID := c.Params("subjectID")
  studentID := c.Params("studentID")

  truancies, err := db.GetTruancies(bson.M{
    "subject.subjectID": subjectID,
   "studentID": studentID,
  }, db.DateSort)
  if err != nil {
    return utils.Error(c, err)
  }

  if len(truancies) == 0 {
    truancies = []models.Truancy {}
  }

  return c.JSON(truancies)
}

// @desc    Get Average Marks by subjectID for multiple students
// @route   GET /api/teacher/averagemarks/:subjectID
// @access  Private
func getAverageMarks(c *fiber.Ctx) error {
  subjectID := c.Params("subjectID")

  averageMarks, err := db.GetAverageMarks(bson.M{
    "subject.subjectID": subjectID,
  }, db.TermSort)
  if err != nil {
    return utils.Error(c, err)
  } 

  if len(averageMarks) == 0 {
    averageMarks = []models.AverageMark {}
  }

  return c.JSON(averageMarks)
}

// @desc    Get final marks by subjectID for one studet
// @route   GET /api/teacher/finalmarks/:subjectID/:studentID
// @access  Private
func getFinalMarks(c *fiber.Ctx) error {
  subjectID := c.Params("subjectID")
  studentID := c.Params("studentID")

  finalMarks, err := db.GetFinalMarks(bson.M{
    "subject.subjectID": subjectID,
    "studentID": studentID,
  }, db.EmptySort)
  if err != nil {
    return utils.Error(c, err)
  }

  if len(finalMarks) == 0 {
    finalMarks = []models.FinalMark {}
  }

  return c.JSON(finalMarks)
}