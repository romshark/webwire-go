package shared

// AuthenticationCredentials represents the credentials of a user
type AuthenticationCredentials struct {
	Name string `json:"name"`
	Password string `json:"pass"`
}
