package service

import (
	"context"
	"fmt"
	api "quiz_app/pkg/api/v1"
	"quiz_app/pkg/logger"
)

type Repository interface {
	CreateSession(context.Context, string, string, []*api.Result) (string, error)
	GetSession(context.Context, string) (*api.Session, error)
	ListSessions(context.Context) ([]*api.Session, error)
}

type Service struct {
	api.StatisticsServer
	repo Repository
}

func (s *Service) CreateSession(ctx context.Context, r *api.CreateSessionRequest) (*api.CreateSessionResponse, error) {
	quiz_id, start_time, end_time, results := r.QuizId, r.StartTime, r.EndTime, r.Results
	session_id, err := s.repo.CreateSession(ctx, quiz_id, start_time, end_time, results)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx, err.Error())
		return nil, err
	}
	logger.GetLoggerFromCtx(ctx).Info(ctx, fmt.Sprintf("new session created with id %s", session_id))
	return &api.CreateSessionResponse{SessionId: session_id}, nil
}

func (s *Service) GetSession(ctx context.Context, r *api.GetSessionRequest) (*api.GetSessionResponse, error) {
	session_id := r.SessionId
	session, err := s.repo.GetSession(ctx, session_id)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx, err.Error())
		return nil, err
	}
	logger.GetLoggerFromCtx(ctx).Info(ctx, fmt.Sprintf("session %s found", session_id))
	return &api.GetSessionResponse{Session: session}, nil
}

func (s *Service) ListSessions(ctx context.Context, r *api.ListSessionsRequest) (*api.ListSessionsResponse, error) {
	sessions, err := s.repo.ListSessions(ctx)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx, err.Error())
		return nil, err
	}
	logger.GetLoggerFromCtx(ctx).Info(ctx, fmt.Sprintf("sessions listed"))
	return &api.ListSessionsResponse{Sessions: sessions}, nil
}
