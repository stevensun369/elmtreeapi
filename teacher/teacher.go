package teacher

import (
	// internal backend package
	"backend-go/db"
	"backend-go/models"
	"backend-go/utils"
	"context"

	// std

	// env

	// fiber
	"github.com/gofiber/fiber/v2"

	// mongodb
	"go.mongodb.org/mongo-driver/bson"
)

// @desc    Updates teacher subject and homeroom grade
// @route   GET /api/teacher/update
// @access  Private
func update(c *fiber.Ctx) error {
  // teacherID
  var teacherID string
  utils.GetLocals(c.Locals("teacherID"), &teacherID)

  teacher, err := db.GetTeacherByID(teacherID)
  if err != nil {
    return utils.Error(c, err)
  }

  tokenString, err := utils.TeacherGenerateToken(teacher.TeacherID, teacher.HomeroomGrade, teacher.SubjectList, teacher.SchoolID)
  if err != nil {
    return utils.Error(c, err)
  } 

  return c.JSON(bson.M{
    "subjectList": teacher.SubjectList,
    "homeroomGrade": teacher.HomeroomGrade,
    "token": tokenString,
  })
}

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

// @desc    Get students that have a subject
// @route   GET /api/teacher/:subjectID
// @access  Private
func getSubjectStudents(c *fiber.Ctx) error {
  subjectID := c.Params("subjectID")

  students, err := db.GetStudents(bson.M{
    "subjectList.subjectID": subjectID,
  }, db.LastNameSort)
  if err != nil {
    utils.Error(c, err)
  }

  if len(students) == 0 {
    students = []models.Student {}
  }

  return c.JSON(students)
}

// @desc    getting timetable
// @route   GET /api/teacher/timetable
// @access  Private
func getPeriods(c *fiber.Ctx) error {
  var subjectList []models.Subject
  utils.GetLocals(c.Locals("subjectList"), &subjectList)
 
  var subjectIDList []string
  for _, subject := range subjectList {
    subjectIDList = append(subjectIDList, subject.SubjectID)
  }
  
  periods, err := db.GetPeriods(bson.M{
    "subject.subjectID": bson.M{"$in": subjectIDList},
  }, db.PeriodSort)
  if err != nil {
    return utils.Error(c, err)
  }

  return c.JSON(periods)
}

// @desc    getting timetable
// @route   GET /api/teacher/school
// @access  Private
func getSchool(c *fiber.Ctx) error {
  var schoolID string
  utils.GetLocals(c.Locals("schoolID"), &schoolID)
 
  var school models.School
  err := db.Schools.FindOne(context.Background(), bson.M{
    "schoolID": schoolID,
  }).Decode(&school)
  if err != nil {
    return utils.Error(c, err)
  }
  
  return c.JSON(school)
}