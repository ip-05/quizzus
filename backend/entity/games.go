package entity

type Game struct {
	Id         uint       `json:"id" gorm:"primary_key"`
	InviteCode string     `json:"inviteCode"`
	Topic      string     `json:"topic"`
	RoundTime  int        `json:"roundTime"`
	Points     float64    `json:"points"`
	Questions  []Question `json:"questions"`
	Owner      string     `json:"ownerId"`
}

type CreateBody struct {
	Topic     string           `json:"topic"`
	RoundTime int              `json:"roundTime"`
	Points    float64          `json:"points"`
	Questions []CreateQuestion `json:"questions"`
}

type UpdateBody struct {
	Topic     string           `json:"topic"`
	RoundTime int              `json:"roundTime"`
	Points    float64          `json:"points"`
	Questions []UpdateQuestion `json:"questions"`
}
