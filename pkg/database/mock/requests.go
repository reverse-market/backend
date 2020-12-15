package mock

import (
	"context"
	"github.com/reverse-market/backend/pkg/database/models"
	"github.com/reverse-market/backend/pkg/simpletime"
	"sort"
	"strings"
	"time"
)

type RequestsRepository struct {
}

var mockRequests = []*models.Request{
	{
		ID:          1,
		UserID:      1,
		Username:    "Ivanov Ivan",
		CategoryID:  1,
		Name:        "name 1",
		ItemName:    "item name 1",
		Description: "description 1",
		Photos:      []string{},
		Price:       100,
		Quantity:    1,
		Date:        simpletime.SimpleTime(time.Date(2020, 12, 14, 0, 0, 0, 0, time.UTC)),
		Finished:    false,
		Tags: []*models.TagInRequest{
			{
				ID:   1,
				Name: "Синий",
			},
			{
				ID:   2,
				Name: "Крансый",
			},
		},
	},
	{
		ID:          2,
		UserID:      1,
		Username:    "Ivanov Ivan",
		CategoryID:  2,
		Name:        "name 2",
		ItemName:    "item name 2",
		Description: "description 2",
		Photos:      []string{},
		Price:       1000,
		Quantity:    1,
		Date:        simpletime.SimpleTime(time.Date(2020, 12, 15, 0, 0, 0, 0, time.UTC)),
		Finished:    false,
		Tags: []*models.TagInRequest{
			{
				ID:   1,
				Name: "Синий",
			},
			{
				ID:   3,
				Name: "Зелёный",
			},
		},
	},
	{
		ID:          3,
		UserID:      2,
		Username:    "Petrov Petr",
		CategoryID:  2,
		Name:        "name 3",
		ItemName:    "item name 3",
		Description: "description 3",
		Photos:      []string{},
		Price:       500,
		Quantity:    1,
		Date:        simpletime.SimpleTime(time.Date(2020, 12, 16, 0, 0, 0, 0, time.UTC)),
		Finished:    false,
		Tags: []*models.TagInRequest{
			{
				ID:   2,
				Name: "Крансый",
			},
			{
				ID:   3,
				Name: "Зелёный",
			},
		},
	},
}

func (rr *RequestsRepository) Add(ctx context.Context, request *models.Request) (int, error) {
	return 4, nil
}

func (rr *RequestsRepository) GetByID(ctx context.Context, id int) (*models.Request, error) {
	if id > 0 && id <= len(mockRequests) {
		return mockRequests[id-1], nil
	} else {
		return nil, models.ErrNoRecord
	}
}

func (rr *RequestsRepository) GetByUserID(ctx context.Context, userID int, search string) ([]*models.Request, error) {
	requests := make([]*models.Request, 0)
	for _, request := range mockRequests {
		if request.UserID == userID &&
			(search == "" || strings.Contains(request.Name, search) || strings.Contains(request.ItemName, search) || strings.Contains(request.Description, search)) {
			requests = append(requests, request)
		}
	}
	return requests, nil
}

func (rr *RequestsRepository) Search(ctx context.Context, filters *models.RequestFilters) ([]*models.Request, error) {
	requests := make([]*models.Request, 0)
	sort.Slice(mockRequests, func(i, j int) bool {
		if filters.SortColumn == "date" {
			if filters.SortDirection == "asc" {
				return time.Time(mockRequests[i].Date).Unix() < time.Time(mockRequests[j].Date).Unix()
			} else {
				return time.Time(mockRequests[i].Date).Unix() >= time.Time(mockRequests[j].Date).Unix()
			}
		} else {
			if filters.SortDirection == "asc" {
				return mockRequests[i].Price < mockRequests[j].Price
			} else {
				return mockRequests[i].Price < mockRequests[j].Price
			}
		}
	})

	for i, request := range mockRequests {
		if i >= filters.Page*filters.Size && i < (filters.Page+1)*filters.Size &&
			(filters.CategoryID == nil || request.CategoryID == *filters.CategoryID) &&
			subslice(filters.Tags, toIDs(request.Tags)) &&
			(filters.PriceFrom == nil || request.Price >= *filters.PriceFrom) &&
			(filters.PriceTo == nil || request.Price <= *filters.PriceTo) &&
			(filters.Search == "" || strings.Contains(request.Name, filters.Search) || strings.Contains(request.ItemName, filters.Search) || strings.Contains(request.Description, filters.Search)) {
			requests = append(requests, request)
		}
	}
	return requests, nil
}

func (rr *RequestsRepository) GetPricesLimits(ctx context.Context) (int, int, error) {
	return 100, 1000, nil
}

func (rr *RequestsRepository) Update(ctx context.Context, request *models.Request) error {
	return nil
}

func (rr *RequestsRepository) Delete(ctx context.Context, id int) error {
	return nil
}
