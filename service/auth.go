package service

import "github.com/luisedmc/ghcmd/domain"

// AuthProvider abstracts GitHub token validation.
type AuthProvider interface {
	ValidateToken(token string) (*domain.User, error)
}

// TokenStore abstracts token persistence.
type TokenStore interface {
	ReadToken() (string, error)
	SaveToken(token string) error
	Close() error
}

// AuthService orchestrates authentication: validating tokens and persisting them.
type AuthService struct {
	auth  AuthProvider
	store TokenStore
}

func NewAuthService(auth AuthProvider, store TokenStore) *AuthService {
	return &AuthService{auth: auth, store: store}
}

// ValidateToken checks whether a token is valid and returns the associated user.
func (s *AuthService) ValidateToken(token string) (*domain.User, error) {
	if token == "" {
		return nil, domain.ErrTokenEmpty
	}
	return s.auth.ValidateToken(token)
}

// ValidateAndStore validates the token, persists it on success, and returns the user.
func (s *AuthService) ValidateAndStore(token string) (*domain.User, error) {
	user, err := s.ValidateToken(token)
	if err != nil {
		return nil, err
	}
	if storeErr := s.store.SaveToken(token); storeErr != nil {
		return nil, storeErr
	}
	return user, nil
}

// LoadStoredToken reads the persisted token from the store.
func (s *AuthService) LoadStoredToken() (string, error) {
	return s.store.ReadToken()
}

// Close releases resources held by the token store.
func (s *AuthService) Close() error {
	return s.store.Close()
}
