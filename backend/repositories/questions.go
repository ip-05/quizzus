package repositories

type Question struct {
	Id      uint     `json:"id" gorm:"primary_key"`
	Name    string   `json:"name"`
	Options []Option `json:"options"`
	GameID  uint     `json:"-"`
}
