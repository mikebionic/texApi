package queries

var GetCompany = `
SELECT 
	id,
   user_id,
   name,
   address,
   phone,
   email,
   logo_url,
   created_at::varchar,
   updated_at::varchar,
   active,
   deleted
FROM tbl_company WHERE id = $1 AND deleted = 0
`

var CreateCompany = `
INSERT INTO tbl_company (user_id, name, address, phone, email, logo_url)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id;
`

var UpdateCompany = `
UPDATE tbl_company
SET name = COALESCE($2, name),
address = COALESCE($3, address),
phone = COALESCE($4, phone),
email = COALESCE($5, email),
logo_url = COALESCE($6, logo_url),
updated_at = NOW()
WHERE id = $1 AND deleted = 0
RETURNING id;
`

var DeleteCompany = `
UPDATE tbl_company
SET deleted = 1, updated_at = NOW()
WHERE id = $1;
`
