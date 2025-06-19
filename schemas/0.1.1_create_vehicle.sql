
CREATE TABLE tbl_vehicle_brand (
   id SERIAL PRIMARY KEY,
   name VARCHAR(100) NOT NULL DEFAULT '',
   country VARCHAR(100) DEFAULT '',
   founded_year INT DEFAULT 0,
   deleted INT DEFAULT 0
);


CREATE TABLE tbl_vehicle_type (
  id SERIAL PRIMARY KEY,
  title_en VARCHAR(100) NOT NULL DEFAULT '',
  desc_en TEXT DEFAULT '',
  title_ru VARCHAR(100) NOT NULL DEFAULT '',
  desc_ru TEXT DEFAULT '',
  title_tk VARCHAR(100) NOT NULL DEFAULT '',
  desc_tk TEXT DEFAULT '',
  title_de VARCHAR(100) NOT NULL DEFAULT '',
  desc_de TEXT DEFAULT '',
  title_ar VARCHAR(100) NOT NULL DEFAULT '',
  desc_ar TEXT DEFAULT '',
  title_es VARCHAR(100) NOT NULL DEFAULT '',
  desc_es TEXT DEFAULT '',
  title_fr VARCHAR(100) NOT NULL DEFAULT '',
  desc_fr TEXT DEFAULT '',
  title_zh VARCHAR(100) NOT NULL DEFAULT '',
  desc_zh TEXT DEFAULT '',
  title_ja VARCHAR(100) NOT NULL DEFAULT '',
  desc_ja TEXT DEFAULT '',
  deleted INT DEFAULT 0
);


CREATE TABLE tbl_vehicle_model (
   id SERIAL PRIMARY KEY,
   vehicle_brand_id INT REFERENCES tbl_vehicle_brand(id) ON DELETE CASCADE DEFAULT 1,
   vehicle_type_id INT REFERENCES tbl_vehicle_type(id) ON DELETE CASCADE DEFAULT 1,
   name VARCHAR(100) NOT NULL DEFAULT '',
   year INT DEFAULT 0,
   feature VARCHAR(255) DEFAULT '',
   deleted INT DEFAULT 0
);


