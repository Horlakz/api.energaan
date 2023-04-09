package auth

import (
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"

	databaseModule "github.com/horlakz/energaan-api/database"
	authModel "github.com/horlakz/energaan-api/database/model/auth"
	"github.com/horlakz/energaan-api/database/repository"
)

type UserRepositoryInterface interface {
	Create(user authModel.User) (authModel.User, error)
	Read(uid uuid.UUID) (authModel.User, error)
	ReadAll(pageable repository.Pageable) ([]authModel.User, repository.Pagination, error)
	Update(user authModel.User) (authModel.User, error)
	Delete(uid uuid.UUID) error
	FindByEmail(email string) (authModel.User, error)
}

type userRepository struct {
	database databaseModule.DatabaseInterface
}

func NewUserRepository(database databaseModule.DatabaseInterface) UserRepositoryInterface {
	return &userRepository{database: database}
}

func (repository *userRepository) Create(user authModel.User) (authModel.User, error) {
	user.Model.Prepare()

	err := repository.database.Connection().Create(&user).Error

	if err != nil {
		return user, err
	}

	return user, nil
}

func (repository *userRepository) Read(uid uuid.UUID) (user authModel.User, err error) {
	err = repository.database.Connection().Model(&authModel.User{}).Where("uuid = ?", uid).First(&user).Error

	if err != nil {
		return user, err
	}

	return user, nil
}

func (repository *userRepository) ReadAll(pageable repository.Pageable) (rows []authModel.User, pagination repository.Pagination, err error) {
	var user authModel.User
	pagination.TotalPages = 1
	pagination.TotalItems = 0
	pagination.CurrentPage = int64(pageable.Page)

	var result *gorm.DB
	var errCount error

	offset := (pageable.Page - 1) * pageable.Size
	searchQuery := repository.database.Connection().Model(&authModel.User{})

	if len(strings.TrimSpace(pageable.Search)) > 0 {
		searchQuery.Where("email LIKE ?", "%"+strings.ToLower(pageable.Search)+"%")
	}

	errCount = searchQuery.Count(&pagination.TotalItems).Error
	paginationQuery := searchQuery.Limit(pageable.Size).Offset(offset).Order(pageable.SortBy + " " + pageable.SortDirection)

	result = paginationQuery.Model(&authModel.User{}).Where(user).Find(&rows)

	if result.Error != nil {
		msg := result.Error
		return nil, pagination, msg
	}

	if errCount != nil {
		return nil, pagination, errCount
	}

	pagination.TotalPages = pagination.TotalItems / int64(pageable.Size)

	return rows, pagination, nil
}

func (repository *userRepository) Update(user authModel.User) (authModel.User, error) {
	var checkRow authModel.User

	err := repository.database.Connection().Model(&authModel.User{}).Where("uuid = ? ", user.UUID.String()).First(&checkRow).Error

	if err != nil {
		return checkRow, err
	}

	err = repository.database.Connection().Model(&checkRow).Updates(user).Error

	if err != nil {
		return user, err
	}

	return user, nil
}

func (repository *userRepository) Delete(uuid uuid.UUID) (err error) {
	var user authModel.User
	err = repository.database.Connection().Model(&authModel.User{}).Where("uuid = ? ", uuid).First(&user).Error

	if err != nil {
		return err
	}

	err = repository.database.Connection().Delete(user).Error

	if err != nil {
		return err
	}

	return nil
}

func (repository *userRepository) FindByEmail(email string) (row authModel.User, err error) {
	err = repository.database.Connection().Model(&authModel.User{}).Where("email = ? ", email).First(&row).Error

	if err != nil {
		return row, err
	}

	return row, nil
}
