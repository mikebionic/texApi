CREATE TABLE tbl_vehicle_brand (
   id SERIAL PRIMARY KEY,
   name VARCHAR(100) NOT NULL DEFAULT '',
   country VARCHAR(100) DEFAULT '',
   founded_year INT DEFAULT 0,
   deleted INT DEFAULT 0
);


CREATE TABLE tbl_vehicle_type (
  id SERIAL PRIMARY KEY,
  type_name VARCHAR(100) NOT NULL DEFAULT '',
  description TEXT DEFAULT '',
  deleted INT DEFAULT 0
);


CREATE TABLE tbl_vehicle_model (
   id SERIAL PRIMARY KEY,
   brand VARCHAR(100) NOT NULL DEFAULT '',
   name VARCHAR(100) NOT NULL DEFAULT '',
   year INT DEFAULT 0,
   vehicle_type_id INT REFERENCES tbl_vehicle_type(id) ON DELETE CASCADE DEFAULT 1,
   feature VARCHAR(255) DEFAULT '',
   deleted INT DEFAULT 0
);



-- Inserting vehicle brands with their countries and founded years
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
     ('MAN', 'Germany', 1758),
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

-- Inserting Expanded Vehicle Types
INSERT INTO tbl_vehicle_type (type_name, description) VALUES
   ('Light Duty Truck', 'Smaller trucks typically used for local deliveries or small cargo transport.'),
   ('Medium Duty Truck', 'Used for deliveries of moderate loads over short to medium distances.'),
   ('Heavy Duty Truck', 'Large trucks often used for long-distance and heavy freight transport.'),
   ('Van', 'Small to medium-sized vehicle, typically used for parcel or courier services.'),
   ('Refrigerated Truck', 'Truck equipped with refrigeration unit for transporting temperature-sensitive goods.'),
   ('Tanker Truck', 'Truck designed for transporting liquid or gaseous cargo.'),
   ('Flatbed Truck', 'Truck with a flat platform for transporting bulky or irregularly shaped cargo.'),
   ('Box Truck', 'Truck with a cargo area that is fully enclosed, typically used for moving or logistics.'),
   ('Trailer', 'A non-motorized vehicle that is towed behind a truck for transporting goods.'),
   ('Pickup Truck', 'Light-duty truck primarily used for personal or small-scale commercial transport.'),
   ('Chassis Cab', 'A truck chassis with no cargo area, often used for customization with various bodies such as flatbeds or boxes.'),
   ('Utility Truck', 'Truck equipped with tools and equipment, used for maintenance or service-related tasks.'),
   ('Dump Truck', 'Truck designed to transport loose materials (like sand, gravel, or demolition waste), typically equipped with a hydraulically operated bed.'),
   ('Tow Truck', 'Specialized vehicle designed for towing broken-down vehicles or other loads.'),
   ('Heavy Haul Truck', 'Truck used for transporting extremely heavy loads, often including specialized trailers.'),
   ('Cement Mixer', 'Truck designed to mix and transport concrete to construction sites.'),
   ('Logging Truck', 'Heavy-duty truck used to transport logs or timber from forest sites to processing mills.'),
   ('Lorry', 'British term for a truck, typically used for transporting goods over long distances.'),
   ('Articulated Lorry', 'A large vehicle consisting of a towing engine (tractor) and a trailer, commonly used for hauling freight.'),
   ('Cargo Van', 'Van designed for transporting goods, often used in delivery services.'),
   ('Minivan', 'Small van typically used for personal or family transport, but can also be used for small goods delivery.'),
   ('Step Van', 'Van with a flat floor and raised cargo area, often used for deliveries or catering services.'),
   ('Car Carrier', 'Truck designed to transport cars and other vehicles, often on a multi-level flatbed platform.'),
   ('Freight Van', 'A van designed specifically for transporting freight or packages, often larger than a standard passenger van.'),
   ('Fire Truck', 'Specialized truck equipped with firefighting equipment, used for emergency response.'),
   ('Ambulance', 'Vehicle designed for transporting sick or injured persons, equipped with medical facilities.'),
   ('Refrigerated Trailer', 'Non-motorized trailer that contains a refrigeration unit for transporting perishable goods.'),
   ('Food Truck', 'Mobile vehicle equipped with facilities for cooking and selling food, often used for catering or street food sales.'),
   ('Mobile Office', 'Truck or trailer that serves as a temporary office or workspace, often used in construction sites or remote locations.'),
   ('Tank Truck', 'Specialized truck for transporting bulk liquids, gases, or hazardous materials, often with a cylindrical tank design.'),
   ('Beverage Truck', 'Truck used to transport bottled beverages, often with specialized storage and cooling features.'),
   ('Bulk Carrier', 'Truck or trailer designed to transport bulk materials like grains, coal, or other loose commodities.'),
   ('Hazmat Truck', 'Truck specially designed and equipped to transport hazardous materials or chemicals.'),
   ('Container Truck', 'Truck that is specifically designed to transport large shipping containers, often used in intermodal transport.'),
   ('High Cube Truck', 'Truck with a higher cargo area, often used for transporting larger or more bulky goods.'),
   ('Flatbed Trailer', 'Non-motorized trailer with a flat platform for transporting heavy or bulky loads that cannot be easily contained within a box.'),
   ('Gooseneck Trailer', 'Type of trailer with a "gooseneck" hitch that connects to the towing vehicle, often used for transporting large loads like RVs, boats, or livestock.'),
   ('Livestock Trailer', 'Specialized trailer for transporting livestock or animals safely.'),
   ('Enclosed Trailer', 'Trailer with an enclosed cargo space, providing weather protection for transported goods.'),
   ('Utility Trailer', 'Non-motorized trailer designed for transporting various types of goods, often used for personal or small-scale commercial use.'),
   ('Furniture Truck', 'Truck designed specifically for moving large furniture items, often with specialized loading equipment and padded cargo areas.'),
   ('Moving Van', 'Van designed for transporting household goods during relocation, often with a ramp or lift for easy loading.'),
   ('Delivery Van', 'Van specifically used for delivering packages or goods to customers, often with a compartment for storage and organization.'),
   ('Limo', 'Luxury vehicle often used for special events, transportation of VIPs, or as part of a commercial limousine service.'),
   ('Car Transporter', 'Large truck or trailer designed to carry multiple cars for transportation between dealerships or auctions.'),
   ('Forklift Truck', 'Truck used for lifting and moving heavy goods or pallets in warehouses, construction sites, and shipping yards.'),
   ('Refrigerated Box Truck', 'Truck with an enclosed cargo area and built-in refrigeration, used for transporting perishable goods.'),
   ('Tandem Truck', 'Truck with two axles at the rear, allowing for greater weight distribution and the ability to haul larger loads.'),
   ('Delivery Truck', 'Truck used for delivering goods to various destinations, often with a box or cargo area that is fully enclosed.'),
   ('Cargo Motorcycle', 'Motorcycle designed for transporting small loads, often used in courier or delivery services.'),
   ('Electric Truck', 'Truck powered by electricity rather than internal combustion engines, designed for transporting goods in a more sustainable way.'),
   ('Hydrogen Truck', 'Truck powered by hydrogen fuel cells, offering an eco-friendly alternative to traditional diesel-powered trucks.'),
   ('Autonomous Truck', 'Self-driving truck capable of operating without human intervention, often used for long-distance freight transportation.'),
   ('Liftgate Truck', 'Truck equipped with a mechanical lift at the rear to assist in loading and unloading heavy items.'),
   ('Box Trailer', 'A trailer with an enclosed cargo space, commonly used for transporting goods in secure, weatherproof conditions.'),
   ('Pole Trailer', 'A trailer designed for carrying long, oversized loads, such as large pipes or timber.'),
   ('Tandem-Trailer', 'A trailer with two or more axles positioned close together for increased weight capacity.'),
   ('Lowboy Trailer', 'A trailer with a low bed height, designed for hauling heavy equipment or oversized cargo.'),
   ('Tilt Trailer', 'Trailer with a tilted bed that can be raised or lowered for easy loading and unloading of equipment or materials.');

