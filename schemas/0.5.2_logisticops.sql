CREATE TYPE payment_method_t AS ENUM ('cash', 'transfer', 'card', 'terminal', 'online', 'coupon');
CREATE TYPE weight_type_t AS ENUM ('kg', 'g', 'lbs', 'oz', 'st', 't', 'tn');
CREATE TYPE response_state_t AS ENUM ('pending', 'accepted', 'declined');
CREATE TYPE currency_t AS ENUM ('USD', 'TMT');

-- ('kg', 'kg', 1)
-- ('grams', 'g', 0.001)
-- ('pounds', 'lbs', 0.453592)
-- ('ounces', 'oz', 0.0283495)
-- ('stones', 'st', 6.35029)
-- ('tonne', 't', 1000)
-- ('tons', 'tn', 907.18474)
-- ('lbs', 'lbs', 0.453592)

-- Мои заявки
CREATE TABLE
    tbl_offer (
        id SERIAL PRIMARY KEY,
        uuid UUID DEFAULT gen_random_uuid (),
        user_id INT REFERENCES tbl_user (id) ON DELETE CASCADE,
        company_id INT REFERENCES tbl_company (id) ON DELETE CASCADE,
        exec_company_id INT NOT NULL DEFAULT 0,
        driver_id INT NOT NULL DEFAULT 0,
        vehicle_id INT NOT NULL DEFAULT 0,
        cargo_id INT NOT NULL DEFAULT 0,
        offer_state state_t NOT NULL DEFAULT 'pending',
        offer_role role_t NOT NULL DEFAULT 'unknown',
        cost_per_km DECIMAL(10, 2) NOT NULL DEFAULT 0.0,
        currency currency_t NOT NULL DEFAULT 'USD',
        from_country VARCHAR(100) NOT NULL DEFAULT '',
        from_region VARCHAR(100) NOT NULL DEFAULT '',
        to_country VARCHAR(100) NOT NULL DEFAULT '',
        to_region VARCHAR(100) NOT NULL DEFAULT '',
        from_address VARCHAR(100) NOT NULL DEFAULT '',
        to_address VARCHAR(100) NOT NULL DEFAULT '',
        sender_contact VARCHAR(100) NOT NULL DEFAULT '',
        recipient_contact VARCHAR(100) NOT NULL DEFAULT '',
        deliver_contact VARCHAR(100) NOT NULL DEFAULT '',
        view_count INT NOT NULL DEFAULT 0,
        validity_start DATE NOT NULL DEFAULT CURRENT_TIMESTAMP,
        validity_end DATE NOT NULL DEFAULT CURRENT_TIMESTAMP,
        delivery_start DATE NOT NULL DEFAULT CURRENT_TIMESTAMP,
        delivery_end DATE NOT NULL DEFAULT CURRENT_TIMESTAMP,
        note TEXT NOT NULL DEFAULT '',
        tax INT DEFAULT 0,
        trade INT DEFAULT 0,
        payment_method payment_method_t NOT NULL DEFAULT 'cash',
        meta TEXT NOT NULL DEFAULT '',
        meta2 TEXT NOT NULL DEFAULT '',
        meta3 TEXT NOT NULL DEFAULT '',
        featured INT NOT NULL DEFAULT 0,
        partner INT NOT NULL DEFAULT 0,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        active INT NOT NULL DEFAULT 0,
        deleted INT NOT NULL DEFAULT 0
    );


CREATE TABLE
    tbl_cargo (
        id SERIAL PRIMARY KEY,
        uuid UUID DEFAULT gen_random_uuid (),
        company_id INT REFERENCES tbl_company (id) ON DELETE CASCADE,
        name VARCHAR(200) NOT NULL DEFAULT '',
        description VARCHAR(1000) NOT NULL DEFAULT '',
        info VARCHAR(1000) NOT NULL DEFAULT '',
        qty INT NOT NULL DEFAULT 0,
        weight INT NOT NULL DEFAULT 0,
        weight_type weight_type_t NOT NULL DEFAULT 'kg',
        meta TEXT NOT NULL DEFAULT '',
        meta2 TEXT NOT NULL DEFAULT '',
        meta3 TEXT NOT NULL DEFAULT '',
        vehicle_type_id INT NOT NULL DEFAULT 0,
        packaging_type_id INT NOT NULL DEFAULT 0,
        gps INT NOT NULL DEFAULT 0,
        photo1_url VARCHAR(200) NOT NULL DEFAULT '',
        photo2_url VARCHAR(200) NOT NULL DEFAULT '',
        photo3_url VARCHAR(200) NOT NULL DEFAULT '',
        docs1_url VARCHAR(200) NOT NULL DEFAULT '',
        docs2_url VARCHAR(200) NOT NULL DEFAULT '',
        docs3_url VARCHAR(200) NOT NULL DEFAULT '',
        note TEXT NOT NULL DEFAULT '',
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        active INT NOT NULL DEFAULT 0,
        deleted INT NOT NULL DEFAULT 0
    );

