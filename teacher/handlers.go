package teacher

import (
	// internal backend package
	"backend-go/db"
	"backend-go/models"
	"math"
	"strconv"
	"strings"

	// internal utils
	"backend-go/utils"

	"context"
	"encoding/json"
	"fmt"

	// jwt
	"github.com/dgrijalva/jwt-go"

	// env
	"os"

	"github.com/joho/godotenv"

	// fiber
	"github.com/gofiber/fiber/v2"

	// bcrypt
	"golang.org/x/crypto/bcrypt"

	// mongodb
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// teacher auth middleware
func teacherMiddleware(c *fiber.Ctx) error {
  var token string

  authHeader := c.Get("Authorization")

  if string(authHeader) != "" && strings.HasPrefix(string(authHeader), "Bearer") {
    token = strings.Fields(string(authHeader))[1]

    // we're just parsing the token: maybe I will put it in the utils
    claims := &utils.Claims{}
    tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface {}, error) {
      return utils.JWTKey, nil
    })

    if err != nil {
      return c.Status(500).SendString(fmt.Sprintf("%v", err))
    }

    if !tkn.Valid {
      return c.Status(500).SendString("token not valid")
    }

    c.Locals("userID", claims.ID)
  }

  if (token == "") {
    return c.Status(500).SendString("no token")
  }

  return c.Next()
}

// @desc   Login teacher
// @route  POST /api/teacher/login
// @access Public
func postLogin(c *fiber.Ctx) error {

  // getting body and unmarshalling it into a body variable
  var body map[string]string
  json.Unmarshal(c.Body(), &body)

  // getting the db collection
  collection, err := db.GetCollection("teachers")
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  // getting the teacher
  var teacher bson.M
  if err = collection.FindOne(context.Background(), bson.M{"cnp": body["cnp"]}).Decode(&teacher); err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  // if there is no teacher, we 500
  if teacher == nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }


  // if the teacher doesn't have a password
  if teacher["password"] == "" {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body["password"]), 10)

    if err != nil {
      return c.Status(500).SendString(fmt.Sprintf("%v", err))
    }

    var modifiedTeacher bson.M
    collection.FindOneAndUpdate(context.Background(), bson.M{"cnp": body["cnp"]}, bson.D{
      {Key: "$set", Value: bson.D{{Key: "password",Value: string(hashedPassword)}}},
    }).Decode(&modifiedTeacher)
    
    // jwt
    tokenString, err := utils.GenerateToken(fmt.Sprintf("%v", modifiedTeacher["teacherID"]))
    if err != nil {
      return c.Status(500).SendString(fmt.Sprintf("%v", err))
    }

    return c.JSON(bson.M{
      "teacherID": modifiedTeacher["teacherID"],
      "firstName": modifiedTeacher["firstName"],
      "lastName": modifiedTeacher["lastName"],
      "cnp": modifiedTeacher["cnp"],
      "password": modifiedTeacher["password"],
      "homeroomGrade": modifiedTeacher["homeroomGrade"],
      "subjectList": modifiedTeacher["subjectList"],   
      "token": tokenString,
    })
  } else {
    hashedPassword := fmt.Sprintf("%v", teacher["password"])

    compareErr := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(body["password"]))

    tokenString, err := utils.GenerateToken(fmt.Sprintf("%v", teacher["teacherID"]))
    if err != nil {
      return c.Status(500).SendString(fmt.Sprintf("%v", err))
    }

    if compareErr == nil {
      return c.JSON(bson.M{
        "teacherID": teacher["teacherID"],
        "firstName": teacher["firstName"],
        "lastName": teacher["lastName"],
        "cnp": teacher["cnp"],
        "password": teacher["password"],
        "homeroomGrade": teacher["homeroomGrade"],
        "subjectList": teacher["subjectList"],
        "token": tokenString,
      })
    } else {
      return c.Status(500).SendString(fmt.Sprintf("%v", compareErr))
    }
  } 
}

// @desc   Get homeroom info
// @route  GET /api/teacher/homeroom/info
// @access Private
func getHomeroomInfo(c *fiber.Ctx) error {
  userID := fmt.Sprintf("%v", c.Locals("userID"))

  // getting a teacher
  teacher, err := models.GetTeacherById(userID)
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  return c.JSON(teacher.HomeroomGrade)
}

