package queries

var GetUser = `
SELECT 
    u.id,
    u.uuid,
    u.username,
    u.password,
    u.email,
    u.phone,
    r.role as role,
    u.role_id,
    u.company_id,
    u.verified,
    u.created_at::varchar,
    u.updated_at::varchar,
    u.active,
    u.deleted,
    u.oauth_provider,
    u.oauth_user_id,
    u.oauth_location,
    u.oauth_access_token,
    u.oauth_access_token_secret,
    u.oauth_refresh_token,
    u.oauth_expires_at::varchar,
    u.oauth_id_token,
    u.refresh_token,
    u.verify_time::varchar,
    u.otp_key
FROM tbl_user u
LEFT JOIN tbl_role r ON u.role_id = r.id`

var CreateUser = `
    INSERT INTO tbl_user (
        username,
        password,
        email,
        phone,
        role_id,
        company_id,
        verified,
        active,
        oauth_provider,
        oauth_user_id,
        oauth_location,
        oauth_access_token,
        oauth_access_token_secret,
        oauth_refresh_token,
        oauth_id_token
    ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15) 
    RETURNING id
`

var UpdateUser = `
    UPDATE tbl_user
    SET
        username = $2,
        password = $3,
        email = $4,
        phone = $5,
        role = $6,
        role_id = $7,
        company_id = $8,
        verified = $9,
        active = $10,
        oauth_provider = $11,
        oauth_user_id = $12,
        oauth_location = $13,
        oauth_access_token = $14,
        oauth_access_token_secret = $15,
        oauth_refresh_token = $16,
        oauth_id_token = $17
    WHERE id = $1
    RETURNING id
`

var ProfileUpdate = `
    UPDATE tbl_user
    SET
        username = COALESCE($2, username),
        password = COALESCE($3, password),
        email = COALESCE($4, email),
        phone = COALESCE($5, phone)
    WHERE id = $1
    RETURNING id
`

var SaveUserWithOTP = `
INSERT INTO tbl_user (
email, phone, role_id, verified, otp_key, verify_time
) VALUES (
CASE WHEN $1 = 'email' THEN $2 ELSE '' END,
CASE WHEN $1 = 'phone' THEN $2 ELSE '' END,
$3,
$4, $5, NOW()
) RETURNING id;
`
var UpdateUserWithOTP = `
UPDATE tbl_user
SET 
    email = CASE WHEN $1 = 'email' THEN $2 ELSE email END,
    phone = CASE WHEN $1 = 'phone' THEN $2 ELSE phone END,
    role_id = $3,
    verified = $4,
    otp_key = $5,
    verify_time = NOW()
WHERE id = $6
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
