package models

type School struct {
	SchoolID  string   `json:"schoolID,omitempty" bson:"schoolID,omitempty"`
	Name      string   `json:"name,omitempty" bson:"name,omitempty"`
	Intervals []string `json:"intervals,omitempty" bson:"intervals,omitempty"`
}