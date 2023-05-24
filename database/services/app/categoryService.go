package service

import (
	"github.com/google/uuid"

	categoryDto "github.com/horlakz/energaan-api/database/dto/app"
	categoryModel "github.com/horlakz/energaan-api/database/model/app"
	categoryRepository "github.com/horlakz/energaan-api/database/repository/app"
)

type categoryService struct {
	categoryRepository categoryRepository.CategoryRespositoryInterface
}

type CategoryServiceInterface interface {
	Create(category categoryDto.CategoryDTO) (categoryDto.CategoryDTO, error)
	Read(slug string) (categoryDto.CategoryDTO, error)
	ReadAll() ([]categoryDto.CategoryDTO, error)
	Update(category categoryDto.CategoryDTO) (categoryDto.CategoryDTO, error)
	Delete(userUUID uuid.UUID, slug string) error
}

func NewCategoryService(categoryRepository categoryRepository.CategoryRespositoryInterface) CategoryServiceInterface {
	return &categoryService{categoryRepository: categoryRepository}
}

func (service *categoryService) ConvertToDTO(category categoryModel.Category) categoryDto.CategoryDTO {
	var categoryDTO categoryDto.CategoryDTO
	categoryDTO.UUID = category.UUID
	categoryDTO.Slug = category.Slug
	categoryDTO.Name = category.Name
	categoryDTO.CreatedAt = category.CreatedAt
	categoryDTO.UpdatedAt = category.UpdatedAt
	categoryDTO.DeletedAt = category.DeletedAt.Time
	categoryDTO.CreatedByID = category.CreatedByID
	categoryDTO.UpdatedByID = category.UpdatedByID
	categoryDTO.DeletedByID = category.DeletedByID
	return categoryDTO
}

func (service *categoryService) ConvertToModel(categoryDTO categoryDto.CategoryDTO) categoryModel.Category {
	var category categoryModel.Category
	category.UUID = categoryDTO.UUID
	category.Slug = categoryDTO.Slug
	category.Name = categoryDTO.Name
	category.CreatedAt = categoryDTO.CreatedAt
	category.UpdatedAt = categoryDTO.UpdatedAt
	category.DeletedAt.Time = categoryDTO.DeletedAt
	category.CreatedByID = categoryDTO.CreatedByID
	category.UpdatedByID = categoryDTO.UpdatedByID
	category.DeletedByID = categoryDTO.DeletedByID
	return category
}

func (service *categoryService) Create(categoryDTO categoryDto.CategoryDTO) (categoryDto.CategoryDTO, error) {
	category := service.ConvertToModel(categoryDTO)
	newRecord, err := service.categoryRepository.Create(category)

	return service.ConvertToDTO(newRecord), err
}

func (service *categoryService) Read(slug string) (categoryDto.CategoryDTO, error) {
	record, err := service.categoryRepository.Read(slug)

	return service.ConvertToDTO(record), err
}

func (service *categoryService) ReadAll() (recordsDto []categoryDto.CategoryDTO, err error) {
	records, err := service.categoryRepository.ReadAll()

	for _, record := range records {
		recordsDto = append(recordsDto, service.ConvertToDTO(record))

	}

	return recordsDto, err
}

func (service *categoryService) Update(categoryDTO categoryDto.CategoryDTO) (categoryDto.CategoryDTO, error) {
	user := service.ConvertToModel(categoryDTO)
	newRecord, err := service.categoryRepository.Update(user)

	return service.ConvertToDTO(newRecord), err
}

func (service *categoryService) Delete(userUUID uuid.UUID, slug string) (err error) {
	rtn := service.categoryRepository.Delete(slug)
	record, _ := service.categoryRepository.Read(slug)
	record.DeletedByID = userUUID

	service.categoryRepository.Update(record)

	return rtn
}
