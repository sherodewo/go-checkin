package repository

import (
	"go-checkin/models"
	"gorm.io/gorm"
)

type DivisiRepository interface {
	FindAll() ([]models.Divisi, error)
	FindAllWhere(operation string, orderType string, orderBy string, limit int, offset int, keyVal map[string]interface{}) ([]models.Divisi, error)
	FindById(id string) (models.Divisi, error)
	FindWhere(email string) (models.Divisi, error)
	//FindBranch(id string) (models.UserDetail, error)
	Save(role models.Divisi) (models.Divisi, error)
	Update(role models.Divisi) (models.Divisi, error)
	Delete(role models.Divisi) error
	Count() (int64, error)
	CountWhere(operation string, keyVal map[string]interface{}) (int64, error)
	DbInstance() *gorm.DB
}

type divisiRepository struct {
	*gorm.DB
}

func NewDivisiRepository(db *gorm.DB) DivisiRepository {
	return &divisiRepository{DB: db}
}
func (r divisiRepository) FindAll() ([]models.Divisi, error) {
	var entities []models.Divisi
	err := r.DB.Find(&entities).Error
	return entities, err
}

func (r divisiRepository) FindAllWhere(operation string, orderType string, orderBy string, limit int, offset int, keyVal map[string]interface{}) ([]models.Divisi, error) {
	var entity []models.Divisi
	q := r.DB.Order(orderBy + " " + orderType).Limit(limit).Offset(offset)

	for k, v := range keyVal {
		switch operation {
		case "and":
			q = q.Where(k, v)
		case "or":
			q = q.Or(k, v)
		}
	}

	err := q.Find(&entity).Error
	return entity, err
}

func (r divisiRepository) FindById(id string) (models.Divisi, error) {
	var entity models.Divisi
	err := r.DB.Where("id = ?", id).First(&entity).Error
	return entity, err
}

func (r divisiRepository) FindWhere(name string) (models.Divisi, error) {
	var entity models.Divisi
	err := r.DB.Where("name = ?", name).Find(&entity).Error
	return entity, err
}

//func (r roleRepository) FindBranch(id string) (models.UserDetail, error) {
//	var entity models.UserDetail
//	err := r.DB.Where("user_id = ?", id).Find(&entity).Error
//	return entity, err
//}

func (r divisiRepository) Save(entity models.Divisi) (models.Divisi, error) {
	err := r.DB.Create(&entity).Error
	return entity, err
}

func (r divisiRepository) Update(entity models.Divisi) (models.Divisi, error) {
	err := r.DB.Model(models.Divisi{ID: entity.ID}).UpdateColumns(&entity).Error
	return entity, err
}

func (r divisiRepository) Delete(entity models.Divisi) error {
	return r.DB.Delete(&entity).Error
}

func (r divisiRepository) Count() (int64, error) {
	var count int64
	err := r.DB.Table("divisi").Count(&count).Error
	return count, err
}

func (r divisiRepository) CountWhere(operation string, keyVal map[string]interface{}) (int64, error) {
	var count int64
	q := r.DB.Model(&models.Divisi{})
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
func (r divisiRepository) DbInstance() *gorm.DB {
	return r.DB
}
