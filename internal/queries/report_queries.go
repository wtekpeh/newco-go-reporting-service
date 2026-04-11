package queries

const TotalUsersQuery = `
	SELECT COUNT(DISTINCT sp.id)
	FROM accounts_staffprofile sp
	LEFT JOIN accounts_branchroleassignment bra ON bra.staff_profile_id = sp.id
	WHERE ($1 = '' OR bra.branch_id::text = $1)
`

const TotalActiveUsersQuery = `
	SELECT COUNT(DISTINCT sp.id)
	FROM accounts_staffprofile sp
	LEFT JOIN accounts_branchroleassignment bra ON bra.staff_profile_id = sp.id
	WHERE sp.is_active = TRUE
	AND ($1 = '' OR bra.branch_id::text = $1)
`
const TotalBranchesQuery = `
	SELECT COUNT(*)
	FROM accounts_branch
`
const TotalBatchesQuery = `
	SELECT COUNT(*)
	FROM cooking_cookbatch
	WHERE ($1 = '' OR branch_id::text = $1)
`
const BatchesThisWeekQuery = `
	SELECT COUNT(*)
	FROM cooking_cookbatch
	WHERE created_at >= date_trunc('week', CURRENT_TIMESTAMP)
	AND ($1 = '' OR branch_id::text = $1)
`
const BatchesThisMonthQuery = `
	SELECT COUNT(*)
	FROM cooking_cookbatch
	WHERE created_at >= date_trunc('month', CURRENT_TIMESTAMP)
	AND ($1 = '' OR branch_id::text = $1)
`
const MostActiveBranchQuery = `
	SELECT b.name, COUNT(cb.id) AS batch_count
	FROM cooking_cookbatch cb
	JOIN accounts_branch b ON b.id = cb.branch_id
	WHERE ($1 = '' OR cb.branch_id::text = $1)
	GROUP BY b.id, b.name
	ORDER BY batch_count DESC, b.name ASC
	LIMIT 1
`
const LargestBranchQuery = `
	SELECT b.name, COUNT(DISTINCT bra.staff_profile_id) AS staff_count
	FROM accounts_branch b
	LEFT JOIN accounts_branchroleassignment bra ON bra.branch_id = b.id
	WHERE ($1 = '' OR b.id::text = $1)
	GROUP BY b.id, b.name
	ORDER BY staff_count DESC, b.name ASC
	LIMIT 1
`
const MostUsedRecipeQuery = `
	SELECT r.name, COUNT(cb.id) AS batch_count
	FROM cooking_cookbatch cb
	JOIN recipes_recipe r ON r.id = cb.recipe_id
	WHERE ($1 = '' OR cb.branch_id::text = $1)
	GROUP BY r.id, r.name
	ORDER BY batch_count DESC, r.name ASC
	LIMIT 1
`
const AverageBatchesPerBranchQuery = `
	SELECT COALESCE(ROUND(COUNT(cb.id)::numeric / NULLIF(COUNT(DISTINCT b.id), 0), 2), 0)
	FROM accounts_branch b
	LEFT JOIN cooking_cookbatch cb ON cb.branch_id = b.id
	WHERE ($1 = '' OR b.id::text = $1)
`
const PeakBatchDayQuery = `
	SELECT 
		TRIM(TO_CHAR(created_at, 'Day')) AS day_name,
		COUNT(*) AS batch_count
	FROM cooking_cookbatch
	WHERE ($1 = '' OR branch_id::text = $1)
	GROUP BY day_name
	ORDER BY batch_count DESC, day_name ASC
	LIMIT 1
`
const RecentBatchesQuery = `
	SELECT
		cb.id,
		r.name AS recipe_name,
		b.name AS branch_name,
		COALESCE(sp.full_name, sp.email, sp.username, sp.keycloak_sub) AS created_by,
		cb.n_people,
		cb.status,
		cb.protein_type,
		cb.created_at
	FROM cooking_cookbatch cb
	JOIN recipes_recipe r ON r.id = cb.recipe_id
	JOIN accounts_branch b ON b.id = cb.branch_id
	JOIN accounts_staffprofile sp ON sp.id = cb.created_by_id
	WHERE ($1 = '' OR DATE(cb.created_at) >= $1::date)
	  AND ($2 = '' OR DATE(cb.created_at) <= $2::date)
	  AND ($3 = '' OR cb.branch_id::text = $3)
	ORDER BY cb.created_at DESC
	LIMIT 10
`
const BatchTrendsQuery = `
	WITH grouped_batches AS (
		SELECT
			CASE
				WHEN $4 = 'week' THEN DATE_TRUNC('week', created_at)
				WHEN $4 = 'month' THEN DATE_TRUNC('month', created_at)
				ELSE DATE_TRUNC('day', created_at)
			END AS grouped_date
		FROM cooking_cookbatch
		WHERE ($1 = '' OR DATE(created_at) >= $1::date)
		  AND ($2 = '' OR DATE(created_at) <= $2::date)
		  AND ($3 = '' OR branch_id::text = $3)
	)
	SELECT
		TO_CHAR(grouped_date, 'YYYY-MM-DD') AS label,
		COUNT(*) AS count
	FROM grouped_batches
	GROUP BY grouped_date
	ORDER BY grouped_date ASC
`
const BranchSummaryQuery = `
	SELECT
		b.id,
		b.name,
		COUNT(DISTINCT bra.staff_profile_id) AS staff_count,
		COUNT(cb.id) AS total_batches
	FROM accounts_branch b
	LEFT JOIN accounts_branchroleassignment bra ON bra.branch_id = b.id
	LEFT JOIN cooking_cookbatch cb ON cb.branch_id = b.id
	WHERE ($1 = '' OR b.id::text = $1)
	GROUP BY b.id, b.name
	ORDER BY total_batches DESC, b.name ASC
`
const RoleDistributionQuery = `
	SELECT role, COUNT(*) AS count
	FROM accounts_branchroleassignment
	WHERE is_active = TRUE
	GROUP BY role
	ORDER BY count DESC, role ASC
`
const UserGrowthQuery = `
	WITH grouped_users AS (
		SELECT
			CASE
				WHEN $4 = 'week' THEN DATE_TRUNC('week', created_at)
				WHEN $4 = 'month' THEN DATE_TRUNC('month', created_at)
				ELSE DATE_TRUNC('day', created_at)
			END AS grouped_date
		FROM accounts_staffprofile
		WHERE ($1 = '' OR DATE(created_at) >= $1::date)
		  AND ($2 = '' OR DATE(created_at) <= $2::date)
		  AND ($3 = '' OR id IN (
				SELECT bra.staff_profile_id
				FROM accounts_branchroleassignment bra
				WHERE bra.branch_id::text = $3
		  ))
	)
	SELECT
		TO_CHAR(grouped_date, 'YYYY-MM-DD') AS label,
		COUNT(*) AS count
	FROM grouped_users
	GROUP BY grouped_date
	ORDER BY grouped_date ASC
`
const GlobalRoleDistributionQuery = `
	SELECT global_role, COUNT(*) AS count
	FROM accounts_staffprofile
	WHERE global_role <> 'none'
	AND ($1 = '' OR id IN (
		SELECT bra.staff_profile_id
		FROM accounts_branchroleassignment bra
		WHERE bra.branch_id::text = $1
	))
	GROUP BY global_role
	ORDER BY count DESC, global_role ASC
`

