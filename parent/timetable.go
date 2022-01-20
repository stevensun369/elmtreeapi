package parent

import (
	"backend-go/db"
	"backend-go/models"
	"backend-go/utils"
	"context"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

// @desc    Get periods
// @route   GET /api/parent/timetable/:studentID
// @access  Private
func getPeriods(c *fiber.Ctx) error {
  studentID := c.Params("studentID")

  student, err := db.GetStudentByID(studentID)
  if err != nil {
    utils.Error(c, err)
  }

  subjectList := student.SubjectList
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

// @desc   Get timetable
// @route  GET /api/parent/timetable/:studentID/teachers
// @access Private
func getPeriodsTeachers(c *fiber.Ctx) error {
  studentID := c.Params("studentID")

  student, err := db.GetStudentByID(studentID)
  if err != nil {
    return utils.Error(c, err)
  }

  subjectList := student.SubjectList
  var subjectIDList []string
  for _, subject := range subjectList {
    subjectIDList = append(subjectIDList, subject.SubjectID)
  }

  teachers, err := db.GetTeachers(bson.M{
    "subjectList.subjectID": bson.M{"$in" : subjectIDList},
  }, db.EmptySort)
  if err != nil {
    return utils.Error(c, err)
  }

  return c.JSON(teachers)
}

// @desc   Get school
// @route  GET /api/parent/school/:studentID/
// @access Private
func getSchool(c *fiber.Ctx) error {
  studentID := c.Params("studentID")

  student, err := db.GetStudentByID(studentID)
  if err != nil {
    return utils.Error(c, err)
  }
  
  var school models.School
  err = db.Schools.FindOne(context.Background(), bson.M{
    "schoolID": student.Grade.SchoolID,
  }).Decode(&school)
  if err != nil {
    return utils.Error(c, err)
  }

  return c.JSON(school)
}
