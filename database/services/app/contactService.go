package service

import (
	contactDto "github.com/horlakz/energaan-api/database/dto/app"
	contactModel "github.com/horlakz/energaan-api/database/model/app"
	"github.com/horlakz/energaan-api/database/repository"
	contactRepository "github.com/horlakz/energaan-api/database/repository/app"
)

type ContactService struct {
	contactRepository contactRepository.ContactRespositoryInterface
}

type ContactServiceInterface interface {
	Create(contact contactDto.ContactDTO) (contactDto.ContactDTO, error)
	Read(email string) (contactDto.ContactDTO, error)
	ReadAll(pageable repository.Pageable) ([]contactDto.ContactDTO, repository.Pagination, error)
	// Update(contact contactDto.ContactDTO) (contactDto.ContactDTO, error)
	// Delete(userUUID uuid.UUID, id uuid.UUID) error
}

func NewContactService(contactRepository contactRepository.ContactRespositoryInterface) ContactServiceInterface {
	return &ContactService{contactRepository: contactRepository}
}

func (service *ContactService) ConvertToDTO(contact contactModel.Contact) contactDto.ContactDTO {
	var contactDTO contactDto.ContactDTO
	contactDTO.UUID = contact.UUID
	contactDTO.FullName = contact.FullName
	contactDTO.Email = contact.Email
	contactDTO.Phone = contact.Phone
	contactDTO.Country = contact.Country
	contactDTO.Message = contact.Message
	contactDTO.CreatedAt = contact.CreatedAt
	contactDTO.UpdatedAt = contact.UpdatedAt
	contactDTO.DeletedAt = contact.DeletedAt.Time
	return contactDTO
}

func (service *ContactService) ConvertToModel(contactDTO contactDto.ContactDTO) contactModel.Contact {
	var contact contactModel.Contact
	contact.UUID = contactDTO.UUID
	contact.FullName = contactDTO.FullName
	contact.Email = contactDTO.Email
	contact.Phone = contactDTO.Phone
	contact.Country = contactDTO.Country
	contact.Message = contactDTO.Message
	contact.CreatedAt = contactDTO.CreatedAt
	contact.UpdatedAt = contactDTO.UpdatedAt
	contact.DeletedAt.Time = contactDTO.DeletedAt
	return contact
}

func (service *ContactService) Create(contactDTO contactDto.ContactDTO) (contactDto.ContactDTO, error) {
	contact := service.ConvertToModel(contactDTO)
	newRecord, err := service.contactRepository.Create(contact)

	return service.ConvertToDTO(newRecord), err
}

func (service *ContactService) Read(email string) (contactDto.ContactDTO, error) {
	record, err := service.contactRepository.Read(email)

	return service.ConvertToDTO(record), err
}

func (service *ContactService) ReadAll(pageable repository.Pageable) (recordsDto []contactDto.ContactDTO, pagination repository.Pagination, err error) {
	records, pagination, err := service.contactRepository.ReadAll(pageable)

	for _, record := range records {
		recordsDto = append(recordsDto, service.ConvertToDTO(record))
	}

	return recordsDto, pagination, err
}

// func (service *ContactService) Update(contactDTO contactDto.ContactDTO) (contactDto.ContactDTO, error) {
// 	user := service.ConvertToModel(contactDTO)
// 	newRecord, err := service.contactRepository.Update(user)

// 	return service.ConvertToDTO(newRecord), err
// }

// func (service *ContactService) Delete(userUUID uuid.UUID, id uuid.UUID) (err error) {
// 	rtn := service.contactRepository.Delete(id)
// 	record, _ := service.contactRepository.Read(id)
// 	record.DeletedByID = userUUID

// 	service.contactRepository.Update(record)

// 	return rtn
// }
