package queries

const GetUnprocessedEvents = `
SELECT
    id,
    action,
    target_type,
    target_id,
    actor_staff_profile_id,
    branch_id,
    message,
    metadata_json,
    created_at,
    processed_at
FROM activity_activityevent
WHERE processed_at IS NULL
ORDER BY created_at ASC
LIMIT $1
`

const MarkEventProcessed = `
UPDATE activity_activityevent
SET processed_at = NOW()
WHERE id = $1
`
