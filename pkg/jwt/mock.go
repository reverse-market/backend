package jwt

const (
	validToken            = "VALID_TOKEN"
	expiredToken          = "EXPIRED_TOKEN"
	tokenEIthWrongSubject = "WRONG_SUBJECT_TOKEN"
)

type MockManager struct {
}

func (m *MockManager) CreateToken(id int) (string, error) {
	return validToken, nil
}

func (m *MockManager) GetIdFromToken(token string) (int, error) {
	switch token {
	case expiredToken:
		return 0, ErrExpiredToken
	case tokenEIthWrongSubject:
		return 0, ErrInvalidSubject
	case validToken:
		return 1, nil
	default:
		return 0, ErrInvalidToken
	}
}