// @desc   Get homeroom students
// @route  GET /api/teacher/homeroom
// @access Private
func getHomeroomStudents(c *fiber.Ctx) error {
  userID := fmt.Sprintf("%v", c.Locals("userID"))

  // getting a teacher
  teacher, err := models.GetTeacherById(userID)
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  // getting the homeroom students
  var students []models.Student
  studentsCollection, err := db.GetCollection("students")
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  cursor, err := studentsCollection.Find(context.Background(), bson.M{
    "grade.gradeID": teacher.HomeroomGrade.GradeID,
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
  userID := fmt.Sprintf("%v", c.Locals("userID"))

  teacher, err := models.GetTeacherById(userID)
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  var averageMarks []models.AverageMark
  averageMarksCollection, err := db.GetCollection("averageMarks")
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  cursor, err := averageMarksCollection.Find(context.Background(), bson.M{
    "grade.gradeID": teacher.HomeroomGrade.GradeID,
  })
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  if err = cursor.All(context.Background(), &averageMarks); err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  return c.JSON(averageMarks)
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

// @desc   Get term marks fro homeroom students
// @route  GET /api/teacher/homeroom/term
// @access Private
func getHomeroomTermMarks(c *fiber.Ctx) error {

  // user id and teacher
  userID := fmt.Sprintf("%v", c.Locals("userID"))
  teacher, err := models.GetTeacherById(userID)
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  // getrting the term marks
  termMarksCollection, err := db.GetCollection("termmarks")
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }
  var termMarks []models.TermMark

  cursor, err := termMarksCollection.Find(context.Background(), bson.M{
    "grade.gradeID": teacher.HomeroomGrade.GradeID,
  })
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  if err = cursor.All(context.Background(), &termMarks); err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  return c.JSON(termMarks)
}

// @desc    Get teacher subjects
// @route   GET /api/teacher/subjects
// @access  Private
func getSubjects(c *fiber.Ctx) error {
  userID := fmt.Sprintf("%v", c.Locals("userID"))
  teacher, err := models.GetTeacherById(userID)
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  return c.JSON(teacher.SubjectList)
}

// @desc    Get all students that a non-homeroom teacher can interact with
// @route   GET /api/teacher/students
// @access  Private
func getStudents(c *fiber.Ctx) error {
  userID := fmt.Sprintf("%v", c.Locals("userID"))
  teacher, err := models.GetTeacherById(userID)
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  var subjectIDList []string
  for _, subject := range teacher.SubjectList {
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

  return c.JSON(students)
}

// @desc    Create mark
// @route   POST /api/teacher/mark
// @access  Private
func createMark(c *fiber.Ctx) error {
  var body map[string]string
  json.Unmarshal(c.Body(), &body)

  value := body["value"]
  date := body["date"]
  subjectID := body["subjectID"]
  studentID := body["studentID"]

  // getting the dateDay and dateMonth
  dateDay := date[:2]
  dateMonth := date[3:5]

  // getting term change
  err := godotenv.Load(".env")
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }
  termDateDay := os.Getenv("TERM_CHANGE_DAY")
  termDateMonth := os.Getenv("TERM_CHANGE_MONTH")

  // deciding the term of the mark
  var isTermOne bool 
  if dateMonth < termDateMonth {
    isTermOne = true
  } else if (dateMonth > termDateMonth) {
    isTermOne = false
  } else { // dateMonth === termDateMonth
    if (dateDay <= termDateDay) {
      isTermOne = true
    } else { // dateDay > termDateDay
      isTermOne = false
    }
  }

  // // temporary: i don't know why it doesn't work
  // fmt.Println(isTermOne)

  // getting and checking if there are averageMarks
  averageMarkTermOne := false
  averageMarkTermTwo := false
  var averageMarks []models.AverageMark

  averageMarksCollection, err := db.GetCollection("averagemarks")
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }
  cursor, err := averageMarksCollection.Find(context.Background(), bson.M{
    "subject.subjectID": subjectID, 
    "studentID": studentID,
  })
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }
  if err = cursor.All(context.Background(), &averageMarks); err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  // now checking for averageMarks for each term
  for _, averageMark := range averageMarks {
    // this is not idiotic:
    // if I would have written averageMarkTermTwo = averageMarks[averageMark].term === 1
    // then it would have been resetted to false when it would have hopped to the other average mark
    if averageMark.Term == 1 {
      averageMarkTermOne = true
    } else if averageMark.Term == 2 {
      averageMarkTermTwo = true
    }
  }
  if averageMarkTermOne && averageMarkTermTwo {
    return c.Status(500).SendString("Media pe ambele semestre a fost încheiată.")
  } else {
    if averageMarkTermOne && isTermOne {
      return c.Status(500).SendString("Media pe primul semestru a fost încheiată.")
    } else if averageMarkTermTwo && !isTermOne {
      return c.Status(500).SendString("Media pe al doilea semestru a fost încheiată.")
    }
  }

  // getting the teacher for the currSubject
  teacherID := fmt.Sprintf("%v", c.Locals("userID"))
  teacher, err := models.GetTeacherById(teacherID)
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  // getting the currSubject
  var currSubject models.Subject
  for _, subject := range teacher.SubjectList {
    if subject.SubjectID == subjectID {
      currSubject = subject
    }
  }

  // getting the marksCollection for generating id
  marksCollection, err := db.GetCollection("marks")
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  // generating id
  var markID = utils.GenID()
  markID = utils.GenID()
  fmt.Println(markID)
  var markGenID models.Mark
  marksCollection.FindOne(context.Background(), bson.M{"markID": markID}).Decode(&markGenID)
  for (markGenID != models.Mark{Subject: models.ShortSubject{}, Grade: models.Grade{}}) {
    markID = utils.GenID()
    marksCollection.FindOne(context.Background(), bson.M{"markID": markID}).Decode(&markGenID)
    fmt.Println(markID)
  } 

  // some prep variables for the mark struct
  valueInt, _ := strconv.Atoi(value) 
  var term int
  if isTermOne {
  term = 1
  } else {
    term = 2
  }
  // filling in all the values to a mark struct
  mark := models.Mark{
    MarkID: markID,
    Value: valueInt,
    DateDay: dateDay,
    DateMonth: dateMonth,
    Subject: models.ShortSubject {
      SubjectID: subjectID,
      Name: currSubject.Name,
    },
    StudentID: studentID,
    Grade: currSubject.Grade,
    Term: term, 
  }

  // inserting the mark
  insertedResult, err := marksCollection.InsertOne(context.Background(), mark)
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  return c.JSON(bson.M{
    "_id": insertedResult.InsertedID,
    "markID": markID,
    "value": valueInt,
    "dateDay": dateDay,
    "dateMonth": dateMonth,
    "subject": bson.M {
      "subjectID": subjectID,
      "name": currSubject.Name,
    },
    "studentID": studentID,
    "grade": currSubject.Grade,
    "term": term, 
  })
}

