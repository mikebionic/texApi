CREATE TYPE entity_t AS ENUM ('individual', 'legal');
CREATE TYPE role_t AS ENUM ('system','admin','sender','carrier','unknown');
CREATE TYPE state_t AS ENUM ('enabled', 'disabled', 'deleted', 'pending', 'archived', 'working');


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
    uuid                      UUID                  DEFAULT gen_random_uuid(),
    username                  VARCHAR(200) NOT NULL DEFAULT '',
    password                  VARCHAR(200) NOT NULL DEFAULT '',
    email                     VARCHAR(100) NOT NULL DEFAULT '',
    phone                     VARCHAR(100) NOT NULL DEFAULT '',
    role                      role_t       NOT NULL DEFAULT 'unknown',
    role_id                   INT REFERENCES tbl_role (id) ON DELETE CASCADE,
    company_id                INT          NOT NULL DEFAULT 0,
    verified                  INT          NOT NULL DEFAULT 0,
    meta                      TEXT         NOT NULL DEFAULT '',
    meta2                     TEXT         NOT NULL DEFAULT '',
    meta3                     TEXT         NOT NULL DEFAULT '',
    created_at                TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,
    updated_at                TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,
    active                    INT                   DEFAULT 1,
    deleted                   INT                   DEFAULT 0,

    oauth_provider            VARCHAR(100) NOT NULL DEFAULT '',
    oauth_user_id             VARCHAR(100) NOT NULL DEFAULT '',
    oauth_location            VARCHAR(200) NOT NULL DEFAULT '',
    oauth_access_token        VARCHAR(500) NOT NULL DEFAULT '',
    oauth_access_token_secret VARCHAR(500) NOT NULL DEFAULT '',
    oauth_refresh_token       VARCHAR(500) NOT NULL DEFAULT '',
    oauth_expires_at          TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    oauth_id_token            VARCHAR(500) NOT NULL DEFAULT '',
    refresh_token             VARCHAR(500) NOT NULL DEFAULT '',
    otp_key                   VARCHAR(20)  NOT NULL DEFAULT '',
    verify_time               TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE tbl_company
(
    id              SERIAL PRIMARY KEY,
    uuid            UUID                  DEFAULT gen_random_uuid(),
    user_id         INT          NOT NULL REFERENCES tbl_user (id) ON DELETE CASCADE,
    role_id         INT          NOT NULL REFERENCES tbl_role (id) ON DELETE CASCADE,
    company_name    VARCHAR(100) NOT NULL DEFAULT '',
    first_name      VARCHAR(100) NOT NULL DEFAULT '',
    last_name       VARCHAR(100) NOT NULL DEFAULT '',
    patronymic_name VARCHAR(100) NOT NULL DEFAULT '',
    phone           VARCHAR(100) NOT NULL DEFAULT '',
    phone2          VARCHAR(100) NOT NULL DEFAULT '',
    phone3          VARCHAR(100) NOT NULL DEFAULT '',
    email           VARCHAR(100) NOT NULL DEFAULT '',
    email2          VARCHAR(100) NOT NULL DEFAULT '',
    email3          VARCHAR(100) NOT NULL DEFAULT '',
    meta            TEXT         NOT NULL DEFAULT '',
    meta2           TEXT         NOT NULL DEFAULT '',
    meta3           TEXT         NOT NULL DEFAULT '',
    address         VARCHAR(200) NOT NULL DEFAULT '',
    country         VARCHAR(200) NOT NULL DEFAULT '',
    country_id      INT          NOT NULL DEFAULT 0,
    city_id         INT          NOT NULL DEFAULT 0,
    image_url       VARCHAR(200) NOT NULL DEFAULT '',
    entity          entity_t     NOT NULL DEFAULT 'individual',
    featured        INT          NOT NULL DEFAULT 0,
    rating          INT          NOT NULL DEFAULT 0,
    partner         INT          NOT NULL DEFAULT 0,
    successful_ops  INT          NOT NULL DEFAULT 0,
    last_active     TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,
    created_at      TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,
    active          INT          NOT NULL DEFAULT 1,
    deleted         INT          NOT NULL DEFAULT 0
);


