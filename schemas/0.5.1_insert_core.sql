
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
('personal_maria', 'pass123', 'maria.driver@gmail.com', '+99365789010', 'carrier', 4, 1, 1, 0),
-- sender (2)
('sender_haknazar', 'pass123', 'haknazar.sender@gmail.com', '+99365789001', 'sender', 3, 1, 1, 0),
('sender_saryyew',  'pass123', 'saryyew.sender@gmail.com',  '+99365789002', 'sender', 3, 1, 1, 0),
-- personal (1)
('personal_mekan',  'pass123', 'mekan.personal@gmail.com',  '+99365789003', 'carrier', 4, 1, 1, 0),
-- fleet (1)
('fleet_dowran',    'pass123', 'dowran.fleet@gmail.com',    '+99365789004', 'carrier', 5, 1, 1, 0),
-- logistics (1)
('logistics_dowlet','pass123', 'dowlet.logistics@gmail.com','+99365789005', 'carrier',6, 1, 1, 0);

INSERT INTO tbl_user (username, password, email, phone, role, role_id, verified, active, deleted, driver_id)
VALUES
('driver-test','test', 'test.driver@gmail.com','+99361234567', 'driver',6, 1, 1, 0, 1);


INSERT INTO tbl_company (user_id,role_id,company_name,first_name,last_name,patronymic_name,address,country,phone,email,image_url,entity,active,deleted) VALUES
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
(14, 3, 'haknazar.sender', 'haknazar', 'sender', '', 'Gorogly street 234, Ashgabat', 'Turkmenistan','+99365789001','haknazar.sender@gmail.com', 'https://images.unsplash.com/photo-1624459294159-598de7ddcae9', 'legal', 1, 0),
(15, 3, 'saryyew.sender',  'saryyew', 'sender', '', 'Yunus Emre street 89, Ashgabat', 'Turkmenistan','+99365789002','saryyew.sender@gmail.com', 'https://images.unsplash.com/photo-1624459294159-598de7ddcae9', 'individual', 1, 0),
-- personal (1)
(16, 4, 'mekan.personal',  'mekan', 'personal', '', 'Yunus Emre street 89, Ashgabat', 'Turkmenistan','+99365789003','mekan.personal@gmail.com', 'https://images.unsplash.com/photo-1624459294159-598de7ddcae9', 'individual', 1, 0),
-- fleet (1)
(17, 5, 'dowran.fleet',    'dowran', 'fleet', '', 'Gorogly street 234, Ashgabat',    'Turkmenistan','+99365789004','dowran.fleet@gmail.com', 'https://images.unsplash.com/photo-1624459294159-598de7ddcae9', 'legal', 1, 0),
-- logistics (1)
(18, 6, 'dowlet.logistics','dowlet', 'logistics', '', 'Atamurat Niyazov street 75, Ashgabat','Turkmenistan','+99365789005','dowlet.logistics@gmail.com', 'https://images.unsplash.com/photo-1624459294159-598de7ddcae9', 'legal', 1, 0);


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


-- Sample Version Inserts for Testing

-- iOS Versions
INSERT INTO tbl_version (
version_number, version_code, title, description, platform, minimal_platform_version,
download_url, file_size, checksum, changelog, release_notes, is_critical_update,
is_beta, auto_update_enabled, rollout_percentage, active_at
) VALUES
-- iOS Production Versions
('1.0.0', 100, 'MyApp iOS Launch', 'Initial release of MyApp for iOS devices', 'ios', 'iOS 14.0',
'https://apps.apple.com/app/myapp/id123456789', 125829120, 'sha256:a1b2c3d4e5f6789012345678901234567890abcdef',
'• Initial release\n• User authentication\n• Basic functionality\n• Push notifications',
'Welcome to MyApp! Create your account and start exploring our amazing features.',
false, false, true, 100, '2024-01-15 10:00:00'),