// @desc    Create truancy
// @route   POST /api/teacher/truancy
// @access  Private
func createTruancy(c *fiber.Ctx) error {
  var body map[string]string
  json.Unmarshal(c.Body(), &body)

  date := body["date"]
  subjectID := body["subjectID"]
  studentID := body["studentID"]

  // getting the dateDay and dateMonth
  dateDay := date[:2]
  dateMonth := date[3:5]

  // getting term change
  err := godotenv.Load(".env")
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }
  termDateDay := os.Getenv("TERM_CHANGE_DAY")
  termDateMonth := os.Getenv("TERM_CHANGE_MONTH")

  // deciding the term of the truancy
  var isTermOne bool 
  if dateMonth < termDateMonth {
    isTermOne = true
  } else if (dateMonth > termDateMonth) {
    isTermOne = false
  } else { // dateMonth === termDateMonth
    if (dateDay <= termDateDay) {
      isTermOne = true
    } else { // dateDay > termDateDay
      isTermOne = false
    }
  }

  // getting the teacher for the currSubject
  teacherID := fmt.Sprintf("%v", c.Locals("userID"))
  teacher, err := models.GetTeacherById(teacherID)
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  // getting the currSubject
  var currSubject models.Subject
  for _, subject := range teacher.SubjectList {
    if subject.SubjectID == subjectID {
      currSubject = subject
    }
  }

  // getting the truancysCollection for generating id
  truancysCollection, err := db.GetCollection("truancies")
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  // generating id
  var truancyID = utils.GenID()
  truancyID = utils.GenID()
  fmt.Println(truancyID)
  var truancyGenID models.Truancy
  truancysCollection.FindOne(context.Background(), bson.M{"truancyID": truancyID}).Decode(&truancyGenID)
  for (truancyGenID != models.Truancy{Subject: models.ShortSubject{}, Grade: models.Grade{}}) {
    truancyID = utils.GenID()
    truancysCollection.FindOne(context.Background(), bson.M{"truancyID": truancyID}).Decode(&truancyGenID)
    fmt.Println(truancyID)
  } 

  // some prep variables for the truancy struct
  var term int
  if isTermOne {
  term = 1
  } else {
    term = 2
  }
  // filling in all the values to a truancy struct
  truancy := models.Truancy{
    TruancyID: truancyID,
    DateDay: dateDay,
    DateMonth: dateMonth,
    Subject: models.ShortSubject {
      SubjectID: subjectID,
      Name: currSubject.Name,
    },
    StudentID: studentID,
    Grade: currSubject.Grade,
    Term: term, 
    Motivated: false,
  }

  // inserting the truancy
  insertedResult, err := truancysCollection.InsertOne(context.Background(), truancy)
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  return c.JSON(bson.M{
    "_id": insertedResult.InsertedID,
    "truancyID": truancyID,
    "dateDay": dateDay,
    "dateMonth": dateMonth,
    "subject": bson.M {
      "subjectID": subjectID,
      "name": currSubject.Name,
    },
    "studentID": studentID,
    "grade": currSubject.Grade,
    "term": term, 
    "motivated": false,
  })
}

