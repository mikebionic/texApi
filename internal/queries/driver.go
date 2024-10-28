package queries

var GetDriver = `
SELECT 
    id,
	company_id,
	first_name,
	last_name,
	patronymic_name,
	phone,
	email,
	avatar_url,
	created_at::varchar,
	updated_at::varchar,
	active,
	deleted
FROM tbl_driver WHERE deleted = 0
`

var CreateDriver = `
INSERT INTO tbl_driver (
    company_id,
	first_name,
	last_name,
	patronymic_name,
	phone,
	email,
	avatar_url)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id;
`

var UpdateDriver = `
UPDATE tbl_driver
SET first_name = COALESCE($2, first_name),
last_name = COALESCE($3, last_name),
patronymic_name = COALESCE($4, patronymic_name),
phone = COALESCE($5, phone),
email = COALESCE($6, email),
avatar_url = COALESCE($7, avatar_url),
updated_at = NOW()
WHERE id = $1 AND deleted = 0
RETURNING id;
`

var DeleteDriver = `
UPDATE tbl_driver
SET deleted = 1, updated_at = NOW()
WHERE id = $1;
`
