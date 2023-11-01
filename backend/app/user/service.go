package user

import "github.com/ip-05/quizzus/entity"

type Repository interface {
	Get(id uint) *entity.User
	GetByGoogle(id string) *entity.User
	GetByDiscord(id string) *entity.User
	GetByTelegram(id string) *entity.User
	Create(e *entity.User) *entity.User
	Update(e *entity.User) *entity.User
	Delete(e *entity.User)
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

	user := s.repo.Create(u)
	return user, nil
}

func (s Service) UpdateUser(id uint, body entity.UpdateUser) (*entity.User, error) {
	user := s.GetUser(id)

	user.Name = body.Name
	user.Picture = body.Picture

	if err := user.Validate(); err != nil {
		return nil, err
	}

	s.repo.Update(user)
	return user, nil
}

func (s Service) DeleteUser(id uint) {
	user := s.GetUser(id)
	s.repo.Delete(user)
}

func (s Service) GetUser(id uint) *entity.User {
	return s.repo.Get(id)
}

func (s Service) GetUserByProvider(id string, provider string) *entity.User {
	switch provider {
	case "google":
		return s.repo.GetByGoogle(id)
	case "discord":
		return s.repo.GetByDiscord(id)
	case "telegram":
		return s.repo.GetByTelegram(id)
	default:
		return nil
	}
}
