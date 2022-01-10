package models

type Period struct {
	PeriodID string `json:"periodID,omitempty" bson:"periodID,omitempty"`
	Day      int    `json:"day,omitempty" bson:"day,omitempty"`
	Interval string `json:"interval,omitempty" bson:"interval,omitempty"`
	Grade    Grade  `json:"grade,omitempty" bson:"grade,omitempty"`

	// modifiable
	Room     string       `json:"room,omitempty" bson:"room,omitempty"`
	Subject  ShortSubject `json:"subject,omitempty" bson:"subject,omitempty"`
	Teacher  ShortTeacher `json:"teacher,omitempty" bson:"teacher,omitempty"`
	Assigned bool         `json:"assigned,omitempty" bson:"assigned,omitempty"`
	Split    string       `json:"split,omitempty" bson:"split,omitempty"`
}

type ShortTeacher struct {
	TeacherID string `json:"teacherID,omitempty" bson:"teacherID,omitempty"`
	FirstName string `json:"firstName,omitempty" bson:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty" bson:"lastName,omitempty"`
}
