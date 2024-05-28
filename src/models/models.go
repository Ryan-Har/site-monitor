package models

type UserInfo struct {
	UUID  string
	Name  string
	Email string
}

// Define a key type to avoid context key collisions
type ContextKey string

const UserInfoKey ContextKey = "userInfo"
