package service

import (
	quoteDto "github.com/horlakz/energaan-api/database/dto/app"
	quoteModel "github.com/horlakz/energaan-api/database/model/app"
	"github.com/horlakz/energaan-api/database/repository"
	quoteRepository "github.com/horlakz/energaan-api/database/repository/app"
)

type QuoteService struct {
	quoteRepository quoteRepository.QuoteRespositoryInterface
}

type QuoteServiceInterface interface {
	Create(quote quoteDto.QuoteDTO) (quoteDto.QuoteDTO, error)
	Read(email string) (quoteDto.QuoteDTO, error)
	CheckEmailServiceTypeExists(email string, serviceType string) (bool, error)
	ReadAll(pageable repository.Pageable) ([]quoteDto.QuoteDTO, repository.Pagination, error)
	// Update(quote quoteDto.QuoteDTO) (quoteDto.QuoteDTO, error)
	// Delete(userUUID uuid.UUID, id uuid.UUID) error
}

func NewQuoteService(quoteRepository quoteRepository.QuoteRespositoryInterface) QuoteServiceInterface {
	return &QuoteService{quoteRepository: quoteRepository}
}

func (service *QuoteService) ConvertToDTO(quote quoteModel.Quote) quoteDto.QuoteDTO {
	var quoteDTO quoteDto.QuoteDTO
	quoteDTO.UUID = quote.UUID
	quoteDTO.FullName = quote.FullName
	quoteDTO.Email = quote.Email
	quoteDTO.ServiceId = quote.ServiceId
	quoteDTO.ServiceType = quote.ServiceType
	quoteDTO.Phone = quote.Phone
	quoteDTO.Country = quote.Country
	quoteDTO.CreatedAt = quote.CreatedAt
	quoteDTO.UpdatedAt = quote.UpdatedAt
	quoteDTO.DeletedAt = quote.DeletedAt.Time
	return quoteDTO
}

func (service *QuoteService) ConvertToModel(quoteDTO quoteDto.QuoteDTO) quoteModel.Quote {
	var quote quoteModel.Quote
	quote.UUID = quoteDTO.UUID
	quote.FullName = quoteDTO.FullName
	quote.Email = quoteDTO.Email
	quote.ServiceId = quoteDTO.ServiceId
	quote.ServiceType = quoteDTO.ServiceType
	quote.Phone = quoteDTO.Phone
	quote.Country = quoteDTO.Country
	quote.CreatedAt = quoteDTO.CreatedAt
	quote.UpdatedAt = quoteDTO.UpdatedAt
	quote.DeletedAt.Time = quoteDTO.DeletedAt
	return quote
}

func (service *QuoteService) Create(quoteDTO quoteDto.QuoteDTO) (quoteDto.QuoteDTO, error) {
	quote := service.ConvertToModel(quoteDTO)
	newRecord, err := service.quoteRepository.Create(quote)

	return service.ConvertToDTO(newRecord), err
}

func (service *QuoteService) Read(email string) (quoteDto.QuoteDTO, error) {
	record, err := service.quoteRepository.Read(email)

	return service.ConvertToDTO(record), err
}

func (service *QuoteService) CheckEmailServiceTypeExists(email string, serviceType string) (bool, error) {
	record, err := service.quoteRepository.CheckEmailServiceTypeExists(email, serviceType)

	if record.Email != "" {
		return true, err
	}

	return false, err
}

func (service *QuoteService) ReadAll(pageable repository.Pageable) (recordsDto []quoteDto.QuoteDTO, pagination repository.Pagination, err error) {
	records, pagination, err := service.quoteRepository.ReadAll(pageable)

	for _, record := range records {
		recordsDto = append(recordsDto, service.ConvertToDTO(record))
	}

	return recordsDto, pagination, err
}

// func (service *QuoteService) Update(quoteDTO quoteDto.QuoteDTO) (quoteDto.QuoteDTO, error) {
// 	user := service.ConvertToModel(quoteDTO)
// 	newRecord, err := service.quoteRepository.Update(user)

// 	return service.ConvertToDTO(newRecord), err
// }

// func (service *QuoteService) Delete(userUUID uuid.UUID, id uuid.UUID) (err error) {
// 	rtn := service.quoteRepository.Delete(id)
// 	record, _ := service.quoteRepository.Read(id)
// 	record.DeletedByID = userUUID

// 	service.quoteRepository.Update(record)

// 	return rtn
// }
