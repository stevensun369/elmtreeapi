package models

type Parent struct {
	ParentID      string   `json:"parentID,omitempty" bson:"parentID,omitempty"`
	StudentIDList []string `json:"studentIDList" bson:"studentIDList"`
	FirstName     string   `json:"firstName,omitempty" bson:"firstName,omitempty"`
	LastName      string   `json:"lastName,omitempty" bson:"lastName,omitempty"`
	CNP           string   `json:"cnp,omitempty" bson:"cnp,omitempty"`
	Password      string   `json:"password,omitempty" bson:"password,omitempty"`
}