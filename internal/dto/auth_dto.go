package dto

// TokenClaims holds the minimal JWT claims we care about for access control.
type TokenClaims struct {
	Sub string `json:"sub"`
}

// StaffAccessRecord is the auth/authorization data loaded from Django-owned tables.
type StaffAccessRecord struct {
	StaffProfileID int64  `db:"staff_profile_id"`
	KeycloakSub    string `db:"keycloak_sub"`
	Email          string `db:"email"`
	Username       string `db:"username"`
	FullName       string `db:"full_name"`
	GlobalRole     string `db:"global_role"`
	IsActive       bool   `db:"is_active"`
}

// AccessContext is what we attach to the request after token + DB resolution.
type AccessContext struct {
	StaffProfileID int64
	KeycloakSub    string
	Email          string
	Username       string
	FullName       string

	GlobalRole  string
	IsExecutive bool

	BranchIDs       []int64
	IsBranchManager bool
}
