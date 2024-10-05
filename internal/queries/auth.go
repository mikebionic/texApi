package queries

var GetUser = `SELECT 
	id,
	uuid,
	username,
	password,
	email,
	fullname,
	phone,
	address,
	role_id,
	verified,
	created_at::varchar,
	updated_at::varchar,
	active,
	deleted
FROM tbl_user`
