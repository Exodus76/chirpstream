package user

type Service interface {
	RegisterNewUser(email, password string) (*User, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) RegisterNewUser(email, password string) (*User, error) {

	newUser := &User{
		Email:    email,
		Password: password,
	}

	//processing logic
	return newUser, nil
}