('1.1.0', 110, 'MyApp iOS Feature Update', 'Major feature update with new capabilities', 'ios', 'iOS 14.0',
'https://apps.apple.com/app/myapp/id123456789', 132841472, 'sha256:b2c3d4e5f6789012345678901234567890abcdef1',
'• Dark mode support\n• Enhanced security\n• Bug fixes\n• Performance improvements\n• New dashboard',
'New features are here! Experience dark mode and improved performance.',
false, false, true, 100, '2024-02-20 14:30:00'),

('1.2.0', 120, 'MyApp iOS Performance Update', 'Critical performance and security updates', 'ios', 'iOS 15.0',
'https://apps.apple.com/app/myapp/id123456789', 138572800, 'sha256:c3d4e5f6789012345678901234567890abcdef12',
'• Critical security patches\n• 50% faster loading times\n• Memory usage optimization\n• Crash fixes',
'Important security and performance improvements. Please update immediately.',
true, false, true, 100, '2024-03-10 09:15:00'),

-- iOS Beta Version
('1.3.0-beta.1', 129, 'MyApp iOS Beta', 'Beta version with experimental features', 'ios', 'iOS 15.0',
'https://testflight.apple.com/join/ABC123DEF', 142606336, 'sha256:d4e5f6789012345678901234567890abcdef123',
'• New AI-powered features\n• Redesigned UI components\n• Advanced analytics\n• Experimental chat bot',
'Try out our latest experimental features! This is a beta version.',
false, true, false, 50, '2024-04-01 16:00:00'),

-- Android Versions
('1.0.0', 1000, 'MyApp Android Launch', 'Initial Android release with core features', 'android', 'Android 8.0',
'https://play.google.com/store/apps/details?id=com.myapp.android', 89478144, 'sha256:e5f6789012345678901234567890abcdef1234',
'• Android launch\n• Material Design UI\n• Google Sign-In integration\n• Offline mode support',
'MyApp is now available on Android! Enjoy native Android experience.',
false, false, true, 100, '2024-01-20 11:00:00'),

('1.1.0', 1100, 'MyApp Android Feature Pack', 'Feature parity with iOS plus Android exclusives', 'android', 'Android 8.0',
'https://play.google.com/store/apps/details?id=com.myapp.android', 95420416, 'sha256:f6789012345678901234567890abcdef12345',
'• Widget support\n• Android Auto integration\n• Adaptive icons\n• Split-screen mode\n• Notification channels',
'Android exclusive features! Widgets, Android Auto, and more.',
false, false, true, 100, '2024-02-25 13:45:00'),

('1.2.1', 1210, 'MyApp Android Hotfix', 'Critical bug fixes and security updates', 'android', 'Android 9.0',
'https://play.google.com/store/apps/details?id=com.myapp.android', 97845248, 'sha256:g789012345678901234567890abcdef123456',
'• Fixed crash on Android 14\n• Security vulnerability patched\n• Battery optimization\n• Background sync improvements',
'Critical fixes for Android users. Update recommended.',
true, false, true, 100, '2024-03-15 10:30:00'),

-- Android Beta
('1.3.0-alpha.2', 1299, 'MyApp Android Alpha', 'Early access alpha build', 'android', 'Android 10.0',
'https://play.google.com/apps/internaltest/ABC123DEF789', 103219200, 'sha256:h89012345678901234567890abcdef1234567',
'• Machine learning features\n• Advanced camera integration\n• Voice commands\n• Gesture navigation',
'Alpha build with cutting-edge features. Expect bugs!',
false, true, false, 25, '2024-04-05 14:20:00'),

-- Web Versions
('1.0.0', 10000, 'MyApp Web Platform Launch', 'Progressive Web App launch', 'web', 'Chrome 90+',
'https://app.myapp.com/', 52428800, 'sha256:i9012345678901234567890abcdef12345678',
'• Progressive Web App\n• Offline functionality\n• Cross-platform sync\n• Responsive design',
'Access MyApp from any browser! Works offline too.',
false, false, true, 100, '2024-01-25 12:00:00'),

