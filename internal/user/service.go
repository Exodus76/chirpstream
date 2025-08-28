package user

type Service interface {
	RegisterNewUser(c *CreateUserRequest) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) RegisterNewUser(c *CreateUserRequest) error {

	_ = &User{
		Email:    c.Email,
		Password: c.Password,
	}

	//processing logic
	return nil
}
