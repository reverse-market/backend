package mock

import (
	"context"
	"github.com/reverse-market/backend/pkg/database/models"
	"github.com/reverse-market/backend/pkg/simpletime"
	"time"
)

type ProposalsRepository struct {
}

var (
	boughtById    = 3
	mockProposals = []*models.Proposal{
		{
			ID:          1,
			UserID:      1,
			Username:    "Petr Petrovich",
			RequestID:   1,
			Name:        "name 1",
			ItemName:    "item name 1",
			Description: "description 1",
			Photos:      []string{},
			Price:       100,
			Quantity:    1,
			Date:        simpletime.SimpleTime(time.Date(2020, 12, 16, 0, 0, 0, 0, time.UTC)),
			BoughtById:  nil,
		},
		{
			ID:          2,
			UserID:      2,
			Username:    "Ivanov Ivan",
			RequestID:   1,
			Name:        "name 2",
			ItemName:    "item name 2",
			Description: "description 2",
			Photos:      []string{},
			Price:       1000,
			Quantity:    1,
			Date:        simpletime.SimpleTime(time.Date(2020, 12, 16, 0, 0, 0, 0, time.UTC)),
			BoughtById:  nil,
		},
		{
			ID:          3,
			UserID:      1,
			Username:    "Petr Petrovich",
			RequestID:   2,
			Name:        "name 3",
			ItemName:    "item name 3",
			Description: "description 3",
			Photos:      []string{},
			Price:       100,
			Quantity:    1,
			Date:        simpletime.SimpleTime(time.Date(2020, 12, 16, 0, 0, 0, 0, time.UTC)),
			BoughtById:  &boughtById,
		},
	}
)

func (pr *ProposalsRepository) Add(ctx context.Context, proposals *models.Proposal) (int, error) {
	return 4, nil
}

func (pr *ProposalsRepository) GetByID(ctx context.Context, id int) (*models.Proposal, error) {
	if id > 0 && id <= len(mockProposals) {
		return mockProposals[id-1], nil
	} else {
		return nil, models.ErrNoRecord
	}
}

func (pr *ProposalsRepository) GetByUserIDSold(ctx context.Context, userID int) ([]*models.Proposal, error) {
	if userID == 1 {
		return []*models.Proposal{mockProposals[2]}, nil
	} else {
		return []*models.Proposal{}, nil
	}
}

func (pr *ProposalsRepository) GetByUserIDBought(ctx context.Context, userID int) ([]*models.Proposal, error) {
	if userID == 3 {
		return []*models.Proposal{mockProposals[2]}, nil
	} else {
		return []*models.Proposal{}, nil
	}
}

func (pr *ProposalsRepository) Update(ctx context.Context, proposals *models.Proposal) error {
	return nil
}

func (pr *ProposalsRepository) Delete(ctx context.Context, id int) error {
	return nil
}

func (pr *ProposalsRepository) GetBestForRequest(ctx context.Context, requestID int) (int, error) {
	switch requestID {
	case 1:
		return 1, nil
	case 2:
		return 3, nil
	default:
		return 0, models.ErrNoRecord
	}
}

func (pr *ProposalsRepository) Buy(ctx context.Context, id int, boughtByID int) error {
	return nil
}
