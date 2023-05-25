package service

import (
	"github.com/google/uuid"

	galleryDto "github.com/horlakz/energaan-api/database/dto/app"
	galleryModel "github.com/horlakz/energaan-api/database/model/app"
	galleryRepository "github.com/horlakz/energaan-api/database/repository/app"
)

type galleryService struct {
	galleryRepository galleryRepository.GalleryRespositoryInterface
}

type GalleryServiceInterface interface {
	Create(gallery galleryDto.GalleryDTO) (galleryDto.GalleryDTO, error)
	Read(id uuid.UUID) (galleryDto.GalleryDTO, error)
	ReadAll() ([]galleryDto.GalleryDTO, error)
	Update(gallery galleryDto.GalleryDTO) (galleryDto.GalleryDTO, error)
	Delete(userUUID uuid.UUID, id uuid.UUID) error
}

func NewGalleryService(galleryRepository galleryRepository.GalleryRespositoryInterface) GalleryServiceInterface {
	return &galleryService{galleryRepository: galleryRepository}
}

func (service *galleryService) ConvertToDTO(gallery galleryModel.Gallery) galleryDto.GalleryDTO {
	var galleryDTO galleryDto.GalleryDTO
	galleryDTO.UUID = gallery.UUID
	galleryDTO.Image = gallery.Image
	galleryDTO.Title = gallery.Title
	galleryDTO.CreatedAt = gallery.CreatedAt
	galleryDTO.UpdatedAt = gallery.UpdatedAt
	galleryDTO.DeletedAt = gallery.DeletedAt.Time
	galleryDTO.CreatedByID = gallery.CreatedByID
	galleryDTO.UpdatedByID = gallery.UpdatedByID
	galleryDTO.DeletedByID = gallery.DeletedByID
	return galleryDTO
}

func (service *galleryService) ConvertToModel(galleryDTO galleryDto.GalleryDTO) galleryModel.Gallery {
	var gallery galleryModel.Gallery
	gallery.UUID = galleryDTO.UUID
	gallery.Image = galleryDTO.Image
	gallery.Title = galleryDTO.Title
	gallery.CreatedAt = galleryDTO.CreatedAt
	gallery.UpdatedAt = galleryDTO.UpdatedAt
	gallery.DeletedAt.Time = galleryDTO.DeletedAt
	gallery.CreatedByID = galleryDTO.CreatedByID
	gallery.UpdatedByID = galleryDTO.UpdatedByID
	gallery.DeletedByID = galleryDTO.DeletedByID
	return gallery
}

func (service *galleryService) Create(galleryDTO galleryDto.GalleryDTO) (galleryDto.GalleryDTO, error) {
	gallery := service.ConvertToModel(galleryDTO)
	newRecord, err := service.galleryRepository.Create(gallery)

	return service.ConvertToDTO(newRecord), err
}

func (service *galleryService) Read(id uuid.UUID) (galleryDto.GalleryDTO, error) {
	record, err := service.galleryRepository.Read(id)

	return service.ConvertToDTO(record), err
}

func (service *galleryService) ReadAll() (recordsDto []galleryDto.GalleryDTO, err error) {
	records, err := service.galleryRepository.ReadAll()

	for _, record := range records {
		recordsDto = append(recordsDto, service.ConvertToDTO(record))

	}

	return recordsDto, err
}

func (service *galleryService) Update(galleryDTO galleryDto.GalleryDTO) (galleryDto.GalleryDTO, error) {
	user := service.ConvertToModel(galleryDTO)
	newRecord, err := service.galleryRepository.Update(user)

	return service.ConvertToDTO(newRecord), err
}

func (service *galleryService) Delete(userUUID uuid.UUID, id uuid.UUID) (err error) {
	rtn := service.galleryRepository.Delete(id)
	record, _ := service.galleryRepository.Read(id)
	record.DeletedByID = userUUID

	service.galleryRepository.Update(record)

	return rtn
}