INSERT INTO tbl_vehicle_brand (name, country, founded_year) VALUES
     ('Ford', 'USA', 1903),
     ('Mercedes-Benz', 'Germany', 1926),
     ('Freightliner', 'USA', 1942),
     ('Volvo', 'Sweden', 1927),
     ('MAN', 'Germany', 1758),
     ('Toyota', 'Japan', 1937),
     ('Hyundai', 'South Korea', 1967),
     ('Isuzu', 'Japan', 1916),
     ('Hino', 'Japan', 1942),
     ('DAF', 'Netherlands', 1928),
     ('Scania', 'Sweden', 1891),
     ('Kenworth', 'USA', 1923),
     ('Iveco', 'Italy', 1975),
     ('Renault', 'France', 1899),
     ('Peterbilt', 'USA', 1939),
     ('Mack', 'USA', 1900),
     ('Chevrolet', 'USA', 1911),
     ('RAM', 'USA', 2009),
     ('Nissan', 'Japan', 1933),
     ('GMC', 'USA', 1901),
     ('Honda', 'Japan', 1948),
     ('Kia', 'South Korea', 1944),
     ('Peugeot', 'France', 1810),
     ('Fiat', 'Italy', 1899),
     ('Chrysler', 'USA', 1925),
     ('Dodge', 'USA', 1900),
     ('Subaru', 'Japan', 1953),
     ('Jaguar', 'UK', 1922),
     ('Land Rover', 'UK', 1948),
     ('Bentley', 'UK', 1919),
     ('Rolls-Royce', 'UK', 1906),
     ('Aston Martin', 'UK', 1913),
     ('Ferrari', 'Italy', 1939),
     ('Lamborghini', 'Italy', 1963),
     ('Maserati', 'Italy', 1914),
     ('Peugeot', 'France', 1810),
     ('Mini', 'Germany', 1959),
     ('BMW', 'Germany', 1916),
     ('Audi', 'Germany', 1909),
     ('Volkswagen', 'Germany', 1937),
     ('Citroën', 'France', 1919),
     ('Porsche', 'Germany', 1931),
     ('Lexus', 'Japan', 1989),
     ('Acura', 'Japan', 1986),
     ('Infinity', 'Japan', 1989),
     ('Chrysler', 'USA', 1925),
     ('Jeep', 'USA', 1941),
     ('Tesla', 'USA', 2003),
     ('Lincoln', 'USA', 1917),
     ('Buick', 'USA', 1899),
     ('Cadillac', 'USA', 1902),
     ('Chevrolet', 'USA', 1911),
     ('Cadillac', 'USA', 1902),
     ('Daihatsu', 'Japan', 1907),
     ('Suzuki', 'Japan', 1909),
     ('Mazda', 'Japan', 1920),
     ('Chery', 'China', 1997),
     ('BYD', 'China', 1995),
     ('Geely', 'China', 1986),
     ('SAIC Motor', 'China', 1958),
     ('Great Wall Motors', 'China', 1984),
     ('Zotye', 'China', 2005),
     ('Foton', 'China', 1996),
     ('Dongfeng Motor', 'China', 1969),
     ('Shacman', 'China', 1998),
     ('Haval', 'China', 2013),
     ('Tata Motors', 'India', 1945),
     ('Mahindra & Mahindra', 'India', 1945),
     ('Maruti Suzuki', 'India', 1981),
     ('Ashok Leyland', 'India', 1948),
     ('BharatBenz', 'India', 2012),
     ('Eicher Motors', 'India', 1948),
     ('Force Motors', 'India', 1958),
     ('Traton', 'Germany', 2018),
     ('Navistar', 'USA', 1902),
     ('Mitsubishi', 'Japan', 1970),
     ('Kawasaki', 'Japan', 1896),
     ('Yamaha', 'Japan', 1953),
     ('Honda', 'Japan', 1948),
     ('Piaggio', 'Italy', 1884),
     ('Suzuki', 'Japan', 1909),
     ('Harley-Davidson', 'USA', 1903),
     ('Indian Motorcycles', 'USA', 1901),
     ('Royal Enfield', 'UK', 1901),
     ('Vespa', 'Italy', 1946),
     ('Ducati', 'Italy', 1926),
     ('BMW Motorrad', 'Germany', 1923),
     ('Buell', 'USA', 1983),
     ('Benelli', 'Italy', 1911),
     ('Peugeot Motorcycles', 'France', 1898),
     ('KTM', 'Austria', 1953),
     ('Aprilia', 'Italy', 1945),
     ('MV Agusta', 'Italy', 1945),
     ('Laverda', 'Italy', 1949),
     ('Zundapp', 'Germany', 1917),
     ('Aermacchi', 'Italy', 1945),
     ('BSA', 'UK', 1861),
     ('Norton', 'UK', 1898),
     ('Matchless', 'UK', 1899),
     ('Rudge', 'UK', 1903),
     ('Sunbeam', 'UK', 1912),
     ('Triumph', 'UK', 1902),
     ('Indian', 'USA', 1901),
     ('Moto Guzzi', 'Italy', 1921),
     ('Guzzi', 'Italy', 1921);


