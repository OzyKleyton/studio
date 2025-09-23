package service

import (
	"github.com/OzyKleyton/studio-api/internal/model"
	"github.com/OzyKleyton/studio-api/internal/model/user"
	"github.com/OzyKleyton/studio-api/internal/repository"
	"github.com/OzyKleyton/studio-api/utils/auth"
	"github.com/OzyKleyton/studio-api/utils/security"
)

type UserService interface {
	CreateUser(userReq *user.UserReq) *model.Response
	FindAllUsers() *model.Response
	FindUserByEmail(email string) *model.Response
	UpdateUser(id uint, userReq *user.UserReq) *model.Response
	DeleteUser(id uint) *model.Response
	Login(userReq user.Login) *model.Response
}

type UserServiceImpl struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &UserServiceImpl{
		repo: repo,
	}
}

func (us *UserServiceImpl) CreateUser(userReq *user.UserReq) *model.Response {
	user := userReq.ToUser()

	if user.Email == "" || user.Username == "" || user.Password == "" {
		return &model.Response{
			Status:  400,
			Message: "Email or Username or Password cannot be empty",
			Data:    nil,
		}
	}

	hash, err := security.EncodePassword(userReq.Password)
	if err != nil {
		return model.NewErrorResponse(err)
	}

	user.Password = string(hash)

	createUser, err := us.repo.Create(user)
	if err != nil {
		return model.NewErrorResponse(err, 500)
	}

	return model.NewSuccessResponse(createUser.ToUserRes())
}

func (us *UserServiceImpl) FindAllUsers() *model.Response {
	users, err := us.repo.FindAll()
	if err != nil {
		return model.NewErrorResponse(err, 404)
	}

	usersResponse := []*user.UserRes{}
	for _, u := range users {
		usersResponse = append(usersResponse, u.ToUserRes())
	}

	return model.NewSuccessResponse(usersResponse)
}

func (us *UserServiceImpl) FindUserByEmail(email string) *model.Response {
	user, err := us.repo.FindByEmail(email)
	if err != nil {
		return model.NewErrorResponse(err, 404)
	}

	return model.NewSuccessResponse(user.ToUserRes())
}

func (us *UserServiceImpl) UpdateUser(id uint, userReq *user.UserReq) *model.Response {
	user, err := us.repo.FindByID(id)
	if err != nil {
		return model.NewErrorResponse(err, 404)
	}

	user.Username = userReq.Username
	user.Email = userReq.Email

	updateUser, err := us.repo.Update(user)
	if err != nil {
		return model.NewErrorResponse(err, 500)
	}

	return model.NewSuccessResponse(updateUser.ToUserRes())
}

func (us *UserServiceImpl) DeleteUser(id uint) *model.Response {
	userID, err := us.repo.FindByID(id)
	if err != nil {
		return model.NewErrorResponse(err, 404)
	}
	_, err = us.repo.Delete(userID.ID)
	if err != nil {
		return model.NewErrorResponse(err, 500)
	}

	return model.NewSuccessResponse(nil)
}

func (us *UserServiceImpl) Login(userReq user.Login) *model.Response {

	if userReq.Email == "" || userReq.Password == "" {
		return &model.Response{
			Status:  400,
			Message: "Email and Password cannot be empty",
			Data:    nil,
		}
	}

	user, err := us.repo.FindByEmail(userReq.Email)
	if err != nil {
		return model.NewErrorResponse(err)
	}

	if hash := security.CompareHashPassword(user.Password, userReq.Password); !hash {
		return &model.Response{
			Status:  400,
			Message: "Password incorrect",
			Data:    nil,
		}
	}

	token, err := auth.GenerateToken(*user)
	if err != nil {
		return model.NewErrorResponse(err)
	}

	loginResponse := map[string]any{
		"token": token,
		"user": map[string]any{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	}

	return model.NewSuccessResponse(loginResponse)
}
