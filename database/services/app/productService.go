package service

import (
	"github.com/google/uuid"

	productDto "github.com/horlakz/energaan-api/database/dto/app"
	productModel "github.com/horlakz/energaan-api/database/model/app"
	"github.com/horlakz/energaan-api/database/repository"
	productRepository "github.com/horlakz/energaan-api/database/repository/app"
)

type productService struct {
	productRepository productRepository.ProductRespositoryInterface
}

type ProductServiceInterface interface {
	Create(product productDto.ProductDTO) (productDto.ProductDTO, error)
	Read(slug string) (productDto.ProductDTO, error)
	ReadByUUID(uuid uuid.UUID) (productDto.ProductDTO, error)
	ReadAll(pageable repository.Pageable, categoryId uuid.UUID) ([]productDto.ProductDTO, repository.Pagination, error)
	Update(product productDto.ProductDTO) (productDto.ProductDTO, error)
	Delete(userUUID uuid.UUID, slug string) error
}

func NewProductService(productRepository productRepository.ProductRespositoryInterface) ProductServiceInterface {
	return &productService{productRepository: productRepository}
}

func (service *productService) ConvertToDTO(product productModel.Product) productDto.ProductDTO {
	var productDTO productDto.ProductDTO
	productDTO.UUID = product.UUID
	productDTO.Slug = product.Slug
	productDTO.Title = product.Title
	productDTO.CategoryID = product.CategoryID
	productDTO.Images = product.Images
	productDTO.Description = product.Description
	productDTO.Features = product.Features
	productDTO.CreatedAt = product.CreatedAt
	productDTO.UpdatedAt = product.UpdatedAt
	productDTO.DeletedAt = product.DeletedAt.Time
	productDTO.CreatedByID = product.CreatedByID
	productDTO.UpdatedByID = product.UpdatedByID
	productDTO.DeletedByID = product.DeletedByID
	return productDTO
}

func (service *productService) ConvertToModel(productDTO productDto.ProductDTO) productModel.Product {
	var product productModel.Product
	product.UUID = productDTO.UUID
	product.Slug = productDTO.Slug
	product.Title = productDTO.Title
	product.CategoryID = productDTO.CategoryID
	product.Images = productDTO.Images
	product.Description = productDTO.Description
	product.Features = productDTO.Features
	product.CreatedAt = productDTO.CreatedAt
	product.UpdatedAt = productDTO.UpdatedAt
	product.DeletedAt.Time = productDTO.DeletedAt
	product.CreatedByID = productDTO.CreatedByID
	product.UpdatedByID = productDTO.UpdatedByID
	product.DeletedByID = productDTO.DeletedByID
	return product
}

func (service *productService) Create(productDTO productDto.ProductDTO) (productDto.ProductDTO, error) {
	product := service.ConvertToModel(productDTO)
	newRecord, err := service.productRepository.Create(product)

	return service.ConvertToDTO(newRecord), err
}

func (service *productService) Read(slug string) (productDto.ProductDTO, error) {
	record, err := service.productRepository.Read(slug)

	return service.ConvertToDTO(record), err
}

func (service *productService) ReadByUUID(uuid uuid.UUID) (productDto.ProductDTO, error) {
	record, err := service.productRepository.ReadByUUID(uuid)

	return service.ConvertToDTO(record), err
}

func (service *productService) ReadAll(pageable repository.Pageable, categoryId uuid.UUID) (recordsDto []productDto.ProductDTO, pagination repository.Pagination, err error) {
	records, pagination, err := service.productRepository.ReadAll(pageable, categoryId)

	for _, record := range records {
		recordsDto = append(recordsDto, service.ConvertToDTO(record))
	}

	return recordsDto, pagination, err
}

func (service *productService) Update(productDTO productDto.ProductDTO) (productDto.ProductDTO, error) {
	user := service.ConvertToModel(productDTO)
	newRecord, err := service.productRepository.Update(user)

	return service.ConvertToDTO(newRecord), err
}

func (service *productService) Delete(userUUID uuid.UUID, slug string) (err error) {
	rtn := service.productRepository.Delete(slug)
	record, _ := service.productRepository.Read(slug)
	record.DeletedByID = userUUID

	service.productRepository.Update(record)

	return rtn
}
