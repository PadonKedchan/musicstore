CREATE TABLE store_info (
    id SERIAL PRIMARY KEY,
    logo_path VARCHAR(255), -- ที่เก็บเส้นทางโลโก้ของร้าน เช่น URL หรือพาธในระบบไฟล์
    store_name VARCHAR(255) NOT NULL,
    description TEXT,
    address VARCHAR(255),
    phone_number VARCHAR(20),
    email VARCHAR(100)
);


INSERT INTO store_info (logo_path, store_name, description, address, phone_number, email) VALUES
('/images/store_logo1.jpg', 'Vinyl Paradise', 'ร้านแผ่นเสียงและอุปกรณ์ดนตรีคุณภาพ นำเข้าจากต่างประเทศ'
, '6 ราชมรรคาใน ตำบลพระปฐมเจดีย์ อำเภอเมืองนครปฐม นครปฐม 73000', '02-123-4567', 'vinylparadise@gmail.com'),

('/images/store_logo2.jpg', 'Melody Master', 'ศูนย์รวมเครื่องดนตรีคุณภาพ'
, '6 ราชมรรคาใน ตำบลพระปฐมเจดีย์ อำเภอเมืองนครปฐม นครปฐม 73000', '02-987-6543', 'melodymaster@gmail.com'),

('/images/store_logo3.jpg', 'Vintage Vinyl', 'ร้านแผ่นเสียงมือสองคุณภาพเยี่ยม ร้านขายเคสโทรศัพท์และเคสไอแพดลายน่ารักสดสัย สีของเคสโทรศัพท์และเคสไอแพดจะมีสีโทนเย็นทุกรูปแบบ 
มีให้เลือกมากมาย สามารถซื้อได้ในราคาย่อมเยา มีให้เลือกหลานรุ่นหลายยี่ห้อ สามารถมาจับจองได้แล้วที่นี่', '6 ราชมรรคาใน ตำบลพระปฐมเจดีย์ อำเภอเมืองนครปฐม นครปฐม 73000', '02-765-4321', 'vintagevinyl@gmail.com'),

('/images/store_logo4.jpg', 'Sound Studio', 'ศูนย์รวมอุปกรณ์สตูดิโอ', '6 ราชมรรคาใน ตำบลพระปฐมเจดีย์ อำเภอเมืองนครปฐม นครปฐม 73000', '02-555-6789', 'soundstudio@gmail.com'),

('/images/store_logo5.jpg', 'Harmony Hub', 'ร้านเครื่องดนตรีครบวงจร',
 '6 ราชมรรคาใน ตำบลพระปฐมเจดีย์ อำเภอเมืองนครปฐม นครปฐม 73000', '02-345-6789', 'harmonyhub@gmail.com');


CREATE TABLE product_info (
    id SERIAL PRIMARY KEY,
    product_name VARCHAR(255) NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    quantity INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    category VARCHAR(100),
    brand VARCHAR(100),
    model VARCHAR(100),
    store_id INT REFERENCES store_info(id),
    is_recommended BOOLEAN DEFAULT FALSE,
    image_path VARCHAR(255) NOT NULL
);




-- Insert Products for Store 1 (ร้านขายแผ่นเสียง)
-- Store 1: ร้านขายแผ่นเสียง (5 unique products, at least 3 recommended)
INSERT INTO product_info (product_name, price, quantity, category, brand, model, store_id, is_recommended, image_path)
VALUES
    ('The Beatles - Abbey Road Vinyl', 1200.00, 15, 'แผ่นเสียง', 'The Beatles', 'Abbey Road', 1, TRUE, '/images/products/Abbey_Road_Vinyl.png'),
    ('Pink Floyd - The Dark Side of the Moon Vinyl', 1500.00, 10, 'แผ่นเสียง', 'Pink Floyd', 'The Dark Side of the Moon', 1, TRUE, '/images/products/Dark_Side_Vinyl.png'),
    ('Nirvana - Nevermind Vinyl', 1000.00, 8, 'แผ่นเสียง', 'Nirvana', 'Nevermind', 1, TRUE, '/images/products/Nevermind_Vinyl.png'),
    ('Led Zeppelin - IV Vinyl', 1800.00, 5, 'แผ่นเสียง', 'Led Zeppelin', 'IV', 1, FALSE, '/images/products/Led_Zeppelin_IV.png'),
    ('AC/DC - Back in Black Vinyl', 1400.00, 12, 'แผ่นเสียง', 'AC/DC', 'Back in Black', 1, FALSE, '/images/products/Back_in_Black_Vinyl.png');

-- Insert Products for Store 2 (ร้านขายกีตาร์)
-- Store 2: ร้านขายกีตาร์ (5 unique products, at least 3 recommended)
INSERT INTO product_info (product_name, price, quantity, category, brand, model, store_id, is_recommended, image_path)
VALUES
    ('Fender Stratocaster Electric Guitar', 25000.00, 10, 'กีตาร์ไฟฟ้า', 'Fender', 'Stratocaster', 2, TRUE, '/images/products/Fender_Stratocaster.png'),
    ('Gibson Les Paul Standard Guitar', 45000.00, 5, 'กีตาร์ไฟฟ้า', 'Gibson', 'Les Paul Standard', 2, TRUE, '/images/products/Gibson_Les_Paul_Standard.png'),
    ('Ibanez RG550 Electric Guitar', 20000.00, 8, 'กีตาร์ไฟฟ้า', 'Ibanez', 'RG550', 2, TRUE, '/images/products/Ibanez_RG550.png'),
    ('Yamaha Pacifica Electric Guitar', 15000.00, 12, 'กีตาร์ไฟฟ้า', 'Yamaha', 'Pacifica 112V', 2, FALSE, '/images/products/Yamaha_Pacifica.png'),
    ('PRS SE Custom 24 Electric Guitar', 25000.00, 6, 'กีตาร์ไฟฟ้า', 'PRS', 'SE Custom 24', 2, FALSE, '/images/products/PRS_SE_Custom24.png');

-- Insert Products for Store 3 (ร้านขายเบส)
-- Store 3: ร้านขายเบส (5 unique products, at least 3 recommended)
INSERT INTO product_info (product_name, price, quantity, category, brand, model, store_id, is_recommended, image_path)
VALUES
    ('Fender Jazz Bass', 25000.00, 10, 'เบสไฟฟ้า', 'Fender', 'Jazz Bass', 3, TRUE, '/images/products/Fender_Jazz_Bass.png'),
    ('Music Man StingRay Bass', 35000.00, 5, 'เบสไฟฟ้า', 'Music Man', 'StingRay', 3, TRUE, '/images/products/Music_Man_StingRay_Bass.png'),
    ('Gibson Thunderbird Bass', 40000.00, 7, 'เบสไฟฟ้า', 'Gibson', 'Thunderbird', 3, TRUE, '/images/products/Gibson_Thunderbird_Bass.png'),
    ('Ibanez SR300E Bass', 15000.00, 8, 'เบสไฟฟ้า', 'Ibanez', 'SR300E', 3, FALSE, '/images/products/Ibanez_SR300E_Bass.png'),
    ('Yamaha TRBX504 Bass', 18000.00, 6, 'เบสไฟฟ้า', 'Yamaha', 'TRBX504', 3, FALSE, '/images/products/Yamaha_TRBX504_Bass.png');

-- Insert Products for Store 4 (ร้านขายกลอง)
-- Store 4: ร้านขายกลอง (5 unique products, at least 3 recommended)
INSERT INTO product_info (product_name, price, quantity, category, brand, model, store_id, is_recommended, image_path)
VALUES
    ('Roland TD-27KV Drum Kit', 89000.00, 8, 'กลองไฟฟ้า', 'Roland', 'TD-27KV', 4, TRUE, '/images/products/Roland_TD-27KV.png'),
    ('Pearl Roadshow Drum Kit', 19000.00, 7, 'กลองชุด', 'Pearl', 'Roadshow', 4, TRUE, '/images/products/Pearl_Roadshow1.png'),
    ('Tama Imperialstar Drum Kit', 28000.00, 5, 'กลองชุด', 'Tama', 'Imperialstar', 4, TRUE, '/images/products/Tama_Imperialstar.png'),
    ('Ludwig Breakbeats Drum Kit', 24000.00, 10, 'กลองชุด', 'Ludwig', 'Breakbeats', 4, FALSE, '/images/products/Ludwig_Breakbeats.png'),
    ('Yamaha Stage Custom Drum Kit', 35000.00, 6, 'กลองชุด', 'Yamaha', 'Stage Custom', 4, FALSE, '/images/products/Yamaha_Stage_Custom.png');

-- Insert Products for Store 5 (ร้านขายลำโพง)
-- Store 5: ร้านขายลำโพง (5 unique products, at least 3 recommended)
INSERT INTO product_info (product_name, price, quantity, category, brand, model, store_id, is_recommended, image_path)
VALUES
    ('JBL Flip 5 Bluetooth Speaker', 5000.00, 20, 'ลำโพงบลูทูธ', 'JBL', 'Flip 5', 5, TRUE, '/images/products/JBL_Flip5.png'),
    ('Bose SoundLink Revolve Bluetooth Speaker', 12000.00, 15, 'ลำโพงบลูทูธ', 'Bose', 'SoundLink Revolve', 5, TRUE, '/images/products/Bose_SoundLink.png'),
    ('Sonos One Smart Speaker', 10000.00, 10, 'ลำโพงสมาร์ท', 'Sonos', 'One', 5, TRUE, '/images/products/Sonos_One.png'),
    ('Marshall Stanmore II Bluetooth Speaker', 9000.00, 8, 'ลำโพงบลูทูธ', 'Marshall', 'Stanmore II', 5, FALSE, '/images/products/Marshall_Stanmore.png'),
    ('Sony SRS-XB43 Bluetooth Speaker', 8000.00, 12, 'ลำโพงบลูทูธ', 'Sony', 'SRS-XB43', 5, FALSE, '/images/products/Sony_SRS_XB43.png');







CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_timestamp
BEFORE UPDATE ON product_info
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();


CREATE TABLE cart (
    id SERIAL PRIMARY KEY,
    store_id INT NOT NULL,  -- อ้างอิงถึงร้านค้าจาก store_info
    product_id INT NOT NULL,  -- รหัสสินค้า
    quantity INT NOT NULL,  -- จำนวนสินค้าในรถเข็น
    added_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,  -- เวลาที่สินค้าได้รับการเพิ่มเข้าไปในรถเข็น
    checked_out_at TIMESTAMP,  -- เวลาเช็คเอาต์ (เมื่อมีการจ่ายเงินหรือทำการเช็คเอาต์)
    status VARCHAR(50) DEFAULT 'in_cart',  -- สถานะของสินค้า (in_cart, ordered)
    FOREIGN KEY (store_id) REFERENCES store_info(id),  -- อ้างอิงถึง store_info(id)
    FOREIGN KEY (product_id) REFERENCES product_info(id)  -- อ้างอิงถึงสินค้าในตาราง product_info
);
