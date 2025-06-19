CREATE TYPE payment_method_t AS ENUM ('cash', 'transfer', 'card','credit','terminal', 'online', 'coupon');
CREATE TYPE weight_type_t AS ENUM ('kg', 'g', 'lbs', 'oz', 'st', 't', 'tn');
CREATE TYPE currency_t AS ENUM (
    'USD',    -- United States Dollar
    'TMT',    -- Turkmenistan Manat
    'EUR',    -- Euro
    'GBP',    -- British Pound Sterling
    'INR',    -- Indian Rupee
    'JPY',    -- Japanese Yen
    'CNY',    -- Chinese Yuan
    'AUD',    -- Australian Dollar
    'CAD',    -- Canadian Dollar
    'CHF',    -- Swiss Franc
    'MXN',    -- Mexican Peso
    'BRL',    -- Brazilian Real
    'RUB',    -- Russian Ruble
    'ZAR',    -- South African Rand
    'SGD',    -- Singapore Dollar
    'KRW',    -- South Korean Won
    'MYR',    -- Malaysian Ringgit
    'PHP',    -- Philippine Peso
    'TRY',    -- Turkish Lira
    'IDR',    -- Indonesian Rupiah
    'AED',    -- United Arab Emirates Dirham
    'SAR',    -- Saudi Riyal
    'THB',    -- Thai Baht
    'SEK',    -- Swedish Krona
    'DKK',    -- Danish Krone
    'NOK',    -- Norwegian Krone
    'HKD',    -- Hong Kong Dollar
    'PLN',    -- Polish Zloty
    'NZD',    -- New Zealand Dollar
    'VND',    -- Vietnamese Dong
    'EGP'     -- Egyptian Pound
    );

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
        trailer_id INT NOT NULL DEFAULT 0,
        vehicle_type_id INT NOT NULL DEFAULT 0,
        cargo_id INT NOT NULL DEFAULT 0,
        packaging_type_id INT NOT NULL DEFAULT 0,
        offer_state state_t NOT NULL DEFAULT 'pending',
        offer_role role_t NOT NULL DEFAULT 'unknown',
        cost_per_km DECIMAL(10, 2) NOT NULL DEFAULT 0.0,
        currency currency_t NOT NULL DEFAULT 'USD',
        from_country_id INT NOT NULL DEFAULT 0,
        from_city_id INT NOT NULL DEFAULT 0,
        to_country_id INT NOT NULL DEFAULT 0,
        to_city_id INT NOT NULL DEFAULT 0,
        distance    INT NOT NULL DEFAULT 0,
        from_country VARCHAR(100) NOT NULL DEFAULT '',
        from_region VARCHAR(100) NOT NULL DEFAULT '',
        to_country VARCHAR(100) NOT NULL DEFAULT '',
        to_region VARCHAR(100) NOT NULL DEFAULT '',
        from_address VARCHAR(100) NOT NULL DEFAULT '',
        to_address VARCHAR(100) NOT NULL DEFAULT '',
        map_url VARCHAR(500) NOT NULL DEFAULT '',
        sender_contact VARCHAR(100) NOT NULL DEFAULT '',
        recipient_contact VARCHAR(100) NOT NULL DEFAULT '',
        deliver_contact VARCHAR(100) NOT NULL DEFAULT '',
        view_count INT NOT NULL DEFAULT 0,
        validity_start DATE NOT NULL DEFAULT CURRENT_TIMESTAMP,
        validity_end DATE NOT NULL DEFAULT CURRENT_TIMESTAMP,
        delivery_start DATE NOT NULL DEFAULT CURRENT_TIMESTAMP,
        delivery_end DATE NOT NULL DEFAULT CURRENT_TIMESTAMP,
        note TEXT NOT NULL DEFAULT '',
        tax INT NOT NULL DEFAULT 0,
        tax_price DECIMAL(10, 2) NOT NULL DEFAULT 0.0,
        trade INT NOT NULL DEFAULT 0,
        discount INT NOT NULL DEFAULT 0,
        payment_method payment_method_t NOT NULL DEFAULT 'cash',
        payment_term VARCHAR NOT NULL DEFAULT '',
        meta TEXT NOT NULL DEFAULT '',
        meta2 TEXT NOT NULL DEFAULT '',
        meta3 TEXT NOT NULL DEFAULT '',
        featured INT NOT NULL DEFAULT 0,
        partner INT NOT NULL DEFAULT 0,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        active INT NOT NULL DEFAULT 1,
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
        active INT NOT NULL DEFAULT 1,
        deleted INT NOT NULL DEFAULT 0
    );


