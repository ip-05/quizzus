package models

type Option struct {
	Id         uint   `json:"id" gorm:"primary_key"`
	Name       string `json:"name"`
	Correct    bool   `json:"correct"`
	QuestionId uint
}
