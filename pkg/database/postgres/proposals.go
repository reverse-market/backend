package postgres

import (
	"context"
	"errors"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/reverse-market/backend/pkg/database/models"
)

type ProposalsRepository struct {
	DB *pgxpool.Pool
}

func (pr *ProposalsRepository) Add(ctx context.Context, proposals *models.Proposal) (int, error) {
	conn, err := pr.DB.Acquire(ctx)
	if err != nil {
		return 0, err
	}
	defer conn.Release()

	stmt := "INSERT INTO proposals (user_id, request_id, description, photos, price, quantity, date) " +
		"VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id"

	var id int
	if err := conn.QueryRow(ctx, stmt,
		proposals.UserID,
		proposals.RequestID,
		proposals.Description,
		proposals.Photos,
		proposals.Price,
		proposals.Quantity,
		proposals.Date,
	).Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (pr *ProposalsRepository) GetByID(ctx context.Context, id int) (*models.Proposal, error) {
	conn, err := pr.DB.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	stmt := "SELECT * FROM proposals_view WHERE id=$1"

	proposal := &models.Proposal{}
	if err := pgxscan.Get(ctx, conn, proposal, stmt, id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}

	return proposal, nil
}

func (pr *ProposalsRepository) GetByUserIDSold(ctx context.Context, userID int) ([]*models.Proposal, error) {
	conn, err := pr.DB.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	stmt := "SELECT * FROM proposals_view WHERE user_id=$1 AND bought_by_id IS NOT NULL"

	proposals := make([]*models.Proposal, 0)
	if err := pgxscan.Select(ctx, conn, &proposals, stmt, userID); err != nil {
		return nil, err
	}

	return proposals, nil
}

func (pr *ProposalsRepository) GetByUserIDBought(ctx context.Context, userID int) ([]*models.Proposal, error) {
	conn, err := pr.DB.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	stmt := "SELECT * FROM proposals_view WHERE bought_by_id = $1"

	proposals := make([]*models.Proposal, 0)
	if err := pgxscan.Select(ctx, conn, &proposals, stmt, userID); err != nil {
		return nil, err
	}

	return proposals, nil
}

func (pr *ProposalsRepository) Update(ctx context.Context, proposals *models.Proposal) error {
	conn, err := pr.DB.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	stmt := "INSERT INTO proposals (user_id=$1, request_id=$2, description=$3, photos=$4," +
		" price=$5, quantity=$6, date=$7) WHERE id=$8"

	if _, err := conn.Exec(ctx, stmt,
		proposals.UserID,
		proposals.RequestID,
		proposals.Description,
		proposals.Photos,
		proposals.Price,
		proposals.Quantity,
		proposals.Date,
		proposals.Date,
	); err != nil {
		return err
	}

	return nil
}

func (pr *ProposalsRepository) Delete(ctx context.Context, id int) error {
	conn, err := pr.DB.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	stmt := "DELETE FROM proposals WHERE id=$1"

	if _, err := conn.Exec(ctx, stmt, id); err != nil {
		return err
	}

	return nil
}

func (pr *ProposalsRepository) GetBestForRequest(ctx context.Context, requestID int) (int, error) {
	conn, err := pr.DB.Acquire(ctx)
	if err != nil {
		return 0, err
	}
	defer conn.Release()

	stmt := "SELECT id FROM proposals WHERE request_id=$1 ORDER BY price LIMIT 1"

	var id int
	if err := conn.QueryRow(ctx, stmt, requestID).Scan(&id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, models.ErrNoRecord
		}
		return 0, err
	}

	return id, nil
}

func (pr *ProposalsRepository) Buy(ctx context.Context, id int, boughtByID int) error {
	conn, err := pr.DB.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	proposal, err := pr.GetByID(ctx, id)
	if err != nil {
		return err
	}

	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	stmt := "UPDATE proposals SET bought_by_id=$1 WHERE id=$2"
	if _, err := tx.Exec(ctx, stmt, boughtByID, id); err != nil {
		return err
	}

	stmt = "DELETE FROM proposals WHERE request_id=$1 AND id <> $2"
	if _, err := tx.Exec(ctx, stmt, proposal.RequestID, id); err != nil {
		return err
	}

	stmt = "UPDATE requests SET finished=true WHERE id=$1"
	if _, err := tx.Exec(ctx, stmt, proposal.RequestID); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