CREATE TABLE tbl_driver
(
    id              SERIAL PRIMARY KEY,
    uuid            UUID                  DEFAULT gen_random_uuid(),
    company_id      INT          NOT NULL REFERENCES tbl_company (id) ON DELETE CASCADE,
    first_name      VARCHAR(100) NOT NULL DEFAULT '',
    last_name       VARCHAR(100) NOT NULL DEFAULT '',
    patronymic_name VARCHAR(100) NOT NULL DEFAULT '',
    phone           VARCHAR(100) NOT NULL DEFAULT '',
    email           VARCHAR(100) NOT NULL DEFAULT '',
    featured        INT          NOT NULL DEFAULT 0,
    rating          INT          NOT NULL DEFAULT 0,
    partner         INT          NOT NULL DEFAULT 0,
    successful_ops  INT          NOT NULL DEFAULT 0,
    image_url       VARCHAR(200) NOT NULL DEFAULT '',
    meta            TEXT         NOT NULL DEFAULT '',
    meta2           TEXT         NOT NULL DEFAULT '',
    meta3           TEXT         NOT NULL DEFAULT '',
    available       INT          NOT NULL DEFAULT 1,
    view_count      INT          NOT NULL DEFAULT 0,
    created_at      TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,
    active          INT          NOT NULL DEFAULT 1,
    deleted         INT          NOT NULL DEFAULT 0
);





CREATE TABLE tbl_packaging_type
(
    id             SERIAL PRIMARY KEY,
    name_ru        VARCHAR(255)   NOT NULL DEFAULT '',
    name_en        VARCHAR(255)   NOT NULL DEFAULT '',
    name_tk        VARCHAR(255)   NOT NULL DEFAULT '',
    category_ru    VARCHAR(255)   NOT NULL DEFAULT '',
    category_en    VARCHAR(255)   NOT NULL DEFAULT '',
    category_tk    VARCHAR(255)   NOT NULL DEFAULT '',
    material       VARCHAR(255)   NOT NULL DEFAULT '',
    dimensions     VARCHAR(255)   NOT NULL DEFAULT '',
    weight         DECIMAL(10, 2) NOT NULL DEFAULT 0.0,
    description_ru TEXT           NOT NULL DEFAULT '',
    description_en TEXT           NOT NULL DEFAULT '',
    description_tk TEXT           NOT NULL DEFAULT '',
    active         INT            NOT NULL DEFAULT 0,
    deleted        INT            NOT NULL DEFAULT 0
);

