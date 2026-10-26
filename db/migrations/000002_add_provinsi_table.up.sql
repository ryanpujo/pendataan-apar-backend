-- -- 1. Buat tabel provinsi
-- CREATE TABLE IF NOT EXISTS provinsi (
--     id SERIAL PRIMARY KEY,
--     name VARCHAR(100) NOT NULL UNIQUE,
--     created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
--     updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
-- );

-- -- 2. Tambahkan kolom provinsi_id ke tabel kota (berelasi ke provinsi.id)
-- ALTER TABLE kota
-- ADD COLUMN IF NOT EXISTS provinsi_id INTEGER REFERENCES provinsi(id) ON DELETE RESTRICT;

-- -- 3. Buat index pada foreign key agar query JOIN lebih cepat
-- CREATE INDEX IF NOT EXISTS idx_kota_provinsi_id ON kota(provinsi_id);