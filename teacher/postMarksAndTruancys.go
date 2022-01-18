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
	"strconv"

	// godotenv

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

  // deciding the term of the mark
  isTermOne := utils.IsTermOne(dateDay, dateMonth)

  // getting and checking if there are averageMarks
  averageMarkTermOne := false
  averageMarkTermTwo := false
  var averageMarks []models.AverageMark

  averageMarks, err := db.GetAverageMarks(bson.M{
    "subject.subjectID": subjectID, 
    "studentID": studentID,
  }, db.EmptySort)
  if err != nil {
    utils.Error(c, err)
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
    return utils.MessageError(c, "Media pe ambele semestre a fost încheiată.")
  } else {
    if averageMarkTermOne && isTermOne {
      return utils.MessageError(c, "Media pe primul semestru a fost încheiată.")
    } else if averageMarkTermTwo && !isTermOne {
      return utils.MessageError(c, "Media pe al doilea semestru a fost încheiată.")
    }
  }

  var subjectList []models.Subject
  utils.GetLocals(c.Locals("subjectList"), &subjectList)

  // getting the currSubject
  var currSubject models.Subject
  for _, subject := range subjectList {
    if subject.SubjectID == subjectID {
      currSubject = subject
    }
  }

  // generating id
  var markID = utils.GenID()
  markID = utils.GenID()
  var markGenID models.Mark
  db.Marks.FindOne(context.Background(), bson.M{"markID": markID}).Decode(&markGenID)
  for (markGenID != models.Mark{Subject: models.ShortSubject{}, Grade: models.Grade{}}) {
    markID = utils.GenID()
    db.Marks.FindOne(context.Background(), bson.M{"markID": markID}).Decode(&markGenID)
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
  insertedResult, err := db.Marks.InsertOne(context.Background(), mark)
  if err != nil {
    utils.Error(c, err)
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

  // deciding the term of the truancy
  isTermOne := utils.IsTermOne(dateDay, dateMonth)

  var subjectList []models.Subject
  utils.GetLocals(c.Locals("subjectList"), &subjectList)

  // getting the currSubject
  var currSubject models.Subject
  for _, subject := range subjectList {
    if subject.SubjectID == subjectID {
      currSubject = subject
    }
  }


  // generating id
  var truancyID = utils.GenID()
  truancyID = utils.GenID()
  var truancyGenID models.Truancy
  db.Truancies.FindOne(context.Background(), bson.M{"truancyID": truancyID}).Decode(&truancyGenID)
  for (truancyGenID != models.Truancy{Subject: models.ShortSubject{}, Grade: models.Grade{}}) {
    truancyID = utils.GenID()
    db.Truancies.FindOne(context.Background(), bson.M{"truancyID": truancyID}).Decode(&truancyGenID)
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
  insertedResult, err := db.Truancies.InsertOne(context.Background(), truancy)
  if err != nil {
    utils.Error(c, err)
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

  // getting the truancy and updating it
  var truancy models.Truancy
  db.Truancies.FindOneAndUpdate(context.Background(), bson.M{"truancyID": body["truancyID"]}, bson.D{
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
  term, _ := strconv.Atoi(body["term"])

  var subjectList []models.Subject
  utils.GetLocals(c.Locals("subjectList"), &subjectList)

  // getting the currSubject
  var currSubject models.Subject
  for _, subject := range subjectList {
    if subject.SubjectID == subjectID {
      currSubject = subject
    }
  }

  // generating id
  var averageMarkID = utils.GenID()
  averageMarkID = utils.GenID()
  var averageMarkGenID models.AverageMark
  db.AverageMarks.FindOne(context.Background(), bson.M{"averageMarkID": averageMarkID}).Decode(&averageMarkGenID)
  for (averageMarkGenID != models.AverageMark{Subject: models.ShortSubject{}, Grade: models.Grade{}}) {
    averageMarkID = utils.GenID()
    db.AverageMarks.FindOne(context.Background(), bson.M{"averageMarkID": averageMarkID}).Decode(&averageMarkGenID)
  } 

  // getting normal marks
  marks, err := db.GetMarks(bson.M{
    "subject.subjectID": subjectID,
    "studentID": studentID,
    "term": term,
  }, db.EmptySort)
  if err != nil {
    utils.Error(c, err)
  }

  // marksValues
  var marksValues []int
  for _, mark := range marks {
    marksValues = append(marksValues, mark.Value)
  }

  // getting final mark
  var finalMark models.FinalMark
  err = db.FinalMarks.FindOne(context.Background(), bson.M{
    "subject.subjectID": subjectID,
    "studentID": studentID,
    "term": term,
  }).Decode(&finalMark)
  if err != nil {
    utils.Error(c, err)
  }
  finalMarkValue := finalMark.Value

  value := utils.CalculateAverageMark(marksValues, finalMarkValue)

  averageMark := models.AverageMark{
    AverageMarkID: averageMarkID,
    Value: value,
    Subject: models.ShortSubject{
      SubjectID: currSubject.SubjectID,
      Name: currSubject.Name,
    },
    StudentID: studentID,
    Grade: currSubject.Grade,
    Term: term,
  }

  insertedResult, err := db.AverageMarks.InsertOne(context.Background(), averageMark)
  if err != nil {
    utils.Error(c, err)
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
    "term": term,
  })
}

// @desc    Create finalMarks
// @route   POST /api/teacher/final
// @access  Private
func createFinalMark(c *fiber.Ctx) error {
  var body map[string]string
  json.Unmarshal(c.Body(), &body)

  value := body["value"]
  subjectID := body["subjectID"]
  studentID := body["studentID"]
  term,_ := strconv.Atoi(body["term"]) 

  // getting and checking if there are averageMarks
  averageMarkTermOne := false
  averageMarkTermTwo := false
  var averageMarks []models.AverageMark

  averageMarks, err := db.GetAverageMarks(bson.M{
    "subject.subjectID": subjectID, 
    "studentID": studentID,
  }, db.EmptySort)
  if err != nil {
    utils.Error(c, err)
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
    return utils.MessageError(c, "Media pe ambele semestre a fost încheiată.")
  } else {
    if averageMarkTermOne && term == 1 {
      return utils.MessageError(c, "Media pe primul semestru a fost încheiată.")
    } else if averageMarkTermTwo && term == 2 {
      return utils.MessageError(c, "Media pe al doilea semestru a fost încheiată.")
    }
  }

  var subjectList []models.Subject
  utils.GetLocals(c.Locals("subjectList"), &subjectList)

  // getting the currSubject
  var currSubject models.Subject
  for _, subject := range subjectList {
    if subject.SubjectID == subjectID {
      currSubject = subject
    }
  }

  // generating id
  var finalMarkID = utils.GenID()
  finalMarkID = utils.GenID()
  var finalMarkGenID models.FinalMark
  db.FinalMarks.FindOne(context.Background(), bson.M{"finalMarkID": finalMarkID}).Decode(&finalMarkGenID)
  for (finalMarkGenID != models.FinalMark{Subject: models.ShortSubject{}, Grade: models.Grade{}}) {
    finalMarkID = utils.GenID()
    db.FinalMarks.FindOne(context.Background(), bson.M{"finalMarkID": finalMarkID}).Decode(&finalMarkGenID)
  } 

  // some prep variables for the finalMark struct
  valueInt, _ := strconv.Atoi(value) 

  // filling in all the values to a finalMark struct
  finalMark := models.FinalMark{
    FinalMarkID: finalMarkID,
    Value: valueInt,
    Subject: models.ShortSubject {
      SubjectID: subjectID,
      Name: currSubject.Name,
    },
    StudentID: studentID,
    Grade: currSubject.Grade,
    Term: term, 
  }

  // inserting the finalMark
  insertedResult, err := db.FinalMarks.InsertOne(context.Background(), finalMark)
  if err != nil {
    utils.Error(c, err)
  }

  return c.JSON(bson.M{
    "_id": insertedResult.InsertedID,
    "finalMarkID": finalMarkID,
    "value": valueInt,
    "subject": bson.M {
      "subjectID": subjectID,
      "name": currSubject.Name,
    },
    "studentID": studentID,
    "grade": currSubject.Grade,
    "term": term, 
  })
}