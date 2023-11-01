package user

import "github.com/ip-05/quizzus/entity"

type Repository interface {
	GetUserById(ID uint) *entity.User
	GetUserByGoogleId(ID string) *entity.User
	GetUserByDiscordId(ID string) *entity.User
	GetUserByTelegramId(ID string) *entity.User
	CreateUser(e *entity.User) *entity.User
	UpdateUser(e *entity.User) *entity.User
	DeleteUser(e *entity.User)
}

type Service struct {
	repo Repository
}

func NewService(userRepo Repository) *Service {
	return &Service{
		repo: userRepo,
	}
}

func (s Service) CreateUser(body *entity.CreateUser) (*entity.User, error) {
	u, err := entity.NewUser(body)
	if err != nil {
		return nil, err
	}

	user := s.repo.CreateUser(u)
	return user, nil
}

func (s Service) UpdateUser(ID uint, body entity.UpdateUser) (*entity.User, error) {
	user := s.GetUserById(ID)

	user.Name = body.Name
	user.Picture = body.Picture

	if err := user.Validate(); err != nil {
		return nil, err
	}

	s.repo.UpdateUser(user)
	return user, nil
}

func (s Service) DeleteUser(ID uint) {
	user := s.GetUserById(ID)
	s.repo.DeleteUser(user)
}

func (s Service) GetUserById(ID uint) *entity.User {
	return s.repo.GetUserById(ID)
}

func (s Service) GetUserByProvider(ID string, provider string) *entity.User {
	switch provider {
	case "google":
		return s.repo.GetUserByGoogleId(ID)
	case "discord":
		return s.repo.GetUserByDiscordId(ID)
	case "telegram":
		return s.repo.GetUserByTelegramId(ID)
	default:
		return nil
	}
}
