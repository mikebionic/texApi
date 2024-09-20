package queries

var GetServices = `
    WITH RECURSIVE parents AS (
        SELECT 
            s.id,
            s.image,
            s.parent_id
        FROM services s
        WHERE s.parent_id IS NULL

        UNION 

        SELECT 
            c.id,
            c.image,
            c.parent_id
        FROM services c 
        INNER JOIN parents p ON p.id = c.parent_id
    ) 
    SELECT
        id,
        JSON_BUILD_OBJECT(
            'tk', GET_SERVICE_TITLE(id, 1),
            'ru', GET_SERVICE_TITLE(id, 2),
            'en', GET_SERVICE_TITLE(id, 3)
        ) AS title,
        image,
        parent_id
    FROM parents
`
var GetService = `
    SELECT
        s.id,
        JSON_BUILD_OBJECT(
            'tk', GET_SERVICE_TITLE(s.id, 1),
            'ru', GET_SERVICE_TITLE(s.id, 2),
            'en', GET_SERVICE_TITLE(s.id, 3)
        ) AS title,
        s.image,
        s.parent_id
    FROM services s
    WHERE s.id = $1
`
var CreateService = `
    INSERT INTO services (parent_id) VALUES (
        CASE WHEN $1 = 0 THEN NULL ELSE $1 END
    ) RETURNING id
`
var CreateServiceTranslates = `
    INSERT INTO service_translates (
        title, languages_id, services_id
    ) VALUES ($1, $2, $7), ($3, $4, $7), ($5, $6, $7)
`
var UpdateService = "UPDATE services SET parent_id = $1 WHERE id = $2"
var UpdateServiceTranslates = `
    UPDATE service_translates SET title = $1
    WHERE languages_id = $2 AND services_id = $3
`
var DeleteService = "DELETE FROM services WHERE id = $1"
var SetServiceImage = "UPDATE services SET image = $1 WHERE id = $2"
var GetServiceList = `
    SELECT
        s.id,
        JSON_BUILD_OBJECT(
            'tk', GET_SERVICE_TITLE(s.id, 1),
            'ru', GET_SERVICE_TITLE(s.id, 2),
            'en', GET_SERVICE_TITLE(s.id, 3)
        ) AS title
    FROM services s
`
