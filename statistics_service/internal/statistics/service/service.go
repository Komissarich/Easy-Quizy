package service

import (
	"context"
	"fmt"
	api "quiz_app/pkg/api/v1"
	"quiz_app/pkg/logger"
)

type Repository interface {
	UpdateStats(context.Context, string, map[string]float32, float32) error

	GetQuizStat(context.Context, string) (*api.QuizStat, error)
	ListQuizzes(context.Context, api.ListQuizzesOption) ([]*api.QuizStat, error)

	GetPlayerStat(context.Context, string) (*api.PlayerStat, error)
	ListPlayers(context.Context, api.ListPlayersOption) ([]*api.PlayerStat, error)

	GetAuthorStat(context.Context, string) (*api.AuthorStat, error)
	ListAuthors(context.Context, api.ListAuthorsOption) ([]*api.AuthorStat, error)
}

type Service struct {
	api.StatisticsServer
	repo Repository
}

func New(ctx context.Context, repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) UpdateStats(ctx context.Context, r *api.UpdateStatsRequest) (*api.UpdateStatsResponse, error) {
	quiz_id, players_score, quiz_rate := r.QuizId, r.PlayersScore, r.QuizRate
	err := s.repo.UpdateStats(ctx, quiz_id, players_score, quiz_rate)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx, err.Error())
		return nil, err
	}
	logger.GetLoggerFromCtx(ctx).Info(ctx, fmt.Sprintf("statistics updated by quiz %s", quiz_id))
	return &api.UpdateStatsResponse{}, nil
}

func (s *Service) GetQuizStat(ctx context.Context, r *api.GetQuizStatRequest) (*api.GetQuizStatResponse, error) {
	quiz_id := r.QuizId
	quiz_stat, err := s.repo.GetQuizStat(ctx, quiz_id)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx, err.Error())
		return nil, err
	}
	logger.GetLoggerFromCtx(ctx).Info(ctx, fmt.Sprintf("quiz %s found", quiz_id))
	return &api.GetQuizStatResponse{Quiz: quiz_stat}, nil
}

func (s *Service) ListQuizzes(ctx context.Context, r *api.ListQuizzesRequest) (*api.ListQuizzesResponse, error) {
	option := r.Option
	result, err := s.repo.ListQuizzes(ctx, option)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx, err.Error())
		return nil, err
	}
	logger.GetLoggerFromCtx(ctx).Info(ctx, "quizzes listed")
	return &api.ListQuizzesResponse{Quizzes: result}, nil
}

func (s *Service) GetPlayerStat(ctx context.Context, r *api.GetPlayerStatRequest) (*api.GetPlayerStatResponse, error) {
	quiz_id := r.UserId
	player_stat, err := s.repo.GetPlayerStat(ctx, quiz_id)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx, err.Error())
		return nil, err
	}
	logger.GetLoggerFromCtx(ctx).Info(ctx, fmt.Sprintf("player %s found", quiz_id))
	return &api.GetPlayerStatResponse{Player: player_stat}, nil
}

func (s *Service) ListPlayers(ctx context.Context, r *api.ListPlayersRequest) (*api.ListPlayersResponse, error) {
	option := r.Option
	result, err := s.repo.ListPlayers(ctx, option)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx, err.Error())
		return nil, err
	}
	logger.GetLoggerFromCtx(ctx).Info(ctx, "players listed")
	return &api.ListPlayersResponse{Players: result}, nil
}

func (s *Service) GetAuthorStat(ctx context.Context, r *api.GetAuthorStatRequest) (*api.GetAuthorStatResponse, error) {
	quiz_id := r.UserId
	author_stat, err := s.repo.GetAuthorStat(ctx, quiz_id)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx, err.Error())
		return nil, err
	}
	logger.GetLoggerFromCtx(ctx).Info(ctx, fmt.Sprintf("author %s found", quiz_id))
	return &api.GetAuthorStatResponse{Author: author_stat}, nil
}

func (s *Service) ListAuthors(ctx context.Context, r *api.ListAuthorsRequest) (*api.ListAuthorsResponse, error) {
	option := r.Option
	result, err := s.repo.ListAuthors(ctx, option)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx, err.Error())
		return nil, err
	}
	logger.GetLoggerFromCtx(ctx).Info(ctx, "authors listed")
	return &api.ListAuthorsResponse{Authors: result}, nil
}
