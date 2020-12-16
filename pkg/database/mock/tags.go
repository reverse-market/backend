package mock

import (
	"context"
	"github.com/reverse-market/backend/pkg/database/models"
	"strings"
)

type TagsRepository struct {
}

var (
	categoryID1 = 1
	categoryID2 = 2
	mockTags    = []*models.Tag{
		{
			ID:         1,
			CategoryID: &categoryID1,
			Name:       "Синий",
		},
		{
			ID:         2,
			CategoryID: &categoryID2,
			Name:       "Красный",
		},
		{
			ID:         3,
			CategoryID: nil,
			Name:       "Зелёный",
		},
	}
)

func (tr *TagsRepository) GetByID(ctx context.Context, id int) (*models.Tag, error) {
	if id > 0 && id <= len(mockTags) {
		return mockTags[id-1], nil
	} else {
		return nil, models.ErrNoRecord
	}
}

func (tr *TagsRepository) GetAll(ctx context.Context, filters *models.TagFilters) ([]*models.Tag, error) {
	tags := make([]*models.Tag, 0)
	for _, tag := range mockTags {
		if (filters.CategoryID == nil || tag.CategoryID == nil || *tag.CategoryID == *filters.CategoryID) &&
			(filters.Search == "" || strings.Contains(strings.ToLower(tag.Name), filters.Search)) {
			tags = append(tags, tag)
		}
	}
	return tags, nil
}
