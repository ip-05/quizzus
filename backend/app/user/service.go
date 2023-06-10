package user

import "github.com/ip-05/quizzus/entity"

type IUserRepo interface {
	Get(id int) *entity.User
	Create(e *entity.User) *entity.User
	Update(id int, e *entity.User) *entity.User
	Delete(e *entity.User)
}

type UserService struct {
	userRepo IUserRepo
}

func NewGameService(userR IUserRepo) *UserService {
	return &UserService{
		userRepo: userR,
	}
}

func (us *UserService) CreateUser() {

}

func (us *UserService) UpdateUser() {

}

func (us *UserService) DeleteUser() {

}

func (us *UserService) GetUser() {

}
