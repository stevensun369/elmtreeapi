package models

type Period struct {
	PeriodID string `json:"periodID,omitempty" bson:"periodID,omitempty"`
	Day      int    `json:"day,omitempty" bson:"day,omitempty"`
	Interval int    `json:"interval,omitempty" bson:"interval,omitempty"`
	Grade    Grade  `json:"grade,omitempty" bson:"grade,omitempty"`

	// modifiable
	Room     string       `json:"room,omitempty" bson:"room,omitempty"`
	Subject  ShortSubject `json:"subject,omitempty" bson:"subject,omitempty"`
	Assigned bool         `json:"assigned,omitempty" bson:"assigned,omitempty"`
}
