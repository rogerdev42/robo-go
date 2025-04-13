package users

import (
	"errors"
	"lesson_05/documentstore"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Service struct {
	coll documentstore.Collection
}

func NewService(coll documentstore.Collection) *Service {
	return &Service{coll}
}

func (s *Service) CreateUser(id, name string) (*User, error) {
	userStruct := User{id, name}
	doc, _ := documentstore.MarshalDocument(userStruct)
	s.coll.Put(*doc)
	return &userStruct, nil
}

func (s *Service) ListUsers() ([]User, error) {
	list := s.coll.List()

	var users []User
	for _, user := range list {
		u := User{}
		err := documentstore.UnmarshalDocument(&user, &u)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (s *Service) GetUser(userID string) (*User, error) {
	user, _ := s.coll.Get(userID)
	if user == nil {
		return nil, ErrUserNotFound
	}
	u := User{}
	err := documentstore.UnmarshalDocument(user, &u)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (s *Service) DeleteUser(userID string) error {
	ok := s.coll.Delete(userID)
	if !ok {
		return ErrUserNotFound
	}
	return nil
}
