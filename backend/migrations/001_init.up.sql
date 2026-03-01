-- migrations/001_init.up.sql

CREATE TABLE pengguna (
    id SERIAL PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    kata_sandi_hash TEXT NOT NULL,
    nama TEXT NOT NULL,
    zona_waktu TEXT NOT NULL,
    dibuat_pada TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);

CREATE TABLE pot (
    id SERIAL PRIMARY KEY,
    id_pengguna INTEGER REFERENCES pengguna(id) ON DELETE CASCADE,
    nama_pot TEXT NOT NULL,
    nama_tanaman TEXT NOT NULL,
    varietas TEXT,
    ukuran_pot TEXT,
    media_tanam TEXT,
    lokasi TEXT,
    tanggal_mulai DATE NOT NULL DEFAULT current_date,
    dibuat_pada TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);

CREATE TABLE siklus (
    id SERIAL PRIMARY KEY,
    id_pot INTEGER REFERENCES pot(id) ON DELETE CASCADE,
    tanggal_mulai DATE NOT NULL,
    status TEXT NOT NULL DEFAULT 'aktif',
    alasan_arsip TEXT,
    diarsipkan_pada TIMESTAMP WITH TIME ZONE
);

CREATE TABLE tahap (
    id SERIAL PRIMARY KEY,
    id_siklus INTEGER REFERENCES siklus(id) ON DELETE CASCADE,
    nama_tahap TEXT NOT NULL,
    urutan INTEGER NOT NULL,
    durasi_hari INTEGER NOT NULL
);

CREATE TABLE tugas (
    id SERIAL PRIMARY KEY,
    id_siklus INTEGER REFERENCES siklus(id) ON DELETE CASCADE,
    judul TEXT NOT NULL,
    tipe TEXT NOT NULL,
    jadwal_json JSONB NOT NULL,
    next_due_at TIMESTAMP WITH TIME ZONE,
    last_completed_at TIMESTAMP WITH TIME ZONE,
    aktif BOOLEAN NOT NULL DEFAULT true,
    dibuat_pada TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);

CREATE TABLE catatan (
    id SERIAL PRIMARY KEY,
    id_siklus INTEGER REFERENCES siklus(id) ON DELETE CASCADE,
    id_pot INTEGER REFERENCES pot(id) ON DELETE CASCADE,
    dibuat_pada TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    catatan_teks TEXT,
    foto_urls_json JSONB,
    id_tugas_terkait INTEGER
);

CREATE TABLE insiden (
    id SERIAL PRIMARY KEY,
    id_siklus INTEGER REFERENCES siklus(id) ON DELETE CASCADE,
    dibuka_pada TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    ditutup_pada TIMESTAMP WITH TIME ZONE,
    jenis TEXT NOT NULL,
    tingkat INTEGER NOT NULL,
    status TEXT NOT NULL DEFAULT 'buka',
    catatan TEXT,
    foto_urls_json JSONB
);

CREATE TABLE panen (
    id SERIAL PRIMARY KEY,
    id_siklus INTEGER REFERENCES siklus(id) ON DELETE CASCADE,
    tanggal DATE NOT NULL,
    jumlah NUMERIC NOT NULL,
    satuan TEXT NOT NULL,
    catatan TEXT,
    grade TEXT
);
