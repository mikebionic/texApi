package queries

const GetDriverList = `
WITH driver_data AS (
	SELECT 
		d.*, COUNT(*) OVER() as total_count
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
					'vehicle_brand_id', v.vehicle_brand_id,
					'numberplate', v.numberplate
				)
			)
			FROM tbl_vehicle v
			WHERE v.company_id = dd.company_id AND v.deleted = 0
		),
		'[]'
	) as assigned_vehicles
FROM driver_data dd
LEFT JOIN tbl_company c ON dd.company_id = c.id;
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
                    'vehicle_type_id', v.vehicle_type_id,
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
    phone, email, image_url, meta, meta2, meta3
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
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
    image_url = COALESCE($7, image_url),
    meta = COALESCE($8, meta),
    meta2 = COALESCE($9, meta2),
    meta3 = COALESCE($10, meta3),
    company_id = COALESCE($11, company_id),
    active = COALESCE($12, active),
    deleted = COALESCE($13, deleted),
    updated_at = NOW()
`

const DeleteDriver = `
UPDATE tbl_driver
SET deleted = 1, updated_at = NOW()
WHERE id = $1;
`
