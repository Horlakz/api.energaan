package app

import (
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"

	databaseModule "github.com/horlakz/energaan-api/database"
	contactModel "github.com/horlakz/energaan-api/database/model/app"
	"github.com/horlakz/energaan-api/database/repository"
)

type ContactRespositoryInterface interface {
	Create(contact contactModel.Contact) (contactModel.Contact, error)
	Read(email string) (contactModel.Contact, error)
	ReadAll(pageable repository.Pageable) ([]contactModel.Contact, repository.Pagination, error)
	Update(contact contactModel.Contact) (contactModel.Contact, error)
	Delete(id uuid.UUID) error
}

type ContactRepository struct {
	database databaseModule.DatabaseInterface
}

func NewContactRepository(database databaseModule.DatabaseInterface) ContactRespositoryInterface {
	return &ContactRepository{database: database}
}

func (repository *ContactRepository) Create(contact contactModel.Contact) (contactModel.Contact, error) {
	contact.Model.Prepare()

	err := repository.database.Connection().Create(&contact).Error

	if err != nil {
		return contact, err
	}

	return contact, nil
}

func (repository *ContactRepository) Read(email string) (contact contactModel.Contact, err error) {
	err = repository.database.Connection().Model(&contactModel.Contact{}).Where("email = ?", email).First(&contact).Error

	if err != nil {
		return contact, err
	}

	return contact, nil
}

func (repository *ContactRepository) ReadAll(pageable repository.Pageable) (rows []contactModel.Contact, pagination repository.Pagination, err error) {
	var contact contactModel.Contact
	pagination.TotalPages = 1
	pagination.TotalItems = 0
	pagination.CurrentPage = int64(pageable.Page)

	var result *gorm.DB
	var errCount error

	offset := (pageable.Page - 1) * pageable.Size
	searchQuery := repository.database.Connection().Model(&contactModel.Contact{})

	if len(strings.TrimSpace(pageable.Search)) > 0 {
		searchQuery.Where("service_type LIKE ?", "%"+strings.ToLower(pageable.Search)+"%")
	}

	errCount = searchQuery.Count(&pagination.TotalItems).Error
	paginationQuery := searchQuery.Limit(pageable.Size).Offset(offset).Order(pageable.SortBy + " " + pageable.SortDirection)

	result = paginationQuery.Model(&contactModel.Contact{}).Where(contact).Find(&rows)

	if result.Error != nil {
		return nil, pagination, result.Error
	}

	if errCount != nil {
		return nil, pagination, errCount
	}

	pagination.TotalPages = pagination.TotalItems / int64(pageable.Size)

	return rows, pagination, nil
}

func (repository *ContactRepository) Update(contact contactModel.Contact) (contactModel.Contact, error) {
	var checkRow contactModel.Contact

	err := repository.database.Connection().Model(&contactModel.Contact{}).Where("uuid = ?", contact.UUID.String()).First(&checkRow).Error

	if err != nil {
		return checkRow, err
	}

	err = repository.database.Connection().Model(&checkRow).Updates(contact).Error

	if err != nil {
		return contact, err
	}

	return contact, nil
}

func (repository *ContactRepository) Delete(id uuid.UUID) (err error) {
	var contact contactModel.Contact

	err = repository.database.Connection().Model(&contactModel.Contact{}).Where("uuid = ?", id).First(&contact).Error

	if err != nil {
		return err
	}

	err = repository.database.Connection().Delete(&contact).Error

	if err != nil {
		return err
	}

	return nil
}
