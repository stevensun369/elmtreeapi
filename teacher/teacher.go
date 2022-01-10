package teacher

import (
	// internal backend package
	"backend-go/db"
	"backend-go/models"
	"backend-go/utils"

	// std
	"context"
	"encoding/json"
	"fmt"

	// env

	// fiber
	"github.com/gofiber/fiber/v2"

	// mongodb
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @desc    Updates teacher subject and homeroom grade
// @route   GET /api/teacher/update
// @access  Private
func update(c *fiber.Ctx) error {
  teacherIDLocals := fmt.Sprintf("%v", c.Locals("teacherID"))
  var teacherID string
  json.Unmarshal([]byte(teacherIDLocals), &teacherID)

  teacher := db.GetTeacherByID(teacherID)

  tokenString, err := utils.TeacherGenerateToken(teacher.TeacherID, teacher.HomeroomGrade, teacher.SubjectList)
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
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
  subjectListLocals := fmt.Sprintf("%v", c.Locals("subjectList"))
  var subjectList []models.Subject
  json.Unmarshal([]byte(subjectListLocals), &subjectList)

  var subjectIDList []string
  for _, subject := range subjectList {
    subjectIDList = append(subjectIDList, subject.SubjectID)
  }

  // getting the students themselves
  var students []models.Student
  studentsCollection, err := db.GetCollection("students")
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }
  
  cursor, err := studentsCollection.Find(context.Background(), bson.M{
    "subjectList.subjectID":  bson.M{"$in": subjectIDList},
  })
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  if err = cursor.All(context.Background(), &students); err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
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
  var students []models.Student 
  studentsCollection, err := db.GetCollection("students")
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }
  
  subjectID := c.Params("subjectID")

  options := options.Find()
  options.SetSort(bson.D{{Key: "lastName", Value: 1}})
  cursor, err := studentsCollection.Find(context.Background(), bson.M{
    "subjectList.subjectID": subjectID,
  }, options)
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  if err = cursor.All(context.Background(), &students); err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  if len(students) == 0 {
    students = []models.Student {}
  }

  return c.JSON(students)
}