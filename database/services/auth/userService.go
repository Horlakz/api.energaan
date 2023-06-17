package service

import (
	"errors"

	"github.com/google/uuid"

	authDtos "github.com/horlakz/energaan-api/database/dto/auth"
	authModels "github.com/horlakz/energaan-api/database/model/auth"
	authRepos "github.com/horlakz/energaan-api/database/repository/auth"
	"github.com/horlakz/energaan-api/helper"
)

type userService struct {
	userRepository authRepos.UserRepositoryInterface
	encrypt        helper.Argon2
}

type UserServiceInterface interface {
	Create(userDTO authDtos.UserDTO) (authDtos.UserDTO, error)
	Read(uid uuid.UUID) (authDtos.UserDTO, error)
	// ReadAll(pageable repository.Pageable) ([]authDtos.UserDTO, repository.Pagination, error)
	// Update(userDTO authDtos.UserDTO) (authDtos.UserDTO, error)
	// Delete(userUUID uuid.UUID, uid uuid.UUID) error
	FindByEmail(email string) (authDtos.UserDTO, error)
	Authenticate(email string, password string) (*authDtos.UserDTO, error)
}

func NewUserService(userRepository authRepos.UserRepositoryInterface) UserServiceInterface {
	return &userService{
		userRepository: userRepository,
	}
}

func (service *userService) ConvertToDTO(user authModels.User) (userDTO authDtos.UserDTO) {
	userDTO.UUID = user.UUID
	userDTO.CreatedAt = user.CreatedAt
	userDTO.UpdatedAt = user.UpdatedAt
	userDTO.DeletedAt = user.DeletedAt.Time
	userDTO.Email = user.Email
	userDTO.FullName = user.FullName

	return userDTO
}

func (service *userService) ConvertToModel(userDTO authDtos.UserDTO) (user authModels.User) {
	user.CreatedAt = userDTO.CreatedAt
	user.UpdatedAt = userDTO.UpdatedAt
	user.DeletedAt.Time = userDTO.DeletedAt
	user.Email = userDTO.Email
	user.FullName = userDTO.FullName
	user.Password = userDTO.Password

	return user
}

func (service *userService) Create(userDTO authDtos.UserDTO) (authDtos.UserDTO, error) {
	user := service.ConvertToModel(userDTO)
	newRecord, err := service.userRepository.Create(user)

	return service.ConvertToDTO(newRecord), err
}

func (service *userService) Read(uid uuid.UUID) (authDtos.UserDTO, error) {
	record, err := service.userRepository.Read(uid)

	return service.ConvertToDTO(record), err
}

// func (service *userService) ReadAll(pageable repository.Pageable) (recordsDto []authDtos.UserDTO, pagination repository.Pagination, err error) {
// 	records, pagination, err := service.userRepository.ReadAll(pageable)

// 	for _, record := range records {
// 		recordsDto = append(recordsDto, service.ConvertToDTO(record))
// 	}

// 	return recordsDto, pagination, err
// }

// func (service *userService) Update(userDTO authDtos.UserDTO) (authDtos.UserDTO, error) {
// 	user := service.ConvertToModel(userDTO)
// 	newRecord, err := service.userRepository.Update(user)

// 	return service.ConvertToDTO(newRecord), err
// }

// func (service *userService) Delete(userUUID uuid.UUID, uid uuid.UUID) (err error) {
// 	rtn := service.userRepository.Delete(uid)
// 	record, _ := service.userRepository.Read(uid)

// 	service.userRepository.Update(record)

// 	return rtn
// }

func (service *userService) FindByEmail(email string) (authDtos.UserDTO, error) {
	record, err := service.userRepository.FindByEmail(email)

	if err != nil {
		err = errors.New("email address not found")
	}

	return service.ConvertToDTO(record), err
}

func (service *userService) Authenticate(email string, password string) (*authDtos.UserDTO, error) {
	user, err := service.userRepository.FindByEmail(email)

	if err != nil {
		return nil, errors.New("user not found")
	}

	passwordValid, _ := service.encrypt.ComparePassword(password, user.Password)

	if !passwordValid {
		return nil, errors.New("password does not match")
	}

	userDTO := service.ConvertToDTO(user)

	return &userDTO, err
}
