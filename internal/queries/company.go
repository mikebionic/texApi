package queries

const GetCompanyWithRelations = `
WITH company_data AS (
    SELECT 
        c.*,
        COUNT(*) OVER() as total_count
    FROM tbl_company c
    WHERE c.deleted = 0
    ORDER BY c.id
    LIMIT $1 OFFSET $2
)
SELECT 
    cd.*,
    json_agg(DISTINCT d.*) FILTER (WHERE d.id IS NOT NULL) as drivers,
    json_agg(DISTINCT v.*) FILTER (WHERE v.id IS NOT NULL) as vehicles
FROM company_data cd
LEFT JOIN tbl_driver d ON cd.id = d.company_id AND d.deleted = 0
LEFT JOIN tbl_vehicle v ON cd.id = v.company_id AND v.deleted = 0
GROUP BY cd.id, cd.total_count;
`

const GetCompanyByID = `
SELECT 
    c.*,
    json_agg(DISTINCT d.*) FILTER (WHERE d.id IS NOT NULL) as drivers,
    json_agg(DISTINCT v.*) FILTER (WHERE v.id IS NOT NULL) as vehicles
FROM tbl_company c
LEFT JOIN tbl_driver d ON c.id = d.company_id AND d.deleted = 0
LEFT JOIN tbl_vehicle v ON c.id = v.company_id AND v.deleted = 0
WHERE c.id = $1 AND c.deleted = 0
GROUP BY c.id;
`

//
//var GetCompany = `
//SELECT
//	id,
//  user_id,
//  name,
//  address,
//  phone,
//  email,
//  logo_url,
//  created_at::varchar,
//  updated_at::varchar,
//  active,
//  deleted
//FROM tbl_company WHERE id = $1 AND deleted = 0
//`
//
//var CreateCompany = `
//INSERT INTO tbl_company (user_id, name, address, phone, email, logo_url)
//VALUES ($1, $2, $3, $4, $5, $6)
//RETURNING id;
//`
//
//var UpdateCompany = `
//UPDATE tbl_company
//SET name = COALESCE($2, name),
//address = COALESCE($3, address),
//phone = COALESCE($4, phone),
//email = COALESCE($5, email),
//logo_url = COALESCE($6, logo_url),
//updated_at = NOW()
//WHERE id = $1 AND deleted = 0
//RETURNING id;
//`

var DeleteCompany = `
UPDATE tbl_company
SET deleted = 1, updated_at = NOW()
WHERE id = $1;
`