// @desc    Get marks
// @route   GET /api/teacher/mark
// @access  Private
func getMarks(c *fiber.Ctx) error {
  subjectID := c.Params("subjectID")
  studentID := c.Params("studentID")

  marksCollection, err := db.GetCollection("marks")
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }
  var marks []models.Mark

  options := options.Find()
  options.SetSort(bson.D{{Key: "dateMonth", Value: 1}, {Key: "dateDay", Value: 1}})
  cursor, err := marksCollection.Find(context.Background(), bson.M{
    "subject.subjectID": subjectID,
    "studentID": studentID,
  }, options)
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }
  if err = cursor.All(context.Background(), &marks); err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  return c.JSON(marks)
}

// @desc    Get marks
// @route   GET /api/teacher/mark
// @access  Private
func getTruancys(c *fiber.Ctx) error {
  subjectID := c.Params("subjectID")
  studentID := c.Params("studentID")

  truancysCollection, err := db.GetCollection("truancies")
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }
  var truancys []models.Truancy

  options := options.Find()
  options.SetSort(bson.D{{Key: "dateMonth", Value: 1}, {Key: "dateDay", Value: 1}})
  cursor, err := truancysCollection.Find(context.Background(), bson.M{
    "subject.subjectID": subjectID,
    "studentID": studentID,
  }, options)
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }
  if err = cursor.All(context.Background(), &truancys); err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }
  return c.JSON(truancys)
}

// @desc    Motivate truancy
// @route   PUT /api/teacher/truancy/motivate
// @access  Private
func motivateTruancy(c *fiber.Ctx) error {
  var body map[string]string
  json.Unmarshal(c.Body(), &body)

  // getting the truancys collection
  truancysCollection, err := db.GetCollection("truancies")
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  // gettint the truancy and updating it
  var truancy models.Truancy
  truancysCollection.FindOneAndUpdate(context.Background(), bson.M{"truancyID": body["truancyID"]}, bson.D{
    {Key: "$set", Value: bson.D{{Key: "motivated", Value: true}}},
  }).Decode(&truancy)

  return c.JSON(truancy)
}

