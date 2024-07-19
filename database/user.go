package database

type User struct {
	UserID     int64
	Name       string
	Authorized bool
	Admin      bool
	Subscribed bool
}