CREATE TYPE response_state_t AS ENUM ('pending', 'accepted', 'declined');
CREATE TABLE
    tbl_offer_response (
        id SERIAL PRIMARY KEY,
        uuid UUID DEFAULT gen_random_uuid (),
        company_id INT REFERENCES tbl_company (id) ON DELETE CASCADE DEFAULT 0,
        offer_id INT REFERENCES tbl_offer (id) ON DELETE CASCADE DEFAULT 0,
        to_company_id INT NOT NULL DEFAULT 0,
        state response_state_t NOT NULL DEFAULT 'pending',
        bid_price DECIMAL(10, 2),
        title VARCHAR(200),
        note VARCHAR(1000),
        reason VARCHAR(1000),
        meta TEXT,
        meta2 TEXT,
        meta3 TEXT,
        value INT,
        rating INT,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        deleted INT DEFAULT 0
    );


INSERT INTO tbl_offer (
    user_id, company_id, driver_id, vehicle_id, cargo_id,
    offer_state, cost_per_km, currency, from_country_id, from_city_id,
    to_country_id, to_city_id, from_country, from_region,
    to_country, to_region, from_address, to_address,
    sender_contact, recipient_contact, deliver_contact,
    view_count, validity_start, validity_end,
    delivery_start, delivery_end, note, tax, trade,
    payment_method, meta, meta2, meta3, featured, partner, offer_role
)
VALUES
    -- Offer 1
    (1, 1, 1, 1, 1, 'enabled', 400.25, 'USD', 1, 1, 2, 2, 'Germany', 'Berlin',
     'Italy', 'Rome', '123 Berlin St, Berlin, Germany', '456 Rome St, Rome, Italy',
     '+49 170 1234567', '+39 06 1234567', '+49 170 7654321',
     10, '2024-11-01', '2025-11-10',
     '2024-11-02', '2025-11-12', 'Urgent transport needed from Berlin to Rome. Special cargo, time-sensitive.', 0, 0,
     'cash', '', '', '', 0, 0, 'carrier'),

    -- Offer 2
    (2, 2, 2, 2, 2, 'pending', 350.50, 'EUR', 3, 3, 4, 4, 'France', 'Paris',
     'Spain', 'Barcelona', '10 Rue de Paris, Paris, France', '20 Carrer de Barcelona, Barcelona, Spain',
     '+33 1 23 45 67 89', '+34 93 123 4567', '+33 1 98 76 54 32',
     15, '2024-11-05', '2025-11-12',
     '2024-11-06', '2025-11-11', 'Looking for reliable transport for goods from Paris to Barcelona. Includes fragile items.', 5, 0,
     'transfer', '', '', '', 0, 0, 'sender'),

    -- Offer 3
    (3, 3, 3, 3, 3, 'enabled', 299.75, 'GBP', 5, 5, 6, 6, 'UK', 'London',
     'Netherlands', 'Amsterdam', '10 Downing St, London, UK', 'Damstraat 12, Amsterdam, Netherlands',
     '+44 20 7946 0958', '+31 20 7946 0987', '+44 20 7946 1234',
     20, '2024-11-07', '2025-02-11',
     '2024-11-08', '2025-11-13', 'Request for timely delivery from London to Amsterdam. Cargo includes electronics and perishable items.', 0, 0,
     'card', '', '', '', 0, 0, 'sender'),

    -- Offer 4
    (1, 1, 4, 4, 4, 'enabled', 100.00, 'PLN', 7, 7, 8, 8, 'Poland', 'Warsaw',
     'Hungary', 'Budapest', '12 Warszawa Rd, Warsaw, Poland', '15 Budapest Str, Budapest, Hungary',
     '+48 22 123 4567', '+36 1 234 5678', '+48 22 765 4321',
     5, '2024-11-10', '2025-11-17',
     '2024-11-11', '2025-11-16', 'Need cargo transport from Warsaw to Budapest. Shipment includes hazardous materials.', 0, 0,
     'terminal', '', '', '', 0, 0, 'carrier'),

    -- Offer 5
    (2, 2, 5, 5, 5, 'enabled', 99.99, 'EUR', 9, 9, 10, 10, 'Belgium', 'Brussels',
     'Austria', 'Vienna', '100 Brussels Ave, Brussels, Belgium', '5 Stephansplatz, Vienna, Austria',
     '+32 2 123 4567', '+43 1 23456789', '+32 2 7654321',
     8, '2024-11-12', '2025-11-20',
     '2024-11-13', '2025-11-19', 'Looking for a driver to transport goods from Brussels to Vienna. The goods are non-perishable and require special handling.', 0, 0,
     'online', '', '', '', 0, 0, 'sender'),

    -- Offer 6
    (3, 3, 6, 6, 6, 'pending', 320.75, 'USD', 11, 11, 12, 12, 'Germany', 'Munich',
     'Austria', 'Salzburg', '21 Munich Blvd, Munich, Germany', '10 Mozartstrasse, Salzburg, Austria',
     '+49 89 123456', '+43 662 123456', '+49 89 654321',
     3, '2024-11-15', '2025-11-22',
     '2024-11-16', '2025-11-23', 'Request for transport of large items from Munich to Salzburg. Cargo includes furniture and appliances.', 0, 0,
     'credit', '', '', '', 0, 0, 'carrier'),

    -- Offer 7
    (4, 4, 7, 7, 7, 'enabled', 500.50, 'EUR', 13, 13, 14, 14, 'Italy', 'Rome',
     'Switzerland', 'Zurich', '22 Via Roma, Rome, Italy', '8 Bahnhofstrasse, Zurich, Switzerland',
     '+39 06 12345678', '+41 44 123 4567', '+39 06 98765432',
     25, '2024-11-20', '2025-12-01',
     '2024-11-21', '2025-11-28', 'Urgent cargo transport needed from Rome to Zurich. Shipment includes medical supplies and documents that need special handling.', 0, 0,
     'online', '', '', '', 0, 0, 'carrier'),

    -- Offer 8
    (4, 4, 8, 8, 8, 'enabled', 150.00, 'USD', 15, 15, 16, 16, 'USA', 'New York',
     'Canada', 'Toronto', '101 Wall Street, New York, USA', '123 King St W, Toronto, Canada',
     '+1 212-555-1234', '+1 416-555-6789', '+1 212-555-4321',
     30, '2024-11-22', '2025-12-05',
     '2024-11-23', '2025-12-02', 'Transport of goods from New York to Toronto. Includes electronics and other sensitive items.', 0, 0,
     'cash', '', '', '', 0, 0, 'carrier'),

    -- Offer 9
    (2, 3, 9, 9, 9, 'pending', 250.00, 'GBP', 17, 17, 18, 18, 'Canada', 'Vancouver',
     'USA', 'Seattle', '555 Vancouver St, Vancouver, Canada', '200 Pine St, Seattle, USA',
     '+1 604-555-7890', '+1 206-555-2345', '+1 604-555-8765',
     40, '2024-11-25', '2025-12-10',
     '2024-11-26', '2025-12-08', 'Goods transport from Vancouver to Seattle. Cargo is perishable, requires refrigerated transport.', 5, 0,
     'transfer', '', '', '', 0, 0, 'sender'),

    -- Offer 10
    (1, 1, 10, 10, 10, 'enabled', 150.50, 'USD', 19, 19, 20, 20, 'Brazil', 'São Paulo',
     'Mexico', 'Mexico City', '10 Avenida Paulista, São Paulo, Brazil', '50 Reforma Ave, Mexico City, Mexico',
     '+55 11 1234-5678', '+52 55 1234-5678', '+55 11 8765-4321',
     50, '2024-11-30', '2025-12-15',
     '2024-12-01', '2025-12-10', 'Transport of large machinery from São Paulo to Mexico City. Requires heavy-duty transport vehicle.', 0, 0,
     'cash', '', '', '', 0, 0, 'carrier'),