('1.1.0', 11000, 'MyApp Web Enhanced', 'Enhanced web experience with new features', 'web', 'Chrome 95+',
'https://app.myapp.com/', 58720256, 'sha256:j012345678901234567890abcdef123456789',
'• WebAssembly integration\n• Enhanced performance\n• Web push notifications\n• File system access\n• Background sync',
'Faster, more powerful web experience with native-like features.',
false, false, true, 100, '2024-03-01 15:30:00'),

('1.2.0', 12000, 'MyApp Web Security Update', 'Important security and compliance updates', 'web', 'Chrome 100+',
'https://app.myapp.com/', 61865984, 'sha256:k12345678901234567890abcdef1234567890',
'• GDPR compliance updates\n• Enhanced encryption\n• Two-factor authentication\n• Security audit fixes',
'Enhanced security and privacy features. Update automatically applied.',
false, false, true, 100, '2024-03-20 08:00:00'),

-- Web Beta
('1.3.0-rc.1', 12950, 'MyApp Web Release Candidate', 'Release candidate with new architecture', 'web', 'Chrome 105+',
'https://beta.myapp.com/', 67108864, 'sha256:l23456789012345678901234567890abcdef01',
'• New React 18 architecture\n• Server-side rendering\n• Edge computing integration\n• Advanced caching',
'Experience our next-generation web platform architecture.',
false, true, false, 75, '2024-04-10 11:45:00');

-- Additional test data for different statuses
-- Deprecated version
INSERT INTO tbl_version (
version_number, version_code, title, description, platform, minimal_platform_version,
download_url, file_size, checksum, changelog, release_notes, is_critical_update,
is_beta, auto_update_enabled, rollout_percentage, active_at, deprecated_at, active
) VALUES
('0.9.0', 90, 'MyApp iOS Beta (Deprecated)', 'Deprecated beta version', 'ios', 'iOS 13.0',
'https://apps.apple.com/app/myapp/id123456789', 115343360, 'sha256:m3456789012345678901234567890abcdef012',
'• Legacy beta features\n• Basic functionality only',
'This version is deprecated. Please update to the latest version.',
false, true, false, 0, '2023-12-01 10:00:00', '2024-01-15 10:00:00', 0);

-- End of life version
INSERT INTO tbl_version (
version_number, version_code, title, description, platform, minimal_platform_version,
download_url, file_size, checksum, changelog, release_notes, is_critical_update,
is_beta, auto_update_enabled, rollout_percentage, active_at, deprecated_at, end_of_life_at, active
) VALUES
('0.8.0', 80, 'MyApp Android Legacy (EOL)', 'End of life version', 'android', 'Android 7.0',
'https://play.google.com/store/apps/details?id=com.myapp.android', 75497472, 'sha256:n456789012345678901234567890abcdef0123',
'• Legacy Android version\n• No longer supported',
'This version is no longer supported. Critical update required.',
true, false, false, 0, '2023-11-01 10:00:00', '2023-12-15 10:00:00', '2024-01-01 00:00:00', 0);



