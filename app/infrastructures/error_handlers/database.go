package error_handlers

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

var pgErrorMessages = map[string]string{
	"23505": "Unique constraint violation (nilai duplikat untuk kolom yang seharusnya unik)",
	"23503": "Foreign key constraint violation (nilai foreign key tidak ada di tabel yang terkait)",
	"23502": "Not null violation (kolom yang tidak boleh kosong tetap menerima nilai NULL)",
	"23514": "Check constraint violation (gagal memenuhi kondisi yang ditentukan dalam CHECK constraint)",
	"22001": "String atau data terlalu panjang untuk kolom",
	"22003": "Data numerik melebihi batas ukuran tipe data",
	"42804": "Tipe data tidak cocok, biasanya terjadi saat ada konversi tipe yang tidak valid",
	"42P01": "Tabel tidak ditemukan, kemungkinan tabel yang dipanggil tidak ada",
	"42703": "Kolom tidak ditemukan, biasanya kesalahan penulisan nama kolom",
	"42601": "Kesalahan sintaks SQL, kesalahan dalam penulisan query SQL",
	"40001": "Deadlock error, terjadi saat dua atau lebih transaksi saling terkunci",
	"42P20": "Tipe data tidak konsisten dengan indeks, ini bisa terjadi pada indeks yang dibangun pada tipe data yang tidak sesuai",
	"39001": "Format data yang tidak valid, seperti format tanggal yang salah",
	"22023": "Tipe data tidak dapat dikonversi, misalnya mencoba untuk mengonversi data string ke tipe data numerik yang tidak valid",
	"22025": "String terlalu panjang untuk tipe data yang ditentukan, misalnya VARCHAR(255) menerima string lebih panjang dari 255 karakter",
	"23000": "Pelanggaran integritas data, sering terjadi saat terjadi masalah dengan constraint yang ditetapkan (misal NOT NULL, UNIQUE, dll)",
	"40003": "Terjadi batasan pada transaksi atau perubahan data, mungkin karena batasan di level transaksi atau query",
	"42000": "Kesalahan SQL terkait permission atau hak akses, terjadi saat pengguna tidak memiliki izin untuk menjalankan query",
	"54000": "Kesalahan pada query yang terlalu besar atau kompleks, bisa terjadi karena batasan dalam resource",
	"55P03": "Transaction in progress, ini berarti ada transaksi yang masih berjalan dan tidak dapat di-interrupt",
	"58000": "Kesalahan internal pada PostgreSQL, biasanya terkait dengan masalah sistem internal PostgreSQL",
	"57014": "Query dibatalkan oleh pengguna, misalnya dengan pg_cancel_backend",
	"0A000": "Fitur tidak tersedia pada versi PostgreSQL yang digunakan",
}

func ParsePGError(err error) string {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if msg, exists := pgErrorMessages[pgErr.Code]; exists {
			return msg
		}
		return "Database errors: " + pgErr.Message
	}
	return "Database errors: " + err.Error()
}
