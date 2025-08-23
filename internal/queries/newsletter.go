package queries

const GetNewsletterList = `
	SELECT 
		n.*,
		COUNT(*) OVER() as total_count
	FROM tbl_newsletter n
`

const GetNewsletterByID = `
	SELECT 
		n.*
	FROM tbl_newsletter n
	WHERE n.id = $1
`

const CreateNewsletter = `
	INSERT INTO tbl_newsletter (
		email, status, first_name, last_name, frequency, 
		ip_address, user_agent, referrer_url, meta, meta2, meta3
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	RETURNING id;
`

const UpdateNewsletter = `
	UPDATE tbl_newsletter
	SET 
		email = COALESCE($2, email),
		status = COALESCE($3, status),
		first_name = COALESCE($4, first_name),
		last_name = COALESCE($5, last_name),
		frequency = COALESCE($6, frequency),
		ip_address = COALESCE($7, ip_address),
		user_agent = COALESCE($8, user_agent),
		referrer_url = COALESCE($9, referrer_url),
		meta = COALESCE($10, meta),
		meta2 = COALESCE($11, meta2),
		meta3 = COALESCE($12, meta3),
		active = COALESCE($13, active),
		deleted = COALESCE($14, deleted),
		unsubscribed_at = CASE 
			WHEN COALESCE($3, status) = 'unsubscribed' AND status != 'unsubscribed' 
			THEN NOW() 
			ELSE unsubscribed_at 
		END,
		updated_at = NOW()
	WHERE id = $1
`

const DeleteNewsletter = `
	UPDATE tbl_newsletter
	SET deleted = 1, active = 0, updated_at = NOW()
	WHERE id = $1
`

const CheckEmailExists = `
	SELECT id FROM tbl_newsletter WHERE email = $1 AND deleted = 0
`
