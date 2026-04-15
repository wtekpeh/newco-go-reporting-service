package queries

const GetNotificationsByRecipientID = `
SELECT
    n.id,
    n.event_id,
    e.action,
    e.target_type,
    e.target_id,
    e.message,
    n.is_read,
    n.read_at,
    n.created_at,
    e.actor_staff_profile_id,
    COALESCE(actor.full_name, '') AS actor_full_name,
    COALESCE(actor.email, '') AS actor_email,
    e.branch_id,
    b.name AS branch_name
FROM activity_notification n
INNER JOIN activity_activityevent e
    ON e.id = n.event_id
INNER JOIN accounts_staffprofile actor
    ON actor.id = e.actor_staff_profile_id
LEFT JOIN accounts_branch b
    ON b.id = e.branch_id
WHERE n.recipient_staff_profile_id = $1
ORDER BY n.created_at DESC
LIMIT $2
`

const GetUnreadNotificationCountByRecipientID = `
SELECT COUNT(*)
FROM activity_notification
WHERE recipient_staff_profile_id = $1
  AND is_read = FALSE
`

const MarkNotificationAsRead = `
UPDATE activity_notification
SET
    is_read = TRUE,
    read_at = NOW()
WHERE id = $1
  AND recipient_staff_profile_id = $2
`