INSERT INTO tbl_vehicle_type (
    title_en, desc_en, title_ru, desc_ru, title_tk, desc_tk,
    title_de, desc_de, title_ar, desc_ar, title_es, desc_es,
    title_fr, desc_fr, title_zh, desc_zh, title_ja, desc_ja
) VALUES
      ('Trailer', 'Unpowered vehicle typically pulled by a truck, used to transport goods.',
       'Прицеп', 'Несамоходное транспортное средство, как правило, буксируемое грузовиком для перевозки грузов.',
       'Treyler', 'Ýük daşamak üçin ýük ulagynyň çekýän, hereketlendirijisiz ulag görnüşi.',
       'Anhänger', 'Nicht angetriebenes Fahrzeug, das typischerweise von einem Lastwagen gezogen wird, um Waren zu transportieren.',
       'مقطورة', 'مركبة غير مزودة بمحرك تُسحب عادةً بواسطة شاحنة لنقل البضائع.',
       'Remolque', 'Vehículo sin motor que generalmente es remolcado por un camión para transportar mercancías.',
       'Remorque', 'Véhicule non motorisé généralement tracté par un camion pour transporter des marchandises.',
       '拖车', '通常由卡车牵引、用于运输货物的无动力车辆。',
       'トレーラー', '通常はトラックに牽引される、貨物輸送用の動力のない車両。'),
      ('Light Duty Truck', 'Smaller trucks typically used for local deliveries or small cargo transport.',
       'Грузовик легкого класса', 'Небольшие грузовики, обычно используемые для местных поставок или перевозки небольших грузов.',
       'Ýeňil ýük ulag', 'Esasan kiçi ýük daşamak üçin ulanylýan kiçi ulaglar.',
       'Leichtlastwagen', 'Kleinere Lastwagen, die typischerweise für lokale Lieferungen oder kleine Transporte verwendet werden.',
       'شاحنة خفيفة', 'شاحنات صغيرة تستخدم عادة للتوصيل المحلي أو نقل الشحنات الصغيرة.',
       'Camión ligero', 'Camiones pequeños que se utilizan normalmente para entregas locales o transporte de pequeñas cargas.',
       'Camion léger', 'Petits camions utilisés généralement pour les livraisons locales ou le transport de petites charges.',
       '轻型卡车', '通常用于本地配送或小型货物运输的小型卡车。',
       '軽トラック', '主に地域配送や小さな荷物の輸送に使用される小型トラック。'),

      ('Medium Duty Truck', 'Used for deliveries of moderate loads over short to medium distances.',
       'Грузовик среднего класса', 'Используется для доставки умеренных грузов на короткие или средние расстояния.',
       'Orta derejeli ýük ulag', 'Orta ýükleri daşamak üçin ulanylýan ulag.',
       'Mittelklasse-LKW', 'Wird für die Lieferung von mittleren Lasten über kurze bis mittlere Distanzen verwendet.',
       'شاحنة متوسطة', 'تستخدم لنقل الحمولات المتوسطة على مسافات قصيرة إلى متوسطة.',
       'Camión mediano', 'Utilizado para entregas de cargas moderadas en distancias cortas a medias.',
       'Camion moyen', 'Utilisé pour des livraisons de charges modérées sur des distances courtes à moyennes.',
       '中型卡车', '用于短距离到中距离运输中等负载的卡车。',
       '中型トラック', '中距離で中型貨物を運ぶのに使用されます。'),

      ('Heavy Duty Truck', 'Large trucks often used for long-distance and heavy freight transport.',
       'Грузовик тяжелого класса', 'Большие грузовики, часто используемые для дальних перевозок тяжелых грузов.',
       'Ağyr ýük ulag', 'Agyr ýükleri uzak aralyga daşamak üçin ulanylýan ulag.',
       'Schwerlastwagen', 'Große Lastwagen, die häufig für Langstreckentransporte schwerer Güter verwendet werden.',
       'شاحنة ثقيلة', 'شاحنات كبيرة تستخدم عادة لنقل البضائع الثقيلة لمسافات طويلة.',
       'Camión pesado', 'Camiones grandes que suelen utilizarse para el transporte de cargas pesadas a larga distancia.',
       'Camion lourd', 'Grands camions souvent utilisés pour le transport de charges lourdes sur de longues distances.',
       '重型卡车', '通常用于长距离和重型货物运输的大型卡车。',
       '大型トラック', '長距離や重い貨物の輸送に使用される大型トラック。'),
      ('Van', 'Small to medium-sized vehicle, typically used for parcel or courier services.',
       'Фургон', 'Маленький или средний автомобиль, обычно используемый для доставки посылок или курьерских услуг.',
       'Minibus', 'Kiçi ýa-da orta ölçegli ulag, adatça sargytlary ýa-da kurýer hyzmatlaryny daşamak üçin ulanylýar.',
       'Lieferwagen', 'Kleines bis mittelgroßes Fahrzeug, typischerweise für Paket- oder Kurierdienste verwendet.',
       'شاحنة صغيرة', 'مركبة صغيرة إلى متوسطة الحجم، تستخدم عادةً في خدمات التوصيل أو البريد السريع.',
       'Furgoneta', 'Vehículo de tamaño pequeño o mediano, típicamente utilizado para servicios de paquetería o mensajería.',
       'Fourgon', 'Véhicule de petite à moyenne taille, généralement utilisé pour des services de colis ou de messagerie.',
       '货车', '小型至中型车辆，通常用于包裹或快递服务。',
       'バン', '小型から中型の車両で、通常は宅配や配送サービスに使用されます。'),

      ('Refrigerated Truck', 'Truck equipped with refrigeration unit for transporting temperature-sensitive goods.',
       'Рефрижератор', 'Грузовик с холодильным оборудованием для перевозки товаров, чувствительных к температуре.',
       'Sowadyjy ýük ulag', 'Temperaturanyň täsirine duýgur harytlary daşamak üçin sowadyjy enjam bilen enjamlaşdyrylan ulag.',
       'Kühlwagen', 'LKW mit Kühleinheit zum Transport temperaturempfindlicher Güter.',
       'شاحنة مبردة', 'شاحنة مجهزة بوحدة تبريد لنقل البضائع الحساسة للحرارة.',
       'Camión refrigerado', 'Camión equipado con unidad de refrigeración para transportar mercancías sensibles a la temperatura.',
       'Camion frigorifique', 'Camion équipé d’une unité de réfrigération pour transporter des marchandises sensibles à la température.',
       '冷藏车', '配备冷藏设备，用于运输对温度敏感的货物的卡车。',
       '冷蔵トラック', '温度に敏感な商品を輸送するために冷却装置を備えたトラック。'),

      ('Tanker Truck', 'Truck designed for transporting liquid or gaseous cargo.',
       'Цистерна', 'Грузовик, предназначенный для перевозки жидких или газообразных грузов.',
       'Suw geçiriji ulag', 'Suwuk ýa-da gazly ýükleri daşamak üçin niýetlenen ulag.',
       'Tankwagen', 'LKW für den Transport von flüssigen oder gasförmigen Gütern.',
       'شاحنة صهريج', 'شاحنة مصممة لنقل الحمولة السائلة أو الغازية.',
       'Camión cisterna', 'Camión diseñado para transportar carga líquida o gaseosa.',
       'Camion-citerne', 'Camion conçu pour transporter des charges liquides ou gazeuses.',
       '罐车', '专为运输液体或气体货物而设计的卡车。',
       'タンクローリー', '液体やガスの貨物を輸送するために設計されたトラック。'),

      ('Flatbed Truck', 'Truck with a flat platform for transporting bulky or irregularly shaped cargo.',
       'Бортовой грузовик', 'Грузовик с плоской платформой для перевозки громоздких или нестандартных грузов.',
       'Platformaly ulag', 'Uly ýa-da adaty däl görnüşli ýükleri daşamak üçin platformaly ulag.',
       'Pritschenwagen', 'LKW mit einer flachen Plattform zum Transport sperriger oder unregelmäßig geformter Güter.',
       'شاحنة مسطحة', 'شاحنة بمنصة مسطحة لنقل البضائع الضخمة أو ذات الشكل غير المنتظم.',
       'Camión plataforma', 'Camión con una plataforma plana para transportar cargas voluminosas o de forma irregular.',
       'Camion à plateau', 'Camion avec une plateforme plate pour transporter des charges volumineuses ou de forme irrégulière.',
       '平板卡车', '配备平坦平台的卡车，用于运输笨重或形状不规则的货物。',
       '平床トラック', 'かさばったり不規則な形の荷物を輸送するための平らなプラットフォーム付きトラック。'),

      ('Box Truck', 'Truck with a cargo area that is fully enclosed, typically used for moving or logistics.',
       'Фургон', 'Грузовик с полностью закрытым грузовым отсеком, обычно используемый для переездов или логистики.',
       'Gutuly ýük ulag', 'Doly ýapyk ýük meýdany bolan ulag, adatça göçürme ýa-da logistika üçin ulanylýar.',
       'Kastenwagen', 'LKW mit vollständig geschlossenem Frachtraum, typischerweise für Umzüge oder Logistik verwendet.',
       'شاحنة بصندوق', 'شاحنة بمنطقة شحن مغلقة بالكامل، تُستخدم عادةً في النقل أو الخدمات اللوجستية.',
       'Camión caja', 'Camión con un área de carga completamente cerrada, generalmente usado para mudanzas o logística.',
       'Camion à caisse', 'Camion avec une zone de chargement entièrement fermée, généralement utilisé pour les déménagements ou la logistique.',
       '厢式货车', '货运区完全封闭的卡车，通常用于搬家或物流。',
       'ボックストラック', '完全に囲まれた貨物エリアを持ち、引っ越しや物流によく使用されるトラック。'),

      ('Trailer', 'A non-motorized vehicle that is towed behind a truck for transporting goods.',
       'Прицеп', 'Немоторизованное транспортное средство, буксируемое грузовиком для перевозки товаров.',
       'Tirkeg', 'Harytlary daşamak üçin ýük ulagynyň yzyna dakylyp çekilýän motoryzasyýa edilmedik ulag.',
       'Anhänger', 'Ein nicht motorisiertes Fahrzeug, das hinter einem LKW gezogen wird, um Güter zu transportieren.',
       'مقطورة', 'مركبة غير مزودة بمحرك يتم سحبها خلف شاحنة لنقل البضائع.',
       'Remolque', 'Vehículo no motorizado que se remolca detrás de un camión para transportar mercancías.',
       'Remorque', 'Véhicule non motorisé remorqué derrière un camion pour transporter des marchandises.',
       '拖车', '一种无动力车辆，被卡车牵引用于运输货物。',
       'トレーラー', 'トラックの後ろに牽引される、エンジンのない車両で、貨物の輸送に使用されます。'),

      ('Pickup Truck', 'Light-duty truck primarily used for personal or small-scale commercial transport.',
       'Пикап', 'Легкий грузовик, используемый в основном для личного или малого коммерческого транспорта.',
       'Käriçik ýük ulag', 'Esasan şahsy ýa-da kiçi göwrümli söwda üçin ulanylýan ýeňil ýük ulag.',
       'Pickup', 'Leichter LKW, der hauptsächlich für den privaten oder kleingewerblichen Transport verwendet wird.',
       'شاحنة صغيرة', 'شاحنة خفيفة تُستخدم بشكل أساسي للنقل الشخصي أو التجاري الصغير.',
       'Camioneta', 'Camión ligero utilizado principalmente para transporte personal o comercial a pequeña escala.',
       'Pick-up', 'Petit camion utilisé principalement pour le transport personnel ou commercial à petite échelle.',
       '皮卡车', '主要用于个人或小规模商业运输的轻型卡车。',
       'ピックアップトラック', '主に個人または小規模な商業輸送に使用される軽量トラック。'),

      ('Chassis Cab', 'A truck chassis with no cargo area, often used for customization with various bodies such as flatbeds or boxes.',
       'Шасси с кабиной', 'Грузовик с шасси без грузового отсека, часто используемый для модификации под различные кузова, такие как платформы или фургоны.',
       'Kabinaly şassi', 'Ýük meýdany bolmadyk, adatça platformalar ýa-da gutular ýaly dürli bedenler üçin özleşdirilýän ulag.',
       'Fahrgestell mit Kabine', 'Ein LKW-Chassis ohne Frachtraum, oft zur Anpassung mit verschiedenen Aufbauten wie Plattformen oder Kasten verwendet.',
       'هيكل شاحنة بكابينة', 'هيكل شاحنة بدون منطقة شحن، يستخدم غالبًا للتخصيص بأجسام مختلفة مثل المنصات أو الصناديق.',
       'Chasis cabina', 'Camión con chasis sin área de carga, a menudo usado para personalización con varias carrocerías como plataformas o cajas.',
       'Châssis cabine', 'Camion avec un châssis sans zone de chargement, souvent utilisé pour des personnalisations avec diverses carrosseries telles que des plateaux ou des boîtes.',
       '底盘驾驶室', '一种没有货物空间的卡车底盘，通常用于定制，如平板车或箱式车。',
       'シャーシキャブ', '貨物エリアのないトラックのシャーシで、フラットベッドやボックスなどのカスタマイズによく使用されます。'),

      ('Utility Truck', 'Truck equipped with tools and equipment, used for maintenance or service-related tasks.',
       'Специальный грузовик', 'Грузовик, оснащённый инструментами и оборудованием, используемый для задач обслуживания или ремонта.',
       'Tehniki hyzmat ulag', 'Enjamlar we gurallar bilen enjamlaşdyrylan, hyzmat ýa-da abatlaýyş işleri üçin ulanylýan ulag.',
       'Werkzeugwagen', 'Ein LKW mit Werkzeugen und Ausrüstung, der für Wartungs- oder Reparaturarbeiten verwendet wird.',
       'شاحنة خدمات', 'شاحنة مجهزة بالأدوات والمعدات، تُستخدم في الصيانة أو المهام المتعلقة بالخدمة.',
       'Camión de servicios', 'Camión equipado con herramientas y equipo, usado para tareas de mantenimiento o relacionadas con servicios.',
       'Camion utilitaire', 'Camion équipé d’outils et d’équipements, utilisé pour des tâches de maintenance ou de service.',
       '维修卡车', '配备工具和设备的卡车，用于维护或相关服务任务。',
       'ユーティリティトラック', 'ツールや機器が装備され、保守やサービス関連の作業に使用されるトラック。'),

      ('Dump Truck', 'Truck designed to transport loose materials (like sand, gravel, or demolition waste), typically equipped with a hydraulically operated bed.',
       'Самосвал', 'Грузовик, предназначенный для перевозки сыпучих материалов (песка, гравия, строительного мусора), оснащённый гидравлически управляемым кузовом.',
       'Gidrawlikli samoswal', 'Esasan çäge, çagyl ýa-da galyndylar ýaly materiallary daşamak üçin gidrawlik enjam bilen enjamlaşdyrylan ulag.',
       'Kipplaster', 'Ein LKW, der zum Transport von Schüttgütern (wie Sand, Kies oder Abbruchmaterial) mit einem hydraulisch betriebenen Kipper ausgestattet ist.',
       'شاحنة قلابة', 'شاحنة مصممة لنقل المواد السائبة مثل الرمل أو الحصى أو مخلفات الهدم، ومزودة بسرير يعمل هيدروليكيًا.',
       'Camión volquete', 'Camión diseñado para transportar materiales sueltos (como arena, grava o escombros), equipado con una cama operada hidráulicamente.',
       'Camion-benne', 'Camion conçu pour transporter des matériaux en vrac (comme le sable, le gravier ou les déchets de démolition), équipé d’une benne hydraulique.',
       '翻斗卡车', '设计用于运输松散材料（如沙子、砾石或拆迁废料）的卡车，配备液压驱动车斗。',
       'ダンプトラック', '砂、砂利、または解体廃材などのバルク材料を運ぶために設計され、油圧作動のベッドを備えたトラック。');


