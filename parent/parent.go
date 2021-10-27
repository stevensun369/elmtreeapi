package parent

import (
	"backend-go/db"
	"backend-go/models"
	"backend-go/utils"
	"context"
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func parentAddStudent(c *fiber.Ctx) error {
  var body map[string]string
  json.Unmarshal(c.Body(), &body)

  // unmarshaling parent id
  var parentID string
  json.Unmarshal([]byte(fmt.Sprintf("%v", c.Locals("parentID"))), &parentID)

  studentsCollection, err := db.GetCollection("students")
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  parentsCollection, err := db.GetCollection("parents")
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  // get student and make sure it is there
  var student models.Student
  studentsCollection.FindOne(context.Background(), bson.M{"cnp": body["studentCNP"], "studentID": body["studentID"]}).Decode(&student)
  if (student.StudentID == "") {
    return c.Status(404).JSON(bson.M{
      "message": "Nu există niciun elev cu datele introduse.",
    })
  }

  newStudentIDListJSON := []byte(fmt.Sprintf("%v", c.Locals("studentIDList")))
  var newStudentIDList []string
  json.Unmarshal(newStudentIDListJSON, &newStudentIDList)
  newStudentIDList = append(newStudentIDList, body["studentID"])

  var parent models.Parent
  if err = parentsCollection.FindOneAndUpdate(context.Background(), bson.M{"parentID": parentID}, bson.M{"$set": bson.M{"studentIDList": newStudentIDList}}).Decode(&parent); err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }
  parent.StudentIDList = append(parent.StudentIDList, body["studentID"])

  var students []models.Student
  options := options.Find()
  options.SetSort(bson.D{{Key: "grade.gradeNumber", Value: 1}, {Key: "grade.gradeLetter", Value: 1}, {Key: "firstName", Value: 1}})
  cursor, err := studentsCollection.Find(context.Background(), bson.M{
    "studentID": bson.M{"$in": newStudentIDList},
  }, options)
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }
  if err = cursor.All(context.Background(), &students); err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  tokenString, err := utils.ParentGenerateToken(parent.ParentID, parent.StudentIDList)
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
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