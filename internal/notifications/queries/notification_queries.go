package queries

const GetExecutiveRecipientIDs = `
SELECT id
FROM accounts_staffprofile
WHERE is_active = TRUE
  AND global_role IN ('boss', 'managing_director')
ORDER BY id
`

const GetBranchManagerRecipientIDsByBranchID = `
SELECT DISTINCT bra.staff_profile_id
FROM accounts_branchroleassignment bra
INNER JOIN accounts_staffprofile sp
    ON sp.id = bra.staff_profile_id
INNER JOIN accounts_branch b
    ON b.id = bra.branch_id
WHERE bra.branch_id = $1
  AND bra.role = 'branch_manager'
  AND bra.is_active = TRUE
  AND sp.is_active = TRUE
  AND b.is_active = TRUE
ORDER BY bra.staff_profile_id
`

const InsertNotification = `
INSERT INTO activity_notification (
    event_id,
    recipient_staff_profile_id,
    is_read,
    read_at,
    emailed_at,
    created_at
)
VALUES ($1, $2, FALSE, NULL, NULL, NOW())
ON CONFLICT (event_id, recipient_staff_profile_id)
DO UPDATE SET recipient_staff_profile_id = activity_notification.recipient_staff_profile_id
RETURNING id
`