-- Ford Models
INSERT INTO tbl_vehicle_model (vehicle_brand_id, name, vehicle_type_id)VALUES
    (1, 'F-150', 1),
    (1, 'F-250 Super Duty', 2),
    (1, 'Transit Van', 4),
    (1, 'F-450 Super Duty', 3),
    (1, 'Super Duty F-550', 3),
    (1, 'Transit 350 HD', 4),
    (1, 'F-650', 3),
    (1, 'F-750', 3),
    (1, 'Ranger', 1),
    (1, 'F-350 Flatbed', 7),
    (1, 'F-550 Super Duty', 3),
    (1, 'F-650 Super Duty', 3);

-- Mercedes-Benz Models
INSERT INTO tbl_vehicle_model (vehicle_brand_id, name, vehicle_type_id)VALUES
    (2, 'Actros 1843', 3),
    (2, 'Sprinter 2500', 4),
    (2, 'Econic', 3),
    (2, 'Atego 1523', 2),
    (2, 'Vito', 4),
    (2, 'Sprinter 3500', 4),
    (2, 'Actros 2545', 3),
    (2, 'Antos', 3),
    (2, 'Atego 1223', 2),
    (2, 'Axor 1824', 3),
    (2, 'Sprinter 516', 4),
    (2, 'Vito Panel Van', 4);

