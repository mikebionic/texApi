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
c.created_at,
c.updated_at,
c.deleted,
 FROM content c WHERE deleted = 0`
