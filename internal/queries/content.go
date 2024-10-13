package queries

var GetContents = `
	SELECT 
	c.id,
	c.uuid,
	c.lang_id,
	c.content_type_id,
	c.title,
	c.slogan,
	c.subtitle,
	c.description,
	c.count,
	c.count_type,
	c.image_url,
	c.video_url,
	c.step,
	c.created_at::VARCHAR,
	c.updated_at::VARCHAR,
	c.active,
	c.deleted
 FROM tbl_content c WHERE c.deleted = 0`

var GetContent = `
    SELECT
	c.id,
	c.uuid,
	c.lang_id,
	c.content_type_id,
	c.title,
	c.slogan,
	c.subtitle,
	c.description,
	c.count,
	c.count_type,
	c.image_url,
	c.video_url,
	c.step,
	c.created_at::VARCHAR,
	c.updated_at::VARCHAR,
	c.active,
	c.deleted
    FROM tbl_content c 
    WHERE c.id = $1
`

var CreateContent = `
    INSERT INTO tbl_content (
	lang_id,
	content_type_id, 
	title,
	slogan,
	subtitle,
	description,
	count,
	count_type,
	image_url,
	video_url,
	step, active)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING id
    `

var UpdateContent = `
    UPDATE tbl_content
    SET
        lang_id = $1,
        content_type_id = $2,
        title = $3,
        slogan = $4,
        subtitle = $5,
        description = $6,
        count = $7,
        count_type = $8,
        image_url = $9,
        video_url = $10,
        step = $11,
        active = $12
    WHERE id = $13
    RETURNING id
`

var DeleteContent = `
	UPDATE tbl_content SET deleted = 1 WHERE id = $1
`