-- Freightliner Models
INSERT INTO tbl_vehicle_model (vehicle_brand_id, name, vehicle_type_id)VALUES
    (3, 'Cascadia 113', 3),
    (3, 'M2 106', 2),
    (3, 'Columbia 120', 3),
    (3, 'Sprinter 2500', 4),
    (3, 'M2 112', 2),
    (3, 'Cascadia 125', 3),
    (3, 'Freightliner M2 106', 2),
    (3, 'Freightliner FLD 120', 3),
    (3, 'Sprinter 3500', 4),
    (3, 'Freightliner 108SD', 3),
    (3, 'Freightliner FLD', 3);

-- Volvo Models
INSERT INTO tbl_vehicle_model (vehicle_brand_id, name, vehicle_type_id)VALUES
    (4, 'FH16', 3),
    (4, 'FMX', 3),
    (4, 'FM', 2),
    (4, 'VNL 670', 3),
    (4, 'VNR', 3),
    (4, 'V70', 4),
    (4, 'V60', 4),
    (4, 'Volvo FH 460', 3),
    (4, 'Volvo L120H', 2),
    (4, 'Volvo EC950F Crawler', 2),
    (4, 'FH', 3),
    (4, 'VNL', 3),
    (4, 'FL', 2),
    (4, 'FE', 2),
    (4, 'L-series', 3),
    (4, 'XC90', 4),
    (4, 'XC60', 4),
    (4, 'S60', 4),
    (4, 'XC40', 4),
    (4, 'FL Electric', 2),
    (4, 'FE Electric', 2),
    (4, 'S90', 4),
    (4, 'XC40 Recharge', 4),
    (4, 'V90 Cross Country', 4),
    (4, 'XC70', 4);


