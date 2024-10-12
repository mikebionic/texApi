CREATE TABLE tbl_role (
    id SERIAL PRIMARY KEY,
    uuid UUID DEFAULT gen_random_uuid(),
    name VARCHAR(100) UNIQUE NOT NULL,
    description TEXT DEFAULT '',
    title VARCHAR(100) NOT NULL DEFAULT '',
    title_ru VARCHAR(100) NOT NULL DEFAULT '',
    subtitle VARCHAR(200) NOT NULL DEFAULT '',
    subtitle_ru VARCHAR(200) NOT NULL DEFAULT ''
);

CREATE TABLE tbl_user (
  id SERIAL PRIMARY KEY,
  uuid UUID DEFAULT gen_random_uuid(),
  username VARCHAR(200) NOT NULL,
  password VARCHAR(200) NOT NULL,
  email VARCHAR(100) NOT NULL DEFAULT '',
  first_name VARCHAR(100) NOT NULL DEFAULT '',
  last_name VARCHAR(100) NOT NULL DEFAULT '',
  nick_name VARCHAR(100) NOT NULL DEFAULT '',
  avatar_url VARCHAR(200) DEFAULT '',
  phone VARCHAR(100) NOT NULL DEFAULT '',
  info_phone VARCHAR(100) NOT NULL DEFAULT '',
  address VARCHAR(200) NOT NULL DEFAULT '',
  role_id INT REFERENCES tbl_role(id) ON DELETE SET NULL DEFAULT 0,
  subrole_id INT NOT NULL DEFAULT 0,
  verified INT NOT NULL DEFAULT 0,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  active INT DEFAULT 1,
  deleted INT DEFAULT 0,

  oauth_provider VARCHAR(100) DEFAULT '',
  oauth_user_id VARCHAR(100) DEFAULT '',
  oauth_location VARCHAR(200) DEFAULT '',
  oauth_access_token VARCHAR(500) DEFAULT '',
  oauth_access_token_secret VARCHAR(500) DEFAULT '',
  oauth_refresh_token VARCHAR(500) DEFAULT '',
  oauth_expires_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  oauth_id_token VARCHAR(500) DEFAULT '',
  refresh_token VARCHAR(500) DEFAULT ''
);


INSERT INTO tbl_role (name, description, title, subtitle, title_ru, subtitle_ru) VALUES
    ('admin', 'Has full access to manage the system','','','',''),
    ('sender', 'Can place orders and track deliveries', 'Sender', 'I am looking for transport','Отправитель', 'Я ищу транспорт'),
    ('carrier', 'Responsible for delivering orders', 'Carrier', 'I am looking for cargo','Перевозчик','Я ищу груз'),
    ('carrier_personal', 'Responsible for delivering orders using their personal vehicle', 'Carrier', 'Personal vehicle','Перевозчик','Личный автотранспорт'),
    ('carrier_owner', 'Responsible for delivering orders with a fleet of vehicles they own', 'Carrier', 'Fleet owner','Перевозчик','Владелец парка автотранспорта'),
    ('carrier_company', 'Responsible for delivering orders through a logistics company', 'Carrier', 'Logistics company','Перевозчик','Логистическая кампания');

INSERT INTO tbl_user (username, password, email, first_name, last_name, nick_name, avatar_url, phone, info_phone, address, role_id, verified, active, deleted)
VALUES
    ('root', 'letmein', 'texlogistics@gmail.com', 'Tex', 'Admin', 'Texy', '', '+0036123456', '+0036123456', 'Ashgabat, Turkmenistan', 1, 1, 1, 0),
    ('customer1', 'password123', 'customer1@example.com', 'Customer', 'One', 'Custy', '', '+123456789', '+123456789', 'Ashgabat, Turkmenistan', 2, 1, 1, 0),
    ('driver1', 'password123', 'driver1@example.com', 'Volodya', '', 'Driver', '', '+123456789', '+123456789', 'Ashgabat, Turkmenistan', 3, 1, 1, 0);
