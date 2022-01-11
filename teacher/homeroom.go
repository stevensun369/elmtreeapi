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
  var homeroomGrade models.Grade
  utils.GetLocals(c.Locals("homeroomGrade"), &homeroomGrade)

  // getting the homeroom students
  students, err := db.GetStudents(bson.M{
    "grade.gradeID": homeroomGrade.GradeID,
  }, db.EmptySort)
  utils.CheckError(c, err)
  
  if len(students) == 0 {
    students = []models.Student {}
  }

  return c.JSON(students)
}

// @desc   Get subjects of a homeroom student
// @route  GET /api/teacher/homeroom/:studentID
// @access Private
func getHomeroomStudentSubjects(c *fiber.Ctx) error {
  studentID := c.Params("studentID")

  student, err := db.GetStudentByID(studentID)
  utils.CheckError(c, err)

  return c.JSON(student.SubjectList)
}

// @desc   Get average marks for the students in the homeroom grade
// @route  GET /api/teacher/homeroom/average
// @access Private
func getHomeroomAverageMarks(c *fiber.Ctx) error {

  // getting the homeroom grade
  var homeroomGrade models.Grade
  utils.GetLocals(c.Locals("homeroomGrade"), &homeroomGrade)

  averageMarks, err := db.GetAverageMarks(bson.M{
    "grade.gradeID": homeroomGrade.GradeID,
  }, db.EmptySort)
  utils.CheckError(c, err)

  if len(averageMarks) == 0 {
    averageMarks = []models.AverageMark {}
  }

  return c.JSON(averageMarks)
}


// @desc   Get term marks fro homeroom students
// @route  GET /api/teacher/homeroom/term
// @access Private
func getHomeroomTermMarks(c *fiber.Ctx) error {

  // getting the homeroom grade
  var homeroomGrade models.Grade
  utils.GetLocals(c.Locals("homeroomGrade"), &homeroomGrade)

  termMarks, err := db.GetTermMarks(bson.M{
    "grade.gradeID": homeroomGrade.GradeID,
  }, db.EmptySort)
  utils.CheckError(c, err)

  if len(termMarks) == 0 {
    termMarks = []models.TermMark {}
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

  // checking that there isn't already a termMark
  var termMarkTest models.TermMark
  db.TermMarks.FindOne(context.Background(), bson.M{
    "studentID": studentID,
    "term": termInt,
  }).Decode(&termMarkTest)

  if (termMarkTest != models.TermMark{Grade: models.Grade{}}) {
    return utils.MessageError(c, "Exista deja o medie pe semestrul " + term)
  }

  // get student and his subjectList
  student, err := db.GetStudentByID(studentID)
  utils.CheckError(c, err)

  averageMarks, err := db.GetAverageMarks(bson.M{
    "studentID": studentID,
    "term": termInt,
  }, db.EmptySort)
  utils.CheckError(c, err)

  if len(averageMarks) != len(student.SubjectList) {
    return utils.MessageError(c, "Nu toate mediile au fost incheiate pe semestrul " + term)
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
  var termMarkGenID models.TermMark
  db.TermMarks.FindOne(context.Background(), bson.M{"termMarkID": termMarkID}).Decode(&termMarkGenID)
  for (termMarkGenID != models.TermMark{Grade: models.Grade{}}) {
    termMarkID = utils.GenID()
    db.TermMarks.FindOne(context.Background(), bson.M{"termMarkID": termMarkID}).Decode(&termMarkGenID)
  } 

  termMark := models.TermMark{
    TermMarkID: termMarkID,
    Value: value,
    StudentID: studentID,
    Grade: student.Grade,
    Term: termInt,
  }

  insertedResult, err := db.TermMarks.InsertOne(context.Background(), termMark)
  utils.CheckError(c, err)

  return c.JSON(bson.M{
    "_id": insertedResult.InsertedID,
    "termMarkID": termMarkID,
    "value": value,
    "studentID": studentID,
    "grade": student.Grade,
    "term": termInt,
  })
}

