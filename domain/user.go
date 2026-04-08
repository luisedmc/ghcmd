package domain

// User represents an authenticated Github user.
type User struct {
	Login       string
	Name        string
	PublicRepos int
}
