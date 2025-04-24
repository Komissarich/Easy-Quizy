package service

import (
	"context"
	api "quizzes/pkg/api/v1"
	v1 "quizzes/pkg/api/v1"

	"google.golang.org/grpc"
)

type Repository interface {
	CreateQuiz(context.Context, string, string, *string, *string, []*v1.CreateQuestion) (string, string, error)
	GetQuiz(context.Context, string) (*v1.GetQuizResponse, error)
	GetQuizByAuthor(context.Context, string) (*v1.GetQuizByAuthorResponse, error)
}
type QuizService struct {
	api.QuizServiceServer
	repo Repository
}

func New(ctx context.Context, repo Repository) *QuizService {
	return &QuizService{repo: repo}
}

func Register(grpcServer *grpc.Server) {
	api.RegisterQuizServiceServer(grpcServer, &QuizService{})
}
func (s *QuizService) CreateQuiz(ctx context.Context, req *api.CreateQuizRequest) (*api.CreateQuizResponse, error) {
	quiz_id, uuID, err := s.repo.CreateQuiz(ctx, req.Name, req.Author, req.ImageId, req.Description, req.Question)
	if err != nil {
		//		logger.GetLoggerFromCtx(ctx).Error(ctx, err.Error())
		return nil, err
	}
	//	logger.GetLoggerFromCtx(ctx).Info(ctx, fmt.Sprintf("created quiz: %s", quiz_id))
	return &api.CreateQuizResponse{QuizId: uuID, ShortId: quiz_id}, nil
}
func (s *QuizService) GetQuiz(ctx context.Context, req *api.GetQuizRequest) (*api.GetQuizResponse, error) {
	resp, err := s.repo.GetQuiz(ctx, req.QuizId)
	if err != nil {
		//		logger.GetLoggerFromCtx(ctx).Error(ctx, err.Error())
		return nil, err
	}
	//	logger.GetLoggerFromCtx(ctx).Info(ctx, fmt.Sprintf("got quiz: %s", req.QuizId))
	return resp, nil
}
func (s *QuizService) GetQuizByAuthor(ctx context.Context, req *api.GetQuizByAuthorRequest) (*api.GetQuizByAuthorResponse, error) {
	resp, err := s.repo.GetQuizByAuthor(ctx, req.Author)
	if err != nil {
		//		logger.GetLoggerFromCtx(ctx).Error(ctx, err.Error())
		return nil, err
	}
	//	logger.GetLoggerFromCtx(ctx).Info(ctx, fmt.Sprintf("got quiz: %s", req.QuizId))
	return resp, nil
}
