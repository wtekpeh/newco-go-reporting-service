package repositories

import (
	"context"
	"errors"

	"newco-go-reporting-service/internal/dto"

	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrStaffProfileNotFound = errors.New("active staff profile not found")

type AccessRepository interface {
	GetActiveStaffByKeycloakSub(ctx context.Context, keycloakSub string) (*dto.StaffAccessRecord, error)
}

type accessRepository struct {
	pool *pgxpool.Pool
}

func NewAccessRepository(pool *pgxpool.Pool) AccessRepository {
	return &accessRepository{pool: pool}
}

func (r *accessRepository) GetActiveStaffByKeycloakSub(ctx context.Context, keycloakSub string) (*dto.StaffAccessRecord, error) {
	const query = `
		SELECT
			id AS staff_profile_id,
			keycloak_sub,
			email,
			username,
			full_name,
			global_role,
			is_active
		FROM accounts_staffprofile
		WHERE keycloak_sub = $1
		  AND is_active = TRUE
		LIMIT 1
	`

	var record dto.StaffAccessRecord

	err := r.pool.QueryRow(ctx, query, keycloakSub).Scan(
		&record.StaffProfileID,
		&record.KeycloakSub,
		&record.Email,
		&record.Username,
		&record.FullName,
		&record.GlobalRole,
		&record.IsActive,
	)
	if err != nil {
		return nil, ErrStaffProfileNotFound
	}

	return &record, nil
}
