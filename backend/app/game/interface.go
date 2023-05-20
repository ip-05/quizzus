package game

type GameRepo interface {
	Get()
	Create()
	Update()
	Delete()
}

type QuestionRepo interface {
	Get()
	Create()
	Update()
	Delete()
}

type OptionRepo interface {
	Get()
	Create()
	Update()
	Delete()
}

type IService interface {
	//methods of service
	CreateGame()
	UpdateGame()
	DeleteGame()
	GetGame()
}
