package queries

var GetContents = `
	SELECT 
	c.id,
	c.uuid,
	c.lang_id,
	c.content_type_id,
	c.title,
	c.subtitle,
	c.description,
	c.image_url,
	c.video_url,
	c.step,
	c.created_at::VARCHAR,
	c.updated_at::VARCHAR,
	c.deleted
 FROM tbl_content c WHERE c.deleted = 0`

var GetContent = `
    SELECT
	c.id,
	c.uuid,
	c.lang_id,
	c.content_type_id,
	c.title,
	c.subtitle,
	c.description,
	c.image_url,
	c.video_url,
	c.step,
	c.created_at::VARCHAR,
	c.updated_at::VARCHAR,
	c.deleted
    FROM tbl_content c 
    WHERE c.id = $1
`

var CreateContent = `
    INSERT INTO tbl_content (
	lang_id,
	content_type_id, 
	title,
	subtitle,
	description,
	image_url,
	video_url,
	step)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id
    `
