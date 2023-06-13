package user

import "github.com/ip-05/quizzus/entity"

type IUserRepo interface {
	Get(id uint) *entity.User
	GetByGoogle(id string) *entity.User
	GetByDiscord(id string) *entity.User
	GetByTelegram(id string) *entity.User
	Create(e *entity.User) *entity.User
	Update(e *entity.User) *entity.User
	Delete(e *entity.User)
}

type UserService struct {
	userRepo IUserRepo
}

func NewUserService(userR IUserRepo) *UserService {
	return &UserService{
		userRepo: userR,
	}
}

func (us *UserService) CreateUser(body *entity.CreateUser) (*entity.User, error) {
	u, err := entity.NewUser(body)
	if err != nil {
		return nil, err
	}

	user := us.userRepo.Create(u)
	return user, nil
}

func (us *UserService) UpdateUser(id uint, body *entity.UpdateUser) (*entity.User, error) {
	user := us.userRepo.Get(id)

	user.Name = body.Name
	user.Picture = body.Picture

	if err := user.Validate(); err != nil {
		return nil, err
	}

	us.userRepo.Update(user)
	return user, nil
}

func (us *UserService) DeleteUser(user *entity.User) {
	us.userRepo.Delete(user)
}

func (us *UserService) GetUser(id uint) *entity.User {
	return us.userRepo.Get(id)
}

func (us *UserService) GetUserByProvider(id string, provider string) *entity.User {
	switch provider {
	case "google":
		return us.userRepo.GetByGoogle(id)
	case "discord":
		return us.userRepo.GetByDiscord(id)
	case "telegram":
		return us.userRepo.GetByTelegram(id)
	default:
		return nil
	}
}