-- Offer 11
(1, 1, 11, 11, 11, 'enabled', 210.50, 'USD', 1, 1, 2, 2, 'USA', 'Los Angeles',
    'Canada', 'Vancouver', '123 Sunset Blvd, Los Angeles, USA', '456 Maple St, Vancouver, Canada',
    '+1 310-555-9876', '+1 604-555-6789', '+1 310-555-5432',
    12, '2024-12-05', '2025-12-20',
    '2024-12-06', '2025-12-18', 'Transporting large containers from Los Angeles to Vancouver. Requires a specialized vehicle.', 0, 0,
    'transfer', '', '', '', 0, 0, 'carrier'),

-- Offer 12
(2, 2, 12, 12, 12, 'pending', 190.30, 'CAD', 3, 3, 4, 4, 'Canada', 'Toronto',
    'USA', 'New York', '100 Bay St, Toronto, Canada', '15 5th Ave, New York, USA',
    '+1 416-555-1234', '+1 212-555-2345', '+1 416-555-8765',
    5, '2024-12-10', '2025-12-30',
    '2024-12-11', '2025-12-25', 'Looking for a carrier for transport of goods from Toronto to New York. Cargo includes furniture.', 0, 0,
    'cash', '', '', '', 0, 0, 'sender'),

-- Offer 13
(3, 3, 13, 13, 13, 'enabled', 320.25, 'EUR', 5, 5, 6, 6, 'Germany', 'Berlin',
    'Austria', 'Vienna', '222 Alexanderplatz, Berlin, Germany', '101 Ringstrasse, Vienna, Austria',
    '+49 30 12345678', '+43 1 23456789', '+49 30 87654321',
    8, '2024-12-12', '2025-12-22',
    '2024-12-13', '2025-12-20', 'Special shipment of electronics from Berlin to Vienna. Requires temperature-controlled transport.', 0, 0,
    'online', '', '', '', 0, 0, 'carrier'),

