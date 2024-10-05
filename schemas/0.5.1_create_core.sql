CREATE TABLE tbl_role (
  id SERIAL PRIMARY KEY,
  uuid UUID DEFAULT gen_random_uuid (),
  name VARCHAR(100) UNIQUE NOT NULL,
  description TEXT DEFAULT ''
);

CREATE TABLE
    tbl_user (
        id SERIAL PRIMARY KEY,
        uuid UUID DEFAULT gen_random_uuid (),
        username VARCHAR(200) NOT NULL,
        password VARCHAR(200) NOT NULL,
        email VARCHAR(100) NOT NULL DEFAULT '',
        fullname VARCHAR(200) NOT NULL DEFAULT '',
        phone VARCHAR(100) NOT NULL DEFAULT '',
        address VARCHAR(200) NOT NULL DEFAULT '',
        role_id INT REFERENCES tbl_role(id) ON DELETE SET NULL DEFAULT 0,
        verified INT NOT NULL DEFAULT 0,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        active INT DEFAULT 1,
        deleted INT DEFAULT 0
    );

INSERT INTO tbl_role (name, description) VALUES
     ('Admin', 'Has full access to manage the system'),
     ('Customer', 'Can place orders and track deliveries'),
     ('Driver', 'Responsible for delivering orders');

INSERT INTO tbl_user (username, password, email, fullname, phone, address, role_id, verified, active, deleted)
VALUES
    ('root', 'PASSWORD_PLACEHOLDER', 'texlogistics@gmail.com', 'Tex Admin', '+0036123456', 'Ashgabat, Turkmenistan', 1, 1, 1, 0),
    ('customer1', 'PASSWORD_PLACEHOLDER', 'customer1@example.com', 'Customer', '+123456789', 'Ashgabat, Turkmenistan', 2, 1, 1, 0),
    ('driver1', 'PASSWORD_PLACEHOLDER', 'driver1@example.com', 'Volodya', '+123456789', 'Ashgabat, Turkmenistan', 3, 1, 1, 0);