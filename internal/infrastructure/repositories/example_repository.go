package repositories

import (
	model "github.com/ilhamfi27/ddd-golang-template/internal/models"
	"gorm.io/gorm"
)

type GormExampleRepository struct {
	db *gorm.DB
}

func NewGormExampleRepository(db *gorm.DB) *GormExampleRepository {
	return &GormExampleRepository{
		db: db,
	}
}

func (r *GormExampleRepository) GetExamples() []model.ExampleModel {
	var examples []model.ExampleModel
	r.db.Find(&examples)
	return examples
}
func (r *GormExampleRepository) GetExampleByID(id uint) model.ExampleModel {
	var example model.ExampleModel
	r.db.First(&example, id)
	return example
}
func (r *GormExampleRepository) UpdateExample(id uint, example model.ExampleModel) error {
	return r.db.Model(&model.ExampleModel{}).Where("id = ?", id).Updates(example).Error
}

func (r *GormExampleRepository) CreateExample(example model.ExampleModel) error {
	return r.db.Create(&example).Error
}

func (r *GormExampleRepository) DeleteExample(id uint) error {
	return r.db.Delete(&model.ExampleModel{}, id).Error
}
