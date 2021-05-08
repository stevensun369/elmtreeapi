package models

import (
	"backend-go/db"
	"context"

	"gopkg.in/mgo.v2/bson"
)

// the password and the subjectList can be empty
type Student struct {
	StudentID   string           `json:"studentID,omitempty" bson:"studentID,omitempty"`
	FirstName   string           `json:"firstName,omitempty" bson:"firstName,omitempty"`
	DadInitials string           `json:"dadInitials,omitempty" bson:"dadInitials,omitempty"`
	LastName    string           `json:"lastName,omitempty" bson:"lastName,omitempty"`
	CNP         string           `json:"cnp,omitempty" bson:"cnp,omitempty"`
	Password    string           `json:"password" bson:"password"`
	Grade       Grade            `json:"grade" bson:"grade"`
	SubjectList []StudentSubject `json:"subjectList" bson:"subjectList"`
}

// we actually define StudentSubject because it doesn't have the grade field.
type StudentSubject struct {
	SubjectID string `json:"subjectID,omitempty" bson:",omitempty"`
	Name      string `json:"name,omitempty" bson:",omitempty"`
}

func GetStudentById(studentID string) (Student, error) {
	// getting a student
	var student Student
  studentsColleciton, err := db.GetCollection("students")
  
  studentsColleciton.FindOne(context.Background(), bson.M{"studentID": studentID}).Decode(&student)

	return student, err
}