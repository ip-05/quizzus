package models

type Game struct {
	Id         uint       `json:"id" gorm:"primary_key"`
	InviteCode string     `json:"inviteCode"`
	Topic      string     `json:"topic"`
	RoundTime  int        `json:"roundTime"`
	Points     float64    `json:"points"`
	Questions  []Question `json:"questions"`
}
