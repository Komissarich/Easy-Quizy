package service

import (
	"context"
	"eazy-quizy-auth/internal/entity"
	"eazy-quizy-auth/internal/repository"
	"eazy-quizy-auth/pkg/logger"
	"errors"
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

	_, err := s.userRepo.FindByUsername(ctx, friendID)
	if err != nil {
		return errors.New("friend user not found")
	}

	alreadyFriends, err := s.friendRepo.CheckFriendship(ctx, userID, friendID)
	if err != nil {
		s.l.Error("Failed to check friendship", zap.Error(err))
		return errors.New("failed to check friendship")
	}

	if alreadyFriends {
		return errors.New("users are already friends")
	}
	friendIDUint, _ := strconv.Atoi(friendID)

	err = s.friendRepo.AddFriend(ctx, userID, uint64(friendIDUint))
	if err != nil {
		s.l.Error("Failed to add friend", zap.Error(err))
		return errors.New("failed to add friend")
	}

	return nil
}

func (s *friendService) RemoveFriend(ctx context.Context, userID uint64, friendID string) error {
	areFriends, err := s.friendRepo.CheckFriendship(ctx, userID, friendID)
	if err != nil {
		s.l.Error("Failed to check friendship", zap.Error(err))
		return errors.New("failed to check friendship")
	}

	if !areFriends {
		return errors.New("users are not friends")
	}

	err = s.friendRepo.RemoveFriend(ctx, userID, friendID)
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
	return s.friendRepo.CheckFriendship(ctx, userID, friendID)
}
