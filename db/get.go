package db

import (
	"backend-go/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// student
func GetStudentByID(studentID string) (models.Student, error) {
	// getting a student
	var student models.Student
  
  err := Students.FindOne(context.Background(), bson.M{"studentID": studentID}).Decode(&student)
	
	return student, err
}

func GetStudentByCNP(cnp string) (models.Student, error) {
	// getting a student
	var student models.Student
  
  err := Students.FindOne(context.Background(), bson.M{"cnp": cnp}).Decode(&student)
	
	return student, err
}

func GetStudentByIDAndCNP(cnp string, studentID string) (models.Student, error) {
	var student models.Student
	
  err := Students.FindOne(context.Background(), bson.M{
		"cnp": cnp, 
		"studentID": studentID,
	}).Decode(&student)

	return student, err
}

// teacher
func GetTeacherByID(teacherID string) (models.Teacher, error) {
	// getting a teacher
	var teacher models.Teacher

  err := Teachers.FindOne(context.Background(), bson.M{"teacherID": teacherID}).Decode(&teacher)
	
	return teacher, err
}

func GetTeacherByCNP(cnp string) (models.Teacher, error) {
	// getting a teacher
	var teacher models.Teacher

  err := Teachers.FindOne(context.Background(), bson.M{"cnp": cnp}).Decode(&teacher)
	
	return teacher, err
}

// marks
func GetMarks(filter interface{}, optionsData interface{}) ([]models.Mark, error) {
	var marks []models.Mark

	options := options.Find()
	options.SetSort(optionsData)

	cursor, err := Marks.Find(context.Background(), filter, options)
	if err != nil {
		return nil, err
	}
	err = cursor.All(context.Background(), &marks)
	return marks, err
}

// truancies
func GetTruancies(filter interface{}, optionsData interface{}) ([]models.Truancy, error) {
	var truancies []models.Truancy

	options := options.Find()
	options.SetSort(optionsData)

	cursor, err := Truancies.Find(context.Background(), filter, options)
	if err != nil {
		return nil, err
	}
	err = cursor.All(context.Background(), &truancies)
	return truancies, err
}

// average marks
func GetAverageMarks(filter interface{}, optionsData interface{}) ([]models.AverageMark, error) {
	var averageMarks []models.AverageMark

	options := options.Find()
	options.SetSort(optionsData)

	cursor, err := AverageMarks.Find(context.Background(), filter, options)
	if err != nil {
		return nil, err
	}
	err = cursor.All(context.Background(), &averageMarks)
	return averageMarks, err
}

// term marks
func GetTermMarks(filter interface{}, optionsData interface{}) ([]models.TermMark, error) {
	var termMarks []models.TermMark

	options := options.Find()
	options.SetSort(optionsData)

	cursor, err := TermMarks.Find(context.Background(), filter, options)
	if err != nil {
		return nil, err
	}
	err = cursor.All(context.Background(), &termMarks)
	return termMarks, err
}

// periods
func GetPeriods(filter interface{}, optionsData interface{}) ([]models.Period, error) {
	var periods []models.Period

	options := options.Find()
	options.SetSort(optionsData)

	cursor, err := Periods.Find(context.Background(), filter, options)
	if err != nil {
		return nil, err
	}
	err = cursor.All(context.Background(), &periods)
	return periods, err
}

// students
func GetStudents(filter interface{}, optionsData interface{}) ([]models.Student, error) {
	var students []models.Student

	options := options.Find()
	options.SetSort(optionsData)

	cursor, err := Students.Find(context.Background(), filter, options)
	if err != nil {
		return nil, err
	}
	err = cursor.All(context.Background(), &students)
	return students, err
}