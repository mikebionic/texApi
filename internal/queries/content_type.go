package queries

var GetContentTypes = `
	SELECT
	id,
	uuid,
	name,
	title,
	description FROM tbl_content_type
`

var GetContentTypesWithContent = `
	SELECT 
-- 		ct.id AS content_type_id,
		ct.uuid AS content_type_uuid,
		ct.name AS content_type_name,
		ct.title AS content_type_title,
		ct.description AS content_type_description,
		c.id,
		c.uuid,
		c.lang_id,
-- 		c.content_type_id AS content_content_type_id,
		c.title,
		c.subtitle,
		c.description,
		c.image_url,
		c.video_url,
		c.step,
		c.created_at::VARCHAR,
		c.updated_at::VARCHAR,
		c.deleted
	FROM 
		tbl_content_type ct
	LEFT JOIN 
		tbl_content c ON c.content_type_id = ct.id
	WHERE 
		c.deleted = 0 OR c.deleted IS NULL
	ORDER BY 
		ct.id, c.step
`
