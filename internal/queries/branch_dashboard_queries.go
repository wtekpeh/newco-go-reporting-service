package queries

const BranchDashboardSummaryQuery = `
	SELECT
		$1::bigint AS branch_id,

		COALESCE((
			SELECT COUNT(DISTINCT bra.staff_profile_id)
			FROM accounts_branchroleassignment bra
			INNER JOIN accounts_branch b
				ON b.id = bra.branch_id
			WHERE bra.branch_id = $1
			  AND bra.is_active = TRUE
			  AND b.is_active = TRUE
		), 0) AS total_staff,

		COALESCE((
			SELECT COUNT(*)
			FROM cooking_cookbatch cb
			WHERE cb.branch_id = $1
		), 0) AS total_batches,

		COALESCE((
			SELECT COUNT(*)
			FROM cooking_cookbatch cb
			WHERE cb.branch_id = $1
			  AND cb.created_at >= date_trunc('week', CURRENT_DATE)
		), 0) AS batches_this_week,

		COALESCE((
			SELECT COUNT(*)
			FROM cooking_cookbatch cb
			WHERE cb.branch_id = $1
			  AND cb.created_at >= date_trunc('month', CURRENT_DATE)
		), 0) AS batches_this_month
	`

const BranchBatchTrendsQuery = `
	WITH grouped_batches AS (
		SELECT
			DATE_TRUNC('day', created_at) AS grouped_date
		FROM cooking_cookbatch
		WHERE branch_id = $1
	)
	SELECT
		TO_CHAR(grouped_date, 'YYYY-MM-DD') AS label,
		COUNT(*) AS count
	FROM grouped_batches
	GROUP BY grouped_date
	ORDER BY grouped_date ASC
`
const BranchRoleDistributionQuery = `
	SELECT role, COUNT(*) AS count
	FROM accounts_branchroleassignment
	WHERE is_active = TRUE
	  AND branch_id = $1
	GROUP BY role
	ORDER BY count DESC, role ASC
`
const BranchRecentBatchesQuery = `
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
	WHERE cb.branch_id = $1
	ORDER BY cb.created_at DESC
	LIMIT 10
`
