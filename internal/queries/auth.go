package queries

var GetUser = `SELECT 
	id,
	uuid,
	username,
	password,
	email,
	first_name,
	last_name,
	nick_name,
	avatar_url,
	phone,
	address,
	role_id,
	verified,
	created_at::varchar,
	updated_at::varchar,
	active,
	deleted,
	oauth_provider,
	oauth_user_id,
	oauth_location,
	oauth_access_token,
	oauth_access_token_secret,
	oauth_refresh_token,
	oauth_expires_at::varchar,
	oauth_id_token
FROM tbl_user`
