package queries

var GetContentTypes = `
	SELECT
	id,
	uuid,
	name,
	title,
	title_ru,
	description,
	parent_id,
	parent_name FROM tbl_content_type
`

var GetContentTypesWithContent = `
SELECT
    ct.id,
    ct.uuid,
    ct.name,
    ct.title,
    ct.title_ru,
    ct.description,
    ct.parent_id,
    ct.parent_name,
    json_agg(
        json_build_object(
            'id', c.id,
            'uuid', c.uuid,
            'lang_id', c.lang_id,
            'content_type_id', c.content_type_id,
            'title', c.title,
            'slogan', c.slogan,
            'subtitle', c.subtitle,
            'description', c.description,
			'count', c.count,
			'count_type', c.count_type,
            'image_url', c.image_url,
            'video_url', c.video_url,
            'step', c.step,
            'created_at', c.created_at::VARCHAR,
            'updated_at', c.updated_at::VARCHAR,
            'active', c.active,
            'deleted', c.deleted
        ) 
    ) FILTER (WHERE c.id != 0 AND (c.lang_id = $1 OR $1 = 0)) AS contents 
FROM
    tbl_content_type ct
LEFT JOIN
    tbl_content c ON c.content_type_id = ct.id AND c.deleted = 0
WHERE
    (ct.id = $2 OR $2 = 0) 
GROUP BY
   ct.id
ORDER BY
   ct.id;
`
