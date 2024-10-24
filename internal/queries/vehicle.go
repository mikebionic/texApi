package queries

var GetVehicles = `
SELECT * FROM tbl_vehicle WHERE deleted = 0 AND ($1::INT IS NULL OR company_id = $1);
`

var GetVehicle = `
SELECT * FROM tbl_vehicle WHERE id = $1 AND deleted = 0;
`

var CreateVehicle = `
INSERT INTO tbl_vehicle (company_id, vehivle_type, brand, vehicle_model, year_of_isse, numberplate, trailer_numberplate, gps_active, photo1_url, photo2_url, photo3_url, docs1_url, docs2_url, docs3_url)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
RETURNING id;
`

var UpdateVehicle = `
UPDATE tbl_vehicle
SET vehivle_type = COALESCE($2, vehivle_type),
brand = COALESCE($3, brand),
vehicle_model = COALESCE($4, vehicle_model),
year_of_isse = COALESCE($5, year_of_isse),
numberplate = COALESCE($6, numberplate),
trailer_numberplate = COALESCE($7, trailer_numberplate),
gps_active = COALESCE($8, gps_active),
photo1_url = COALESCE($9, photo1_url),
photo2_url = COALESCE($10, photo2_url),
photo3_url = COALESCE($11, photo3_url),
docs1_url = COALESCE($12, docs1_url),
docs2_url = COALESCE($13, docs2_url),
docs3_url = COALESCE($14, docs3_url),
updated_at = NOW()
WHERE id = $1 AND deleted = 0
RETURNING id;
`

var DeleteVehicle = `
UPDATE tbl_vehicle
SET deleted = 1, updated_at = NOW()
WHERE id = $1;
`
