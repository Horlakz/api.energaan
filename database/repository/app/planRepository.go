package app

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"

	databaseModule "github.com/horlakz/energaan-api/database"
	planModel "github.com/horlakz/energaan-api/database/model/app"
	"github.com/horlakz/energaan-api/database/repository"
)

type PlanRespositoryInterface interface {
	Create(plan planModel.Plan) (planModel.Plan, error)
	Read(slug string) (planModel.Plan, error)
	ReadByUUID(uuid uuid.UUID) (planModel.Plan, error)
	ReadAll(pageable repository.Pageable) ([]planModel.Plan, repository.Pagination, error)
	Update(plan planModel.Plan) (planModel.Plan, error)
	Delete(slug string) error
}

type planRepository struct {
	database databaseModule.DatabaseInterface
}

func NewPlanRepository(database databaseModule.DatabaseInterface) PlanRespositoryInterface {
	return &planRepository{database: database}
}

func (repository *planRepository) Create(plan planModel.Plan) (planModel.Plan, error) {
	plan.Model.Prepare()

	err := repository.database.Connection().Create(&plan).Error

	if err != nil {
		return plan, err
	}

	return plan, nil
}

func (repository *planRepository) Read(slug string) (plan planModel.Plan, err error) {
	err = repository.database.Connection().Model(&planModel.Plan{}).Where("slug = ?", slug).First(&plan).Error

	if err != nil {
		return plan, err
	}

	return plan, nil
}

func (repository *planRepository) ReadByUUID(uuid uuid.UUID) (plan planModel.Plan, err error) {
	err = repository.database.Connection().Model(&planModel.Plan{}).Where("uuid = ?", uuid).First(&plan).Error

	if err != nil {
		return plan, err
	}

	return plan, nil
}

func (repository *planRepository) ReadAll(pageable repository.Pageable) (rows []planModel.Plan, pagination repository.Pagination, err error) {
	var plan planModel.Plan
	pagination.TotalPages = 1
	pagination.TotalItems = 0
	pagination.CurrentPage = int64(pageable.Page)

	var result *gorm.DB
	var errCount error

	offset := (pageable.Page - 1) * pageable.Size
	searchQuery := repository.database.Connection().Model(&planModel.Plan{})

	if len(strings.TrimSpace(pageable.Search)) > 0 {
		searchQuery.Where("title LIKE ?", "%"+strings.ToLower(pageable.Search)+"%")
	}

	errCount = searchQuery.Count(&pagination.TotalItems).Error
	paginationQuery := searchQuery.Limit(pageable.Size).Offset(offset).Order(pageable.SortBy + " " + pageable.SortDirection)

	result = paginationQuery.Model(&planModel.Plan{}).Where(plan).Find(&rows)

	if result.Error != nil {
		return nil, pagination, result.Error
	}

	if errCount != nil {
		return nil, pagination, errCount
	}

	pagination.TotalPages = pagination.TotalItems / int64(pageable.Size)

	return rows, pagination, nil
}

func (repository *planRepository) Update(plan planModel.Plan) (planModel.Plan, error) {
	var checkRow planModel.Plan

	err := repository.database.Connection().Model(&planModel.Plan{}).Where("uuid = ?", plan.UUID.String()).First(&checkRow).Error

	if err != nil {
		return checkRow, err
	}

	err = repository.database.Connection().Model(&checkRow).Updates(plan).Error

	if err != nil {
		return plan, err
	}

	return plan, nil
}

func (repository *planRepository) Delete(slug string) (err error) {
	var plan planModel.Plan

	err = repository.database.Connection().Model(&planModel.Plan{}).Where("slug = ?", slug).First(&plan).Error

	if err != nil {
		return err
	}

	err = repository.database.Connection().Delete(&plan).Error

	fmt.Println(err, " this is the error")

	if err != nil {
		return err
	}

	return nil
}
