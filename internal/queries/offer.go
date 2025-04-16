package queries

const CreateOffer = `
INSERT INTO tbl_offer (
    user_id,
    company_id,
    driver_id,
    vehicle_id,
    cargo_id,
    cost_per_km,
    currency,
    from_country_id,
    from_city_id,
    to_country_id,
    to_city_id,
    from_country,
    from_region,
    to_country,
    to_region,
    from_address,
    to_address,
    sender_contact,
    recipient_contact,
    deliver_contact,
    validity_start,
    validity_end,
    delivery_start,
    delivery_end,
    note,
    tax,
    tax_price,
    trade,
    discount,
    payment_method,
    meta,
    meta2,
    meta3,
    offer_role,
    exec_company_id,
	vehicle_type_id,
	packaging_type_id,
	distance,
	map_url,
	payment_term
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34, $35, $36, $37, $38, $39, $40)
RETURNING id;
`
const UpdateOffer = `
UPDATE tbl_offer
SET
    company_id = COALESCE($2,company_id),
    exec_company_id = COALESCE($3,exec_company_id),
    driver_id = COALESCE($4,driver_id),
    vehicle_id = COALESCE($5,vehicle_id),
    vehicle_type_id = COALESCE($6,vehicle_type_id),
    cargo_id = COALESCE($7,cargo_id),
    packaging_type_id = COALESCE($8,packaging_type_id),
    offer_state = COALESCE($9,offer_state),
    offer_role = COALESCE($10,offer_role),
    cost_per_km = COALESCE($11,cost_per_km),
    currency = COALESCE($12,currency),
    from_country_id = COALESCE($13,from_country_id),
    from_city_id = COALESCE($14,from_city_id),
    to_country_id = COALESCE($15,to_country_id),
    to_city_id = COALESCE($16,to_city_id),
    distance = COALESCE($17,distance),
    from_country = COALESCE($18,from_country),
    from_region = COALESCE($19,from_region),
    to_country = COALESCE($20,to_country),
    to_region = COALESCE($21,to_region),
    from_address = COALESCE($22,from_address),
    to_address = COALESCE($23,to_address),
    map_url = COALESCE($24,map_url),
    sender_contact = COALESCE($25,sender_contact),
    recipient_contact = COALESCE($26,recipient_contact),
    deliver_contact = COALESCE($27,deliver_contact),
    view_count = COALESCE($28,view_count),
    validity_start = COALESCE($29,validity_start),
    validity_end = COALESCE($30,validity_end),
    delivery_start = COALESCE($31,delivery_start),
    delivery_end = COALESCE($32,delivery_end),
    note = COALESCE($33,note),
    tax = COALESCE($34,tax),
    tax_price = COALESCE($35,tax_price),
    trade = COALESCE($36,trade),
    discount = COALESCE($37,discount),
	payment_method = COALESCE($38,payment_method),
	payment_term = COALESCE($39,payment_term),
	meta = COALESCE($40,meta),
	meta2 = COALESCE($41,meta2),
	meta3 = COALESCE($42,meta3),
	featured = COALESCE($43,featured),
	partner = COALESCE($44,partner),
	active = COALESCE($45,active),
	deleted = COALESCE($46,deleted),
    updated_at = NOW()
WHERE id = $1 AND company_id = $2
`

const DeleteOffer = `
UPDATE tbl_offer
SET deleted = 1, updated_at = NOW()
WHERE id = $1 AND company_id = $2;
`
