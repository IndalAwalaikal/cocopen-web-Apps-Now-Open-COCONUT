CREATE TABLE users (
    id_user INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    full_name VARCHAR(100) NOT NULL DEFAULT,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    profile_picture VARCHAR(255) DEFAULT 'default.jpg',
    role ENUM('user', 'admin') DEFAULT 'user',
    is_verified BOOLEAN DEFAULT FALSE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE password_resets (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    token VARCHAR(255) NOT NULL,
    expires_at DATETIME NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id_user)
);

CREATE TABLE email_verification_tokens (
    id_email_verification_token INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    token VARCHAR(255) NOT NULL UNIQUE,
    expires_at DATETIME NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_user_id (user_id),
    FOREIGN KEY (user_id) REFERENCES users(id_user) ON DELETE CASCADE
);

CREATE TABLE pendaftar (
    id_pendaftar INT AUTO_INCREMENT PRIMARY KEY,
    nama_lengkap VARCHAR(100) NOT NULL,
    asal_kampus VARCHAR(100) NOT NULL,
    prodi VARCHAR(100) NOT NULL,
    semester VARCHAR(20) NOT NULL,
    no_wa VARCHAR(15) NOT NULL,
    domisili VARCHAR(100) NOT NULL,
    alamat_sekarang TEXT NOT NULL,
    tinggal_dengan VARCHAR(100) NOT NULL,
    alasan_masuk TEXT NOT NULL,
    pengetahuan_coconut TEXT,
    foto_path VARCHAR(255) NOT NULL DEFAULT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    status ENUM('pending', 'diterima', 'ditolak') DEFAULT 'pending',
    user_id INT,
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id_user) ON DELETE SET NULL,
    INDEX idx_user_id (user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE jadwal (
    id_jadwal INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id_user) ON DELETE CASCADE,
    pendaftar_id INT,
    FOREIGN KEY (pendaftar_id) REFERENCES pendaftar(id_pendaftar) ON DELETE SET NULL,
    tanggal DATE NOT NULL,
    jam_mulai TIME NOT NULL,
    jam_selesai TIME NOT NULL,
    tempat VARCHAR(255) NOT NULL,
    jenis_jadwal ENUM('pribadi', 'umum') DEFAULT 'pribadi',
    konfirmasi_jadwal ENUM('belum', 'dikonfirmasi', 'ditolak') DEFAULT 'belum',
    catatan TEXT DEFAULT NULL,
    pengajuan_perubahan BOOLEAN DEFAULT FALSE,
    alasan_perubahan TEXT DEFAULT NULL,
    tanggal_diajukan DATE NULL,
    jam_mulai_diajukan TIME NULL,
    jam_selesai_diajukan TIME NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_user_id (user_id),
    INDEX idx_pendaftar_id (pendaftar_id),
    INDEX idx_tanggal (tanggal),
    INDEX idx_konfirmasi (konfirmasi_jadwal)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


-- 1. Tabel soal_test: bank soal pilihan ganda
CREATE TABLE soal_test (
    id_soal INT PRIMARY KEY AUTO_INCREMENT,
    nomor INT NOT NULL,
    pertanyaan TEXT NOT NULL,
    pilihan_a TEXT NOT NULL,
    pilihan_b TEXT NOT NULL,
    pilihan_c TEXT NOT NULL,
    pilihan_d TEXT NOT NULL,
    jawaban_benar CHAR(1) NOT NULL CHECK (jawaban_benar IN ('A', 'B', 'C', 'D')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY unique_nomor (nomor)
);

-- 2. Tabel test: konfigurasi tes (admin atur durasi, aktif/tidak)
CREATE TABLE test (
    id_test INT PRIMARY KEY AUTO_INCREMENT,
    judul VARCHAR(255) NOT NULL DEFAULT 'Tes Seleksi',
    deskripsi TEXT,
    durasi_menit INT NOT NULL DEFAULT 30, -- contoh: 30 menit
    waktu_mulai TIMESTAMP NULL, -- opsional: kapan tes dibuka
    waktu_selesai TIMESTAMP NULL, -- opsional: kapan tes ditutup
    aktif BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- 3. Tabel hasil_test: hasil ujian per user
CREATE TABLE hasil_test (
    id_hasil INT PRIMARY KEY AUTO_INCREMENT,
    user_id INT NOT NULL,
    pendaftar_id INT NOT NULL,
    id_test INT NOT NULL,
    skor_benar INT NOT NULL DEFAULT 0,
    skor_salah INT NOT NULL DEFAULT 0,
    nilai DECIMAL(5,2) NOT NULL DEFAULT 0.00,
    waktu_mulai TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    waktu_selesai TIMESTAMP, -- waktu saat user submit
    durasi_menit INT AS (TIMESTAMPDIFF(MINUTE, waktu_mulai, waktu_selesai)) STORED, -- durasi pakai
    FOREIGN KEY (user_id) REFERENCES users(id_user),
    FOREIGN KEY (pendaftar_id) REFERENCES pendaftar(id_pendaftar),
    FOREIGN KEY (id_test) REFERENCES test(id_test),
    UNIQUE KEY unique_user_per_test (user_id, id_test)
);

-- 4. Tabel jawaban_user: jawaban per soal
CREATE TABLE jawaban_user (
    id_jawaban INT PRIMARY KEY AUTO_INCREMENT,
    id_hasil INT NOT NULL,
    id_soal INT NOT NULL,
    jawaban_user CHAR(1) NOT NULL CHECK (jawaban_user IN ('A', 'B', 'C', 'D')),
    is_benar BOOLEAN NOT NULL,
    FOREIGN KEY (id_hasil) REFERENCES hasil_test(id_hasil) ON DELETE CASCADE,
    FOREIGN KEY (id_soal) REFERENCES soal_test(id_soal),
    UNIQUE KEY unique_jawaban_soal (id_hasil, id_soal)
);