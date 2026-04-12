package repositories

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type BranchRoleRecord struct {
	BranchID int64
	Role     string
}

type BranchAccessRepository interface {
	GetActiveBranchRolesByStaffProfileID(ctx context.Context, staffProfileID int64) ([]BranchRoleRecord, error)
}

type branchAccessRepository struct {
	pool *pgxpool.Pool
}

func NewBranchAccessRepository(pool *pgxpool.Pool) BranchAccessRepository {
	return &branchAccessRepository{pool: pool}
}

func (r *branchAccessRepository) GetActiveBranchRolesByStaffProfileID(
	ctx context.Context,
	staffProfileID int64,
) ([]BranchRoleRecord, error) {
	const query = `
		SELECT
			bra.branch_id,
			bra.role
		FROM accounts_branchroleassignment bra
		INNER JOIN accounts_branch b
			ON b.id = bra.branch_id
		WHERE bra.staff_profile_id = $1
		  AND bra.is_active = TRUE
		  AND b.is_active = TRUE
		ORDER BY bra.branch_id, bra.role
	`

	rows, err := r.pool.Query(ctx, query, staffProfileID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []BranchRoleRecord

	for rows.Next() {
		var record BranchRoleRecord
		if err := rows.Scan(&record.BranchID, &record.Role); err != nil {
			return nil, err
		}
		records = append(records, record)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return records, nil
}