INSERT INTO tbl_packaging_type (name_ru, name_en, category_ru, category_en, material, dimensions, weight, description_ru, description_en)
VALUES
    ('Картонная коробка', 'Cardboard Box', 'Вторичная упаковка', 'Secondary Packaging', 'Картон', '30x30x30 см', 1.50, 'Общепринятая упаковка для товаров, таких как электроника, книги и мелкие предметы.', 'Common packaging for electronics, books, and small items'),
    ('Деревянный ящик', 'Wooden Crate', 'Третичная упаковка', 'Tertiary Packaging', 'Дерево', '120x120x120 см', 15.00, 'Мощная деревянная упаковка для перевозки крупногабаритного оборудования или материалов.', 'Heavy-duty wooden crates for large machinery or equipment'),
    ('Пленка Shrink', 'Shrink Wrap', 'Вторичная упаковка', 'Secondary Packaging', 'Пластик', 'N/A', 0.10, 'Пластиковая пленка, которая оборачивает товары, создавая защиту для транспортировки.', 'Plastic film used to wrap around goods for secure transport'),
    ('Палета', 'Pallet', 'Третичная упаковка', 'Tertiary Packaging', 'Дерево', '120x80 см', 25.00, 'Палета для укладки товаров для транспортировки с использованием погрузчика.', 'Pallet used to stack goods for easier transport with a forklift'),
    ('Тетра Пак', 'Tetra Pak', 'Первичная упаковка', 'Primary Packaging', 'Картон', '200x150x100 мм', 0.25, 'Упаковка для жидких продуктов, таких как молоко или сок.', 'Packaging used for liquid products like milk or juice'),
    ('Картонная коробка с клапаном', 'Flap Carton Box', 'Вторичная упаковка', 'Secondary Packaging', 'Картон', '40x40x40 см', 3.00, 'Коробка с клапаном для упаковки средних товаров.', 'Flap carton box used for medium-sized goods packaging'),
    ('Гофрированный картон', 'Corrugated Cardboard', 'Вторичная упаковка', 'Secondary Packaging', 'Картон', 'N/A', 0.50, 'Гофрированный картон, использующийся для упаковки хрупких товаров.', 'Corrugated cardboard used for packing fragile goods'),
    ('Картонный контейнер', 'Cardboard Container', 'Третичная упаковка', 'Tertiary Packaging', 'Картон', '120x80 см', 10.00, 'Контейнеры из картона для транспортировки больших объемов товаров.', 'Cardboard containers used for bulk goods transport'),
    ('Мешок', 'Bag', 'Первичная упаковка', 'Primary Packaging', 'Пластик', '30x40 см', 0.15, 'Мешки для упаковки сыпучих товаров, таких как зерно, уголь и порошки.', 'Bags used for packaging bulk goods such as grain, coal, and powders'),
    ('Пластиковая бутылка', 'Plastic Bottle', 'Первичная упаковка', 'Primary Packaging', 'Пластик', '1 литр', 0.25, 'Пластиковая бутылка для напитков и жидких продуктов.', 'Plastic bottle used for beverages and liquid products'),
    ('Пластиковый контейнер', 'Plastic Container', 'Первичная упаковка', 'Primary Packaging', 'Пластик', '500 мл', 0.20, 'Пластиковый контейнер для упаковки продуктов питания или бытовых товаров.', 'Plastic container used for food or household goods'),
    ('Металлическая банка', 'Metal Can', 'Первичная упаковка', 'Primary Packaging', 'Металл', '500 мл', 0.30, 'Металлическая банка для упаковки напитков или консервированных продуктов.', 'Metal can used for packaging beverages or canned food'),
    ('Стеклянная бутылка', 'Glass Bottle', 'Первичная упаковка', 'Primary Packaging', 'Стекло', '1 литр', 0.45, 'Стеклянная бутылка для упаковки напитков, таких как соки и вино.', 'Glass bottle used for packaging beverages like juices and wine'),
    ('Вакуумная упаковка', 'Vacuum Packaging', 'Первичная упаковка', 'Primary Packaging', 'Пластик', 'N/A', 0.30, 'Упаковка, в которой удален воздух, используемая для хранения продуктов или товаров.', 'Packaging where air is removed, used for storing food or goods'),
    ('Упаковка с регулируемой атмосферой', 'Modified Atmosphere Packaging (MAP)', 'Первичная упаковка', 'Primary Packaging', 'Пластик', 'N/A', 0.50, 'Упаковка с контролируемым составом воздуха для продления срока хранения продуктов.', 'Packaging with controlled air composition for extended product shelf life'),
    ('Стретч-пленка', 'Stretch Film', 'Вторичная упаковка', 'Secondary Packaging', 'Пластик', 'N/A', 0.10, 'Пленка, обвивающая паллеты для их закрепления и защиты в процессе транспортировки.', 'Film used to wrap around pallets for securing and protecting goods during transport'),
    ('Пластиковая пленка', 'Plastic Wrap', 'Вторичная упаковка', 'Secondary Packaging', 'Пластик', 'N/A', 0.05, 'Пленка для упаковки продуктов или товаров в небольших количествах.', 'Film used for wrapping small quantities of products or goods'),
    ('Обертка для продукции', 'Product Wrap', 'Вторичная упаковка', 'Secondary Packaging', 'Пластик', 'N/A', 0.10, 'Пленка или бумага, используемая для упаковки отдельных единиц продукции.', 'Plastic or paper wrap used for individual product packaging'),
    ('Деревянный контейнер', 'Wooden Container', 'Третичная упаковка', 'Tertiary Packaging', 'Дерево', '150x120x100 см', 20.00, 'Деревянный контейнер для транспортировки крупногабаритных и тяжёлых товаров.', 'Wooden container used for transporting oversized and heavy goods'),
    ('Металлический контейнер', 'Metal Container', 'Третичная упаковка', 'Tertiary Packaging', 'Металл', '200x150x150 см', 30.00, 'Металлические контейнеры для транспортировки опасных грузов или химикатов.', 'Metal containers used for transporting hazardous goods or chemicals'),
    ('Изотермическая упаковка', 'Isothermal Packaging', 'Первичная упаковка', 'Primary Packaging', 'Пластик/Термопласт', '30x30x30 см', 0.80, 'Упаковка, поддерживающая температуру для чувствительных к температуре продуктов, таких как медикаменты или еда.', 'Packaging that maintains temperature for temperature-sensitive products like medications or food'),
    ('Пакет с клапаном', 'Valve Bag', 'Первичная упаковка', 'Primary Packaging', 'Пластик', '50x50 см', 0.50, 'Мешок с клапаном для упаковки порошков или сыпучих товаров, таких как цемент или химикаты.', 'Bag with a valve for packaging powders or bulk goods like cement or chemicals'),
    ('Картонная коробка для электроники', 'Electronics Cardboard Box', 'Вторичная упаковка', 'Secondary Packaging', 'Картон', '25x25x25 см', 1.00, 'Коробка, используемая для упаковки электроники, такой как телевизоры или компьютеры.', 'Cardboard box used for packaging electronics like televisions or computers');



