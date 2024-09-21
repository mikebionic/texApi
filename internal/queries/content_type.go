package queries

var GetContentTypes = `
	SELECT
	id,
	uuid,
	name,
	title,
	description FROM tbl_content_type
`
