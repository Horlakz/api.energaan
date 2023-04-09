package service

import (
	"github.com/google/uuid"

	planDto "github.com/horlakz/energaan-api/database/dto/app"
	planModel "github.com/horlakz/energaan-api/database/model/app"
	"github.com/horlakz/energaan-api/database/repository"
	planRepository "github.com/horlakz/energaan-api/database/repository/app"
)

type planService struct {
	planRepository planRepository.PlanRespositoryInterface
}

type PlanServiceInterface interface {
	Create(plan planDto.PlanDTO) (planDto.PlanDTO, error)
	Read(slug string) (planDto.PlanDTO, error)
	ReadAll(pageable repository.Pageable) ([]planDto.PlanDTO, repository.Pagination, error)
	Update(plan planDto.PlanDTO) (planDto.PlanDTO, error)
	Delete(userUUID uuid.UUID, slug string) error
}

func NewPlanService(planRepository planRepository.PlanRespositoryInterface) PlanServiceInterface {
	return &planService{planRepository: planRepository}
}

func (service *planService) ConvertToDTO(plan planModel.Plan) planDto.PlanDTO {
	var planDTO planDto.PlanDTO
	planDTO.UUID = plan.UUID
	planDTO.Slug = plan.Slug
	planDTO.Title = plan.Title
	planDTO.Image = plan.Image
	planDTO.Description = plan.Description
	planDTO.Price = plan.Price
	planDTO.CreatedAt = plan.CreatedAt
	planDTO.UpdatedAt = plan.UpdatedAt
	planDTO.DeletedAt = plan.DeletedAt.Time
	planDTO.CreatedByID = plan.CreatedByID
	planDTO.UpdatedByID = plan.UpdatedByID
	planDTO.DeletedByID = plan.DeletedByID
	return planDTO
}

func (service *planService) ConvertToModel(planDTO planDto.PlanDTO) planModel.Plan {
	var plan planModel.Plan
	plan.UUID = planDTO.UUID
	plan.Slug = planDTO.Slug
	plan.Title = planDTO.Title
	plan.Image = planDTO.Image
	plan.Description = planDTO.Description
	plan.Price = planDTO.Price
	plan.CreatedAt = planDTO.CreatedAt
	plan.UpdatedAt = planDTO.UpdatedAt
	plan.DeletedAt.Time = planDTO.DeletedAt
	plan.CreatedByID = planDTO.CreatedByID
	plan.UpdatedByID = planDTO.UpdatedByID
	plan.DeletedByID = planDTO.DeletedByID
	return plan
}

func (service *planService) Create(planDTO planDto.PlanDTO) (planDto.PlanDTO, error) {
	plan := service.ConvertToModel(planDTO)
	newRecord, err := service.planRepository.Create(plan)

	return service.ConvertToDTO(newRecord), err
}

func (service *planService) Read(slug string) (planDto.PlanDTO, error) {
	record, err := service.planRepository.Read(slug)

	return service.ConvertToDTO(record), err
}

func (service *planService) ReadAll(pageable repository.Pageable) (recordsDto []planDto.PlanDTO, pagination repository.Pagination, err error) {
	records, pagination, err := service.planRepository.ReadAll(pageable)

	for _, record := range records {
		recordsDto = append(recordsDto, service.ConvertToDTO(record))
	}

	return recordsDto, pagination, err
}

func (service *planService) Update(planDTO planDto.PlanDTO) (planDto.PlanDTO, error) {
	user := service.ConvertToModel(planDTO)
	newRecord, err := service.planRepository.Update(user)

	return service.ConvertToDTO(newRecord), err
}

func (service *planService) Delete(userUUID uuid.UUID, slug string) (err error) {
	rtn := service.planRepository.Delete(slug)
	record, _ := service.planRepository.Read(slug)
	record.DeletedByID = userUUID

	service.planRepository.Update(record)

	return rtn
}
