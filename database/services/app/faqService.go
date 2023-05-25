package service

import (
	"github.com/google/uuid"

	faqDto "github.com/horlakz/energaan-api/database/dto/app"
	faqModel "github.com/horlakz/energaan-api/database/model/app"
	faqRepository "github.com/horlakz/energaan-api/database/repository/app"
)

type FaqService struct {
	faqRepository faqRepository.FaqRespositoryInterface
}

type FaqServiceInterface interface {
	Create(faq faqDto.FaqDTO) (faqDto.FaqDTO, error)
	Read(id uuid.UUID) (faqDto.FaqDTO, error)
	ReadAll() ([]faqDto.FaqDTO, error)
	Update(faq faqDto.FaqDTO) (faqDto.FaqDTO, error)
	Delete(userUUID uuid.UUID, id uuid.UUID) error
}

func NewFaqService(faqRepository faqRepository.FaqRespositoryInterface) FaqServiceInterface {
	return &FaqService{faqRepository: faqRepository}
}

func (service *FaqService) ConvertToDTO(faq faqModel.Faq) faqDto.FaqDTO {
	var faqDTO faqDto.FaqDTO
	faqDTO.UUID = faq.UUID
	faqDTO.Title = faq.Title
	faqDTO.Description = faq.Description
	faqDTO.CreatedAt = faq.CreatedAt
	faqDTO.UpdatedAt = faq.UpdatedAt
	faqDTO.DeletedAt = faq.DeletedAt.Time
	faqDTO.CreatedByID = faq.CreatedByID
	faqDTO.UpdatedByID = faq.UpdatedByID
	faqDTO.DeletedByID = faq.DeletedByID
	return faqDTO
}

func (service *FaqService) ConvertToModel(faqDTO faqDto.FaqDTO) faqModel.Faq {
	var faq faqModel.Faq
	faq.UUID = faqDTO.UUID
	faq.Title = faqDTO.Title
	faq.Description = faqDTO.Description
	faq.CreatedAt = faqDTO.CreatedAt
	faq.UpdatedAt = faqDTO.UpdatedAt
	faq.DeletedAt.Time = faqDTO.DeletedAt
	faq.CreatedByID = faqDTO.CreatedByID
	faq.UpdatedByID = faqDTO.UpdatedByID
	faq.DeletedByID = faqDTO.DeletedByID
	return faq
}

func (service *FaqService) Create(faqDTO faqDto.FaqDTO) (faqDto.FaqDTO, error) {
	faq := service.ConvertToModel(faqDTO)
	newRecord, err := service.faqRepository.Create(faq)

	return service.ConvertToDTO(newRecord), err
}

func (service *FaqService) Read(id uuid.UUID) (faqDto.FaqDTO, error) {
	record, err := service.faqRepository.Read(id)

	return service.ConvertToDTO(record), err
}

func (service *FaqService) ReadAll() (recordsDto []faqDto.FaqDTO, err error) {
	records, err := service.faqRepository.ReadAll()

	for _, record := range records {
		recordsDto = append(recordsDto, service.ConvertToDTO(record))

	}

	return recordsDto, err
}

func (service *FaqService) Update(faqDTO faqDto.FaqDTO) (faqDto.FaqDTO, error) {
	user := service.ConvertToModel(faqDTO)
	newRecord, err := service.faqRepository.Update(user)

	return service.ConvertToDTO(newRecord), err
}

func (service *FaqService) Delete(userUUID uuid.UUID, id uuid.UUID) (err error) {
	rtn := service.faqRepository.Delete(id)
	record, _ := service.faqRepository.Read(id)
	record.DeletedByID = userUUID

	service.faqRepository.Update(record)

	return rtn
}