-- Inserting models used for Logistics and Transportation

-- Ford Models
INSERT INTO tbl_vehicle_model (brand, name, vehicle_type_id, feature) VALUES
    ('Ford', 'F-150', 1, 'Light-duty pickup truck used for cargo and small deliveries'),
    ('Ford', 'F-250 Super Duty', 2, 'Medium-duty truck, ideal for heavier loads and short-distance hauling'),
    ('Ford', 'Transit Van', 4, 'Van used for transporting small goods or packages in urban areas'),
    ('Ford', 'F-450 Super Duty', 3, 'Heavy-duty truck with a large towing capacity for long hauls'),
    ('Ford', 'Super Duty F-550', 3, 'Heavy-duty truck ideal for construction and freight logistics'),
    ('Ford', 'Transit 350 HD', 4, 'Heavy-duty version of the Transit van for bulk deliveries'),
    ('Ford', 'F-650', 3, 'Medium-heavy duty truck, often used in construction and freight transport'),
    ('Ford', 'F-750', 3, 'Heavy-duty truck used for industrial freight transport'),
    ('Ford', 'Ranger', 1, 'Compact pickup truck, commonly used for local deliveries'),
    ('Ford', 'F-350 Flatbed', 7, 'Flatbed truck used for transporting bulky or oversized goods'),
    ('Ford', 'F-550 Super Duty', 3, 'Used for carrying large loads and heavy-duty hauling'),
    ('Ford', 'F-650 Super Duty', 3, 'Truck for larger cargo hauling and freight logistics');

