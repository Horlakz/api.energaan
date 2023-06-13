package app

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"

	databaseModule "github.com/horlakz/energaan-api/database"
	productModel "github.com/horlakz/energaan-api/database/model/app"
	"github.com/horlakz/energaan-api/database/repository"
)

type ProductRespositoryInterface interface {
	Create(product productModel.Product) (productModel.Product, error)
	Read(slug string) (productModel.Product, error)
	ReadByUUID(uuid uuid.UUID) (productModel.Product, error)
	ReadAll(pageable repository.Pageable, categoryId uuid.UUID) ([]productModel.Product, repository.Pagination, error)
	Update(product productModel.Product) (productModel.Product, error)
	Delete(slug string) error
}

type productRepository struct {
	database databaseModule.DatabaseInterface
}

func NewProductRepository(database databaseModule.DatabaseInterface) ProductRespositoryInterface {
	return &productRepository{database: database}
}

func (repository *productRepository) Create(product productModel.Product) (productModel.Product, error) {
	product.Model.Prepare()

	err := repository.database.Connection().Create(&product).Error

	if err != nil {
		return product, err
	}

	return product, nil
}

func (repository *productRepository) Read(slug string) (product productModel.Product, err error) {
	err = repository.database.Connection().Model(&productModel.Product{}).Where("slug = ?", slug).First(&product).Error

	if err != nil {
		return product, err
	}

	return product, nil
}

func (repository *productRepository) ReadByUUID(uuid uuid.UUID) (product productModel.Product, err error) {
	err = repository.database.Connection().Model(&productModel.Product{}).Where("uuid = ?", uuid).First(&product).Error

	if err != nil {
		return product, err
	}

	return product, nil
}

func (repository *productRepository) ReadAll(pageable repository.Pageable, categoryId uuid.UUID) (rows []productModel.Product, pagination repository.Pagination, err error) {
	var product productModel.Product
	pagination.TotalPages = 1
	pagination.TotalItems = 0
	pagination.CurrentPage = int64(pageable.Page)

	var result *gorm.DB
	var errCount error

	offset := (pageable.Page - 1) * pageable.Size
	searchQuery := repository.database.Connection().Model(&productModel.Product{})

	if len(strings.TrimSpace(pageable.Search)) > 0 {
		searchQuery.Where("title LIKE ?", "%"+strings.ToLower(pageable.Search)+"%")
	}

	if categoryId != uuid.Nil {
		searchQuery.Where("category_id = ?", categoryId)
	}

	errCount = searchQuery.Count(&pagination.TotalItems).Error
	paginationQuery := searchQuery.Limit(pageable.Size).Offset(offset).Order(pageable.SortBy + " " + pageable.SortDirection)

	result = paginationQuery.Model(&productModel.Product{}).Where(product).Find(&rows)

	if result.Error != nil {
		return nil, pagination, result.Error
	}

	if errCount != nil {
		return nil, pagination, errCount
	}

	pagination.TotalPages = pagination.TotalItems / int64(pageable.Size)

	return rows, pagination, nil
}

func (repository *productRepository) Update(product productModel.Product) (productModel.Product, error) {
	var checkRow productModel.Product

	err := repository.database.Connection().Model(&productModel.Product{}).Where("uuid = ?", product.UUID.String()).First(&checkRow).Error

	if err != nil {
		return checkRow, err
	}

	err = repository.database.Connection().Model(&checkRow).Updates(product).Error

	if err != nil {
		return product, err
	}

	return product, nil
}

func (repository *productRepository) Delete(slug string) (err error) {
	var product productModel.Product

	err = repository.database.Connection().Model(&productModel.Product{}).Where("slug = ?", slug).First(&product).Error

	if err != nil {
		return err
	}

	err = repository.database.Connection().Delete(&product).Error

	fmt.Println(err, " this is the error")

	if err != nil {
		return err
	}

	return nil
}
