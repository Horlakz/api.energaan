package app

import (
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"

	databaseModule "github.com/horlakz/energaan-api/database"
	quoteModel "github.com/horlakz/energaan-api/database/model/app"
	"github.com/horlakz/energaan-api/database/repository"
)

type QuoteRespositoryInterface interface {
	Create(quote quoteModel.Quote) (quoteModel.Quote, error)
	Read(email string) (quoteModel.Quote, error)
	CheckEmailServiceTypeExists(email string, serviceType string) (quoteModel.Quote, error)
	ReadAll(pageable repository.Pageable) ([]quoteModel.Quote, repository.Pagination, error)
	Update(quote quoteModel.Quote) (quoteModel.Quote, error)
	Delete(id uuid.UUID) error
}

type QuoteRepository struct {
	database databaseModule.DatabaseInterface
}

func NewQuoteRepository(database databaseModule.DatabaseInterface) QuoteRespositoryInterface {
	return &QuoteRepository{database: database}
}

func (repository *QuoteRepository) Create(quote quoteModel.Quote) (quoteModel.Quote, error) {
	quote.Model.Prepare()

	err := repository.database.Connection().Create(&quote).Error

	if err != nil {
		return quote, err
	}

	return quote, nil
}

func (repository *QuoteRepository) Read(email string) (quote quoteModel.Quote, err error) {
	err = repository.database.Connection().Model(&quoteModel.Quote{}).Where("email = ?", email).First(&quote).Error

	if err != nil {
		return quote, err
	}

	return quote, nil
}

func (repository *QuoteRepository) CheckEmailServiceTypeExists(email string, serviceType string) (quote quoteModel.Quote, err error) {
	err = repository.database.Connection().Model(&quoteModel.Quote{}).Where("email = ? AND service_type = ?", email, serviceType).First(&quote).Error

	if err != nil {
		return quote, err
	}

	return quote, nil
}

func (repository *QuoteRepository) ReadAll(pageable repository.Pageable) (rows []quoteModel.Quote, pagination repository.Pagination, err error) {
	var quote quoteModel.Quote
	pagination.TotalPages = 1
	pagination.TotalItems = 0
	pagination.CurrentPage = int64(pageable.Page)

	var result *gorm.DB
	var errCount error

	offset := (pageable.Page - 1) * pageable.Size
	searchQuery := repository.database.Connection().Model(&quoteModel.Quote{})

	if len(strings.TrimSpace(pageable.Search)) > 0 {
		searchQuery.Where("service_type LIKE ?", "%"+strings.ToLower(pageable.Search)+"%")
	}

	errCount = searchQuery.Count(&pagination.TotalItems).Error
	paginationQuery := searchQuery.Limit(pageable.Size).Offset(offset).Order(pageable.SortBy + " " + pageable.SortDirection)

	result = paginationQuery.Model(&quoteModel.Quote{}).Where(quote).Find(&rows)

	if result.Error != nil {
		return nil, pagination, result.Error
	}

	if errCount != nil {
		return nil, pagination, errCount
	}

	pagination.TotalPages = pagination.TotalItems / int64(pageable.Size)

	return rows, pagination, nil
}

func (repository *QuoteRepository) Update(quote quoteModel.Quote) (quoteModel.Quote, error) {
	var checkRow quoteModel.Quote

	err := repository.database.Connection().Model(&quoteModel.Quote{}).Where("uuid = ?", quote.UUID.String()).First(&checkRow).Error

	if err != nil {
		return checkRow, err
	}

	err = repository.database.Connection().Model(&checkRow).Updates(quote).Error

	if err != nil {
		return quote, err
	}

	return quote, nil
}

func (repository *QuoteRepository) Delete(id uuid.UUID) (err error) {
	var quote quoteModel.Quote

	err = repository.database.Connection().Model(&quoteModel.Quote{}).Where("uuid = ?", id).First(&quote).Error

	if err != nil {
		return err
	}

	err = repository.database.Connection().Delete(&quote).Error

	if err != nil {
		return err
	}

	return nil
}
