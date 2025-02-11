package queries

const GetCargoList = `
    SELECT 
        c.*,
        COUNT(*) OVER() as total_count
    FROM tbl_cargo c
`

const GetCargoByID = `
SELECT 
    c.*
FROM tbl_cargo c
WHERE c.id = $1
`

const CreateCargo = `
INSERT INTO tbl_cargo (
    company_id, name, description, info, qty, weight, meta, meta2, meta3, 
    vehicle_type_id, packaging_type_id, gps, photo1_url, photo2_url, photo3_url, 
    docs1_url, docs2_url, docs3_url, note, weight_type
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)
RETURNING id;
`

const UpdateCargo = `
UPDATE tbl_cargo
SET 
    name = COALESCE($2, name),
    description = COALESCE($3, description),
    info = COALESCE($4, info),
    qty = COALESCE($5, qty),
    weight = COALESCE($6, weight),
    meta = COALESCE($7, meta),
    meta2 = COALESCE($8, meta2),
    meta3 = COALESCE($9, meta3),
    vehicle_type_id = COALESCE($10, vehicle_type_id),
    packaging_type_id = COALESCE($11, packaging_type_id),
    gps = COALESCE($12, gps),
    photo1_url = COALESCE($13, photo1_url),
    photo2_url = COALESCE($14, photo2_url),
    photo3_url = COALESCE($15, photo3_url),
    docs1_url = COALESCE($16, docs1_url),
    docs2_url = COALESCE($17, docs2_url),
    docs3_url = COALESCE($18, docs3_url),
    note = COALESCE($19, note),
    active = COALESCE($20, active),
    deleted = COALESCE($21, deleted),
    weight_type = COALESCE($22, weight_type),
    updated_at = NOW()
WHERE id = $1 
`

const DeleteCargo = `
UPDATE tbl_cargo
SET deleted = 1, updated_at = NOW()
WHERE id = $1
`

const GetDetailedCargoList = `
    SELECT 
        c.*,
        COUNT(*) OVER() as total_count,
        
        comp.id as "company.id",
        comp.uuid as "company.uuid",
        comp.user_id as "company.user_id",
        comp.role_id as "company.role_id",
        comp.company_name as "company.company_name",
        comp.first_name as "company.first_name",
        comp.last_name as "company.last_name",
        comp.patronymic_name as "company.patronymic_name",
        comp.phone as "company.phone",
        comp.phone2 as "company.phone2",
        comp.phone3 as "company.phone3",
        comp.email as "company.email",
        comp.email2 as "company.email2",
        comp.email3 as "company.email3",
        comp.meta as "company.meta",
        comp.meta2 as "company.meta2",
        comp.meta3 as "company.meta3",
        comp.address as "company.address",
        comp.country as "company.country",
        comp.country_id as "company.country_id",
        comp.city_id as "company.city_id",
        comp.image_url as "company.image_url",
        comp.entity as "company.entity",
        comp.featured as "company.featured",
        comp.rating as "company.rating",
        comp.partner as "company.partner",
        comp.successful_ops as "company.successful_ops",
        comp.created_at as "company.created_at",
        comp.updated_at as "company.updated_at",
        comp.active as "company.active",
        comp.deleted as "company.deleted",
        
        vt.id as "vehicle_type.id",
        vt.title_en as "vehicle_type.title_en",
        vt.desc_en as "vehicle_type.desc_en",
        vt.title_ru as "vehicle_type.title_ru",
        vt.desc_ru as "vehicle_type.desc_ru",
        vt.title_tk as "vehicle_type.title_tk",
        vt.desc_tk as "vehicle_type.desc_tk",
        vt.title_de as "vehicle_type.title_de",
        vt.desc_de as "vehicle_type.desc_de",
        vt.title_ar as "vehicle_type.title_ar",
        vt.desc_ar as "vehicle_type.desc_ar",
        vt.title_es as "vehicle_type.title_es",
        vt.desc_es as "vehicle_type.desc_es",
        vt.title_fr as "vehicle_type.title_fr",
        vt.desc_fr as "vehicle_type.desc_fr",
        vt.title_zh as "vehicle_type.title_zh",
        vt.desc_zh as "vehicle_type.desc_zh",
        vt.title_ja as "vehicle_type.title_ja",
        vt.desc_ja as "vehicle_type.desc_ja",
        vt.deleted as "vehicle_type.deleted",
        
        pt.id as "packaging_type.id",
        pt.name_ru as "packaging_type.name_ru",
        pt.name_en as "packaging_type.name_en",
        pt.name_tk as "packaging_type.name_tk",
        pt.category_ru as "packaging_type.category_ru",
        pt.category_en as "packaging_type.category_en",
        pt.category_tk as "packaging_type.category_tk",
        pt.material as "packaging_type.material",
        pt.dimensions as "packaging_type.dimensions",
        pt.weight as "packaging_type.weight",
        pt.description_ru as "packaging_type.description_ru",
        pt.description_en as "packaging_type.description_en",
        pt.description_tk as "packaging_type.description_tk",
        pt.active as "packaging_type.active",
        pt.deleted as "packaging_type.deleted"
    FROM tbl_cargo c
    LEFT JOIN tbl_company comp ON c.company_id = comp.id
    LEFT JOIN tbl_vehicle_type vt ON c.vehicle_type_id = vt.id
    LEFT JOIN tbl_packaging_type pt ON c.packaging_type_id = pt.id
`
