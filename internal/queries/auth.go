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
	info_phone,
	address,
	role_id,
	subrole_id,
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
	oauth_id_token,
	refresh_token
FROM tbl_user`

var CreateUser = `
    INSERT INTO tbl_user (
		username,
		password,
		email,
		first_name,
		last_name,
		nick_name,
		avatar_url,
		phone,
		info_phone,
		address,
		role_id,
		subrole_id,
		verified,
		active,
		oauth_provider,
		oauth_user_id,
		oauth_location,
		oauth_access_token,
		oauth_access_token_secret,
		oauth_refresh_token,
		oauth_id_token
    ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21) 
      RETURNING id
`
