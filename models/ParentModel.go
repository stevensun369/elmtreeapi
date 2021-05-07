package models

type Parent struct {
	ParentID      string   `json:"parentID,omitempty"`
	StudentIDList []string `json:"studentIDList"`
	FirstName     string   `json:"firstName,omitempty"`
	LastName      string   `json:"lastName,omitempty"`
	CNP           string   `json:"cnp,omitempty"`
	Password      string   `json:"password,omitempty"`
}