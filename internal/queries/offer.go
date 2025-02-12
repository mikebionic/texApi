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
    driver_id = COALESCE($2, driver_id),
    vehicle_id = COALESCE($3, vehicle_id),
    cargo_id = COALESCE($4, cargo_id),
    cost_per_km = COALESCE($5, cost_per_km),
    currency = COALESCE($6, currency),
    from_country_id = COALESCE($7, from_country_id),
    from_city_id = COALESCE($8, from_city_id),
    to_country_id = COALESCE($9, to_country_id),
    to_city_id = COALESCE($10, to_city_id),
    from_country = COALESCE($11, from_country),
    from_region = COALESCE($12, from_region),
    to_country = COALESCE($13, to_country),
    to_region = COALESCE($14, to_region),
    from_address = COALESCE($15, from_address),
    to_address = COALESCE($16, to_address),
    sender_contact = COALESCE($17, sender_contact),
    recipient_contact = COALESCE($18, recipient_contact),
    deliver_contact = COALESCE($19, deliver_contact),
    validity_start = COALESCE($20, validity_start),
    validity_end = COALESCE($21, validity_end),
    delivery_start = COALESCE($22, delivery_start),
    delivery_end = COALESCE($23, delivery_end),
    note = COALESCE($24, note),
    tax = COALESCE($25, tax),
    tax_price = COALESCE($26, tax_price),
    trade = COALESCE($27, trade),
    discount = COALESCE($28, discount),
    payment_method = COALESCE($29, payment_method),
    meta = COALESCE($30, meta),
    meta2 = COALESCE($31, meta2),
    meta3 = COALESCE($32, meta3),
    active = COALESCE($33, active),
    deleted = COALESCE($34, deleted),
    exec_company_id = COALESCE($35, exec_company_id),
    offer_state = COALESCE($36, offer_state),
    offer_role = COALESCE($37, offer_role),
	vehicle_type_id = COALESCE($38,vehicle_type_id),
	packaging_type_id = COALESCE($39,packaging_type_id),
	distance = COALESCE($40,distance),
	map_url = COALESCE($41,map_url),
	payment_term = COALESCE($42,payment_term),
    updated_at = NOW()
WHERE id = $1 AND company_id = $43
`

const DeleteOffer = `
UPDATE tbl_offer
SET deleted = 1, updated_at = NOW()
WHERE id = $1 AND company_id = $2;
`
