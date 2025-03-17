package service

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"task-manager/dto/request"
	"task-manager/dto/response"
	"task-manager/model"
	"task-manager/repository"
	"task-manager/util"
)

type AuthService interface {
	Register(req request.RegisterRequest) (*response.UserResponse, error)
	Login(req request.LoginRequest) (*response.TokenResponse, error)
	GetUserByID(id uint) (*response.UserResponse, error)
}

type authService struct {
	userRepo repository.UserRepository
	db       *gorm.DB
}

func NewAuthService(userRepo repository.UserRepository, db *gorm.DB) AuthService {
	return &authService{
		userRepo: userRepo,
		db:       db,
	}
}

func (s *authService) Register(req request.RegisterRequest) (*response.UserResponse, error) {
	existingUser, err := s.userRepo.FindByUsername(req.Username)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("username already taken")
	}

	existingUser, err = s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("email already registered")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		userRepoTx := repository.NewUserRepository(tx)
		if err := userRepoTx.Create(user); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &response.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}

func (s *authService) Login(req request.LoginRequest) (*response.TokenResponse, error) {
	user, err := s.userRepo.FindByUsername(req.Username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid username or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	token, expiration, err := util.GenerateJWTToken(user.ID, user.Username)
	if err != nil {
		return nil, err
	}

	return &response.TokenResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   expiration,
		UserID:      user.ID,
	}, nil
}

func (s *authService) GetUserByID(id uint) (*response.UserResponse, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	return &response.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}