package model

// Repository represents a Github repository
type Repository struct {
	Name        string
	Owner       string
	OwnerURL    string
	Description string
	URL         string
	Stars       int
	Forks       int
	Language    string
	OpenIssues  int
	CreatedAt   string
	License     string
}
