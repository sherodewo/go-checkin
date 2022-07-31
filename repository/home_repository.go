package repository

import (
	"fmt"
	"go-checkin/models"
	"gorm.io/gorm"
)

type HomeRepository interface {
	FindAll() ([]models.Presence, error)
	FindAllByParam(sql string) ([]models.PresenceDatatable, error)
	FindAllWhere(operation string, orderType string, orderBy string, limit int, offset int, keyVal map[string]interface{}) ([]models.PresenceDatatable, error)
	FindById(id string) (models.Presence, error)
	FindWhere(email string) (models.Presence, error)
	Count() (int64, error)
	CountWhere(operation string, keyVal map[string]interface{}) (int64, error)
	DbInstance() *gorm.DB
}

type homeRepository struct {
	*gorm.DB
}

func NewHomeRepository(db *gorm.DB) HomeRepository {
	return &homeRepository{DB: db}
}

func (r homeRepository) FindAll() ([]models.Presence, error) {
	var entities []models.Presence
	err := r.DB.Find(&entities).Error
	return entities, err
}

func (r homeRepository) FindAllWhere(operation string, orderType string, orderBy string, limit int, offset int, keyVal map[string]interface{}) ([]models.PresenceDatatable, error) {
	var entity []models.PresenceDatatable
	sql := fmt.Sprint("SELECT a.*, b.name from presence a LEFT JOIN user b on b.user_id = a.user_id")

	q := r.DB.Raw(sql).Order(orderBy + " " + orderType).Limit(limit).Offset(offset)

	for k, v := range keyVal {
		switch operation {
		case "and":
			q = q.Where(k, v)
		case "or":
			q = q.Or(k, v)
		}
	}

	err := q.Scan(&entity).Error
	return entity, err
}

func (r homeRepository) FindById(id string) (models.Presence, error) {
	var entity models.Presence
	err := r.DB.Where("id = ?", id).First(&entity).Error
	return entity, err
}

func (r homeRepository) FindWhere(name string) (models.Presence, error) {
	var entity models.Presence
	err := r.DB.Where("name = ?", name).Find(&entity).Error
	return entity, err
}

func (r homeRepository) Count() (int64, error) {
	var count int64
	err := r.DB.Table("user_role").Count(&count).Error
	return count, err
}

func (r homeRepository) CountWhere(operation string, keyVal map[string]interface{}) (int64, error) {
	var count int64
	q := r.DB.Model(&models.Presence{})
	for k, v := range keyVal {
		switch operation {
		case "and":
			q = q.Where(k, v)
		case "or":
			q = q.Or(k, v)
		}
	}

	err := q.Count(&count).Error
	return count, err
}

func (r homeRepository) FindAllByParam(sql string) ([]models.PresenceDatatable, error) {
	var entities []models.PresenceDatatable
	err := r.DB.Raw(sql).
		Find(&entities).Error
	return entities, err
}

func (r homeRepository) DbInstance() *gorm.DB {
	return r.DB
}
