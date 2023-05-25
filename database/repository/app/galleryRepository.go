package app

import (
	"gorm.io/gorm"

	databaseModule "github.com/horlakz/energaan-api/database"
	galleryModel "github.com/horlakz/energaan-api/database/model/app"
)

type GalleryRespositoryInterface interface {
	Create(gallery galleryModel.Gallery) (galleryModel.Gallery, error)
	Read(slug string) (galleryModel.Gallery, error)
	ReadAll() ([]galleryModel.Gallery, error)
	Update(gallery galleryModel.Gallery) (galleryModel.Gallery, error)
	Delete(slug string) error
}

type GalleryRepository struct {
	database databaseModule.DatabaseInterface
}

func NewgalleryRepository(database databaseModule.DatabaseInterface) GalleryRespositoryInterface {
	return &GalleryRepository{database: database}
}

func (repository *GalleryRepository) Create(gallery galleryModel.Gallery) (galleryModel.Gallery, error) {
	gallery.Model.Prepare()

	err := repository.database.Connection().Create(&gallery).Error

	if err != nil {
		return gallery, err
	}

	return gallery, nil
}

func (repository *GalleryRepository) Read(slug string) (gallery galleryModel.Gallery, err error) {
	err = repository.database.Connection().Model(&galleryModel.Gallery{}).Where("slug = ?", slug).First(&gallery).Error

	if err != nil {
		return gallery, err
	}

	return gallery, nil
}

func (repository *GalleryRepository) ReadAll() (rows []galleryModel.Gallery, err error) {
	var gallery galleryModel.Gallery

	var result *gorm.DB
	var errCount error

	result = repository.database.Connection().Model(&galleryModel.Gallery{}).Where(gallery).Find(&rows)

	if result.Error != nil {
		return nil, result.Error
	}

	if errCount != nil {
		return nil, errCount
	}

	return rows, nil
}

func (repository *GalleryRepository) Update(gallery galleryModel.Gallery) (galleryModel.Gallery, error) {
	var checkRow galleryModel.Gallery

	err := repository.database.Connection().Model(&galleryModel.Gallery{}).Where("uuid = ?", gallery.UUID.String()).First(&checkRow).Error

	if err != nil {
		return checkRow, err
	}

	err = repository.database.Connection().Model(&checkRow).Updates(gallery).Error

	if err != nil {
		return gallery, err
	}

	return gallery, nil
}

func (repository *GalleryRepository) Delete(slug string) (err error) {
	var gallery galleryModel.Gallery

	err = repository.database.Connection().Model(&galleryModel.Gallery{}).Where("slug = ?", slug).First(&gallery).Error

	if err != nil {
		return err
	}

	err = repository.database.Connection().Delete(&gallery).Error

	if err != nil {
		return err
	}

	return nil
}
