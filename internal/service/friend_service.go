package service

import (
	"context"
	"eazy-quizy-auth/internal/entity"
	"eazy-quizy-auth/internal/repository"
	"errors"
	"log"
)

type FriendService interface {
	AddFriend(ctx context.Context, userID, friendID string) error
	RemoveFriend(ctx context.Context, userID, friendID string) error
	GetFriends(ctx context.Context, userID string) ([]*entity.User, error)
	CheckFriendship(ctx context.Context, userID, friendID string) (bool, error)
}

type friendService struct {
	friendRepo repository.FriendRepository
	userRepo   repository.UserRepository
}

func NewFriendService(friendRepo repository.FriendRepository, userRepo repository.UserRepository) *friendService {
	return &friendService{
		friendRepo: friendRepo,
		userRepo:   userRepo,
	}
}

func (s *friendService) AddFriend(ctx context.Context, userID, friendID string) error {
	if userID == friendID {
		return errors.New("cannot add yourself as a friend")
	}

	_, err := s.userRepo.FindByID(ctx, friendID)
	if err != nil {
		return errors.New("friend user not found")
	}

	alreadyFriends, err := s.friendRepo.CheckFriendship(ctx, userID, friendID)
	if err != nil {
		log.Printf("Error checking friendship: %v", err)
		return errors.New("failed to check friendship")
	}

	if alreadyFriends {
		return errors.New("users are already friends")
	}

	err = s.friendRepo.AddFriend(ctx, userID, friendID)
	if err != nil {
		log.Printf("Error adding friend: %v", err)
		return errors.New("failed to add friend")
	}

	return nil
}

func (s *friendService) RemoveFriend(ctx context.Context, userID, friendID string) error {
	areFriends, err := s.friendRepo.CheckFriendship(ctx, userID, friendID)
	if err != nil {
		log.Printf("Error checking friendship: %v", err)
		return errors.New("failed to check friendship")
	}

	if !areFriends {
		return errors.New("users are not friends")
	}

	err = s.friendRepo.RemoveFriend(ctx, userID, friendID)
	if err != nil {
		log.Printf("Error removing friend: %v", err)
		return errors.New("failed to remove friend")
	}

	return nil
}

// GetFriends возвращает список друзей пользователя
func (s *friendService) GetFriends(ctx context.Context, userID string) ([]*entity.User, error) {
	friendIDs, err := s.friendRepo.GetFriendIDs(ctx, userID)
	if err != nil {
		log.Printf("Error getting friend IDs: %v", err)
		return nil, errors.New("failed to get friend list")
	}

	var friends []*entity.User
	for _, id := range friendIDs {
		friend, err := s.userRepo.FindByID(ctx, id)
		if err != nil {
			continue
		}
		friends = append(friends, friend)
	}

	return friends, nil
}

func (s *friendService) CheckFriendship(ctx context.Context, userID, friendID string) (bool, error) {
	return s.friendRepo.CheckFriendship(ctx, userID, friendID)
}
