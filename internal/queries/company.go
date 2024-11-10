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
GROUP BY 
    cd.id, 
    cd.uuid,
    cd.user_id,
    cd.role_id,
    cd.company_name,
    cd.first_name,
    cd.last_name,
    cd.patronymic_name,
    cd.phone,
    cd.phone2,
    cd.phone3,
    cd.email,
    cd.email2,
    cd.email3,
    cd.meta,
    cd.meta2,
    cd.meta3,
    cd.address,
    cd.country,
    cd.country_id,
    cd.city_id,
    cd.image_url,
    cd.entity,
    cd.featured,
    cd.rating,
    cd.partner,
    cd.successful_ops,
    cd.created_at,
    cd.updated_at,
    cd.active,
    cd.deleted,
    cd.total_count;
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

var CreateCompany = `
INSERT INTO tbl_company (
    user_id, role_id, company_name, first_name, last_name, patronymic_name, phone, phone2, phone3, email, email2, email3, 
    meta, meta2, meta3, address, country, country_id, city_id, image_url, entity
)
VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21
)
RETURNING id;
`

var UpdateUserCompany = `
UPDATE tbl_user SET company_id = $1 WHERE id = $2;
`

var UpdateCompany = `
UPDATE tbl_company
SET
company_name = COALESCE($2, company_name),
first_name = COALESCE($3, first_name),
last_name = COALESCE($4, last_name),
patronymic_name = COALESCE($5, patronymic_name),
phone = COALESCE($6, phone),
phone2 = COALESCE($7, phone2),
phone3 = COALESCE($8, phone3),
email = COALESCE($9, email),
email2 = COALESCE($10, email2),
email3 = COALESCE($11, email3),
meta = COALESCE($12, meta),
meta2 = COALESCE($13, meta2),
meta3 = COALESCE($14, meta3),
address = COALESCE($15, address),
country = COALESCE($16, country),
country_id = COALESCE($17, country_id),
city_id = COALESCE($18, city_id),
image_url = COALESCE($19, image_url),
entity = COALESCE($20, entity),
user_id = COALESCE($21, user_id),
role_id = COALESCE($22, role_id),
active = COALESCE($23, active),
deleted = COALESCE($24, deleted),
updated_at = NOW()
`

var DeleteCompany = `
UPDATE tbl_company
SET deleted = 1, updated_at = NOW()
WHERE id = $1;
`
