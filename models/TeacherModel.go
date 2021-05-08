package models

import (
	"backend-go/db"
	"context"

	"gopkg.in/mgo.v2/bson"
)

// the password and the subjectList can be empty
type Teacher struct {
	TeacherID     string    `json:"teacherID,omitempty" bson:"teacherID,omitempty"`
	FirstName     string    `json:"firstName,omitempty" bson:"firstName,omitempty"`
	LastName      string    `json:"lastName,omitempty" bson:"lastName,omitempty"`
	CNP           string    `json:"cnp,omitempty" bson:"cnp,omitempty"`
	Password      string    `json:"password" bson:"password"`
	HomeroomGrade Grade     `json:"homeroomGrade" bson:"homeroomGrade"`
	SubjectList   []Subject `json:"subjectList" bson:"subjectList"`
}

func GetTeacherById(teacherID string) (Teacher, error) {
	// getting a teacher
	var teacher Teacher
  teachersColleciton, err := db.GetCollection("teachers")
  
  teachersColleciton.FindOne(context.Background(), bson.M{"teacherID": teacherID}).Decode(&teacher)

	return teacher, err
}
