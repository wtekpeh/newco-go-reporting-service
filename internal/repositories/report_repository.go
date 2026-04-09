package repositories

import (
	"context"
	"newco-go-reporting-service/internal/queries"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ReportRepository struct {
	DB *pgxpool.Pool
}

func NewReportRepository(db *pgxpool.Pool) *ReportRepository {
	return &ReportRepository{
		DB: db,
	}
}

func (r *ReportRepository) TotalUsers() (int, error) {
	var count int

	query := queries.TotalUsersQuery

	err := r.DB.QueryRow(context.Background(), query).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *ReportRepository) TotalActiveUsers() (int, error) {
	var count int

	query := queries.TotalActiveUsersQuery

	err := r.DB.QueryRow(context.Background(), query).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