-- Mercedes-Benz Models
INSERT INTO tbl_vehicle_model (brand, name, vehicle_type_id, feature) VALUES
    ('Mercedes-Benz', 'Actros 1843', 3, 'Heavy-duty truck for long-distance hauling of freight'),
    ('Mercedes-Benz', 'Sprinter 2500', 4, 'Van used for parcel deliveries in urban environments'),
    ('Mercedes-Benz', 'Econic', 3, 'City delivery truck with a low entry for easy loading'),
    ('Mercedes-Benz', 'Atego 1523', 2, 'Medium-duty truck designed for regional deliveries and freight'),
    ('Mercedes-Benz', 'Vito', 4, 'Small delivery van ideal for urban logistics'),
    ('Mercedes-Benz', 'Sprinter 3500', 4, 'Larger cargo van used for high-volume urban deliveries'),
    ('Mercedes-Benz', 'Actros 2545', 3, 'Long-distance freight truck for carrying heavy cargo'),
    ('Mercedes-Benz', 'Antos', 3, 'Heavy-duty truck for transporting goods over long distances'),
    ('Mercedes-Benz', 'Atego 1223', 2, 'Light to medium truck for regional logistics'),
    ('Mercedes-Benz', 'Axor 1824', 3, 'Heavy-duty truck for industrial freight transport'),
    ('Mercedes-Benz', 'Sprinter 516', 4, 'Cargo van with high capacity for long deliveries'),
    ('Mercedes-Benz', 'Vito Panel Van', 4, 'Panel van for transporting goods and small logistics');

-- Freightliner Models
INSERT INTO tbl_vehicle_model (brand, name, vehicle_type_id, feature) VALUES
    ('Freightliner', 'Cascadia 113', 3, 'Heavy-duty truck designed for long-distance freight transport'),
    ('Freightliner', 'M2 106', 2, 'Medium-duty truck used for regional deliveries and transportation'),
    ('Freightliner', 'Columbia 120', 3, 'Long-haul truck for carrying freight across the country'),
    ('Freightliner', 'Sprinter 2500', 4, 'Van for parcel deliveries in urban environments'),
    ('Freightliner', 'M2 112', 2, 'Truck used for regional hauling of medium-weight cargo'),
    ('Freightliner', 'Cascadia 125', 3, 'Heavy-duty truck used for freight transport across regions'),
    ('Freightliner', 'Freightliner M2 106', 2, 'Commercial truck used for urban deliveries and regional transportation'),
    ('Freightliner', 'Freightliner FLD 120', 3, 'Legacy truck used for transporting freight'),
    ('Freightliner', 'Sprinter 3500', 4, 'Larger cargo van used for high-volume urban deliveries'),
    ('Freightliner', 'Freightliner 108SD', 3, 'Heavy-duty truck for heavy freight transport'),
    ('Freightliner', 'Freightliner FLD', 3, 'Older model used for regional and long-distance freight transport');

