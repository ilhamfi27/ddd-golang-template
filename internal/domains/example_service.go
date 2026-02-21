package domains

import (
	"context"

	"github.com/ilhamfi27/ddd-golang-template/internal/application/dto"
	"github.com/ilhamfi27/ddd-golang-template/internal/infrastructure/repositories"
)

type ExampleService struct {
	repo *repositories.GormExampleRepository
}

func NewExampleService(repo *repositories.GormExampleRepository) *ExampleService {
	return &ExampleService{
		repo: repo,
	}
}

func (s *ExampleService) Hello(ctx context.Context, data dto.ParseExampleDto) (map[string]interface{}, error) {
	examples := s.repo.GetExamples(ctx)
	return map[string]interface{}{
		"name":     data.Name,
		"examples": examples,
	}, nil
}