-- MAN Models
INSERT INTO tbl_vehicle_model (vehicle_brand_id, name, vehicle_type_id)VALUES
    (5, 'TGX', 3),
    (5, 'TGM', 2),
    (5, 'TGL', 2),
    (5, 'MAN TGE', 4),
    (5, 'MAN TGS', 3),
    (5, 'MAN TGE 3.180', 4),
    (5, 'MAN CLA', 3),
    (5, 'MAN TGS 18.440', 3),
    (5, 'MAN TGE 4x4', 4),
    (5, 'MAN L2000', 2),
    (5, 'MAN TGS 26.440', 3);

-- Toyota Models
INSERT INTO tbl_vehicle_model (vehicle_brand_id, name, vehicle_type_id)VALUES
    (6, 'Hilux', 1),
    (6, 'Dyna 150', 2),
    (6, 'Proace', 4),
    (6, 'HiAce', 4),
    (6, 'Coaster', 3),
    (6, 'Toyota Tacoma', 1),
    (6, 'Tundra', 1),
    (6, 'Toyota 4Runner', 4),
    (6, 'Proace City', 4),
    (6, 'Toyota Tundra CrewMax', 1),
    (6, 'Toyota HiAce 4x4', 4);

-- Hyundai Models
INSERT INTO tbl_vehicle_model (vehicle_brand_id, name, vehicle_type_id)VALUES
    (7, 'Hyundai Porter 2', 4),
    (7, 'Hyundai Mighty', 2),
    (7, 'Hyundai Xcient', 3),
    (7, 'Hyundai HD170', 2),
    (7, 'Hyundai HD78', 2),
    (7, 'Hyundai Santa Fe', 4),
    (7, 'Hyundai Tucson', 4),
    (7, 'Hyundai H350', 4),
    (7, 'Hyundai Elantra', 4),
    (7, 'Hyundai Staria', 4);

