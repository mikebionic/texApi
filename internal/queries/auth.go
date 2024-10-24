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
	refresh_token,
	verify_time::varchar
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

var UpdateUser = `
    UPDATE tbl_user
    SET
        username = $2,
        password = $3,
        email = $4,
        first_name = $5,
        last_name = $6,
        nick_name = $7,
        avatar_url = $8,
        phone = $9,
        info_phone = $10,
        address = $11,
        role_id = $12,
        subrole_id = $13,
        verified = $14,
        active = $15,
        oauth_provider = $16,
        oauth_user_id = $17,
        oauth_location = $18,
        oauth_access_token = $19,
        oauth_access_token_secret = $20,
        oauth_refresh_token = $21,
        oauth_id_token = $22
    WHERE id = $1
    RETURNING id;
`

var SaveUserWithOTP = `
INSERT INTO tbl_user (
email, phone, role_id, verified, otp_key, verify_time
) VALUES (
CASE WHEN $1 = 'email' THEN $2 ELSE '' END,
CASE WHEN $1 = 'phone' THEN $2 ELSE '' END,
$3,
0, $4, NOW()
) RETURNING id;
`
var UpdateUserWithOTP = `
UPDATE tbl_user
SET 
    email = CASE WHEN $1 = 'email' THEN $2 ELSE email END,
    phone = CASE WHEN $1 = 'phone' THEN $2 ELSE phone END,
    role_id = $3,
    otp_key = $4,
    verify_time = NOW(),
    verified = 0
WHERE id = $5
RETURNING id;
`

var GetOTPInfo = `
SELECT id, otp_key, verify_time
FROM tbl_user
WHERE 
    (CASE WHEN $1 = 'email' THEN email ELSE phone END) = $2
LIMIT 1;
`

var VerifyUserByID = `
UPDATE tbl_user
SET
	verified = 1,
	verify_time = NOW()
WHERE id = $1
RETURNING id;
`