CREATE TABLE tbl_vehicle
(
    id                  SERIAL PRIMARY KEY,
    uuid                UUID                                                    DEFAULT gen_random_uuid(),
    company_id          INT          NOT NULL REFERENCES tbl_company (id) ON DELETE CASCADE,
    vehicle_type_id     INT REFERENCES tbl_vehicle_type (id) ON DELETE CASCADE  DEFAULT 1,
    vehicle_brand_id    INT REFERENCES tbl_vehicle_brand (id) ON DELETE CASCADE DEFAULT 1,
    vehicle_model_id    INT REFERENCES tbl_vehicle_model (id) ON DELETE CASCADE DEFAULT 1,
    year_of_issue       VARCHAR(10)  NOT NULL                                   DEFAULT '',
    mileage             INT          NOT NULL                                   DEFAULT 0,
    numberplate         VARCHAR(20)  NOT NULL                                   DEFAULT '',
    trailer_numberplate VARCHAR(20)  NOT NULL                                   DEFAULT '',
    gps                 INT          NOT NULL                                   DEFAULT 0,
    photo1_url          VARCHAR(200) NOT NULL                                   DEFAULT '',
    photo2_url          VARCHAR(200) NOT NULL                                   DEFAULT '',
    photo3_url          VARCHAR(200) NOT NULL                                   DEFAULT '',
    docs1_url           VARCHAR(200) NOT NULL                                   DEFAULT '',
    docs2_url           VARCHAR(200) NOT NULL                                   DEFAULT '',
    docs3_url           VARCHAR(200) NOT NULL                                   DEFAULT '',
    view_count          INT          NOT NULL                                   DEFAULT 0,
    meta                TEXT         NOT NULL                                   DEFAULT '',
    meta2               TEXT         NOT NULL                                   DEFAULT '',
    meta3               TEXT         NOT NULL                                   DEFAULT '',
    available           INT          NOT NULL                                   DEFAULT 1,
    created_at          TIMESTAMP    NOT NULL                                   DEFAULT CURRENT_TIMESTAMP,
    updated_at          TIMESTAMP    NOT NULL                                   DEFAULT CURRENT_TIMESTAMP,
    active              INT          NOT NULL                                   DEFAULT 1,
    deleted             INT          NOT NULL                                   DEFAULT 0
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
    ('root', 'PASSWORD_PLACEHOLDER', 'texlogistics@gmail.com', '+0036123456', 'admin', 2, 1, 1, 0),
    ('customer1', 'PASSWORD_PLACEHOLDER', 'customer1@example.com', '+123456789', 'sender', 3, 1, 1, 0),
    ('driver1', 'PASSWORD_PLACEHOLDER', 'driver1@example.com', '+123456789', 'carrier', 4, 1, 1, 0),
    ('sender_anna', 'PASSWORD_PLACEHOLDER', 'anna.logistics@gmail.com', '+99365789001', 'sender', 3, 1, 1, 0),
    ('sender_mikhail', 'PASSWORD_PLACEHOLDER', 'mikhail.trans@gmail.com', '+99365789002', 'sender', 3, 1, 1, 0),
    ('personal_ivan', 'PASSWORD_PLACEHOLDER', 'ivan.driver@gmail.com', '+99365789003', 'carrier', 4, 1, 1, 0),
    ('personal_elena', 'PASSWORD_PLACEHOLDER', 'elena.driver@gmail.com', '+99365789004', 'carrier', 4, 1, 1, 0),
    ('fleet_boris', 'PASSWORD_PLACEHOLDER', 'boris.fleet@gmail.com', '+99365789005', 'carrier', 5, 1, 1, 0),
    ('fleet_dmitry', 'PASSWORD_PLACEHOLDER', 'dmitry.fleet@gmail.com', '+99365789006', 'carrier', 5, 1, 1, 0),
    ('company_sergei', 'PASSWORD_PLACEHOLDER', 'sergei.logistics@gmail.com', '+99365789007', 'carrier', 6, 1, 1, 0),
    ('company_natalia', 'PASSWORD_PLACEHOLDER', 'natalia.logistics@gmail.com', '+99365789008', 'carrier', 6, 1, 1, 0),
    ('sender_alex', 'PASSWORD_PLACEHOLDER', 'alex.cargo@gmail.com', '+99365789009', 'sender', 3, 1, 1, 0),
    ('personal_maria', 'PASSWORD_PLACEHOLDER', 'maria.driver@gmail.com', '+99365789010', 'carrier', 4, 1, 1, 0),
    -- sender (2)
    ('sender_hz', 'PASSWORD_PLACEHOLDER', 'hz.sender@gmail.com', '+99365789001', 'sender', 3, 1, 1, 0),
    ('sender_sw',  'PASSWORD_PLACEHOLDER', 'sw.sender@gmail.com',  '+99365789002', 'sender', 3, 1, 1, 0),
    -- personal (1)
    ('personal_mn',  'PASSWORD_PLACEHOLDER', 'mn.personal@gmail.com',  '+99365789003', 'carrier', 4, 1, 1, 0),
    -- fleet (1)
    ('fleet_dr',    'PASSWORD_PLACEHOLDER', 'dr.fleet@gmail.com',    '+99365789004', 'carrier', 5, 1, 1, 0),
    -- logistics (1)
    ('logistics_dl','PASSWORD_PLACEHOLDER', 'dl.logistics@gmail.com','+99365789005', 'carrier',6, 1, 1, 0);


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
    (13, 4, 'Maria Express Delivery', 'Maria', 'Kozlova', 'Sergeevna', 'Magtymguly avenue 198, Ashgabat', 'Turkmenistan', '+99365789010', 'maria.driver@gmail.com', 'https://images.unsplash.com/photo-1624459294159-598de7ddcae9', 'individual', 1, 0),
    -- sender (2)
    (14, 3, 'hz.sender', 'hz', 'sender', '', 'Gorogly street 234, Ashgabat', 'Turkmenistan','+99365789001','hz.sender@gmail.com', 'https://images.unsplash.com/photo-1624459294159-598de7ddcae9', 'legal', 1, 0),
    (15, 3, 'sw.sender',  'sw', 'sender', '', 'Yunus Emre street 89, Ashgabat', 'Turkmenistan','+99365789002','sw.sender@gmail.com', 'https://images.unsplash.com/photo-1624459294159-598de7ddcae9', 'individual', 1, 0),
    -- personal (1)
    (16, 4, 'mn.personal',  'mn', 'personal', '', 'Yunus Emre street 89, Ashgabat', 'Turkmenistan','+99365789003','mn.personal@gmail.com', 'https://images.unsplash.com/photo-1624459294159-598de7ddcae9', 'individual', 1, 0),
    -- fleet (1)
    (17, 5, 'dr.fleet',    'dr', 'fleet', '', 'Gorogly street 234, Ashgabat',    'Turkmenistan','+99365789004','dr.fleet@gmail.com', 'https://images.unsplash.com/photo-1624459294159-598de7ddcae9', 'legal', 1, 0),
    -- logistics (1)
    (18, 6, 'dl.logistics','dl', 'logistics', '', 'Atamurat Niyazov street 75, Ashgabat','Turkmenistan','+99365789005','dl.logistics@gmail.com', 'https://images.unsplash.com/photo-1624459294159-598de7ddcae9', 'legal', 1, 0);


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

UPDATE tbl_user SET company_id = 14 WHERE id = 14;
UPDATE tbl_user SET company_id = 15 WHERE id = 15;
UPDATE tbl_user SET company_id = 16 WHERE id = 16;
UPDATE tbl_user SET company_id = 17 WHERE id = 17;
UPDATE tbl_user SET company_id = 18 WHERE id = 18;


INSERT INTO tbl_driver (company_id, first_name, last_name, patronymic_name, phone, email, image_url, active, deleted) VALUES
(1, 'John', 'Doe', 'Smith', '+1234567890', 'john.doe@logisticscorp.com', 'http://example.com/avatar1.png', 1, 0),
(1, 'Jane', 'Doe', 'Johnson', '+0987654321', 'jane.doe@logisticscorp.com', 'http://example.com/avatar2.png', 1, 0),
(2, 'Michael', 'Brown', 'Williams', '+1122334455', 'michael.brown@fastmovers.com', 'http://example.com/avatar3.png', 1, 0),
(2, 'Emily', 'Clark', 'Davis', '+2233445566', 'emily.clark@fastmovers.com', 'http://example.com/avatar4.png', 1, 0),
(3, 'Robert', 'Lee', 'Martin', '+3344556677', 'robert.lee@speedydeliveries.com', 'http://example.com/avatar5.png', 1, 0),
(3, 'Anna', 'Taylor', 'Thompson', '+4455667788', 'anna.taylor@speedydeliveries.com', 'http://example.com/avatar6.png', 1, 0);


INSERT INTO tbl_vehicle (company_id, vehicle_type_id, vehicle_brand_id, vehicle_model_id, year_of_issue, mileage, numberplate, trailer_numberplate, gps, photo1_url, photo2_url, photo3_url, docs1_url, docs2_url, docs3_url, active, deleted) VALUES
(1, 1, 4,12, '2019', 1023954, 'ABC123', 'TRAIL123', 1, 'https://img.linemedia.com/img/s/dump-truck-Volvo-FH16-750-6x4-Retarder-Full-Steel---1730104405751076482_big--24102810295850812600.jpg', 'https://img.linemedia.com/img/s/dump-truck-Volvo-FH16-750-6x4-Retarder-Full-Steel---1730104406528835478_big--24102810295850812600.jpg', 'https://img.linemedia.com/img/s/dump-truck-Volvo-FH16-750-6x4-Retarder-Full-Steel---1730104404137616440_big--24102810295850812600.jpg', 'http://example.com/vehicle1_docs1.pdf', '', '', 1, 0),
(1, 5, 3, 45, '2020', 1022234,'XYZ456', 'TRAIL456', 1, 'https://img.linemedia.com/img/s/coach-bus-Mercehttpdes-Benz-Sprinter-518---1729426654848855287_big--24102015102897570700.jpg', 'https://img.linemedia.com/img/s/coach-bus-Mercedes-Benz-Sprinter-518---1729426656178478265_big--24102015102897570700.jpg', 'https://img.linemedia.com/img/s/coach-bus-Mercedes-Benz-Sprinter-518---1729426657288003993_big--24102015102897570700.jpg', 'http://example.com/vehicle2_docs1.pdf', '', '', 1, 0),
(2, 2, 5, 66, '2018', 23954, 'LMN789', 'TRAIL789', 0, 'https://img.linemedia.com/img/s/forestry-equipment-wood-chipper-Jenz-MAN-TGS-33-500-HEM-583-R-Palfinger-Epsilon-S110F101---1721826471689124800_big--24072415525385274700.jpg', 'https://img.linemedia.com/img/s/forestry-equipment-wood-chipper-Jenz-MAN-TGS-33-500-HEM-583-R-Palfinger-Epsilon-S110F101---1721826472190064857_big--24072415525385274700.jpg', 'https://img.linemedia.com/img/s/forestry-equipment-wood-chipper-Jenz-MAN-TGS-33-500-HEM-583-R-Palfinger-Epsilon-S110F101---1721826472610252577_big--24072415525385274700.jpg', 'http://example.com/vehicle3_docs1.pdf', '', '', 1, 0),
(2, 2, 11, 88, '2017', 96954,'GHI321', '', 1, 'http://example.com/vehicle4_photo1.png', '', '', 'http://example.com/vehicle4_docs1.pdf', '', '', 1, 0),
(3, 2, 10, 52, '2021', 403954, 'JKL654', 'TRAIL654', 1, 'http://example.com/vehicle5_photo1.png', 'http://example.com/vehicle5_photo2.png', '', 'http://example.com/vehicle5_docs1.pdf', '', '', 1, 0),
(3, 5, 1, 5, '2022',  53954, 'MNO987', '', 0, 'http://example.com/vehicle6_photo1.png', '', '', 'http://example.com/vehicle6_docs1.pdf', '', '', 1, 0);