-- Isuzu Models
INSERT INTO tbl_vehicle_model (vehicle_brand_id, name, vehicle_type_id)VALUES
    (8, 'Isuzu NPR', 2),
    (8, 'Isuzu FTR', 3),
    (8, 'Isuzu ELF', 2),
    (8, 'Isuzu D-Max', 1),
    (8, 'Isuzu NQR', 2),
    (8, 'Isuzu NRR', 2),
    (8, 'Isuzu Giga', 3),
    (8, 'Isuzu MU-X', 4),
    (8, 'Isuzu V-Cross', 1),
    (8, 'Isuzu F-Series', 3);
-- Hino Models
INSERT INTO tbl_vehicle_model (vehicle_brand_id, name, vehicle_type_id) VALUES
    (9, 'Hino 300', 2),
    (9, 'Hino 500', 3),
    (9, 'Hino 700', 3),
    (9, 'Hino Dutro', 2),
    (9, 'Hino XL Series', 3),
    (9, 'Hino 300 Series', 2),
    (9, 'Hino Ranger', 3),
    (9, 'Hino L-Series', 2),
    (9, 'Hino Hybrid', 2),
    (9, 'Hino 500 Series', 3);
-- DAF Models
INSERT INTO tbl_vehicle_model (vehicle_brand_id, name, vehicle_type_id)VALUES
    (10, 'DAF XF', 3),
    (10, 'DAF CF', 3),
    (10, 'DAF LF', 2),
    (10, 'DAF CF Electric', 3),
    (10, 'DAF XF 105', 3),
    (10, 'DAF LF 18t', 2),
    (10, 'DAF CF 85', 3),
    (10, 'DAF CF 75', 3),
    (10, 'DAF XF 530', 3),
    (10, 'DAF CF 480', 3);
-- Scania Models
INSERT INTO tbl_vehicle_model (vehicle_brand_id, name, vehicle_type_id)VALUES
    (11, 'Scania R Series', 3),
    (11, 'Scania P Series', 2),
    (11, 'Scania G Series', 3),
    (11, 'Scania L Series', 2),
    (11, 'Scania S Series', 3),
    (11, 'Scania P 230', 2),
    (11, 'Scania V8', 3),
    (11, 'Scania R500', 3),
    (11, 'Scania Citywide', 3),
    (11, 'Scania G440', 3);
-- Kenworth Models
INSERT INTO tbl_vehicle_model (vehicle_brand_id, name, vehicle_type_id)VALUES
    (12, 'Kenworth T680', 3),
    (12, 'Kenworth W900', 3),
    (12, 'Kenworth T880', 3),
    (12, 'Kenworth K270', 2),
    (12, 'Kenworth T800', 3),
    (12, 'Kenworth T370', 2),
    (12, 'Kenworth C500', 3),
    (12, 'Kenworth T300', 2),
    (12, 'Kenworth T680 Advantage', 3),
    (12, 'Kenworth K370', 2);
-- Iveco Models
INSERT INTO tbl_vehicle_model (vehicle_brand_id, name, vehicle_type_id)VALUES
    (13, 'Iveco Stralis', 3),
    (13, 'Iveco Daily', 4),
    (13, 'Iveco Trakker', 3),
    (13, 'Iveco Eurocargo', 2),
    (13, 'Iveco S-Way', 3),
    (13, 'Iveco Hi-Way', 3),
    (13, 'Iveco Magelys', 3),
    (13, 'Iveco Eurotech', 3),
    (13, 'Iveco 682', 2),
    (13, 'Iveco TurboStar', 3);
