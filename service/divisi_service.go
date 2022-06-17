package service

import (
	"go-checkin/dto"
	"go-checkin/models"
	"go-checkin/repository"
	"gorm.io/gorm"
	"time"
)

type DivisiService struct {
	DivisiRepository repository.DivisiRepository
}

func NewDivisiService(repository repository.DivisiRepository) *DivisiService {
	return &DivisiService{
		DivisiRepository: repository,
	}
}

func (s *DivisiService) QueryDatatable(searchValue string, orderType string, orderBy string, limit int, offset int) (
	recordTotal int64, recordFiltered int64, data []models.Divisi, err error) {
	recordTotal, err = s.DivisiRepository.Count()

	if searchValue != "" {
		recordFiltered, err = s.DivisiRepository.CountWhere("or", map[string]interface{}{
			"name LIKE ?": "%" + searchValue + "%",
			"id LIKE ?":   "%" + searchValue + "%",
		})

		data, err = s.DivisiRepository.FindAllWhere("or", orderType, "created_at", limit, offset, map[string]interface{}{
			"name LIKE ?": "%" + searchValue + "%",
			"id LIKE ?":   "%" + searchValue + "%",
		})
		return recordTotal, recordFiltered, data, err
	}

	recordFiltered, err = s.DivisiRepository.CountWhere("or", map[string]interface{}{
		"1 =?": 1,
	})

	data, err = s.DivisiRepository.FindAllWhere("or", orderType, "created_at", limit, offset, map[string]interface{}{
		"1= ?": 1,
	})
	return recordTotal, recordFiltered, data, err
}
func (s *DivisiService) SaveDivisi(dto dto.DivisiDto) (*models.Divisi, error) {
	entity := models.Divisi{
		Name:        dto.Name,
		Description: dto.Description,
		CreatedAt:   time.Now(),
	}
	data, err := s.DivisiRepository.Save(entity)
	return &data, err
}
func (s *DivisiService) FindUserById(id string) (*models.Divisi, error) {
	data, err := s.DivisiRepository.FindById(id)

	return &data, err
}
func (s *DivisiService) DeleteDivisi(id string) error {
	entity := models.Divisi{
		ID: id,
	}
	err := s.DivisiRepository.Delete(entity)
	if err != nil {
		return err
	} else {
		return nil
	}
}
func (s *DivisiService) UpdateDivisi(id string, dto dto.DivisiDto) (*models.Divisi, error) {
	var entity models.Divisi
	entity.ID = id
	entity.Name = dto.Name
	entity.Description = dto.Description

	data, err := s.DivisiRepository.Update(entity)

	return &data, err
}
func (s *DivisiService) GetDbInstance() *gorm.DB {
	return s.DivisiRepository.DbInstance()
}