// @desc    Create AverageMark
// @route   POST /api/teacher/average
// @access  Private
func createAverageMark(c *fiber.Ctx) error {
  var body map[string]string
  json.Unmarshal(c.Body(), &body)
  subjectID := body["subjectID"]
  studentID := body["studentID"]
  term := body["term"]
  termInt, err := strconv.Atoi(term)
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }


  // getting the teacher for the currSubject
  teacherID := fmt.Sprintf("%v", c.Locals("userID"))
  teacher, err := models.GetTeacherById(teacherID)
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  // getting the currSubject
  var currSubject models.Subject
  for _, subject := range teacher.SubjectList {
    if subject.SubjectID == subjectID {
      currSubject = subject
    }
  }

  // getting the averageMark collection
  averageMarksCollection, err := db.GetCollection("averagemarks")
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  // generating id
  var averageMarkID = utils.GenID()
  averageMarkID = utils.GenID()
  fmt.Println(averageMarkID)
  var averageMarkGenID models.AverageMark
  averageMarksCollection.FindOne(context.Background(), bson.M{"averageMarkID": averageMarkID}).Decode(&averageMarkGenID)
  for (averageMarkGenID != models.AverageMark{Subject: models.ShortSubject{}, Grade: models.Grade{}}) {
    averageMarkID = utils.GenID()
    averageMarksCollection.FindOne(context.Background(), bson.M{"averageMarkID": averageMarkID}).Decode(&averageMarkGenID)
    fmt.Println(averageMarkGenID)
  } 

  // getting the markscollection
  marksCollection, err := db.GetCollection("marks")
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }
    
  // getting all the related marks
  var marks []models.Mark
  cursor, err := marksCollection.Find(context.Background(), bson.M{
    "subject.subjectID": subjectID,
    "studentID": studentID,
    "term": termInt,
  })
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }
  if err = cursor.All(context.Background(), &marks); err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  markValuesCounter := 0
  marksCounter := 0
  for _, mark := range marks {
    markValuesCounter += mark.Value
    marksCounter += 1
  }

  // it rounds to the closest integer - cause that's how it works in real life
  value :=  int(math.Round(float64(markValuesCounter) / float64(marksCounter)))

  averageMark := models.AverageMark{
    AverageMarkID: averageMarkID,
    Value: value,
    Subject: models.ShortSubject{
      SubjectID: currSubject.SubjectID,
      Name: currSubject.Name,
    },
    StudentID: studentID,
    Grade: currSubject.Grade,
    Term: termInt,
  }

  insertedResult, err := averageMarksCollection.InsertOne(context.Background(), averageMark)
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  return c.JSON(bson.M{
    "_id": insertedResult.InsertedID,
    "averageMarkID": averageMarkID,
    "value": value,
    "subject": bson.M {
      "subjectID": currSubject.SubjectID,
      "name": currSubject.Name,
    },
    "studentID": studentID,
    "grade": currSubject.Grade,
    "term": termInt,
  })
}

// @desc    Get Average Mark
// @route   POST /api/teacher/average/:subjectID/:studentID
// @access  Private
func getAverageMarks(c *fiber.Ctx) error {
  subjectID := c.Params("subjectID")
  studentID := c.Params("studentID")

  averageMarksCollection, err := db.GetCollection("averagemarks")
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }
  var averageMarks []models.AverageMark

  options := options.Find()
  options.SetSort(bson.D{{Key: "term", Value: 1}})
  cursor, err := averageMarksCollection.Find(context.Background(), bson.M{
    "subject.subjectID": subjectID,
    "studentID": studentID,
  }, options)
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }
  if err = cursor.All(context.Background(), &averageMarks); err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  return c.JSON(averageMarks)
}

// @desc    Get Average Marks by subjectID for multiple students
// @route   POST /api/teacher/average/:subjectID
// @access  Private
func getAverageMarksSubject(c *fiber.Ctx) error {
  subjectID := c.Params("subjectID")

  averageMarksCollection, err := db.GetCollection("averagemarks")
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }
  var averageMarks []models.AverageMark

  options := options.Find()
  options.SetSort(bson.D{{Key: "term", Value: 1}})
  cursor, err := averageMarksCollection.Find(context.Background(), bson.M{
    "subject.subjectID": subjectID,
  }, options)
  if err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }
  if err = cursor.All(context.Background(), &averageMarks); err != nil {
    return c.Status(500).SendString(fmt.Sprintf("%v", err))
  }

  return c.JSON(averageMarks)
}