-- Мои отклики

CREATE TABLE
    tbl_response (
        id SERIAL PRIMARY KEY,
        uuid UUID DEFAULT gen_random_uuid (),
        company_id INT REFERENCES tbl_company (id) ON DELETE CASCADE DEFAULT 0,
        offer_id INT REFERENCES tbl_offer (id) ON DELETE CASCADE DEFAULT 0,
--         tut nado obyasnit mekanu
        response_company_id INT REFERENCES tbl_company (id) ON DELETE CASCADE DEFAULT 0,
        state response_state_t NOT NULL DEFAULT 'pending',
        title VARCHAR(200) NOT NULL DEFAULT '',
        note VARCHAR(1000) NOT NULL DEFAULT '',
        reason VARCHAR(1000) NOT NULL DEFAULT '',
        meta TEXT NOT NULL DEFAULT '',
        meta2 TEXT NOT NULL DEFAULT '',
        meta3 TEXT NOT NULL DEFAULT '',
        value INT NOT NULL DEFAULT 0,
        rating INT NOT NULL DEFAULT 0,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        deleted INT
    );

-- Insert into tbl_offer
INSERT INTO tbl_offer (
    user_id, company_id, driver_id, vehicle_id, cargo_id,
    offer_state, cost_per_km, currency, from_country, from_region,
    to_country, to_region, from_address, to_address,
    sender_contact, recipient_contact, deliver_contact,
    view_count, validity_start, validity_end,
    delivery_start, delivery_end, note, tax, trade,
    payment_method, meta, meta2, meta3, featured, partner,offer_role
)
VALUES
    (1, 1, 1, 1, 1, 'enabled', 400.25, 'USD', 'Germany', 'Berlin',
     'Italy', 'Rome', 'Berlin Address', 'Rome Address', 'Sender Contact',
     'Recipient Contact', 'Delivery Contact', 10, '2024-11-01', '2025-11-10',
     '2024-11-02', '2025-11-12', 'Urgent transport needed from Berlin to Rome.', 0, 0,
     'cash', '', '', '', 0, 0,'carrier'),

    (2, 2, 2, 2, 2, 'pending', 350.50, 'EUR', 'France', 'Paris',
     'Spain', 'Barcelona', 'Paris Address', 'Barcelona Address', 'Sender Contact',
     'Recipient Contact', 'Delivery Contact', 15, '2024-11-05', '2025-11-12',
     '2024-11-06', '2025-11-11', 'Looking for reliable transport for goods from Paris to Barcelona.', 0, 0,
     'transfer', '', '', '', 0, 0,'sender'),

    (3, 3, 3, 3, 3, 'enabled', 299.75, 'GBP', 'UK', 'London',
     'Netherlands', 'Amsterdam', 'London Address', 'Amsterdam Address', 'Sender Contact',
     'Recipient Contact', 'Delivery Contact', 20, '2024-11-07', '2025-11-14',
     '2024-11-08', '2025-11-13', 'Request for timely delivery from London to Amsterdam.', 0, 0,
     'card', '', '', '', 0, 0,'sender'),

    (1, 1, 4, 4, 4, 'enabled', 100.00, 'PLN', 'Poland', 'Warsaw',
     'Hungary', 'Budapest', 'Warsaw Address', 'Budapest Address', 'Sender Contact',
     'Recipient Contact', 'Delivery Contact', 5, '2024-11-10', '2025-11-17',
     '2024-11-11', '2025-11-16', 'Need cargo transport from Warsaw to Budapest.', 0, 0,
     'terminal', '', '', '', 0, 0,'carrier'),

    (2, 2, 5, 5, 5, 'enabled', 99.99, 'EUR', 'Belgium', 'Brussels',
     'Austria', 'Vienna', 'Brussels Address', 'Vienna Address', 'Sender Contact',
     'Recipient Contact', 'Delivery Contact', 8, '2024-11-12', '2025-11-20',
     '2024-11-13', '2025-11-19', 'Looking for a driver to transport goods from Brussels to Vienna.', 0, 0,
     'online', '', '', '', 0, 0,'sender');
