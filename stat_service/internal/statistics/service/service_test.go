package service_test

import (
	"testing"
)

//  service.UpdateStats(ctx, &api.UpdateStatsRequest{QuizId: "123", AuthorId: "234", PlayersScore: map[string]float32{"ivan": 4.5, "masha": 2.5}, QuizRate: 4.5})
// 	service.UpdateStats(ctx, &api.UpdateStatsRequest{QuizId: "435", AuthorId: "789", PlayersScore: map[string]float32{"ivan": 4.5, "igor": 3.5}, QuizRate: 2.5})
// 	service.UpdateStats(ctx, &api.UpdateStatsRequest{QuizId: "2657", AuthorId: "234", PlayersScore: map[string]float32{"igor": 1.5, "masha": 2.5}, QuizRate: 1.5})

// 	service.GetQuizStat(ctx, &api.GetQuizStatRequest{QuizId: "123"})
// 	service.ListQuizzes(ctx, &api.ListQuizzesRequest{Option: 0})
// 	service.ListQuizzes(ctx, &api.ListQuizzesRequest{Option: 1})
// 	service.GetAuthorStat(ctx, &api.GetAuthorStatRequest{UserId: "234"})
// 	service.ListAuthors(ctx, &api.ListAuthorsRequest{Option: 0})
// 	service.ListAuthors(ctx, &api.ListAuthorsRequest{Option: 1})
// 	service.ListAuthors(ctx, &api.ListAuthorsRequest{Option: 2})
// 	service.GetPlayerStat(ctx, &api.GetPlayerStatRequest{UserId: "ivan"})
// 	service.ListPlayers(ctx, &api.ListPlayersRequest{Option: 0})
// 	service.ListPlayers(ctx, &api.ListPlayersRequest{Option: 1})
// 	service.ListPlayers(ctx, &api.ListPlayersRequest{Option: 2})

func TestUpdateStats(t *testing.T) {

}

func TestGetQuizStat(t *testing.T) {

}

func TestListQuizzes(t *testing.T) {

}

func TestGetPlayerStat(t *testing.T) {

}

func TestListPlayers(t *testing.T) {

}

func TestGetAuthorStat(t *testing.T) {

}

func TestListAuthors(t *testing.T) {

}
