package users

type Service interface {
	ListUsers() (*[]UserResponse, error)
	//GetUserByID()
	//DeleteUserByID()
	//UpdateUserByID()
}
type UserService struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &UserService{repo}
}

func (s *UserService) ListUsers() (*[]UserResponse, error) {

	listUsers, err := s.repo.ListUsers()
	if err != nil {
		return nil, err
	}
	return listUsers, nil

}
