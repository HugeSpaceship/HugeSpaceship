package db

type User struct {
	ID         int
	Name       string
	AvatarHash [20]byte
	Bio        string
}
