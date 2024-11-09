CREATE TYPE entity_t AS ENUM ('individual', 'legal');
CREATE TYPE role_t AS ENUM ('system','admin','sender','carrier','unknown');
CREATE TYPE state_t AS ENUM ('enabled', 'disabled', 'deleted');


CREATE TABLE tbl_role (
    id SERIAL PRIMARY KEY,
    role  role_t NOT NULL DEFAULT 'unknown',
    uuid UUID DEFAULT gen_random_uuid(),
    name VARCHAR(100) UNIQUE NOT NULL,
    description TEXT DEFAULT '',
    title VARCHAR(100) NOT NULL DEFAULT '',
    title_ru VARCHAR(100) NOT NULL DEFAULT '',
    subtitle VARCHAR(200) NOT NULL DEFAULT '',
    subtitle_ru VARCHAR(200) NOT NULL DEFAULT ''
);


CREATE TABLE tbl_user
(
    id                        SERIAL PRIMARY KEY,
    uuid                      UUID                                           DEFAULT gen_random_uuid(),
    username                  VARCHAR(200) NOT NULL                          DEFAULT '',
    password                  VARCHAR(200) NOT NULL                          DEFAULT '',
    email                     VARCHAR(100) NOT NULL                          DEFAULT '',
    phone                     VARCHAR(100) NOT NULL                          DEFAULT '',
    role                      role_t       NOT NULL                          DEFAULT 'unknown',
    role_id                   INT REFERENCES tbl_role (id) ON DELETE CASCADE,
    company_id                INT          NOT NULL                          DEFAULT 0,
    verified                  INT          NOT NULL                          DEFAULT 0,
    created_at                TIMESTAMP                                      DEFAULT CURRENT_TIMESTAMP,
    updated_at                TIMESTAMP                                      DEFAULT CURRENT_TIMESTAMP,
    active                    INT                                            DEFAULT 1,
    deleted                   INT                                            DEFAULT 0,

    oauth_provider            VARCHAR(100) NOT NULL                          DEFAULT '',
    oauth_user_id             VARCHAR(100) NOT NULL                          DEFAULT '',
    oauth_location            VARCHAR(200) NOT NULL                          DEFAULT '',
    oauth_access_token        VARCHAR(500) NOT NULL                          DEFAULT '',
    oauth_access_token_secret VARCHAR(500) NOT NULL                          DEFAULT '',
    oauth_refresh_token       VARCHAR(500) NOT NULL                          DEFAULT '',
    oauth_expires_at          TIMESTAMP    NOT NULL                          DEFAULT CURRENT_TIMESTAMP,
    oauth_id_token            VARCHAR(500) NOT NULL                          DEFAULT '',
    refresh_token             VARCHAR(500) NOT NULL                          DEFAULT '',
    otp_key                   VARCHAR(20)  NOT NULL                          DEFAULT '',
    verify_time               TIMESTAMP    NOT NULL                          DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE tbl_company
(
    id              SERIAL PRIMARY KEY,
    uuid            UUID                                                             DEFAULT gen_random_uuid(),
    user_id         INT          NOT NULL REFERENCES tbl_user (id) ON DELETE CASCADE,
    role_id         INT          NOT NULL REFERENCES tbl_role (id) ON DELETE CASCADE,
    company_name    VARCHAR(100) NOT NULL                                            DEFAULT '',
    first_name      VARCHAR(100) NOT NULL                                            DEFAULT '',
    last_name       VARCHAR(100) NOT NULL                                            DEFAULT '',
    patronymic_name VARCHAR(100) NOT NULL                                            DEFAULT '',
    phone           VARCHAR(100) NOT NULL                                            DEFAULT '',
    phone2           VARCHAR(100) NOT NULL                                            DEFAULT '',
    phone3           VARCHAR(100) NOT NULL                                            DEFAULT '',
    email           VARCHAR(100) NOT NULL                                            DEFAULT '',
    email2          VARCHAR(100) NOT NULL                                            DEFAULT '',
    email3           VARCHAR(100) NOT NULL                                            DEFAULT '',
    meta           TEXT NOT NULL                                            DEFAULT '',
    meta2           TEXT NOT NULL                                            DEFAULT '',
    meta3           TEXT NOT NULL                                            DEFAULT '',
    address         VARCHAR(200) NOT NULL                                            DEFAULT '',
    country         VARCHAR(200) NOT NULL                                            DEFAULT '',
    image_url       VARCHAR(200)                                                     DEFAULT '',
    entity          entity_t     NOT NULL                                            DEFAULT 'individual',
    created_at      TIMESTAMP                                                        DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP                                                        DEFAULT CURRENT_TIMESTAMP,
    active          INT                                                              DEFAULT 1,
    deleted         INT                                                              DEFAULT 0
);


CREATE TABLE tbl_driver
(
    id              SERIAL PRIMARY KEY,
    uuid            UUID                  DEFAULT gen_random_uuid(),
    company_id      INT                 NOT NULL  REFERENCES tbl_company (id) ON DELETE CASCADE,
    first_name      VARCHAR(100) NOT NULL DEFAULT '',
    last_name       VARCHAR(100) NOT NULL DEFAULT '',
    patronymic_name VARCHAR(100) NOT NULL DEFAULT '',
    phone           VARCHAR(100) NOT NULL DEFAULT '',
    email           VARCHAR(100) NOT NULL DEFAULT '',
    avatar_url      VARCHAR(200)          DEFAULT '',
    created_at      TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,
    active          INT                   DEFAULT 1,
    deleted         INT                   DEFAULT 0
);

CREATE TABLE tbl_vehicle
(
    id                  SERIAL PRIMARY KEY,
    uuid                UUID                  DEFAULT gen_random_uuid(),
    company_id      INT              NOT NULL     REFERENCES tbl_company (id) ON DELETE CASCADE,
    vehicle_type        VARCHAR(100) NOT NULL DEFAULT '',
    brand               VARCHAR(100) NOT NULL DEFAULT '',
    vehicle_model       VARCHAR(100) NOT NULL DEFAULT '',
    year_of_issue       VARCHAR(10)  NOT NULL DEFAULT '',
    mileage             INT                   DEFAULT 0,
    numberplate         VARCHAR(20)  NOT NULL DEFAULT '',
    trailer_numberplate VARCHAR(20)  NOT NULL DEFAULT '',
    gps_active          INT                   DEFAULT 0,
    photo1_url          VARCHAR(200)          DEFAULT '',
    photo2_url          VARCHAR(200)          DEFAULT '',
    photo3_url          VARCHAR(200)          DEFAULT '',
    docs1_url           VARCHAR(200)          DEFAULT '',
    docs2_url           VARCHAR(200)          DEFAULT '',
    docs3_url           VARCHAR(200)          DEFAULT '',
    created_at          TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,
    updated_at          TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,
    active              INT                   DEFAULT 1,
    deleted             INT                   DEFAULT 0
);


INSERT INTO tbl_role (role, name, description, title, subtitle, title_ru, subtitle_ru) VALUES
   ('system', 'system', 'System level access','','','',''),
   ('admin', 'admin', 'Has full access to manage the system','','','',''),
   ('sender', 'sender', 'Can place orders and track deliveries', 'Sender', 'I am looking for transport','Отправитель', 'Я ищу транспорт'),
   ('carrier', 'carrier_personal', 'Responsible for delivering orders using their personal vehicle', 'Carrier', 'Personal vehicle','Перевозчик','Личный автотранспорт'),
   ('carrier', 'carrier_owner', 'Responsible for delivering orders with a fleet of vehicles they own', 'Carrier', 'Fleet owner','Перевозчик','Владелец парка автотранспорта'),
   ('carrier', 'carrier_company', 'Responsible for delivering orders through a logistics company', 'Carrier', 'Logistics company','Перевозчик','Логистическая кампания');

INSERT INTO tbl_user (username, password, email, phone, role, role_id, verified, active, deleted)
VALUES
    ('root', 'letmein', 'texlogistics@gmail.com', '+0036123456', 'admin', 2, 1, 1, 0),
    ('customer1', 'password123', 'customer1@example.com', '+123456789', 'sender', 3, 1, 1, 0),
    ('driver1', 'password123', 'driver1@example.com', '+123456789', 'carrier', 4, 1, 1, 0),
    ('sender_anna', 'pass123', 'anna.logistics@gmail.com', '+99365789001', 'sender', 3, 1, 1, 0),
    ('sender_mikhail', 'pass123', 'mikhail.trans@gmail.com', '+99365789002', 'sender', 3, 1, 1, 0),
    ('personal_ivan', 'pass123', 'ivan.driver@gmail.com', '+99365789003', 'carrier', 4, 1, 1, 0),
    ('personal_elena', 'pass123', 'elena.driver@gmail.com', '+99365789004', 'carrier', 4, 1, 1, 0),
    ('fleet_boris', 'pass123', 'boris.fleet@gmail.com', '+99365789005', 'carrier', 5, 1, 1, 0),
    ('fleet_dmitry', 'pass123', 'dmitry.fleet@gmail.com', '+99365789006', 'carrier', 5, 1, 1, 0),
    ('company_sergei', 'pass123', 'sergei.logistics@gmail.com', '+99365789007', 'carrier', 6, 1, 1, 0),
    ('company_natalia', 'pass123', 'natalia.logistics@gmail.com', '+99365789008', 'carrier', 6, 1, 1, 0),
    ('sender_alex', 'pass123', 'alex.cargo@gmail.com', '+99365789009', 'sender', 3, 1, 1, 0),
    ('personal_maria', 'pass123', 'maria.driver@gmail.com', '+99365789010', 'carrier', 4, 1, 1, 0);

INSERT INTO tbl_company (
    user_id,
    role_id,
    company_name,
    first_name,
    last_name,
    patronymic_name,
    address,
    country,
    phone,
    email,
    image_url,
    entity,
    active,
    deleted
) VALUES
    (1, 2, 'Logistics Corp', 'Tex', 'Admin', '', '123 Main St, Cityville', 'Turkmenistan', '+1234567890', 'info@logisticscorp.com', 'https://images.unsplash.com/photo-1635070041078-e363dbe005cb', 'legal', 1, 0),
    (2, 3, 'Fast Movers', 'Customer', 'One', '', '456 Elm St, Townsville', 'Turkmenistan', '+0987654321', 'contact@fastmovers.com', 'https://images.unsplash.com/photo-1635070041078-e363dbe005cb', 'legal', 1, 0),
    (3, 4, 'Speedy Deliveries', 'Volodya', 'Driver', '', '789 Oak St, Villageville', 'Turkmenistan', '+1122334455', 'support@speedydeliveries.com', 'https://images.unsplash.com/photo-1635070041078-e363dbe005cb', 'individual', 1, 0),
    (4, 3, 'Anna Logistics Solutions', 'Anna', 'Petrova', 'Mikhailovna', 'Magtymguly avenue 142, Ashgabat', 'Turkmenistan', '+99365789001', 'anna.logistics@gmail.com', 'https://images.unsplash.com/photo-1560179707-f14e90ef3623', 'legal', 1, 0),
    (5, 3, 'Mikhail Transit Hub', 'Mikhail', 'Ivanov', 'Sergeevich', 'Andaliba street 54, Ashgabat', 'Turkmenistan', '+99365789002', 'mikhail.trans@gmail.com', 'https://images.unsplash.com/photo-1623259838743-9f1e884fba89', 'legal', 1, 0),
    (6, 4, 'Personal Delivery Service', 'Ivan', 'Smirnov', 'Alexandrovich', 'Garashsyzlyk avenue 32, Ashgabat', 'Turkmenistan', '+99365789003', 'ivan.driver@gmail.com', 'https://images.unsplash.com/photo-1601628828688-632f38a5a7d0', 'individual', 1, 0),
    (7, 4, 'Elena Express', 'Elena', 'Volkova', 'Dmitrievna', 'Atamurat Niyazov street 75, Ashgabat', 'Turkmenistan', '+99365789004', 'elena.driver@gmail.com', 'https://images.unsplash.com/photo-1554768804-50c1e2b50a6e', 'individual', 1, 0),
    (8, 5, 'Boris Fleet Management', 'Boris', 'Kuznetsov', 'Ivanovich', 'Oguzhan street 127, Ashgabat', 'Turkmenistan', '+99365789005', 'boris.fleet@gmail.com', 'https://images.unsplash.com/photo-1586528116311-ad8dd3c8310d', 'legal', 1, 0),
    (9, 5, 'Dmitry Transportation Co', 'Dmitry', 'Sokolov', 'Petrovich', 'Yunus Emre street 89, Ashgabat', 'Turkmenistan', '+99365789006', 'dmitry.fleet@gmail.com', 'https://images.unsplash.com/photo-1570449942860-bb66578b6e69', 'legal', 1, 0),
    (10, 6, 'Sergei Logistics Group', 'Sergei', 'Popov', 'Mikhailovich', 'Gorogly street 234, Ashgabat', 'Turkmenistan', '+99365789007', 'sergei.logistics@gmail.com', 'https://images.unsplash.com/photo-1566576912321-d58ddd7a6088', 'legal', 1, 0),
    (11, 6, 'Natalia Cargo Systems', 'Natalia', 'Morozova', 'Andreevna', 'A. Niyazov street 156, Ashgabat', 'Turkmenistan', '+99365789008', 'natalia.logistics@gmail.com', 'https://images.unsplash.com/photo-1635070041078-e363dbe005cb', 'legal', 1, 0),
    (12, 3, 'Alex Cargo Solutions', 'Alexander', 'Lebedev', 'Vladimirovich', 'Bitarap street 67, Ashgabat', 'Turkmenistan', '+99365789009', 'alex.cargo@gmail.com', 'https://images.unsplash.com/photo-1542744173-8e7e53415bb0', 'legal', 1, 0),
    (13, 4, 'Maria Express Delivery', 'Maria', 'Kozlova', 'Sergeevna', 'Magtymguly avenue 198, Ashgabat', 'Turkmenistan', '+99365789010', 'maria.driver@gmail.com', 'https://images.unsplash.com/photo-1624459294159-598de7ddcae9', 'individual', 1, 0);

UPDATE tbl_user SET company_id = 1 WHERE id = 1;
UPDATE tbl_user SET company_id = 2 WHERE id = 2;
UPDATE tbl_user SET company_id = 3 WHERE id = 3;
UPDATE tbl_user SET company_id = 4 WHERE id = 4;
UPDATE tbl_user SET company_id = 5 WHERE id = 5;
UPDATE tbl_user SET company_id = 6 WHERE id = 6;
UPDATE tbl_user SET company_id = 7 WHERE id = 7;
UPDATE tbl_user SET company_id = 8 WHERE id = 8;
UPDATE tbl_user SET company_id = 9 WHERE id = 9;
UPDATE tbl_user SET company_id = 10 WHERE id = 10;
UPDATE tbl_user SET company_id = 11 WHERE id = 11;
UPDATE tbl_user SET company_id = 12 WHERE id = 12;
UPDATE tbl_user SET company_id = 13 WHERE id = 13;


INSERT INTO tbl_driver (company_id, first_name, last_name, patronymic_name, phone, email, avatar_url, active, deleted) VALUES
(1, 'John', 'Doe', 'Smith', '+1234567890', 'john.doe@logisticscorp.com', 'http://example.com/avatar1.png', 1, 0),
(1, 'Jane', 'Doe', 'Johnson', '+0987654321', 'jane.doe@logisticscorp.com', 'http://example.com/avatar2.png', 1, 0),
(2, 'Michael', 'Brown', 'Williams', '+1122334455', 'michael.brown@fastmovers.com', 'http://example.com/avatar3.png', 1, 0),
(2, 'Emily', 'Clark', 'Davis', '+2233445566', 'emily.clark@fastmovers.com', 'http://example.com/avatar4.png', 1, 0),
(3, 'Robert', 'Lee', 'Martin', '+3344556677', 'robert.lee@speedydeliveries.com', 'http://example.com/avatar5.png', 1, 0),
(3, 'Anna', 'Taylor', 'Thompson', '+4455667788', 'anna.taylor@speedydeliveries.com', 'http://example.com/avatar6.png', 1, 0);


-- Insert mock data into tbl_vehicle (Each company has multiple vehicles)
-- Drivers can have multiple vehicles, or vehicles can be unassigned.
INSERT INTO tbl_vehicle (company_id, vehicle_type, brand, vehicle_model, year_of_issue, mileage, numberplate, trailer_numberplate, gps_active, photo1_url, photo2_url, photo3_url, docs1_url, docs2_url, docs3_url, active, deleted) VALUES
(1, 'Truck', 'Volvo', 'FH16', '2019', 1023954, 'ABC123', 'TRAIL123', 1, 'https://img.linemedia.com/img/s/dump-truck-Volvo-FH16-750-6x4-Retarder-Full-Steel---1730104405751076482_big--24102810295850812600.jpg', 'https://img.linemedia.com/img/s/dump-truck-Volvo-FH16-750-6x4-Retarder-Full-Steel---1730104406528835478_big--24102810295850812600.jpg', 'https://img.linemedia.com/img/s/dump-truck-Volvo-FH16-750-6x4-Retarder-Full-Steel---1730104404137616440_big--24102810295850812600.jpg', 'http://example.com/vehicle1_docs1.pdf', '', '', 1, 0),
(1, 'Van', 'Mercedes', 'Sprinter', '2020', 1022234,'XYZ456', 'TRAIL456', 1, 'https://img.linemedia.com/img/s/coach-bus-Mercedes-Benz-Sprinter-518---1729426654848855287_big--24102015102897570700.jpg', 'https://img.linemedia.com/img/s/coach-bus-Mercedes-Benz-Sprinter-518---1729426656178478265_big--24102015102897570700.jpg', 'https://img.linemedia.com/img/s/coach-bus-Mercedes-Benz-Sprinter-518---1729426657288003993_big--24102015102897570700.jpg', 'http://example.com/vehicle2_docs1.pdf', '', '', 1, 0),
(2, 'Truck', 'MAN', 'TGS', '2018', 23954, 'LMN789', 'TRAIL789', 0, 'https://img.linemedia.com/img/s/forestry-equipment-wood-chipper-Jenz-MAN-TGS-33-500-HEM-583-R-Palfinger-Epsilon-S110F101---1721826471689124800_big--24072415525385274700.jpg', 'https://img.linemedia.com/img/s/forestry-equipment-wood-chipper-Jenz-MAN-TGS-33-500-HEM-583-R-Palfinger-Epsilon-S110F101---1721826472190064857_big--24072415525385274700.jpg', 'https://img.linemedia.com/img/s/forestry-equipment-wood-chipper-Jenz-MAN-TGS-33-500-HEM-583-R-Palfinger-Epsilon-S110F101---1721826472610252577_big--24072415525385274700.jpg', 'http://example.com/vehicle3_docs1.pdf', '', '', 1, 0),
(2, 'Truck', 'Scania', 'R500', '2017', 96954,'GHI321', '', 1, 'http://example.com/vehicle4_photo1.png', '', '', 'http://example.com/vehicle4_docs1.pdf', '', '', 1, 0),
(3, 'Truck', 'DAF', 'XF', '2021', 403954, 'JKL654', 'TRAIL654', 1, 'http://example.com/vehicle5_photo1.png', 'http://example.com/vehicle5_photo2.png', '', 'http://example.com/vehicle5_docs1.pdf', '', '', 1, 0),
(3, 'Van', 'Ford', 'Transit', '2022',  53954, 'MNO987', '', 0, 'http://example.com/vehicle6_photo1.png', '', '', 'http://example.com/vehicle6_docs1.pdf', '', '', 1, 0);
