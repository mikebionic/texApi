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
		lang_id = COALESCE($1, lang_id),
		content_type_id = COALESCE($2, content_type_id),
		title = COALESCE($3, title),
		slogan = COALESCE($4, slogan),
		subtitle = COALESCE($5, subtitle),
		description = COALESCE($6, description),
		count = COALESCE($7, count),
		count_type = COALESCE($8, count_type),
		image_url = COALESCE($9, image_url),
		video_url = COALESCE($10, video_url),
		step = COALESCE($11, step),
		active = COALESCE($12, active)
	WHERE id = $13
	RETURNING id
`

var DeleteContent = `
	UPDATE tbl_content SET deleted = 1 WHERE id = $1
`
