package models

type Period struct {
	Interval  string  `json:"interval,omitempty" bson:"interval,omitempty"`
	Day       int     `json:"day,omitempty" bson:"day,omitempty"`
	Room      string  `json:"room,omitempty" bson:"room,omitempty"`
	Subject   Subject `json:"subject,omitempty" bson:"subject,omitempty"`
	TeacherID string  `json:"teacherID,omitempty" bson:"teacherID,omitempty"`
}