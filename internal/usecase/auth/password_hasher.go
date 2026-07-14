package auth

// MaxPasswordBytes — ограничение bcrypt: символы после 72 байт
const MaxPasswordBytes = 72

// PasswordHasher — порт для хэширования и проверки паролей.
type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hash, password string) error
}
