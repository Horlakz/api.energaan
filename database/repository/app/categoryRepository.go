package app

import (
	"gorm.io/gorm"

	databaseModule "github.com/horlakz/energaan-api/database"
	categoryModel "github.com/horlakz/energaan-api/database/model/app"
)

type CategoryRespositoryInterface interface {
	Create(category categoryModel.Category) (categoryModel.Category, error)
	Read(slug string) (categoryModel.Category, error)
	ReadAll() ([]categoryModel.Category, error)
	Update(category categoryModel.Category) (categoryModel.Category, error)
	Delete(slug string) error
}

type categoryRepository struct {
	database databaseModule.DatabaseInterface
}

func NewCategoryRepository(database databaseModule.DatabaseInterface) CategoryRespositoryInterface {
	return &categoryRepository{database: database}
}

func (repository *categoryRepository) Create(category categoryModel.Category) (categoryModel.Category, error) {
	category.Model.Prepare()

	err := repository.database.Connection().Create(&category).Error

	if err != nil {
		return category, err
	}

	return category, nil
}

func (repository *categoryRepository) Read(slug string) (category categoryModel.Category, err error) {
	err = repository.database.Connection().Model(&categoryModel.Category{}).Where("slug = ?", slug).First(&category).Error

	if err != nil {
		return category, err
	}

	return category, nil
}

func (repository *categoryRepository) ReadAll() (rows []categoryModel.Category, err error) {
	var category categoryModel.Category

	var result *gorm.DB
	var errCount error

	result = repository.database.Connection().Model(&categoryModel.Category{}).Where(category).Find(&rows)

	if result.Error != nil {
		return nil, result.Error
	}

	if errCount != nil {
		return nil, errCount
	}

	return rows, nil
}

func (repository *categoryRepository) Update(category categoryModel.Category) (categoryModel.Category, error) {
	var checkRow categoryModel.Category

	err := repository.database.Connection().Model(&categoryModel.Category{}).Where("uuid = ?", category.UUID.String()).First(&checkRow).Error

	if err != nil {
		return checkRow, err
	}

	err = repository.database.Connection().Model(&checkRow).Updates(category).Error

	if err != nil {
		return category, err
	}

	return category, nil
}

func (repository *categoryRepository) Delete(slug string) (err error) {
	var category categoryModel.Category

	err = repository.database.Connection().Model(&categoryModel.Category{}).Where("slug = ?", slug).First(&category).Error

	if err != nil {
		return err
	}

	err = repository.database.Connection().Delete(&category).Error

	if err != nil {
		return err
	}

	return nil
}
