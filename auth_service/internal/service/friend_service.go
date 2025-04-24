package service

import (
	"context"
	"eazy-quizy-auth/internal/entity"
	"eazy-quizy-auth/internal/repository"
	"eazy-quizy-auth/pkg/logger"
	"errors"
	"fmt"
	"strconv"

	"go.uber.org/zap"
)

type FriendService interface {
	AddFriend(ctx context.Context, userID uint64, friendID string) error
	RemoveFriend(ctx context.Context, userID uint64, friendID string) error

	GetFriends(ctx context.Context, userID uint64) ([]*entity.User, error)

	CheckFriendship(ctx context.Context, userID uint64, friendID string) (bool, error)
}

type friendService struct {
	friendRepo repository.FriendRepository
	userRepo   repository.UserRepository
	l          *logger.Logger
}

func NewFriendService(friendRepo repository.FriendRepository, userRepo repository.UserRepository, l *logger.Logger) *friendService {
	return &friendService{
		friendRepo: friendRepo,
		userRepo:   userRepo,
		l:          l,
	}
}

func (s *friendService) AddFriend(ctx context.Context, userID uint64, friendID string) error {
	userFriend, err := s.userRepo.FindByUsername(ctx, friendID)
	if err != nil {
		return errors.New("friend user not found")
	}

	alreadyFriends := s.friendRepo.CheckFriendship(ctx, userID, userFriend.ID)
	if !alreadyFriends {
		s.l.Error("Failed to add friend", zap.Error(errors.New("users are already friends")))
		return errors.New("users are already friends")
	}

	err = s.friendRepo.AddFriend(ctx, userID, userFriend.ID)
	if err != nil {
		s.l.Error("Failed to add friend", zap.Error(err))
		return errors.New("failed to add friend")
	}

	return nil
}

func (s *friendService) RemoveFriend(ctx context.Context, userID uint64, friendID string) error {
	userFriend, err := s.userRepo.FindByUsername(ctx, friendID)
	if err != nil {
		return errors.New("friend user not found")
	}

	areFriends := s.friendRepo.CheckFriendship(ctx, userID, userFriend.ID)
	if !areFriends {
		s.l.Error("Failed to remove friend", zap.Error(errors.New("users are not friends")))
		return errors.New("users are not friends")
	}

	if !areFriends {
		return errors.New("users are not friends")
	}

	err = s.friendRepo.RemoveFriend(ctx, userID, userFriend.ID)
	if err != nil {
		s.l.Error("Failed to remove friend", zap.Error(err))
		return errors.New("failed to remove friend")
	}

	return nil
}

func (s *friendService) GetFriends(ctx context.Context, userID uint64) ([]*entity.User, error) {
	friendIDs, err := s.friendRepo.GetFriendIDs(ctx, userID)
	if err != nil {
		s.l.Error("Failed to get friend list", zap.Error(err))
		return nil, errors.New("failed to get friend list")
	}

	var friends []*entity.User
	for _, id := range friendIDs {
		friendsID, err := strconv.Atoi(id)
		if err != nil {
			s.l.Error("Failed to get friend list", zap.Error(err))
		}
		friend, err := s.userRepo.FindByID(ctx, uint64(friendsID))
		if err != nil {
			continue
		}
		friends = append(friends, friend)
	}

	return friends, nil
}

func (s *friendService) CheckFriendship(ctx context.Context, userID uint64, friendID string) (bool, error) {
	userFriend, err := s.userRepo.FindByUsername(ctx, friendID)
	if err != nil {
		return false, fmt.Errorf("failed to check friendship: %v", err)
	}

	if ok := s.friendRepo.CheckFriendship(ctx, userID, userFriend.ID); !ok {
		return false, fmt.Errorf("failed to check friendship: %v", ok)
	}

	return true, nil
}