-- Volvo Models
INSERT INTO tbl_vehicle_model (brand, name, vehicle_type_id, feature) VALUES
    ('Volvo', 'FH16', 3, 'Heavy-duty truck designed for long-haul freight transport'),
    ('Volvo', 'FMX', 3, 'Heavy-duty truck built for off-road conditions and construction transport'),
    ('Volvo', 'FM', 2, 'Medium-duty truck for regional freight and deliveries'),
    ('Volvo', 'VNL 670', 3, 'Long-distance truck designed for freight hauling with comfort features'),
    ('Volvo', 'VNR', 3, 'Compact heavy-duty truck optimized for urban freight and regional distribution'),
    ('Volvo', 'V70', 4, 'Sporty wagon used for transporting smaller cargo or equipment'),
    ('Volvo', 'V60', 4, 'Mid-size estate car, used for small, last-mile cargo deliveries'),
    ('Volvo', 'Volvo FH 460', 3, 'Heavy-duty truck for carrying large freight across distances'),
    ('Volvo', 'Volvo L120H', 2, 'Large wheel loader, often used for material handling in logistics and construction'),
    ('Volvo', 'Volvo EC950F Crawler', 2, 'Construction truck used for carrying heavy equipment and materials'),
    ('Volvo', 'FH', 3, 'Heavy-duty truck for long-haul freight and cross-border logistics'),
    ('Volvo', 'VNL', 3, 'Heavy-duty truck used primarily for long-distance freight across North America'),
    ('Volvo', 'FL', 2, 'Medium-duty truck used for urban and regional logistics, deliveries, and waste collection'),
    ('Volvo', 'FE', 2, 'Medium-duty truck for carrying goods within urban environments and short regional distances'),
    ('Volvo', 'L-series', 3, 'Heavy-duty truck designed for large freight transport and long-haul logistics'),
    ('Volvo', 'XC90', 4, 'Luxury SUV used for transporting smaller high-value goods and VIP logistics'),
    ('Volvo', 'XC60', 4, 'Mid-size SUV used for small freight deliveries in urban environments'),
    ('Volvo', 'S60', 4, 'Sedan used for high-end goods and small cargo transport in urban logistics'),
    ('Volvo', 'XC40', 4, 'Compact SUV used for regional deliveries and transporting small goods'),
    ('Volvo', 'FL Electric', 2, 'Electric medium-duty truck used for environmentally friendly urban deliveries'),
    ('Volvo', 'FE Electric', 2, 'Electric truck used for urban freight transport, reducing emissions in cities'),
    ('Volvo', 'S90', 4, 'Luxury sedan used for transporting small amounts of valuable goods'),
    ('Volvo', 'XC40 Recharge', 4, 'Electric compact SUV used for regional logistics and light deliveries'),
    ('Volvo', 'V90 Cross Country', 4, 'Wagon used for transporting goods in smaller quantities over regional distances'),
    ('Volvo', 'XC70', 4, 'Crossover SUV used for transporting goods in smaller quantities, including luxury logistics');


-- MAN Models
INSERT INTO tbl_vehicle_model (brand, name, vehicle_type_id, feature) VALUES
    ('MAN', 'TGX', 3, 'Heavy-duty truck for long-distance transport and freight hauling'),
    ('MAN', 'TGM', 2, 'Medium-duty truck used for regional distribution and logistics'),
    ('MAN', 'TGL', 2, 'Light-duty truck for city logistics and urban freight transport'),
    ('MAN', 'MAN TGE', 4, 'Cargo van for deliveries and commercial transport of medium-sized cargo'),
    ('MAN', 'MAN TGS', 3, 'Heavy-duty truck designed for construction and freight transport'),
    ('MAN', 'MAN TGE 3.180', 4, 'Large cargo van designed for transporting heavier parcels or equipment'),
    ('MAN', 'MAN CLA', 3, 'Freight transport truck used for regional delivery and long-distance transport'),
    ('MAN', 'MAN TGS 18.440', 3, 'High-performance heavy-duty truck for carrying large loads and long-distance travel'),
    ('MAN', 'MAN TGE 4x4', 4, 'Off-road cargo van for deliveries in challenging environments'),
    ('MAN', 'MAN L2000', 2, 'Medium-duty truck for urban and regional transport of goods'),
    ('MAN', 'MAN TGS 26.440', 3, 'Truck for heavy-duty hauling of goods in industrial and commercial logistics');