-- Insert into tbl_response
INSERT INTO tbl_response (
    company_id, response_company_id, offer_id, state, title, note, reason,
    value, rating
)
VALUES
    (2, 3, 1, 'declined', 'Response Title 1', 'Response note 1', 'Reason 1', 0, 0),
    (3, 2, 2, 'accepted', 'Response Title 2', 'Response note 2', 'Reason 2', 0, 0),
    (1, 2, 2, 'pending', 'Response Title 3', 'Response note 3', 'Reason 3', 0, 0),
    (2, 3, 3, 'declined', 'Response Title 4', 'Response note 4', 'Reason 4', 0, 0),
    (3, 4, 4, 'declined', 'Response Title 5', 'Response note 5', 'Reason 5', 0, 0),
    (1, 3, 3, 'accepted', 'Response Title 6', 'Response note 6', 'Reason 6', 0, 0),
    (2, 3, 3, 'declined', 'Response Title 7', 'Response note 7', 'Reason 7',0, 0),
    (1, 2, 2, 'pending', 'Response Title 8', 'Response note 8', 'Reason 8',0, 0),
    (2, 1, 1, 'pending', 'Response Title 9', 'Response note 9', 'Reason 9',0, 0),
    (1, 2, 2, 'accepted', 'Response Title 10', 'Response note 10', 'Reason 10',0, 0),
    (2, 1, 1, 'declined', 'Response Title 11', 'Response note 11', 'Reason 11',0, 0),
    (1, 3, 3, 'declined', 'Response Title 12', 'Response note 12', 'Reason 12',0, 0),
    (2, 3, 3, 'declined', 'Response Title 13', 'Response note 13', 'Reason 13',0, 0),
    (3, 1, 1, 'pending', 'Response Title 14', 'Response note 14', 'Reason 14',0, 0),
    (2, 1, 1, 'accepted', 'Response Title 15', 'Response note 15', 'Reason 15',0, 0),
    (1, 3, 3, 'accepted', 'Response Title 16', 'Response note 16', 'Reason 16',0, 0),
    (2, 1, 1, 'pending', 'Response Title 17', 'Response note 17', 'Reason 17',0, 0),
    (3, 2, 2, 'declined', 'Response Title 18', 'Response note 18', 'Reason 18',0, 0),
    (1, 2, 2, 'declined', 'Response Title 19', 'Response note 19', 'Reason 19',0, 0);

INSERT INTO tbl_cargo (
    company_id, name, description, info, qty, weight, meta, meta2, meta3,
    vehicle_type_id, packaging_type_id, gps, photo1_url, photo2_url, photo3_url,
    docs1_url, docs2_url, docs3_url, note, active, deleted
)
VALUES
    (1, 'Electronic Components', 'Various electronic components for manufacturing.', 'Sensitive components requiring careful handling.', 100, 250, 'fragile', '', '', 1, 2, 1,
     'https://example.com/photo1.jpg', 'https://example.com/photo2.jpg', 'https://example.com/photo3.jpg',
     'https://example.com/docs1.pdf', '', '', 'Ensure careful handling during transport.', 1, 0),

    (2, 'Furniture', 'Office furniture including desks, chairs, and filing cabinets.', 'Furniture for a corporate office relocation.', 50, 1200, 'sturdy', '', '', 2, 3, 0,
     'https://example.com/furniture_photo1.jpg', 'https://example.com/furniture_photo2.jpg', '',
     '', '', '', 'Fragile items, handle with care.', 1, 0),

    (3, 'Clothing and Apparel', 'Seasonal clothing and fashion apparel for retail.', 'Bulk shipment for clothing store.', 500, 800, 'perishable', '', '', 3, 1, 0,
     'https://example.com/clothing_photo1.jpg', 'https://example.com/clothing_photo2.jpg', '',
     'https://example.com/clothing_docs.pdf', '', '', 'Seasonal fashion, to be delivered by the end of the week.', 1, 0),


    (4, 'Machinery Parts', 'Heavy machinery parts, including gears, motors, and components.', 'Used for industrial machinery repair.', 20, 5000, 'heavy', '', '', 4, 4, 1,
     'https://example.com/machinery_photo1.jpg', 'https://example.com/machinery_photo2.jpg', 'https://example.com/machinery_photo3.jpg',
     '', '', '', 'Handle with extreme caution due to weight and size.', 1, 0),

    (5, 'Food Products', 'Packaged food products including snacks, canned goods, and beverages.', 'Requires temperature-controlled environment.', 1000, 2000, 'temperature-sensitive', '', '', 5, 2, 1,
     'https://example.com/food_photo1.jpg', '', '',
     'https://example.com/food_docs.pdf', '', '', 'Store in a temperature-controlled environment during transport.', 1, 0);