const ActiveStatusDistributionQuery = `
	SELECT
		CASE
			WHEN is_active = TRUE THEN 'Active'
			ELSE 'Inactive'
		END AS status,
		COUNT(*) AS count
	FROM accounts_staffprofile
	WHERE ($1 = '' OR id IN (
		SELECT bra.staff_profile_id
		FROM accounts_branchroleassignment bra
		WHERE bra.branch_id::text = $1
	))
	GROUP BY is_active
	ORDER BY status ASC
`
const BatchStatusSummaryQuery = `
	SELECT
		status,
		COUNT(*) AS count
	FROM cooking_cookbatch
	WHERE ($1 = '' OR branch_id::text = $1)
	GROUP BY status
	ORDER BY count DESC, status ASC
`
const BranchTrendsFlatQuery = `
	SELECT
		b.id AS branch_id,
		b.name AS branch_name,
		TO_CHAR(DATE(cb.created_at), 'YYYY-MM-DD') AS label,
		COUNT(cb.id) AS count
	FROM cooking_cookbatch cb
	JOIN accounts_branch b ON b.id = cb.branch_id
	WHERE ($1 = '' OR cb.branch_id::text = $1)
	GROUP BY b.id, b.name, DATE(cb.created_at)
	ORDER BY b.name ASC, DATE(cb.created_at) ASC
`
