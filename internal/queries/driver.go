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
	image_url,
	created_at::varchar,
	updated_at::varchar,
	active,
	deleted
FROM tbl_driver WHERE deleted = 0
`

//var CreateDriver = `
//INSERT INTO tbl_driver (
//   company_id,
//	first_name,
//	last_name,
//	patronymic_name,
//	phone,
//	email,
//	avatar_url)
//VALUES ($1, $2, $3, $4, $5, $6, $7)
//RETURNING id;
//`
//
//var UpdateDriver = `
//UPDATE tbl_driver
//SET first_name = COALESCE($2, first_name),
//last_name = COALESCE($3, last_name),
//patronymic_name = COALESCE($4, patronymic_name),
//phone = COALESCE($5, phone),
//email = COALESCE($6, email),
//avatar_url = COALESCE($7, avatar_url),
//updated_at = NOW()
//WHERE id = $1 AND deleted = 0
//RETURNING id;
//`
//
//var DeleteDriver = `
//UPDATE tbl_driver
//SET deleted = 1, updated_at = NOW()
//WHERE id = $1;
//`

const GetDriverList = `
WITH driver_data AS (
    SELECT 
        d.*,
        COUNT(*) OVER() as total_count
    FROM tbl_driver d
    WHERE d.deleted = 0
    ORDER BY d.id
    LIMIT $1 OFFSET $2
)
SELECT 
    dd.*,
    json_build_object(
        'id', c.id,
        'company_name', c.company_name,
        'country', c.country
    ) as company,
    COALESCE(
        (
            SELECT json_agg(
                json_build_object(
                    'id', v.id,
                    'vehicle_type', v.vehicle_type,
                    'numberplate', v.numberplate
                )
            )
            FROM tbl_vehicle v
            WHERE v.company_id = dd.company_id AND v.deleted = 0
        ),
        '[]'
    ) as assigned_vehicles
FROM driver_data dd
LEFT JOIN tbl_company c ON dd.company_id = c.id
GROUP BY dd.id, dd.total_count, c.id, c.company_name, c.country;
`

const GetDriverByID = `
SELECT 
    d.*,
    json_build_object(
        'id', c.id,
        'company_name', c.company_name,
        'country', c.country
    ) as company,
    COALESCE(
        (
            SELECT json_agg(
                json_build_object(
                    'id', v.id,
                    'vehicle_type', v.vehicle_type,
                    'numberplate', v.numberplate
                )
            )
            FROM tbl_vehicle v
            WHERE v.company_id = d.company_id AND v.deleted = 0
        ),
        '[]'
    ) as assigned_vehicles
FROM tbl_driver d
LEFT JOIN tbl_company c ON d.company_id = c.id
WHERE d.id = $1 AND d.deleted = 0;
`

const CreateDriver = `
INSERT INTO tbl_driver (
    company_id, first_name, last_name, patronymic_name,
    phone, email, image_url
) VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id;
`

const UpdateDriver = `
UPDATE tbl_driver
SET 
    first_name = COALESCE($2, first_name),
    last_name = COALESCE($3, last_name),
    patronymic_name = COALESCE($4, patronymic_name),
    phone = COALESCE($5, phone),
    email = COALESCE($6, email),
    featured = COALESCE($7, featured),
    rating = COALESCE($8, rating),
    partner = COALESCE($9, partner),
    image_url = COALESCE($10, image_url),
    active = COALESCE($11, active),
    updated_at = NOW()
WHERE id = $1 AND deleted = 0
RETURNING id;
`

const DeleteDriver = `
UPDATE tbl_driver
SET deleted = 1, updated_at = NOW()
WHERE id = $1;
`