-- Toyota Models
INSERT INTO tbl_vehicle_model (brand, name, vehicle_type_id, feature) VALUES
    ('Toyota', 'Hilux', 1, 'Light-duty pickup truck used for cargo and small deliveries'),
    ('Toyota', 'Dyna 150', 2, 'Medium-duty truck for transporting goods in urban areas and small freight hauling'),
    ('Toyota', 'Proace', 4, 'Van used for deliveries and commercial transport of medium-sized cargo'),
    ('Toyota', 'HiAce', 4, 'Versatile van for transporting goods and small equipment'),
    ('Toyota', 'Coaster', 3, 'Minibus designed for transporting people and small cargo'),
    ('Toyota', 'Toyota Tacoma', 1, 'Compact pickup truck for transporting smaller cargo loads'),
    ('Toyota', 'Tundra', 1, 'Full-size pickup truck for hauling heavy loads and transporting goods'),
    ('Toyota', 'Toyota 4Runner', 4, 'SUV used for off-road cargo transport and delivery'),
    ('Toyota', 'Proace City', 4, 'Compact van for urban deliveries and cargo transport'),
    ('Toyota', 'Toyota Tundra CrewMax', 1, 'Full-size pickup truck used for heavy cargo and freight transport'),
    ('Toyota', 'Toyota HiAce 4x4', 4, 'All-wheel-drive van used for deliveries in rural and off-road areas');

-- Hyundai Models
INSERT INTO tbl_vehicle_model (brand, name, vehicle_type_id, feature) VALUES
    ('Hyundai', 'Hyundai Porter 2', 4, 'Light-duty van designed for city deliveries and small cargo'),
    ('Hyundai', 'Hyundai Mighty', 2, 'Medium-duty truck ideal for regional deliveries and logistics'),
    ('Hyundai', 'Hyundai Xcient', 3, 'Heavy-duty truck for long-distance freight transport'),
    ('Hyundai', 'Hyundai HD170', 2, 'Medium-duty truck for both urban and regional freight delivery'),
    ('Hyundai', 'Hyundai HD78', 2, 'Light-duty truck used for transporting goods within cities and local areas'),
    ('Hyundai', 'Hyundai Santa Fe', 4, 'SUV used occasionally for transporting smaller goods'),
    ('Hyundai', 'Hyundai Tucson', 4, 'Compact SUV for small deliveries or VIP transport'),
    ('Hyundai', 'Hyundai H350', 4, 'Van used for commercial deliveries and transporting goods'),
    ('Hyundai', 'Hyundai Elantra', 4, 'Compact car used for small goods or last-mile delivery'),
    ('Hyundai', 'Hyundai Staria', 4, 'Large van designed for transporting goods and small equipment');

-- Isuzu Models
INSERT INTO tbl_vehicle_model (brand, name, vehicle_type_id, feature) VALUES
    ('Isuzu', 'Isuzu NPR', 2, 'Medium-duty truck commonly used for city deliveries and freight'),
    ('Isuzu', 'Isuzu FTR', 3, 'Heavy-duty truck used for transporting bulk goods over long distances'),
    ('Isuzu', 'Isuzu ELF', 2, 'Light-duty truck designed for small cargo delivery and city logistics'),
    ('Isuzu', 'Isuzu D-Max', 1, 'Pickup truck used for smaller cargo and regional deliveries'),
    ('Isuzu', 'Isuzu NQR', 2, 'Medium-duty truck with a wide cargo capacity for both urban and regional hauling'),
    ('Isuzu', 'Isuzu NRR', 2, 'Medium-duty truck used for transporting heavy goods in city and regional environments'),
    ('Isuzu', 'Isuzu Giga', 3, 'Heavy-duty truck used for large freight and bulk transport'),
    ('Isuzu', 'Isuzu MU-X', 4, 'SUV used for smaller logistics or high-value goods transport'),
    ('Isuzu', 'Isuzu V-Cross', 1, 'Pickup truck used for regional freight and transporting goods off-road'),
    ('Isuzu', 'Isuzu F-Series', 3, 'Heavy-duty truck designed for long-haul logistics and freight transport');
