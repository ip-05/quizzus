package entity

type Option struct {
	Id         uint   `json:"id" gorm:"primary_key"`
	Name       string `json:"name"`
	Correct    bool   `json:"correct"`
	QuestionID uint   `json:"-"`
}

type CreateOption struct {
	Name    string `json:"name"`
	Correct bool   `json:"correct"`
}

type UpdateOption struct {
	Id      uint   `json:"id"`
	Name    string `json:"name"`
	Correct bool   `json:"correct"`
}
