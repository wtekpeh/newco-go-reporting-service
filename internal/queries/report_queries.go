package queries

const TotalUsersQuery = `
	SELECT COUNT(*) 
	FROM accounts_staffprofile
`

const TotalActiveUsersQuery = `
	SELECT COUNT(*)
	FROM accounts_staffprofile
	WHERE is_active = TRUE
`
const TotalBranchesQuery = `
	SELECT COUNT(*)
	FROM accounts_branch
`
