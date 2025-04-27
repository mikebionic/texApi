package queries

var GetUser = `
SELECT 
    id,
    uuid,
    username,
    password,
    email,
    phone,
    role,
    role_id,
    company_id,
    verified,
    meta,
    meta2,
    meta3,
    refresh_token,
    otp_key,
    verify_time::varchar,
    created_at::varchar,
    updated_at::varchar,
    active,
    deleted
FROM tbl_user`

var CreateUser = `
    INSERT INTO tbl_user (
        username,
        password,
        email,
        phone,
        role,
        role_id,
        company_id,
        verified,
        meta,
        meta2,
        meta3,
        otp_key,
        refresh_token,
        active
    ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14) 
    RETURNING id
`
var UpdateUser = `
    UPDATE tbl_user
    SET
        username = COALESCE($2,username),
        password = COALESCE($3,password),
        email = COALESCE($4,email),
        phone = COALESCE($5,phone),
        role = COALESCE($6,role),
        role_id = COALESCE($7,role_id),
        company_id = COALESCE($8,company_id),
        verified = COALESCE($9,verified),
        active = COALESCE($10,active),
        updated_at = CURRENT_TIMESTAMP
    WHERE id = $1
    RETURNING id
`

var UserUpdate = `
    UPDATE tbl_user
    SET
        username = COALESCE($2, username),
        password = COALESCE($3, password),
        email = COALESCE($4, email),
        phone = COALESCE($5, phone),
        updated_at = CURRENT_TIMESTAMP
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
	$4, $5, $6, NOW()
) RETURNING id;
`
var UpdateUserWithOTP = `
UPDATE tbl_user
SET 
    email = CASE WHEN $1 = 'email' THEN $2 ELSE email END,
    phone = CASE WHEN $1 = 'phone' THEN $2 ELSE phone END,
    role = $3,
    role_id = $4,
    verified = $5,
    otp_key = $6,
    verify_time = NOW()
WHERE id = $7
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
    verify_time = CURRENT_TIMESTAMP,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING id;
`