--
-- -- Nissan Models
-- INSERT INTO tbl_vehicle_model (vehicle_brand_id, name, vehicle_type_id)VALUES
--     (14, 'Nissan Navara', 1),
--     (14, 'Nissan NV3500', 4),
--     (14, 'Nissan Frontier', 1),
--     (14, 'Nissan NV200', 4),
--     (14, 'Nissan Titan XD', 1),
--     (14, 'Nissan Atlas', 2);
-- -- Mitsubishi Fuso Models
-- INSERT INTO tbl_vehicle_model (vehicle_brand_id, name, vehicle_type_id)VALUES
--     (15, 'Mitsubishi Fuso Canter', 2),
--     (15, 'Mitsubishi Fuso Fighter', 3),
--     (15, 'Mitsubishi Fuso Rosa', 4),
--     (15, 'Mitsubishi Fuso Super Great', 3),
--     (15, 'Mitsubishi Fuso FE', 2);
-- -- Peugeot Models
-- INSERT INTO tbl_vehicle_model (vehicle_brand_id, name, vehicle_type_id)VALUES
--     (16, 'Peugeot Boxer', 4),
--     (16, 'Peugeot Partner', 4),
--     (16, 'Peugeot Expert', 4),
--     (16, 'Peugeot 508 SW', 4);
-- -- Citroën Models
-- INSERT INTO tbl_vehicle_model (vehicle_brand_id, name, vehicle_type_id)
-- VALUES
--     (17, 'Citroën Berlingo', 4),
--     (17, 'Citroën Jumpy', 4),
--     (17, 'Citroën Jumper', 4),
--     (17, 'Citroën C4 Picasso', 4);
-- -- Renault Trucks Models
-- INSERT INTO tbl_vehicle_model (vehicle_brand_id, name, vehicle_type_id)
-- VALUES
--     (18, 'Renault Master', 4),
--     (18, 'Renault Trafic', 4),
--     (18, 'Renault Trucks T', 3),
--     (18, 'Renault Trucks K', 3),
--     (18, 'Renault Trucks C', 3),
--     (18, 'Renault Trucks D', 2);
-- -- Honda Models
-- INSERT INTO tbl_vehicle_model (vehicle_brand_id, name, vehicle_type_id)VALUES
--     (19, 'Honda Ridgeline', 1),
--     (19, 'Honda Acty', 4),
--     (19, 'Honda CR-V', 4);
-- -- Dongfeng Models
-- INSERT INTO tbl_vehicle_model (vehicle_brand_id, name, vehicle_type_id)VALUES
--     (20, 'Dongfeng KX', 2),
--     (20, 'Dongfeng Tianlong', 3),
--     (20, 'Dongfeng DF', 2),
--     (20, 'Dongfeng EQ', 3),
--     (20, 'Dongfeng Fengshen', 4);
-- -- Tata Motors Models
-- INSERT INTO tbl_vehicle_model (vehicle_brand_id, name, vehicle_type_id)VALUES
--     (21, 'Tata LPT', 3),
--     (21, 'Tata Ultra', 2),
--     (21, 'Tata Ace', 1),
--     (21, 'Tata SFC', 2),
--     (21, 'Tata Winger', 4);
-- -- Mahindra Models
-- INSERT INTO tbl_vehicle_model (vehicle_brand_id, name, vehicle_type_id)VALUES
--     (22, 'Mahindra Bolero Pik-Up', 1),
--     (22, 'Mahindra Scorpio', 4),
--     (22, 'Mahindra XUV700', 4),
--     (22, 'Mahindra Treo', 4),
--     (22, 'Mahindra 4110', 2);
-- -- FAW Models
-- INSERT INTO tbl_vehicle_model (vehicle_brand_id, name, vehicle_type_id)VALUES
--     (23, 'FAW Jiefang', 3),
--     (23, 'FAW T3', 2),
--     (23, 'FAW V2', 4),
--     (23, 'FAW Qingling', 2);
-- -- Zhengzhou Yutong Models
-- INSERT INTO tbl_vehicle_model (vehicle_brand_id, name, vehicle_type_id)VALUES
--     (24, 'Yutong ZK6122H', 3),
--     (24, 'Yutong ZK6105', 3),
--     (24, 'Yutong T7', 3);
-- -- Great Wall Motors Models
-- INSERT INTO tbl_vehicle_model (vehicle_brand_id, name, vehicle_type_id)VALUES
--     (25, 'Great Wall Wingle', 1),
--     (25, 'Great Wall Steed', 1),
--     (25, 'Great Wall Voleex', 4);
-- -- Kia Motors Models
-- INSERT INTO tbl_vehicle_model (vehicle_brand_id, name, vehicle_type_id)VALUES
--     (26, 'Kia Bongo', 2),
--     (26, 'Kia K2500', 2),
--     (26, 'Kia Soul EV', 4);
--
