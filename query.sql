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
    jam_selesai_diajukan TIME NULL;
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_user_id (user_id),
    INDEX idx_pendaftar_id (pendaftar_id),
    INDEX idx_tanggal (tanggal),
    INDEX idx_konfirmasi (konfirmasi_jadwal)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


root:LsCMppNyBvMVaDrilcsXgEHvYOccZDYN@tcp(turntable.proxy.rlwy.net:22858)/railway?parseTime=true&loc=Local
