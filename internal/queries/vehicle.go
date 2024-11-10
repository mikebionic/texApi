package queries

// Vehicle Brand queries
var GetVehicleBrand = `
SELECT id, name, country, founded_year, deleted
FROM tbl_vehicle_brand 
WHERE deleted = 0
`

var CreateVehicleBrand = `
INSERT INTO tbl_vehicle_brand (name, country, founded_year)
VALUES ($1, $2, $3)
RETURNING id;
`

var UpdateVehicleBrand = `
UPDATE tbl_vehicle_brand
SET name = COALESCE($2, name),
    country = COALESCE($3, country),
    founded_year = COALESCE($4, founded_year)
WHERE id = $1 AND deleted = 0
RETURNING id;
`

var DeleteVehicleBrand = `
UPDATE tbl_vehicle_brand
SET deleted = 1
WHERE id = $1;
`

// Vehicle Type queries
var GetVehicleType = `
SELECT id, type_name, description, deleted
FROM tbl_vehicle_type 
WHERE deleted = 0
`

var CreateVehicleType = `
INSERT INTO tbl_vehicle_type (type_name, description)
VALUES ($1, $2)
RETURNING id;
`

var UpdateVehicleType = `
UPDATE tbl_vehicle_type
SET type_name = COALESCE($2, type_name),
    description = COALESCE($3, description)
WHERE id = $1 AND deleted = 0
RETURNING id;
`

var DeleteVehicleType = `
UPDATE tbl_vehicle_type
SET deleted = 1
WHERE id = $1;
`

// Vehicle Model queries
var GetVehicleModel = `
SELECT m.id, m.name, m.year, m.brand, m.vehicle_type_id, 
       t.type_name as vehicle_type, m.feature, m.deleted
FROM tbl_vehicle_model m
LEFT JOIN tbl_vehicle_type t ON t.id = m.vehicle_type_id
WHERE m.deleted = 0
`

var CreateVehicleModel = `
INSERT INTO tbl_vehicle_model (name, year, brand, vehicle_type_id, feature)
VALUES ($1, $2, $3, $4, $5)
RETURNING id;
`

var UpdateVehicleModel = `
UPDATE tbl_vehicle_model
SET name = COALESCE($2, name),
    year = COALESCE($3, year),
    brand = COALESCE($4, brand),
    vehicle_type_id = COALESCE($5, vehicle_type_id),
    feature = COALESCE($6, feature)
WHERE id = $1 AND deleted = 0
RETURNING id;
`

var DeleteVehicleModel = `
UPDATE tbl_vehicle_model
SET deleted = 1
WHERE id = $1;
`

const GetVehicleList = `
WITH vehicle_data AS (
    SELECT 
        v.*,
        COUNT(*) OVER() as total_count
    FROM tbl_vehicle v
    WHERE v.deleted = 0
    ORDER BY v.id
    LIMIT $1 OFFSET $2
)
SELECT 
    vd.id, vd.uuid, vd.company_id, vd.vehicle_type,
    vd.vehicle_brand_id, vd.vehicle_model_id, vd.year_of_issue,
    vd.mileage, vd.numberplate, vd.trailer_numberplate,
    vd.gps_active, vd.photo1_url, vd.photo2_url,
    vd.photo3_url, vd.docs1_url, vd.docs2_url,
    vd.docs3_url, vd.view_count, vd.created_at,
    vd.updated_at, vd.active, vd.deleted, vd.total_count,
    json_build_object(
        'id', c.id,
        'company_name', c.company_name,
        'country', c.country
    ) AS company,
    json_build_object(
        'id', vb.id,
        'name', vb.name,
        'country', vb.country,
        'founded_year', vb.founded_year
    ) AS brand,
    json_build_object(
        'id', vm.id,
        'name', vm.name,
        'year', vm.year,
        'vehicle_type', t.type_name
    ) AS model
FROM vehicle_data vd
LEFT JOIN tbl_company c ON vd.company_id = c.id
LEFT JOIN tbl_vehicle_brand vb ON vd.vehicle_brand_id = vb.id
LEFT JOIN tbl_vehicle_model vm ON vd.vehicle_model_id = vm.id
LEFT JOIN tbl_vehicle_type t ON t.id = vm.vehicle_type_id  -- Adding the join to vehicle type table
GROUP BY 
    vd.id, vd.uuid, vd.company_id, vd.vehicle_type,
    vd.vehicle_brand_id, vd.vehicle_model_id, vd.year_of_issue,
    vd.mileage, vd.numberplate, vd.trailer_numberplate,
    vd.gps_active, vd.photo1_url, vd.photo2_url,
    vd.photo3_url, vd.docs1_url, vd.docs2_url,
    vd.docs3_url, vd.view_count, vd.created_at,
    vd.updated_at, vd.active, vd.deleted, vd.total_count,
    c.id, c.company_name, c.country,
    vb.id, vb.name, vb.country, vb.founded_year,
    vm.id, vm.name, vm.year, t.type_name;

`

const CreateVehicle = `
INSERT INTO tbl_vehicle (
    company_id, vehicle_type, vehicle_brand_id, vehicle_model_id,
    year_of_issue, mileage, numberplate, trailer_numberplate,
    gps_active, photo1_url, photo2_url, photo3_url,
    docs1_url, docs2_url, docs3_url, created_at, updated_at,
    active, deleted, uuid
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12,
    $13, $14, $15, NOW(), NOW(), 1, 0, gen_random_uuid()
)
RETURNING id;
`

const UpdateVehicle = `
UPDATE tbl_vehicle
SET 
    vehicle_type = COALESCE($2, vehicle_type),
    vehicle_brand_id = COALESCE($3, vehicle_brand_id),
    vehicle_model_id = COALESCE($4, vehicle_model_id),
    year_of_issue = COALESCE($5, year_of_issue),
    mileage = COALESCE($6, mileage),
    numberplate = COALESCE($7, numberplate),
    trailer_numberplate = COALESCE($8, trailer_numberplate),
    gps_active = COALESCE($9, gps_active),
    photo1_url = COALESCE($10, photo1_url),
    photo2_url = COALESCE($11, photo2_url),
    photo3_url = COALESCE($12, photo3_url),
    docs1_url = COALESCE($13, docs1_url),
    docs2_url = COALESCE($14, docs2_url),
    docs3_url = COALESCE($15, docs3_url),
    active = COALESCE($16, active),
    company_id = COALESCE($17, company_id),
    deleted = COALESCE($18, deleted),
    updated_at = NOW()
`

const DeleteVehicle = `
UPDATE tbl_vehicle
SET deleted = 1, updated_at = NOW(), active = 0
WHERE id = $1;
`

const GetVehicleByID = `
SELECT 
    v.*,
    json_build_object(
        'id', c.id,
        'company_name', c.company_name,
        'country', c.country
    ) as company,
    json_build_object(
        'id', vb.id,
        'name', vb.name
    ) as brand,
    json_build_object(
        'id', vm.id,
        'name', vm.name
    ) as model
FROM tbl_vehicle v
LEFT JOIN tbl_company c ON v.company_id = c.id
LEFT JOIN tbl_vehicle_brand vb ON v.vehicle_brand_id = vb.id
LEFT JOIN tbl_vehicle_model vm ON v.vehicle_model_id = vm.id
WHERE v.id = $1 AND v.deleted = 0;
`