INSERT INTO tbl_organization (
name, description_en, description_ru, description_tk, email, image_url, logo_url, icon_url, banner_url,
website_url, about_text, refund_text, delivery_text, contact_text, terms_conditions, privacy_policy,
address1, address2, address3, address4, address_title1, address_title2, address_title3, address_title4,
contact_phone1, contact_phone2, contact_phone3, contact_phone4, contact_title1, contact_title2, contact_title3, contact_title4,
meta, meta2, meta3
) VALUES (
'TEX Logistics',
'TEX is a revolutionary logistics platform connecting carriers, shippers, and companies worldwide. Easily find freight, optimize routes, and grow your network.',
'TEX - это революционная платформа в сфере логистики, соединяющая перевозчиков, отправителей и компании по всему миру. Легко находите грузы, оптимизируйте маршруты и расширяйте сеть.',
'TEX - bu dünýä logistikasyndaky üstünlikli täzelik. Bu platforma ýük daşamagy, müşderileri we kompaniýalary birleşdirýär.',
'admin@texexpress.pro',
'https://images.unsplash.com/photo-1551232864-3f00be7be8f5',
'https://cdn-icons-png.flaticon.com/512/2250/2250246.png',
'https://cdn-icons-png.flaticon.com/512/7436/7436887.png',
'https://images.unsplash.com/photo-1600650684452-8a96e7bcb7be',
'https://www.texexpress.pro',

'{
"en": "TEX is more than a logistics platform – it’s a global network. With 5,000+ companies, 20,000+ users, and 50,000+ successful deals, we transform freight management into seamless digital collaboration.",
"ru": "TEX - это не просто платформа, это глобальная сеть. С 5,000+ компаниями, 20,000+ пользователями и 50,000+ успешными сделками, мы делаем управление грузоперевозками цифровым и эффективным.",
"tk": "TEX diňe platforma däl - bu dünýä logistikasy üçin döredilen hyzmatdaşlyk ulgamydyr. 5000+ kompaniýa, 20000+ ulanyjy we 50000+ üstünlikli geleşik bilen işiňizi ýokarlandyryň."
}',

'{
"en": "Due to the nature of our service, returns are generally not applicable. For issues with transactions or agreements, please contact support.",
"ru": "Из-за характера наших услуг возврат, как правило, не предусмотрен. При возникновении проблем с транзакциями или соглашениями свяжитесь с нашей службой поддержки.",
"tk": "Hyzmat görnüşimiz sebäpli yzyna gaýtarmak mümkin däl. Töleg ýa-da şertnama bilen bagly mesele ýüze çyksa, bize ýüz tutuň."
}',

'{
"en": "Our platform facilitates freight transportation across countries and regions. Timeframes depend on the agreement between parties and selected carrier.",
"ru": "Наша платформа обеспечивает грузоперевозки по странам и регионам. Сроки зависят от условий соглашений и выбранных перевозчиков.",
"tk": "Platformamyz ýurtlar we sebitler boýunça ýük daşamagy üpjün edýär. Eltip bermegiň wagty saýlanan ýükçi bilen baglaşylan şertnama baglydyr."
}',

'{
"en": "Need help? Reach out to us via email or through our platform chat support. Our team is here to assist you.",
"ru": "Нужна помощь? Свяжитесь с нами по электронной почте или в чате платформы. Мы всегда рады помочь.",
"tk": "Kömek gerekmi? Bize e-poçta ýa-da sahypadaky söhbetdeşlik arkaly ýüz tutuň. Biz kömek etmäge taýýar."
}',

'{
"en": "By accessing or using TEX Logistic, you agree to be bound by the following terms and conditions:\n\n1. All users must provide accurate company and contact information.\n2. TEX is not liable for contractual breaches between users but facilitates secure communication and deal-making.\n3. All logistics transactions must comply with national and international trade laws.\n4. Any dispute arising from usage of the platform must be resolved via arbitration under Turkmenistan law.\n5. Platform services are subject to availability and may change without notice.\n\nBy continuing, you confirm you understand and accept these terms.",
"ru": "Используя платформу TEX Logistic, вы соглашаетесь со следующими условиями:\n\n1. Все пользователи обязаны предоставлять достоверную информацию о компании и контактах.\n2. TEX не несёт ответственности за нарушение договоров между пользователями, но обеспечивает безопасную коммуникацию и заключение сделок.\n3. Все логистические операции должны соответствовать национальным и международным торговым законам.\n4. Споры, возникающие в связи с использованием платформы, подлежат арбитражу по законодательству Туркменистана.\n5. Услуги платформы могут изменяться без предварительного уведомления.\n\nПродолжая работу с платформой, вы подтверждаете согласие с этими условиями.",
"tk": "TEX Logistic platformasyny ulanmak bilen, aşakdaky şertlere razy bolýarsyňyz:\n\n1. Ulanyjylar dogry kompaniýa we aragatnaşyk maglumatlaryny bermelidir.\n2. TEX ulanyjylar arasyndaky şertnamalaryň bozulmagyna jogapkär däl, diňe howpsuz aragatnaşyk we geleşik döretmek mümkinçiliklerini üpjün edýär.\n3. Ýük daşama amallary milli we halkara söwda kanunlaryna laýyk bolmalydyr.\n4. Ulanyşdan ýüze çykýan dawalara Türkmenistanyň kanunlaryna laýyklykda arbitraž arkaly serediler.\n5. Platformanyň hyzmatlary öňünden duýduryşsyz üýtgäp biler.\n\nPlatformany ulanmak bilen, bu şertleri kabul edýändigiňizi tassyklaýarsyňyz."
}',

