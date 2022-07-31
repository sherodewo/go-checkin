package service

import (
	"fmt"
	"go-checkin/dto"
	"go-checkin/models"
	"go-checkin/repository"
	"gorm.io/gorm"
)

type HomeService struct {
	HomeRepository repository.HomeRepository
}

func NewHomeService(repository repository.HomeRepository) *HomeService {
	return &HomeService{
		HomeRepository: repository,
	}
}

func (s *HomeService) QueryDatatable(searchValue string, orderType string, orderBy string, limit int, offset int) (
	recordTotal int64, recordFiltered int64, data []models.PresenceDatatable, err error) {
	recordTotal, err = s.HomeRepository.Count()

	if searchValue != "" {
		recordFiltered, err = s.HomeRepository.CountWhere("or", map[string]interface{}{
			"location LIKE ?": "%" + searchValue + "%",
			"id LIKE ?":       "%" + searchValue + "%",
		})

		data, err = s.HomeRepository.FindAllWhere("or", orderType, "created_at", limit, offset, map[string]interface{}{
			"location LIKE ?": "%" + searchValue + "%",
			"id LIKE ?":       "%" + searchValue + "%",
		})
		return recordTotal, recordFiltered, data, err
	}

	recordFiltered, err = s.HomeRepository.CountWhere("or", map[string]interface{}{
		"1 =?": 1,
	})

	data, err = s.HomeRepository.FindAllWhere("or", orderType, "created_at", limit, offset, map[string]interface{}{
		"1= ?": 1,
	})
	return recordTotal, recordFiltered, data, err
}

func (s *HomeService) GetDbInstance() *gorm.DB {
	return s.HomeRepository.DbInstance()
}

func (s *HomeService) GetAllPresence(req dto.Excel) ([]models.PresenceDatatable, error) {
	var sql string
	var times string
	sql = fmt.Sprint("SELECT a.*, b.name from presence a LEFT JOIN user b on b.user_id = a.user_id")

	if req.Start == req.End {
		times = fmt.Sprintf(" DATE(a.created_at) = CURDATE()")
		req.Start = ""
		req.End = ""
	}
	if req.ID != "" && (req.Start == "" || req.End == "") {
		fmt.Println("MASUK 1")
		sql = fmt.Sprintf("SELECT a.*, b.name from presence a LEFT JOIN user b on b.user_id = a.user_id WHERE b.user_id = '%s' AND %s", req.ID, times)
	} else if req.ID == "" && (req.Start == "" || req.End == "") {
		sql = fmt.Sprintf("SELECT a.*, b.name from presence a LEFT JOIN user b on b.user_id = b.user_id WHERE %s ", times)
	}
	if req.ID != "" && (req.Start != "" || req.End != "") {
		fmt.Println("MASUK 2")
		sql += fmt.Sprintf(" WHERE b.user_id = '%s' AND a.created_at BETWEEN '%s' AND '%s'", req.ID, req.Start, req.End)
	}
	if req.ID == "" && (req.Start != "" || req.End != "") {
		fmt.Println("MASUK 3")
		sql += fmt.Sprintf(" WHERE a.created_at BETWEEN '%s' AND '%s'", req.Start, req.End)
	}

	data, err := s.HomeRepository.FindAllByParam(sql)
	if err != nil {
		return data, err
	}

	return data, nil
}
