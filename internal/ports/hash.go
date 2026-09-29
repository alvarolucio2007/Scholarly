package ports

type PasswordHasher interface {
	Hash(password string) (string, error)
	Check(password, hashedPassword string) (bool, error)
}