-- Offer 14
(4, 4, 14, 14, 14, 'enabled', 150.75, 'GBP', 7, 7, 8, 8, 'UK', 'Manchester',
    'France', 'Paris', '10 High St, Manchester, UK', '45 Rue Lafayette, Paris, France',
    '+44 161 555 2345', '+33 1 70 55 78 90', '+44 161 555 9876',
    20, '2024-12-15', '2025-12-25',
    '2024-12-16', '2025-12-24', 'Transporting luxury goods from Manchester to Paris. Delivery needs to be discreet and secure.', 0, 0,
    'credit', '', '', '', 0, 0, 'carrier'),

-- Offer 15
(1, 1, 15, 15, 15, 'pending', 275.50, 'AUD', 9, 9, 10, 10, 'Australia', 'Sydney',
    'New Zealand', 'Auckland', '20 George St, Sydney, Australia', '120 Queen St, Auckland, New Zealand',
    '+61 2 5550 1234', '+64 9 555 9876', '+61 2 5550 4321',
    30, '2024-12-18', '2025-12-28',
    '2024-12-19', '2025-12-26', 'Transport required for perishable goods from Sydney to Auckland. Must ensure cold chain integrity.', 0, 0,
    'cash', '', '', '', 0, 0, 'sender'),

-- Offer 16
(2, 2, 16, 16, 16, 'enabled', 230.00, 'USD', 11, 11, 12, 12, 'Mexico', 'Mexico City',
    'USA', 'Los Angeles', '123 Avenida Reforma, Mexico City, Mexico', '789 Santa Monica Blvd, Los Angeles, USA',
    '+52 55 1234-5678', '+1 310-555-2345', '+52 55 9876-5432',
    18, '2024-12-20', '2025-12-30',
    '2024-12-21', '2025-12-29', 'Transport of electronic components from Mexico City to Los Angeles. Requires secure handling.', 0, 0,
    'transfer', '', '', '', 0, 0, 'carrier'),

-- Offer 17
(3, 3, 17, 17, 17, 'enabled', 500.00, 'JPY', 13, 13, 14, 14, 'Japan', 'Tokyo',
    'South Korea', 'Seoul', '7 Shibuya, Tokyo, Japan', '15 Gangnam, Seoul, South Korea',
    '+81 3 1234-5678', '+82 2 555-1234', '+81 3 9876-5432',
    35, '2024-12-25', '2025-12-31',
    '2024-12-26', '2025-12-30', 'Shipping of machinery from Tokyo to Seoul. Must be delivered in time for a factory upgrade project.', 0, 0,
    'transfer', '', '', '', 0, 0, 'carrier'),

