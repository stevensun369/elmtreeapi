package teacher

import (
	"github.com/gofiber/fiber/v2"

	// internal packages
	"backend-go/db"
	"backend-go/models"
	"backend-go/utils"

	// std
	"context"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"strconv"

	// godotenv
	"github.com/joho/godotenv"

	// mongodb
	"go.mongodb.org/mongo-driver/bson"
)

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

  subjectListLocals := fmt.Sprintf("%v", c.Locals("subjectList"))
  var subjectList []models.Subject
  json.Unmarshal([]byte(subjectListLocals), &subjectList)

  // getting the currSubject
  var currSubject models.Subject
  for _, subject := range subjectList {
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
  var markGenID models.Mark
  marksCollection.FindOne(context.Background(), bson.M{"markID": markID}).Decode(&markGenID)
  for (markGenID != models.Mark{Subject: models.ShortSubject{}, Grade: models.Grade{}}) {
    markID = utils.GenID()
    marksCollection.FindOne(context.Background(), bson.M{"markID": markID}).Decode(&markGenID)
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

  subjectListLocals := fmt.Sprintf("%v", c.Locals("subjectList"))
  var subjectList []models.Subject
  json.Unmarshal([]byte(subjectListLocals), &subjectList)

  // getting the currSubject
  var currSubject models.Subject
  for _, subject := range subjectList {
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
  var truancyGenID models.Truancy
  truancysCollection.FindOne(context.Background(), bson.M{"truancyID": truancyID}).Decode(&truancyGenID)
  for (truancyGenID != models.Truancy{Subject: models.ShortSubject{}, Grade: models.Grade{}}) {
    truancyID = utils.GenID()
    truancysCollection.FindOne(context.Background(), bson.M{"truancyID": truancyID}).Decode(&truancyGenID)
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


  subjectListLocals := fmt.Sprintf("%v", c.Locals("subjectList"))
  var subjectList []models.Subject
  json.Unmarshal([]byte(subjectListLocals), &subjectList)

  // getting the currSubject
  var currSubject models.Subject
  for _, subject := range subjectList {
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
  var averageMarkGenID models.AverageMark
  averageMarksCollection.FindOne(context.Background(), bson.M{"averageMarkID": averageMarkID}).Decode(&averageMarkGenID)
  for (averageMarkGenID != models.AverageMark{Subject: models.ShortSubject{}, Grade: models.Grade{}}) {
    averageMarkID = utils.GenID()
    averageMarksCollection.FindOne(context.Background(), bson.M{"averageMarkID": averageMarkID}).Decode(&averageMarkGenID)
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