'{
"en": "We take your privacy seriously. TEX Logistic collects only necessary data to provide logistics services:\n\n1. We store company registration, contact details, and usage logs.\n2. We never share personal or company data with third parties without explicit consent.\n3. Data is stored securely in compliance with international standards.\n4. Users may request data removal or review at any time.\n\nYour data is used solely to match you with logistics partners, improve service quality, and ensure transaction transparency.",
"ru": "Мы серьёзно относимся к вашей конфиденциальности. TEX Logistic собирает только необходимые данные для предоставления логистических услуг:\n\n1. Мы храним регистрационные данные компании, контактную информацию и логи использования.\n2. Мы не передаём личные или корпоративные данные третьим лицам без вашего согласия.\n3. Данные хранятся безопасно в соответствии с международными стандартами.\n4. Пользователи могут в любой момент запросить удаление или просмотр данных.\n\nВаши данные используются исключительно для подбора логистических партнёров, улучшения качества сервиса и прозрачности сделок.",
"tk": "TEX Logistic şahsy maglumatlaryňyzyň howpsuzlygyny möhüm hasaplaýar. Biz diňe zerur maglumatlary ýygnap ulanýarys:\n\n1. Kompaniýanyň hasaba alyş maglumatlary, aragatnaşyk maglumatlary we ulanma ýazgylary saklanylýar.\n2. Şahsy ýa-da kompaniýa maglumatlary siziň razylygyňyz bolmazdan üçünji taraplara berilmeýär.\n3. Maglumatlar halkara howpsuzlyk ülňülerine laýyklykda ygtybarly saklanylýar.\n4. Ulanyjylar islendik wagtda maglumatlaryny pozmagy ýa-da gözden geçirmegi haýyş edip biler.\n\nToplanan maglumatlar diňe hyzmatdaş tapmak, hyzmaty gowulandyrmak we geleşikleriň aç-açanlygyny üpjün etmek üçin ulanylýar."
}',

'Gorogly köçesi 100, Aşgabat, Türkmenistan',
'Türkmenabat şäher şahamçasy',
'Mary şäher ofisi',
'Daşoguz logistika merkezi',
'Baş edara - Aşgabat',
'Şahamça - Türkmenabat',
'Şahamça - Mary',
'Şahamça - Daşoguz',

'+993 12 34-56-78',
'+993 65 78-90-12',
'+993 63 12-34-56',
'+993 64 56-78-90',

'Hyzmat bölümi - Aşgabat',
'Tizlik hyzmaty - Türkmenabat',
'Ýük ugrukdyryjy - Mary',
'Konsultasiýa - Daşoguz',

'{"industry":"logistics", "region":"Turkmenistan", "languages":["tk","ru","en"]}',
'{"partners":5000, "active_users":20000, "transactions":50000, "vehicle_shipments":5000}',
'{"features":["freight matching", "carrier tracking", "route optimization", "real-time notifications"]}'
);