-- -- Hino Models
-- INSERT INTO tbl_vehicle_model (brand, name, vehicle_type_id, feature) VALUES
--     (9, 'Hino 300', 2, 'Medium-duty truck ideal for small to medium-size freight transport'),
--     (9, 'Hino 500', 3, 'Heavy-duty truck designed for long-distance freight hauling'),
--     (9, 'Hino 700', 3, 'Heavy-duty truck used for high-capacity freight transport'),
--     (9, 'Hino Dutro', 2, 'Light-duty truck used for transporting goods within cities and local areas'),
--     (9, 'Hino XL Series', 3, 'Heavy-duty truck for large-scale freight logistics'),
--     (9, 'Hino 300 Series', 2, 'Medium-duty truck with wide cargo area for urban and regional delivery'),
--     (9, 'Hino Ranger', 3, 'Heavy-duty truck for transporting goods over long distances'),
--     (9, 'Hino L-Series', 2, 'Light-duty truck used for logistics and regional deliveries'),
--     (9, 'Hino Hybrid', 2, 'Medium-duty truck with hybrid engine, ideal for urban and regional freight deliveries'),
--     (9, 'Hino 500 Series', 3, 'Heavy-duty truck used for international and cross-border freight hauling');
-- -- DAF Models
-- INSERT INTO tbl_vehicle_model (brand, name, vehicle_type_id, feature) VALUES
--     (10, 'DAF XF', 3, 'Heavy-duty truck for long-haul freight transport'),
--     (10, 'DAF CF', 3, 'Heavy-duty truck used for regional transport and logistics'),
--     (10, 'DAF LF', 2, 'Medium-duty truck used for urban and regional deliveries'),
--     (10, 'DAF CF Electric', 3, 'Electric heavy-duty truck designed for regional freight logistics'),
--     (10, 'DAF XF 105', 3, 'Older model used for long-haul freight and intercity transport'),
--     (10, 'DAF LF 18t', 2, 'Medium-duty truck for transporting goods in urban environments'),
--     (10, 'DAF CF 85', 3, 'Heavy-duty truck for freight transport with high engine performance'),
--     (10, 'DAF CF 75', 3, 'Lightweight truck used for freight transport in regional logistics'),
--     (10, 'DAF XF 530', 3, 'Long-haul truck used for cross-border logistics and large freight hauling'),
--     (10, 'DAF CF 480', 3, 'Heavy-duty truck with high carrying capacity for transporting bulk goods');
-- -- Scania Models
-- INSERT INTO tbl_vehicle_model (brand, name, vehicle_type_id, feature) VALUES
--     (11, 'Scania R Series', 3, 'Heavy-duty truck for long-distance and international freight transport'),
--     (11, 'Scania P Series', 2, 'Medium-duty truck used for urban deliveries and regional freight hauling'),
--     (11, 'Scania G Series', 3, 'Heavy-duty truck for long-haul and high-efficiency logistics'),
--     (11, 'Scania L Series', 2, 'Light-duty truck designed for city logistics and small cargo transport'),
--     (11, 'Scania S Series', 3, 'Top-end heavy-duty truck used for freight transport over long distances'),
--     (11, 'Scania P 230', 2, 'Medium-duty truck designed for regional logistics and deliveries'),
--     (11, 'Scania V8', 3, 'High-performance truck with large engine capacity for hauling heavy freight'),
--     (11, 'Scania R500', 3, 'Heavy-duty truck used for cross-country freight and logistics'),
--     (11, 'Scania Citywide', 3, 'City bus with cargo capabilities for inner-city freight transport'),
--     (11, 'Scania G440', 3, 'Heavy-duty truck used for international freight and transport logistics');
-- -- Kenworth Models
-- INSERT INTO tbl_vehicle_model (brand, name, vehicle_type_id, feature) VALUES
--     (12, 'Kenworth T680', 3, 'Long-haul heavy-duty truck used for regional freight transport'),
--     (12, 'Kenworth W900', 3, 'Classic heavy-duty truck for long-distance and bulk freight transport'),
--     (12, 'Kenworth T880', 3, 'Heavy-duty truck designed for construction and freight logistics'),
--     (12, 'Kenworth K270', 2, 'Medium-duty truck used for urban logistics and regional deliveries'),
--     (12, 'Kenworth T800', 3, 'Heavy-duty truck used for cross-border and long-haul freight transport'),
--     (12, 'Kenworth T370', 2, 'Light-duty truck for regional transport and local deliveries'),
--     (12, 'Kenworth C500', 3, 'Heavy-duty truck with off-road capabilities for construction and heavy freight'),
--     (12, 'Kenworth T300', 2, 'Medium-duty truck for local freight and urban logistics'),
--     (12, 'Kenworth T680 Advantage', 3, 'Long-distance truck for freight transport with fuel-efficient features'),
--     (12, 'Kenworth K370', 2, 'Medium-duty truck designed for urban deliveries and city logistics');
-- -- Iveco Models
-- INSERT INTO tbl_vehicle_model (brand, name, vehicle_type_id, feature) VALUES
--     (13, 'Iveco Stralis', 3, 'Heavy-duty truck for long-distance freight and logistics'),
--     (13, 'Iveco Daily', 4, 'Van used for medium cargo deliveries and commercial logistics'),
--     (13, 'Iveco Trakker', 3, 'Heavy-duty truck designed for construction and off-road logistics'),
--     (13, 'Iveco Eurocargo', 2, 'Medium-duty truck used for urban and regional freight transport'),
--     (13, 'Iveco S-Way', 3, 'Next-gen heavy-duty truck for fuel efficiency and long-haul logistics'),
--     (13, 'Iveco Hi-Way', 3, 'Long-distance truck for freight transport and hauling'),
--     (13, 'Iveco Magelys', 3, 'Coach used for transporting large groups or as a logistics support vehicle'),
--     (13, 'Iveco Eurotech', 3, 'Heavy-duty truck for regional freight and inter-city transport'),
--     (13, 'Iveco 682', 2, 'Medium-duty truck for regional deliveries and logistics'),
--     (13, 'Iveco TurboStar', 3, 'Truck used for high-performance freight hauling');
--
-- -- Nissan Models
-- INSERT INTO tbl_vehicle_model (brand, name, vehicle_type_id, feature) VALUES
--     (14, 'Nissan Navara', 1, 'Compact pickup truck used for regional and small cargo transport'),
--     (14, 'Nissan NV3500', 4, 'Cargo van used for urban deliveries and light freight'),
--     (14, 'Nissan Frontier', 1, 'Pickup truck used for urban and regional freight transport'),
--     (14, 'Nissan NV200', 4, 'Small van used for city logistics and delivery services'),
--     (14, 'Nissan Titan XD', 1, 'Heavy-duty pickup truck used for transporting larger cargo and regional logistics'),
--     (14, 'Nissan Atlas', 2, 'Medium-duty truck for transporting bulk goods and regional deliveries');
-- -- Mitsubishi Fuso Models
-- INSERT INTO tbl_vehicle_model (brand, name, vehicle_type_id, feature) VALUES
--     (15, 'Mitsubishi Fuso Canter', 2, 'Medium-duty truck commonly used for city and regional deliveries'),
--     (15, 'Mitsubishi Fuso Fighter', 3, 'Heavy-duty truck for long-distance freight transport'),
--     (15, 'Mitsubishi Fuso Rosa', 4, 'Minibus with cargo space used for transporting goods in urban settings'),
--     (15, 'Mitsubishi Fuso Super Great', 3, 'Heavy-duty truck used for transporting goods across long distances'),
--     (15, 'Mitsubishi Fuso FE', 2, 'Light-duty truck used for smaller cargo transport and urban logistics');
-- -- Peugeot Models
-- INSERT INTO tbl_vehicle_model (brand, name, vehicle_type_id, feature) VALUES
--     (16, 'Peugeot Boxer', 4, 'Van used for city logistics, deliveries, and transporting medium-sized goods'),
--     (16, 'Peugeot Partner', 4, 'Compact van used for urban deliveries and transporting small cargo'),
--     (16, 'Peugeot Expert', 4, 'Medium-sized van for business deliveries and logistics'),
--     (16, 'Peugeot 508 SW', 4, 'Estate car used for high-value small cargo and city deliveries');
-- -- Citroën Models
-- INSERT INTO tbl_vehicle_model (brand, name, vehicle_type_id, feature) VALUES
--     (17, 'Citroën Berlingo', 4, 'Van used for transporting small goods and deliveries in urban environments'),
--     (17, 'Citroën Jumpy', 4, 'Medium-sized van for transporting goods and business logistics'),
--     (17, 'Citroën Jumper', 4, 'Large van for business deliveries and transporting larger loads'),
--     (17, 'Citroën C4 Picasso', 4, 'Compact car with cargo space for light deliveries');
-- -- Renault Trucks Models
-- INSERT INTO tbl_vehicle_model (brand, name, vehicle_type_id, feature) VALUES
--     (18, 'Renault Master', 4, 'Cargo van for urban deliveries and transporting small to medium goods'),
--     (18, 'Renault Trafic', 4, 'Van designed for light cargo transport and regional deliveries'),
--     (18, 'Renault Trucks T', 3, 'Heavy-duty truck used for long-distance freight and logistics'),
--     (18, 'Renault Trucks K', 3, 'Off-road truck for logistics in construction and rugged environments'),
--     (18, 'Renault Trucks C', 3, 'Heavy-duty truck for urban and regional freight transport'),
--     (18, 'Renault Trucks D', 2, 'Medium-duty truck used for regional freight logistics and transportation');
-- -- Honda Models
-- INSERT INTO tbl_vehicle_model (brand, name, vehicle_type_id, feature) VALUES
--     (19, 'Honda Ridgeline', 1, 'Light-duty pickup truck used for regional and city deliveries'),
--     (19, 'Honda Acty', 4, 'Mini truck used for small deliveries in urban and rural areas'),
--     (19, 'Honda CR-V', 4, 'Compact SUV used for transporting small goods in urban settings');
-- -- Dongfeng Models
-- INSERT INTO tbl_vehicle_model (brand, name, vehicle_type_id, feature) VALUES
--     (20, 'Dongfeng KX', 2, 'Medium-duty truck used for logistics and transporting bulk goods'),
--     (20, 'Dongfeng Tianlong', 3, 'Heavy-duty truck for long-haul freight and logistics'),
--     (20, 'Dongfeng DF', 2, 'Light-duty truck used for urban deliveries and regional freight transport'),
--     (20, 'Dongfeng EQ', 3, 'Heavy-duty truck for cross-border and long-distance logistics'),
--     (20, 'Dongfeng Fengshen', 4, 'Cargo van for transporting small and medium-sized freight');
-- -- Tata Motors Models
-- INSERT INTO tbl_vehicle_model (brand, name, vehicle_type_id, feature) VALUES
--     (21, 'Tata LPT', 3, 'Heavy-duty truck used for transporting large goods and long-distance logistics'),
--     (21, 'Tata Ultra', 2, 'Medium-duty truck for regional transport and logistics'),
--     (21, 'Tata Ace', 1, 'Light-duty pickup truck used for small cargo deliveries in cities'),
--     (21, 'Tata SFC', 2, 'Medium-duty truck used for carrying freight and construction materials'),
--     (21, 'Tata Winger', 4, 'Van used for transporting goods in smaller quantities over short distances');
-- -- Mahindra Models
-- INSERT INTO tbl_vehicle_model (brand, name, vehicle_type_id, feature) VALUES
--     (22, 'Mahindra Bolero Pik-Up', 1, 'Light-duty pickup truck for small deliveries and regional transport'),
--     (22, 'Mahindra Scorpio', 4, 'SUV used for transporting small cargo and high-value goods'),
--     (22, 'Mahindra XUV700', 4, 'SUV used for delivering goods and executive transport'),
--     (22, 'Mahindra Treo', 4, 'Electric van used for city deliveries and urban logistics'),
--     (22, 'Mahindra 4110', 2, 'Medium-duty truck used for regional logistics and deliveries');
-- -- FAW Models
-- INSERT INTO tbl_vehicle_model (brand, name, vehicle_type_id, feature) VALUES
--     (23, 'FAW Jiefang', 3, 'Heavy-duty truck used for long-haul logistics and transporting bulk goods'),
--     (23, 'FAW T3', 2, 'Medium-duty truck for urban deliveries and regional logistics'),
--     (23, 'FAW V2', 4, 'Van used for transporting small goods in cities and rural areas'),
--     (23, 'FAW Qingling', 2, 'Light-duty truck used for regional deliveries and small freight transport');
-- -- Zhengzhou Yutong Models
-- INSERT INTO tbl_vehicle_model (brand, name, vehicle_type_id, feature) VALUES
--     (24, 'Yutong ZK6122H', 3, 'Bus with cargo capabilities, often used for large-scale urban freight logistics'),
--     (24, 'Yutong ZK6105', 3, 'Heavy-duty truck for logistics and freight transport across regions'),
--     (24, 'Yutong T7', 3, 'Long-distance freight transport vehicle with a focus on fuel efficiency');
-- -- Great Wall Motors Models
-- INSERT INTO tbl_vehicle_model (brand, name, vehicle_type_id, feature) VALUES
--     (25, 'Great Wall Wingle', 1, 'Pickup truck for small freight transport and urban logistics'),
--     (25, 'Great Wall Steed', 1, 'Pickup truck for regional and cross-border freight deliveries'),
--     (25, 'Great Wall Voleex', 4, 'Small van for urban logistics and business deliveries');
-- -- Kia Motors Models
-- INSERT INTO tbl_vehicle_model (brand, name, vehicle_type_id, feature) VALUES
--     (26, 'Kia Bongo', 2, 'Medium-duty truck used for regional logistics and transporting goods'),
--     (26, 'Kia K2500', 2, 'Light-duty truck used for city deliveries and transporting small freight'),
--     (26, 'Kia Soul EV', 4, 'Electric compact vehicle used for small goods and light delivery');
--
