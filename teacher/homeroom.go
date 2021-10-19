package teacher

import (
	"github.com/gofiber/fiber/v2"

	// internal pacakges
	"backend-go/db"
	"backend-go/models"
	"backend-go/utils"

	// std
	"context"
	"encoding/json"
	"fmt"
	"math"
	"strconv"

	// mongodb
	"go.mongodb.org/mongo-driver/bson"
)

// @desc   Get homeroom students
// @route  GET /api/teacher/homeroom
// @access Private
func getHomeroomStudents(c *fiber.Ctx) error {

  // getting the homeroom grade
  homeroomGradeLocals := fmt.Sprintf("%v", c.Locals("homeroomGrade"))
  var homeroomGrade models.Grade
  json.Unmarshal([]byte(homeroomGradeLocals), &homeroomGrade)

  // getting the homeroom students
  var students []models.Student
  studentsCollection, err := db.GetCollection("students")
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  cursor, err := studentsCollection.Find(context.Background(), bson.M{
    "grade.gradeID": homeroomGrade.GradeID,
  })
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  if err = cursor.All(context.Background(), &students); err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  return c.JSON(students)
}

// @desc   Get subjects of a homeroom student
// @route  GET /api/teacher/homeroom/:studentID
// @access Private
func getHomeroomStudentSubjects(c *fiber.Ctx) error {
  studentID := c.Params("studentID")

  // getting the students
  var student models.Student
  studentsCollection, err := db.GetCollection("students")
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  studentsCollection.FindOne(context.Background(), bson.M{
    "studentID": studentID,
  }).Decode(&student)

  return c.JSON(student.SubjectList)
}

// @desc   Get average marks for the students in the homeroom grade
// @route  GET /api/teacher/homeroom/average
// @access Private
func getHomeroomAverageMarks(c *fiber.Ctx) error {

  // getting the homeroom grade
  homeroomGradeLocals := fmt.Sprintf("%v", c.Locals("homeroomGrade"))
  var homeroomGrade models.Grade
  json.Unmarshal([]byte(homeroomGradeLocals), &homeroomGrade)


  var averageMarks []models.AverageMark
  averageMarksCollection, err := db.GetCollection("averagemarks")
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  cursor, err := averageMarksCollection.Find(context.Background(), bson.M{
    "grade.gradeID": homeroomGrade.GradeID,
  })
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  if err = cursor.All(context.Background(), &averageMarks); err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  return c.JSON(averageMarks)
}


// @desc   Get term marks fro homeroom students
// @route  GET /api/teacher/homeroom/term
// @access Private
func getHomeroomTermMarks(c *fiber.Ctx) error {

  // getting the homeroom grade
  homeroomGradeLocals := fmt.Sprintf("%v", c.Locals("homeroomGrade"))
  var homeroomGrade models.Grade
  json.Unmarshal([]byte(homeroomGradeLocals), &homeroomGrade)

  // getting the term marks
  termMarksCollection, err := db.GetCollection("termmarks")
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }
  var termMarks []models.TermMark

  cursor, err := termMarksCollection.Find(context.Background(), bson.M{
    "grade.gradeID": homeroomGrade.GradeID,
  })
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  if err = cursor.All(context.Background(), &termMarks); err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  return c.JSON(termMarks)
}
  
// @desc   Create term mark for the students at term
// @route  POST /api/teacher/homeroom/term
// @access Private
func createHomeroomTermMark(c *fiber.Ctx) error {
  // getting studentID and term from body
  var body map[string]string
  json.Unmarshal(c.Body(), &body)
  studentID := body["studentID"]
  term := body["term"]
  termInt, _ := strconv.Atoi(term)

  // term Marks Collection
  termMarksCollection, err := db.GetCollection("termmarks")
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  // checking that there isn't already a termMark
  var termMarkTest models.TermMark
  termMarksCollection.FindOne(context.Background(), bson.M{
    "studentID": studentID,
    "term": termInt,
  }).Decode(&termMarkTest)

  if (termMarkTest != models.TermMark{Grade: models.Grade{}}) {
    return c.Status(500).SendString("Exista deja o medie pe semestrul " + term)
  }

  // get student and his subjectList
  student, err := models.GetStudentById(studentID)
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  // get averageMarks
  var averageMarks []models.AverageMark
  averageMarksCollection, err := db.GetCollection("averagemarks")
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }
  
  cursor, err := averageMarksCollection.Find(context.Background(), bson.M{
    "studentID": studentID,
    "term": termInt,
  })
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  if err = cursor.All(context.Background(), &averageMarks); err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  if len(averageMarks) != len(student.SubjectList) {
    return c.Status(500).SendString("Nu toate mediile au fost incheiate pe semestrul" + term)
  }

  var value float64 = 0
  for _, averageMark := range averageMarks {
    value += float64(averageMark.Value)
  }
  // still have to find out how it is going to work: whether I have to round, or to floor. 
  // I'll go with floor, for now. I hope I remember
  value = math.Floor(float64(value) / float64(len(student.SubjectList)) * 100) / 100

  // generating id
  var termMarkID = utils.GenID()
  termMarkID = utils.GenID()
  fmt.Println(termMarkID)
  var termMarkGenID models.TermMark
  termMarksCollection.FindOne(context.Background(), bson.M{"termMarkID": termMarkID}).Decode(&termMarkGenID)
  for (termMarkGenID != models.TermMark{Grade: models.Grade{}}) {
    termMarkID = utils.GenID()
    termMarksCollection.FindOne(context.Background(), bson.M{"termMarkID": termMarkID}).Decode(&termMarkGenID)
    fmt.Println(termMarkID)
  } 

  termMark := models.TermMark{
    TermMarkID: termMarkID,
    Value: value,
    StudentID: studentID,
    Grade: student.Grade,
    Term: termInt,
  }

  insertedResult, err := termMarksCollection.InsertOne(context.Background(), termMark)
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  return c.JSON(bson.M{
    "_id": insertedResult.InsertedID,
    "termMarkID": termMarkID,
    "value": value,
    "studentID": studentID,
    "grade": student.Grade,
    "term": termInt,
  })
}

