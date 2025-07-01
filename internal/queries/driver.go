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
					'company_id',v.company_id,
					'vehicle_type_id',v.vehicle_type_id,
					'vehicle_brand_id',v.vehicle_brand_id,
					'vehicle_model_id',v.vehicle_model_id,
					'year_of_issue',v.year_of_issue,
					'mileage',v.mileage,
					'numberplate',v.numberplate,
					'trailer_numberplate',v.trailer_numberplate,
					'gps',v.gps,
					'photo1_url',v.photo1_url,
					'photo2_url',v.photo2_url,
					'photo3_url',v.photo3_url,
					'docs1_url',v.docs1_url,
					'docs2_url',v.docs2_url,
					'docs3_url',v.docs3_url,
					'view_count',v.view_count,
					'meta',v.meta,
					'meta2',v.meta2,
					'meta3',v.meta3,
					'available',v.available
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

const GetFilteredDriverList = `
WITH driver_data AS (
    SELECT 
        d.*, 
        COUNT(*) OVER() as total_count
    FROM tbl_driver d
    WHERE :whereClause
    ORDER BY :orderBy
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
GROUP BY 
    dd.id, dd.uuid, dd.company_id, dd.first_name, dd.last_name, dd.patronymic_name, 
    dd.phone, dd.email, dd.featured, dd.rating, dd.partner, dd.successful_ops, 
    dd.image_url, dd.meta, dd.meta2, dd.meta3, dd.created_at, dd.updated_at, dd.active, dd.deleted, dd.total_count, 
    c.id, c.company_name, c.country;
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
    block_reason = COALESCE($12, block_reason),
    active = COALESCE($13, active),
    deleted = COALESCE($14, deleted),
    updated_at = NOW()
`

const DeleteDriver = `
UPDATE tbl_driver
SET deleted = 1, updated_at = NOW()
WHERE id = $1`

const CreateDriverUser = `
INSERT INTO tbl_user (
    username, password, email, phone, role, role_id, 
    verified, active, deleted, driver_id
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING id;
`

const UpdateDriverUser = `
UPDATE tbl_user
SET 
    email = COALESCE($1, email),
    phone = COALESCE($2, phone),
    active = COALESCE($3, active),
    deleted = COALESCE($4, deleted),
    updated_at = NOW()
WHERE driver_id = $5 AND deleted = 0;
`

const DeleteDriverUser = `
UPDATE tbl_user
SET deleted = 1, updated_at = NOW()
WHERE driver_id = $1;
`
