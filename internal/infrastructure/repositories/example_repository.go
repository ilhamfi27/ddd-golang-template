package repositories

import (
	"context"

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

func (r *GormExampleRepository) GetExamples(ctx context.Context) []model.ExampleModel {
	var examples []model.ExampleModel
	r.db.WithContext(ctx).Find(&examples)
	return examples
}
func (r *GormExampleRepository) GetExampleByID(ctx context.Context, id uint) model.ExampleModel {
	var example model.ExampleModel
	r.db.WithContext(ctx).First(&example, id)
	return example
}
func (r *GormExampleRepository) UpdateExample(ctx context.Context, id uint, example model.ExampleModel) error {
	return r.db.WithContext(ctx).Model(&model.ExampleModel{}).Where("id = ?", id).Updates(example).Error
}

func (r *GormExampleRepository) CreateExample(ctx context.Context, example model.ExampleModel) error {
	return r.db.WithContext(ctx).Create(&example).Error
}

func (r *GormExampleRepository) DeleteExample(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.ExampleModel{}, id).Error
}