-- Offer 18
(4, 4, 18, 18, 18, 'enabled', 125.50, 'BRL', 17, 17, 18, 18, 'Brazil', 'Rio de Janeiro',
    'Argentina', 'Buenos Aires', '50 Copacabana Blvd, Rio de Janeiro, Brazil', '100 Av. 9 de Julio, Buenos Aires, Argentina',
    '+55 21 9876-5432', '+54 11 1234-5678', '+55 21 1234-9876',
    40, '2024-12-28', '2025-12-31',
    '2024-12-29', '2025-12-31', 'Transport of raw materials from Rio de Janeiro to Buenos Aires. Requires heavy-duty trucks.', 0, 0,
    'transfer', '', '', '', 0, 0, 'carrier'),

-- Offer 19
(1, 1, 19, 19, 19, 'enabled', 110.00, 'USD', 19, 19, 20, 20, 'Italy', 'Rome',
    'Germany', 'Berlin', '123 Via Roma, Rome, Italy', '500 Unter den Linden, Berlin, Germany',
    '+39 06 1234-5678', '+49 30 1234-5678', '+39 06 8765-4321',
    50, '2025-01-05', '2025-12-31',
    '2025-01-06', '2025-12-30', 'Shipment of industrial equipment from Rome to Berlin. Requires cranes for unloading.', 5, 0,
    'credit', '', '', '', 0, 0, 'carrier'),

-- Offer 20
(2, 2, 20, 20, 20, 'pending', 295.00, 'GBP', 1, 1, 2, 2, 'Spain', 'Madrid',
    'Portugal', 'Lisbon', '15 Gran Via, Madrid, Spain', '200 Avenida da Liberdade, Lisbon, Portugal',
    '+34 91 123 4567', '+351 21 123 4567', '+34 91 765 4321',
    60, '2025-01-10', '2025-12-15',
    '2025-01-12', '2025-12-14', 'Transport of hazardous materials from Madrid to Lisbon. Requires safety measures during transport.', 0, 0,
    'cash', '', '', '', 0, 0, 'sender');

INSERT INTO tbl_offer_response (
    company_id, to_company_id, offer_id, state, bid_price, title, note, reason
)
VALUES
    (2, 3, 1, 'declined', 1200.50, 'Not a fit for us', 'We appreciate the offer but won’t proceed.', 'Price is too high.'),
    (3, 2, 2, 'accepted', 1500.00, 'Happy to move forward', 'We accept the offer and look forward to working together.', ''),
    (1, 2, 2, 'pending', 1100.75, 'Need more time', 'We are reviewing the offer and will get back soon.', ''),
    (2, 3, 3, 'declined', 1300.00, 'Not within budget', 'The proposal is interesting, but we cannot afford it.', 'Too expensive.'),
    (3, 4, 4, 'declined', 1400.00, 'Not a strategic fit', 'We appreciate it, but it’s not aligned with our goals.', ''),
    (1, 3, 3, 'accepted', 1250.25, 'Let’s proceed', 'We accept and will begin discussions soon.', ''),
    (2, 3, 3, 'declined', 1350.00, 'Not at this time', 'We’re unable to proceed with this offer currently.', ''),
    (1, 2, 2, 'pending', 1180.00, 'Considering the offer', 'We need a few more days to review.', ''),
    (2, 1, 1, 'pending', 1225.00, 'Evaluating internally', 'We are discussing this with our team.', ''),
    (1, 2, 2, 'accepted', 1450.50, 'Great opportunity', 'Excited to move forward with this deal.', ''),
    (2, 1, 1, 'declined', 1150.00, 'Not suitable', 'Unfortunately, this offer doesn’t align with our needs.', ''),
    (1, 3, 3, 'declined', 1275.00, 'Different priorities', 'We are focusing on other opportunities at this time.', ''),
    (2, 3, 3, 'declined', 1290.00, 'Budget constraints', 'We can’t afford this deal within our budget.', ''),
    (3, 1, 1, 'pending', 1195.00, 'Awaiting approval', 'Our management team is reviewing this.', ''),
    (2, 1, 1, 'accepted', 1480.00, 'Excited to collaborate', 'We are happy with the terms and ready to proceed.', ''),
    (1, 3, 3, 'accepted', 1340.00, 'Looking forward', 'Let’s finalize the next steps soon.', ''),
    (2, 1, 1, 'pending', 1235.00, 'Considering options', 'We are comparing this with other offers.', ''),
    (3, 2, 2, 'declined', 1190.00, 'Not viable for us', 'This deal doesn’t align with our needs.', ''),
    (1, 2, 2, 'declined', 1260.00, 'Different direction', 'We are taking another approach to our strategy.', '');

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
