package users

import (
	"errors"
	"fmt"
	"lesson_05/documentstore"
)

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrUserCreateFailed = errors.New("failed to create user")
	ErrUserListFailed   = errors.New("failed to list users")
	ErrUserGetFailed    = errors.New("failed to get user")
	ErrUserDeleteFailed = errors.New("failed to delete user")
)

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Service struct {
	coll documentstore.Collection
}

func NewService() *Service {
	store := documentstore.NewStore()
	coll, _ := store.CreateCollection("users", &documentstore.CollectionConfig{PrimaryKey: "ID"})
	return &Service{*coll}
}

func (s *Service) CreateUser(id, name string) (*User, error) {
	userStruct := User{id, name}
	doc, err := documentstore.MarshalDocument(userStruct)
	if err != nil {
		return nil, fmt.Errorf(ErrUserCreateFailed.Error()+": %w", err)
	}

	err = s.coll.Put(*doc)
	if err != nil {
		return nil, fmt.Errorf(ErrUserCreateFailed.Error()+": %w", err)
	}

	return &userStruct, nil
}

func (s *Service) ListUsers() ([]User, error) {
	list := s.coll.List()

	var users []User
	for _, user := range list {
		u := User{}
		err := documentstore.UnmarshalDocument(&user, &u)
		if err != nil {
			return nil, fmt.Errorf(ErrUserListFailed.Error()+": %w", err)
		}
		users = append(users, u)
	}
	return users, nil
}

func (s *Service) GetUser(userID string) (*User, error) {
	user, err := s.coll.Get(userID)
	if err != nil {
		return nil, fmt.Errorf(ErrUserNotFound.Error()+": %w", err)
	}
	u := User{}
	err = documentstore.UnmarshalDocument(user, &u)
	if err != nil {
		return nil, fmt.Errorf(ErrUserGetFailed.Error()+": %w", err)
	}
	return &u, nil
}

func (s *Service) DeleteUser(userID string) error {
	err := s.coll.Delete(userID)
	if err != nil {
		return fmt.Errorf(ErrUserDeleteFailed.Error()+": %w", err)
	}
	return nil
}
