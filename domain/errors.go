package domain

import "errors"

var (
	// Database errors
	ErrDBInsertFailed = errors.New("failed to insert github token into database")

	// Token errors
	ErrTokenEmpty       = errors.New("github token is empty")
	ErrTokenTest        = errors.New("error validating github token")
	ErrTokenInvalid     = errors.New("invalid github token")
	ErrTokenForbidden   = errors.New("github token lacks required permissions")
	ErrTokenRateLimited = errors.New("github API rate limit exceeded")
	ErrTokenServerError = errors.New("github API server error")

	// Search errors
	ErrSearchNotFound = errors.New("repository not found")

	// Create errors
	ErrRepoAlreadyExists = errors.New("repository already exists")
	ErrRepoCreateFailed  = errors.New("repository creation failed")
	ErrRepoUnauthorized  = errors.New("token lacks permission to create repositories")
)
