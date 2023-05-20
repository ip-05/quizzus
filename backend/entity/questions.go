package entity

type Question struct {
	Id      uint     `json:"id" gorm:"primary_key"`
	Name    string   `json:"name"`
	Options []Option `json:"options"`
	GameID  uint     `json:"-"`
}

type CreateQuestion struct {
	Name    string         `json:"name"`
	Options []CreateOption `json:"options"`
}

type UpdateQuestion struct {
	Id      uint           `json:"id"`
	Name    string         `json:"name"`
	Options []UpdateOption `json:"options"`
}
