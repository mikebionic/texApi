package queries

var GetPackagingTypes = `
    SELECT 
        id, name_ru, name_en, name_tk,
        category_ru, category_en, category_tk,
        material, dimensions, weight,
        description_ru, description_en, description_tk,
        active, deleted
    FROM tbl_packaging_type 
    WHERE deleted = 0`

var GetPackagingType = `
    SELECT 
        id, name_ru, name_en, name_tk,
        category_ru, category_en, category_tk,
        material, dimensions, weight,
        description_ru, description_en, description_tk,
        active, deleted
    FROM tbl_packaging_type 
    WHERE id = $1 AND deleted = 0`

var CreatePackagingType = `
    INSERT INTO tbl_packaging_type (
        name_ru, name_en, name_tk,
        category_ru, category_en, category_tk,
        material, dimensions, weight,
        description_ru, description_en, description_tk,
        active
    ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
    RETURNING id`

var UpdatePackagingType = `
    UPDATE tbl_packaging_type
    SET
        name_ru = COALESCE($1, name_ru),
        name_en = COALESCE($2, name_en),
        name_tk = COALESCE($3, name_tk),
        category_ru = COALESCE($4, category_ru),
        category_en = COALESCE($5, category_en),
        category_tk = COALESCE($6, category_tk),
        material = COALESCE($7, material),
        dimensions = COALESCE($8, dimensions),
        weight = COALESCE($9, weight),
        description_ru = COALESCE($10, description_ru),
        description_en = COALESCE($11, description_en),
        description_tk = COALESCE($12, description_tk),
        active = COALESCE($13, active)
    WHERE id = $14
    RETURNING id`

var DeletePackagingType = `
    UPDATE tbl_packaging_type SET deleted = 1 WHERE id = $1`
