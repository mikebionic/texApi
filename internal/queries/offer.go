package queries

const GetMyOfferList = `
WITH offer_data AS (
    SELECT 
        o.*,
        COUNT(*) OVER() as total_count
    FROM tbl_offer o
    WHERE o.deleted = 0 AND o.company_id = $1 AND o.offer_state = $4
    ORDER BY o.id
    LIMIT $2 OFFSET $3
)
SELECT 
    od.*,
    json_build_object(
        'id', c.id,
        'company_name', c.company_name,
        'country', c.country
    ) as company,
    json_build_object(
        'id', d.id,
        'first_name', d.first_name,
        'last_name', d.last_name,
        'image_url', d.image_url
    ) as assigned_driver,
    json_build_object(
        'id', v.id,
        'vehicle_type', v.vehicle_type,
        'numberplate', v.numberplate
    ) as assigned_vehicle,
    json_build_object(
        'id', c2.id,
        'name', c2.name,
        'description', c2.description,
        'info', c2.info
    ) as cargo
FROM offer_data od
LEFT JOIN tbl_company c ON od.company_id = c.id
LEFT JOIN tbl_driver d ON od.driver_id = d.id
LEFT JOIN tbl_vehicle v ON od.vehicle_id = v.id
LEFT JOIN tbl_cargo c2 ON od.cargo_id = c2.id
GROUP BY 
    od.id, od.uuid, od.user_id, od.company_id, od.exec_company_id,
    od.driver_id, od.vehicle_id, od.cargo_id, od.offer_state, 
    od.offer_role, od.cost_per_km, od.currency, 
    od.from_country_id, od.from_city_id, od.to_country_id, od.to_city_id, 
    od.from_country, od.from_region, od.to_country, od.to_region, od.from_address, 
    od.to_address, od.sender_contact, od.recipient_contact, 
    od.deliver_contact, od.view_count, od.validity_start, 
    od.validity_end, od.delivery_start, od.delivery_end, od.note, 
    od.tax, od.tax_price, od.trade, od.discount, od.payment_method, od.meta, od.meta2, 
    od.meta3, od.featured, od.partner, od.created_at, od.updated_at, 
    od.active, od.deleted, od.total_count, 
    c.id, c.company_name, c.country,
    d.id, d.first_name, d.last_name, d.image_url,
    v.id, v.vehicle_type, v.numberplate,
    c2.id, c2.name, c2.description, c2.info;

`

const GetOfferList = `
    SELECT 
        o.*,
        COUNT(*) OVER() as total_count
    FROM tbl_offer o
    WHERE 
        o.validity_end > CURRENT_TIMESTAMP
        AND o.delivery_end > CURRENT_TIMESTAMP
`

const GetOfferByID = `
SELECT 
    o.*,
    json_build_object(
        'id', c.id,
        'company_name', c.company_name,
        'country', c.country
    ) as company,
    json_build_object(
        'id', d.id,
        'first_name', d.first_name,
        'last_name', d.last_name,
        'image_url', d.image_url
    ) as assigned_driver,
    json_build_object(
        'id', v.id,
        'vehicle_type', v.vehicle_type,
        'numberplate', v.numberplate
    ) as assigned_vehicle,
    json_build_object(
        'id', c2.id,
        'name', c2.name,
        'description', c2.description,
        'info', c2.info
    ) as cargo
FROM tbl_offer o
LEFT JOIN tbl_company c ON o.company_id = c.id
LEFT JOIN tbl_driver d ON o.driver_id = d.id
LEFT JOIN tbl_vehicle v ON o.vehicle_id = v.id
LEFT JOIN tbl_cargo c2 ON o.cargo_id = c2.id
WHERE o.id = $1 AND o.deleted = 0;
`
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
    exec_company_id
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34)
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
    updated_at = NOW()
WHERE (id = $1 AND company_id = $36) AND (active = 1 AND deleted = 0)
RETURNING id;
`

const DeleteOffer = `
UPDATE tbl_offer
SET deleted = 1, updated_at = NOW()
WHERE id = $1 AND company_id = $2;
`
