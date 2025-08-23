-- MySQL dump 10.13  Distrib 8.0.43, for Linux (x86_64)
--
-- Host: 127.0.0.1    Database: dbCocopen
-- ------------------------------------------------------
-- Server version	5.5.5-10.4.32-MariaDB

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `email_verification_tokens`
--

DROP TABLE IF EXISTS `email_verification_tokens`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `email_verification_tokens` (
  `id_email_verification_token` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL,
  `token` varchar(255) NOT NULL,
  `expires_at` datetime NOT NULL,
  `created_at` datetime DEFAULT current_timestamp(),
  PRIMARY KEY (`id_email_verification_token`),
  UNIQUE KEY `token` (`token`),
  KEY `idx_user_id` (`user_id`),
  CONSTRAINT `email_verification_tokens_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id_user`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=13 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `email_verification_tokens`
--

LOCK TABLES `email_verification_tokens` WRITE;
/*!40000 ALTER TABLE `email_verification_tokens` DISABLE KEYS */;
/*!40000 ALTER TABLE `email_verification_tokens` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `hasil_test`
--

DROP TABLE IF EXISTS `hasil_test`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `hasil_test` (
  `id_hasil` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL,
  `pendaftar_id` int(11) NOT NULL,
  `id_test` int(11) NOT NULL,
  `skor_benar` int(11) NOT NULL DEFAULT 0,
  `skor_salah` int(11) NOT NULL DEFAULT 0,
  `nilai` decimal(5,2) NOT NULL DEFAULT 0.00,
  `waktu_mulai` timestamp NOT NULL DEFAULT current_timestamp(),
  `waktu_selesai` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
  `durasi_menit` int(11) GENERATED ALWAYS AS (timestampdiff(MINUTE,`waktu_mulai`,`waktu_selesai`)) STORED,
  PRIMARY KEY (`id_hasil`),
  UNIQUE KEY `unique_user_per_test` (`user_id`,`id_test`),
  KEY `pendaftar_id` (`pendaftar_id`),
  KEY `id_test` (`id_test`),
  CONSTRAINT `hasil_test_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id_user`),
  CONSTRAINT `hasil_test_ibfk_2` FOREIGN KEY (`pendaftar_id`) REFERENCES `pendaftar` (`id_pendaftar`),
  CONSTRAINT `hasil_test_ibfk_3` FOREIGN KEY (`id_test`) REFERENCES `test` (`id_test`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `hasil_test`
--

LOCK TABLES `hasil_test` WRITE;
/*!40000 ALTER TABLE `hasil_test` DISABLE KEYS */;
/*!40000 ALTER TABLE `hasil_test` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `jadwal`
--

DROP TABLE IF EXISTS `jadwal`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `jadwal` (
  `id_jadwal` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL,
  `pendaftar_id` int(11) DEFAULT NULL,
  `tanggal` date NOT NULL,
  `jam_mulai` time NOT NULL,
  `jam_selesai` time NOT NULL,
  `tempat` varchar(255) NOT NULL,
  `konfirmasi_jadwal` enum('belum','dikonfirmasi','ditolak') DEFAULT 'belum',
  `catatan` text DEFAULT NULL,
  `pengajuan_perubahan` tinyint(1) DEFAULT 0,
  `alasan_perubahan` text DEFAULT NULL,
  `created_at` datetime DEFAULT current_timestamp(),
  `updated_at` datetime DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  `tanggal_diajukan` date DEFAULT NULL,
  `jam_mulai_diajukan` time DEFAULT NULL,
  `jam_selesai_diajukan` time DEFAULT NULL,
  `jenis_jadwal` enum('pribadi','umum') DEFAULT 'pribadi',
  PRIMARY KEY (`id_jadwal`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_pendaftar_id` (`pendaftar_id`),
  KEY `idx_tanggal` (`tanggal`),
  KEY `idx_konfirmasi` (`konfirmasi_jadwal`),
  CONSTRAINT `jadwal_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id_user`) ON DELETE CASCADE,
  CONSTRAINT `jadwal_ibfk_2` FOREIGN KEY (`pendaftar_id`) REFERENCES `pendaftar` (`id_pendaftar`) ON DELETE SET NULL
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `jadwal`
--

LOCK TABLES `jadwal` WRITE;
/*!40000 ALTER TABLE `jadwal` DISABLE KEYS */;
INSERT INTO `jadwal` VALUES (4,19,NULL,'2025-04-20','14:00:00','15:00:00','Ruang Meeting 2','dikonfirmasi','Wawancara lanjutan',0,NULL,'2025-08-21 02:59:06','2025-08-21 03:04:12',NULL,NULL,NULL,'umum'),(5,19,3,'2025-04-15','09:00:00','10:30:00','Ruang Wawancara A','belum','Bawa dokumen asli',0,NULL,'2025-08-21 21:48:12','2025-08-21 21:48:12',NULL,NULL,NULL,'pribadi');
/*!40000 ALTER TABLE `jadwal` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `jawaban_user`
--

DROP TABLE IF EXISTS `jawaban_user`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `jawaban_user` (
  `id_jawaban` int(11) NOT NULL AUTO_INCREMENT,
  `id_hasil` int(11) NOT NULL,
  `id_soal` int(11) NOT NULL,
  `jawaban_user` char(1) NOT NULL CHECK (`jawaban_user` in ('A','B','C','D')),
  `is_benar` tinyint(1) NOT NULL,
  PRIMARY KEY (`id_jawaban`),
  UNIQUE KEY `unique_jawaban_soal` (`id_hasil`,`id_soal`),
  KEY `id_soal` (`id_soal`),
  CONSTRAINT `jawaban_user_ibfk_1` FOREIGN KEY (`id_hasil`) REFERENCES `hasil_test` (`id_hasil`) ON DELETE CASCADE,
  CONSTRAINT `jawaban_user_ibfk_2` FOREIGN KEY (`id_soal`) REFERENCES `soal_test` (`id_soal`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `jawaban_user`
--

LOCK TABLES `jawaban_user` WRITE;
/*!40000 ALTER TABLE `jawaban_user` DISABLE KEYS */;
/*!40000 ALTER TABLE `jawaban_user` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `password_resets`
--

DROP TABLE IF EXISTS `password_resets`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `password_resets` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL,
  `token` varchar(255) NOT NULL,
  `expires_at` datetime NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  PRIMARY KEY (`id`),
  KEY `user_id` (`user_id`),
  CONSTRAINT `password_resets_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id_user`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `password_resets`
--

LOCK TABLES `password_resets` WRITE;
/*!40000 ALTER TABLE `password_resets` DISABLE KEYS */;
INSERT INTO `password_resets` VALUES (1,7,'5096a234b2483f0d07f4ea6569bc889d3dcdd3f9944213bdcabe6ae0482ff9db','2025-08-09 17:55:30','2025-08-08 09:55:30');
/*!40000 ALTER TABLE `password_resets` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `pendaftar`
--

DROP TABLE IF EXISTS `pendaftar`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `pendaftar` (
  `id_pendaftar` int(11) NOT NULL AUTO_INCREMENT,
  `nama_lengkap` varchar(100) NOT NULL,
  `asal_kampus` varchar(100) NOT NULL,
  `prodi` varchar(100) NOT NULL,
  `semester` varchar(20) NOT NULL,
  `no_wa` varchar(15) NOT NULL,
  `domisili` varchar(100) NOT NULL,
  `alamat_sekarang` text NOT NULL,
  `tinggal_dengan` varchar(100) NOT NULL,
  `alasan_masuk` text NOT NULL,
  `pengetahuan_coconut` text DEFAULT NULL,
  `created_at` datetime DEFAULT current_timestamp(),
  `updated_at` datetime DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  `status` enum('pending','diterima','ditolak') DEFAULT 'pending',
  `foto_path` varchar(255) NOT NULL DEFAULT '',
  `user_id` int(11) DEFAULT NULL,
  PRIMARY KEY (`id_pendaftar`),
  KEY `idx_user_id` (`user_id`),
  CONSTRAINT `pendaftar_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id_user`) ON DELETE SET NULL
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `pendaftar`
--

LOCK TABLES `pendaftar` WRITE;
/*!40000 ALTER TABLE `pendaftar` DISABLE KEYS */;
INSERT INTO `pendaftar` VALUES (3,'morgan','unm','informatika','3','085757106358','btp','btp','kos','p','p','2025-08-16 06:45:48','2025-08-16 06:48:39','ditolak','',33),(4,'morgan','unm','informatika','3','085757106358','btp','btp','kos','p','p','2025-08-16 06:45:51','2025-08-16 06:46:40','ditolak','',33),(5,'nawat','ekf','njdj','2','567895678956','ksbda','jkbafe','Kos','n','d','2025-08-17 01:00:00','2025-08-17 01:00:00','pending','1755363600697069624.png',33);
/*!40000 ALTER TABLE `pendaftar` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `soal_test`
--

DROP TABLE IF EXISTS `soal_test`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `soal_test` (
  `id_soal` int(11) NOT NULL AUTO_INCREMENT,
  `nomor` int(11) NOT NULL,
  `pertanyaan` text NOT NULL,
  `pilihan_a` text NOT NULL,
  `pilihan_b` text NOT NULL,
  `pilihan_c` text NOT NULL,
  `pilihan_d` text NOT NULL,
  `jawaban_benar` char(1) NOT NULL CHECK (`jawaban_benar` in ('A','B','C','D')),
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (`id_soal`),
  UNIQUE KEY `unique_nomor` (`nomor`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `soal_test`
--

LOCK TABLES `soal_test` WRITE;
/*!40000 ALTER TABLE `soal_test` DISABLE KEYS */;
/*!40000 ALTER TABLE `soal_test` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `test`
--

DROP TABLE IF EXISTS `test`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `test` (
  `id_test` int(11) NOT NULL AUTO_INCREMENT,
  `judul` varchar(255) NOT NULL DEFAULT 'Tes Seleksi',
  `deskripsi` text DEFAULT NULL,
  `durasi_menit` int(11) NOT NULL DEFAULT 30,
  `waktu_mulai` timestamp NULL DEFAULT NULL,
  `waktu_selesai` timestamp NULL DEFAULT NULL,
  `aktif` tinyint(1) DEFAULT 1,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (`id_test`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `test`
--

LOCK TABLES `test` WRITE;
/*!40000 ALTER TABLE `test` DISABLE KEYS */;
/*!40000 ALTER TABLE `test` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `users`
--

DROP TABLE IF EXISTS `users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `users` (
  `id_user` int(11) NOT NULL AUTO_INCREMENT,
  `username` varchar(50) NOT NULL,
  `password` varchar(255) DEFAULT NULL,
  `role` enum('user','admin') DEFAULT 'user',
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  `email` varchar(255) NOT NULL,
  `is_verified` tinyint(1) DEFAULT 0,
  `full_name` varchar(100) NOT NULL DEFAULT '',
  `profile_picture` varchar(255) DEFAULT 'default.jpg',
  PRIMARY KEY (`id_user`),
  UNIQUE KEY `username` (`username`),
  UNIQUE KEY `email` (`email`)
) ENGINE=InnoDB AUTO_INCREMENT=35 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `users`
--

LOCK TABLES `users` WRITE;
/*!40000 ALTER TABLE `users` DISABLE KEYS */;
INSERT INTO `users` VALUES (7,'morgan','$2a$10$2pK7IwNmNHYV4jxlCD1I7O4d8bF8DiJKGfdndf4SM2Mq65R3yKIV2','user','2025-08-06 05:47:04','2025-08-23 09:51:05','indalawalaikal@gmail.com',1,'morgan','default.jpg'),(18,'fajrul','$2a$10$pNuQ2C9qgdiEtFGHFZml/eZjKIEwxoNmOZ2gENSv2e3Pi0kEnz7Um','user','2025-08-08 17:33:12','2025-08-23 09:51:38','indalawalaikal0055@gmail.com',0,'fajrul','default.jpg'),(19,'admin1','$2a$12$8UjwK9AY/b.9Q3F5l1BafOyO0V7KU5NAncg1OBuXbp.MG6U3H6Jmy','admin','2025-08-08 17:36:01','2025-08-15 17:02:12','',1,'admin1','default.jpg'),(27,'inferorum','$2a$10$HYkN4yHuCUmtNRNH3VhMs.8.0bTMYc.d8ztaYYdoMRNnR648Tjyfu','user','2025-08-11 13:10:18','2025-08-23 09:51:49','inferorum@gmail.com',0,'inferorum','default.jpg'),(33,'indal','$2a$10$kYuDwjAHoVhF0ytRogTz2uvDaxRY9lkOjMo80lC7/iIhVTH5J07OW','user','2025-08-15 17:56:06','2025-08-15 22:06:49','indalawalaikal05@gmail.com',1,'indal awaluddin','1755295609599036122.png'),(34,'nopal','$2a$10$YR1RZxTeuTP.UGD4A7.DUekaSpeYY9pJaN67G6WLo16dSxzxAGQPi','user','2025-08-17 13:37:13','2025-08-17 13:37:58','ngondokgaming@gmail.com',1,'nopal','default.jpg');
/*!40000 ALTER TABLE `users` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2025-08-24  2:52:45
