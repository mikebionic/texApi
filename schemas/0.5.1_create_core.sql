CREATE TYPE entity_t AS ENUM ('individual', 'legal');
CREATE TYPE role_t AS ENUM ('inactive','admin', 'sender', 'carrier');

CREATE TABLE tbl_role (
    id SERIAL PRIMARY KEY,
    role  role_t NOT NULL DEFAULT 'inactive',
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
    uuid                      UUID                                                     DEFAULT gen_random_uuid(),
    username                  VARCHAR(200) NOT NULL                                    DEFAULT '',
    password                  VARCHAR(200) NOT NULL                                    DEFAULT '',
    email                     VARCHAR(100) NOT NULL                                    DEFAULT '',
    info_email                VARCHAR(100) NOT NULL                                    DEFAULT '',
    first_name                VARCHAR(100) NOT NULL                                    DEFAULT '',
    last_name                 VARCHAR(100) NOT NULL                                    DEFAULT '',
    nick_name                 VARCHAR(100) NOT NULL                                    DEFAULT '',
    avatar_url                VARCHAR(200)                                             DEFAULT '',
    phone                     VARCHAR(100) NOT NULL                                    DEFAULT '',
    info_phone                VARCHAR(100) NOT NULL                                    DEFAULT '',
    address                   VARCHAR(200) NOT NULL                                    DEFAULT '',
    entity                    entity_t     NOT NULL                                    DEFAULT 'individual',
    role                        role_t NOT NULL DEFAULT 'inactive',
    role_id                   INT          REFERENCES tbl_role (id) ON DELETE SET NULL DEFAULT 0,
    verified                  INT          NOT NULL                                    DEFAULT 0,
    created_at                TIMESTAMP                                                DEFAULT CURRENT_TIMESTAMP,
    updated_at                TIMESTAMP                                                DEFAULT CURRENT_TIMESTAMP,
    active                    INT                                                      DEFAULT 1,
    deleted                   INT                                                      DEFAULT 0,

    oauth_provider            VARCHAR(100)                                             DEFAULT '',
    oauth_user_id             VARCHAR(100)                                             DEFAULT '',
    oauth_location            VARCHAR(200)                                             DEFAULT '',
    oauth_access_token        VARCHAR(500)                                             DEFAULT '',
    oauth_access_token_secret VARCHAR(500)                                             DEFAULT '',
    oauth_refresh_token       VARCHAR(500)                                             DEFAULT '',
    oauth_expires_at          TIMESTAMP                                                DEFAULT CURRENT_TIMESTAMP,
    oauth_id_token            VARCHAR(500)                                             DEFAULT '',
    refresh_token             VARCHAR(500)                                             DEFAULT '',
    otp_key                   VARCHAR(20)                                              DEFAULT '',
    verify_time               TIMESTAMP                                                DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE tbl_company
(
    id         SERIAL PRIMARY KEY,
    uuid       UUID                                                     DEFAULT gen_random_uuid(),
    user_id    INT          REFERENCES tbl_user (id) ON DELETE SET NULL DEFAULT 0,
    name       VARCHAR(100) NOT NULL                                    DEFAULT '',
    address    VARCHAR(200) NOT NULL                                    DEFAULT '',
    country    VARCHAR(200) NOT NULL                                    DEFAULT '',
    phone      VARCHAR(100) NOT NULL                                    DEFAULT '',
    email      VARCHAR(100) NOT NULL                                    DEFAULT '',
    logo_url   VARCHAR(200)                                             DEFAULT '',
    created_at TIMESTAMP                                                DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP                                                DEFAULT CURRENT_TIMESTAMP,
    active     INT                                                      DEFAULT 1,
    deleted    INT                                                      DEFAULT 0
);

-- TODO: can we save the foreign key with zero value?
CREATE TABLE tbl_driver
(
    id              SERIAL PRIMARY KEY,
    uuid            UUID                  DEFAULT gen_random_uuid(),
    company_id      INT                   REFERENCES tbl_company (id) ON DELETE SET NULL DEFAULT 0,
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
    company_id      INT                   REFERENCES tbl_company (id) ON DELETE SET NULL DEFAULT 0,
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

CREATE TABLE tbl_vehicle_brand (
    name               VARCHAR(100) PRIMARY KEY
);


INSERT INTO tbl_role (name, description, title, subtitle, title_ru, subtitle_ru) VALUES
    ('admin', 'Has full access to manage the system','','','',''),
    ('sender', 'Can place orders and track deliveries', 'Sender', 'I am looking for transport','Отправитель', 'Я ищу транспорт'),
    ('carrier_personal', 'Responsible for delivering orders using their personal vehicle', 'Carrier', 'Personal vehicle','Перевозчик','Личный автотранспорт'),
    ('carrier_owner', 'Responsible for delivering orders with a fleet of vehicles they own', 'Carrier', 'Fleet owner','Перевозчик','Владелец парка автотранспорта'),
    ('carrier_company', 'Responsible for delivering orders through a logistics company', 'Carrier', 'Logistics company','Перевозчик','Логистическая кампания');

INSERT INTO tbl_user (username, password, email, first_name, last_name, nick_name, avatar_url, phone, info_phone, address, role_id, verified, active, deleted)
VALUES
    ('root', 'PASSWORD_PLACEHOLDER', 'texlogistics@gmail.com', 'Tex', 'Admin', 'Texy', '', '+0036123456', '+0036123456', 'Ashgabat, Turkmenistan', 1, 1, 1, 0),
    ('customer1', 'PASSWORD_PLACEHOLDER', 'customer1@example.com', 'Customer', 'One', 'Custy', '', '+123456789', '+123456789', 'Ashgabat, Turkmenistan', 2, 1, 1, 0),
    ('driver1', 'PASSWORD_PLACEHOLDER', 'driver1@example.com', 'Volodya', '', 'Driver', '', '+123456789', '+123456789', 'Ashgabat, Turkmenistan', 3, 1, 1, 0);

-- Insert mock data into tbl_company
INSERT INTO tbl_company (user_id, name, address, phone, email, logo_url, active, deleted) VALUES
(1, 'Logistics Corp', '123 Main St, Cityville', '+1234567890', 'info@logisticscorp.com', 'http://example.com/logo1.png', 1, 0),
(2, 'Fast Movers', '456 Elm St, Townsville', '+0987654321', 'contact@fastmovers.com', 'http://example.com/logo2.png', 1, 0),
(3, 'Speedy Deliveries', '789 Oak St, Villageville', '+1122334455', 'support@speedydeliveries.com', 'http://example.com/logo3.png', 1, 0);


-- Insert mock data into tbl_driver (Each company has multiple drivers)
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
