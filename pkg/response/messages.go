package response

const (
	// Success Messages
	SuccessLogin    = "Login berhasil, selamat datang kembali"
	SuccessRegister = "Registrasi akun berhasil"
	SuccessFetch    = "Data berhasil dimuat"
	SuccessUpdate   = "Data berhasil diperbarui"
	SuccessDelete   = "Data berhasil dihapus"
	SuccessRefresh  = "Sesi berhasil diperbarui"
	SuccessCreated  = "Data berhasil ditambahkan"

	// Error Messages
	ErrInvalidRequest     = "Parameter permintaan tidak valid"
	ErrInvalidCredentials = "Kombinasi email atau password tidak sesuai"
	ErrEmailExists        = "Alamat email sudah terdaftar"
	ErrUsernameExists     = "Username sudah digunakan"
	ErrInternalServer     = "Terjadi gangguan pada server, silakan coba beberapa saat lagi"
	ErrUnauthorized       = "Akses ditolak, silakan login terlebih dahulu"
	ErrUserNotFound       = "Pengguna tidak ditemukan"
	ErrInvalidToken       = "Sesi tidak valid atau telah berakhir"
	ErrResourceNotFound   = "Data yang Anda cari tidak ditemukan"
)
