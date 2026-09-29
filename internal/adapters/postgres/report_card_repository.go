package postgres

import (
	"context"
	"database/sql"

	"github.com/alvarolucio2007/Scholarly/internal/domain"
	"github.com/alvarolucio2007/Scholarly/internal/ports"
)

var _ ports.ReportCardRepository = (*ReportCardRepo)(nil)

type ReportCardRepo struct {
	db *sql.DB
}

func NewReportCardRepo(db *sql.DB) *ReportCardRepo {
	return &ReportCardRepo{db: db}
}

func (r *ReportCardRepo) View(ctx context.Context) ([]*domain.ReportCard, error) {
	return nil, nil
}
