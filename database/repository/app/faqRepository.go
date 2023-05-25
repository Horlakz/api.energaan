package app

import (
	"gorm.io/gorm"

	databaseModule "github.com/horlakz/energaan-api/database"
	faqModel "github.com/horlakz/energaan-api/database/model/app"
)

type FaqRespositoryInterface interface {
	Create(faq faqModel.Faq) (faqModel.Faq, error)
	Read(slug string) (faqModel.Faq, error)
	ReadAll() ([]faqModel.Faq, error)
	Update(faq faqModel.Faq) (faqModel.Faq, error)
	Delete(slug string) error
}

type FaqRepository struct {
	database databaseModule.DatabaseInterface
}

func NewFaqRepository(database databaseModule.DatabaseInterface) FaqRespositoryInterface {
	return &FaqRepository{database: database}
}

func (repository *FaqRepository) Create(faq faqModel.Faq) (faqModel.Faq, error) {
	faq.Model.Prepare()

	err := repository.database.Connection().Create(&faq).Error

	if err != nil {
		return faq, err
	}

	return faq, nil
}

func (repository *FaqRepository) Read(slug string) (faq faqModel.Faq, err error) {
	err = repository.database.Connection().Model(&faqModel.Faq{}).Where("slug = ?", slug).First(&faq).Error

	if err != nil {
		return faq, err
	}

	return faq, nil
}

func (repository *FaqRepository) ReadAll() (rows []faqModel.Faq, err error) {
	var faq faqModel.Faq

	var result *gorm.DB
	var errCount error

	result = repository.database.Connection().Model(&faqModel.Faq{}).Where(faq).Find(&rows)

	if result.Error != nil {
		return nil, result.Error
	}

	if errCount != nil {
		return nil, errCount
	}

	return rows, nil
}

func (repository *FaqRepository) Update(faq faqModel.Faq) (faqModel.Faq, error) {
	var checkRow faqModel.Faq

	err := repository.database.Connection().Model(&faqModel.Faq{}).Where("uuid = ?", faq.UUID.String()).First(&checkRow).Error

	if err != nil {
		return checkRow, err
	}

	err = repository.database.Connection().Model(&checkRow).Updates(faq).Error

	if err != nil {
		return faq, err
	}

	return faq, nil
}

func (repository *FaqRepository) Delete(slug string) (err error) {
	var faq faqModel.Faq

	err = repository.database.Connection().Model(&faqModel.Faq{}).Where("slug = ?", slug).First(&faq).Error

	if err != nil {
		return err
	}

	err = repository.database.Connection().Delete(&faq).Error

	if err != nil {
		return err
	}

	return nil
}
