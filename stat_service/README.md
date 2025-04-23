### Statistics service

This document describes the gRPC/REST API endpoints for the Statistics service. All endpoints are versioned under `/v1/stats`.

---

##### Table of Contents
1. [Updating Statistics](#updating-statistics)
2. [Quiz Statistics](#quiz-statistics)
3. [Player Statistics](#player-statistics)
4. [Author Statistics](#author-statistics)

---

#### Updating Statistics

##### `UpdateStats`
Updates statistics for a quiz session.

**HTTP Method**: `POST /v1/stats/update`  
**gRPC Method**: `Statistics.UpdateStats`

###### Request:

```protobuf
message UpdateStatsRequest {
    string quiz_id = 1;               // Unique quiz identifier
    map<string, float> players_score = 2; // Map of player IDs to their scores
    float quiz_rate = 3;             // Average quiz rating (0.0–5.0)
}
```
###### Response:
```protobuf
message UpdateStatsResponse {} // Empty on success
```
###### Example:
```bash
curl -X POST "http://localhost:8080/v1/stats/update" \
  -H "Content-Type: application/json" \
  -d '{
    "quiz_id": "quiz_123",
    "players_score": {"player_1": 95.5, "player_2": 88.0},
    "quiz_rate": 4.7
  }'
```
#### Quiz Statistics
##### `GetQuizStat`
Retrieves statistics for a specific quiz.

**HTTP Method**: `GET /v1/stats/quiz/{quiz_id}`

**gRPC Method**: `Statistics.GetQuizStat`

###### Request:

```protobuf
message GetQuizStatRequest {
    string quiz_id = 1; // Quiz ID
}
```
###### Response:
```protobuf
message QuizStat {
    string quiz_id = 1;
    string author_id = 2;
    int32 num_sessions = 3; // Total sessions
    float avg_rate = 4;      // Average rating (0.0–5.0)
}
```
###### Example:
```bash
curl "http://localhost:8080/v1/stats/quiz/quiz_123"
```

##### `ListQuizzes`

Lists quizzes sorted by a specified option.

**HTTP Method**: `GET /v1/stats/quizzes/{option}`

**gRPC Method**: `Statistics.ListQuizzes`

###### Options:
```protobuf
enum ListQuizzesOption {
    AVG_RATE = 0;      // Sort by average rating (desc)
    NUM_SESSIONS = 1;  // Sort by session count (desc)
}
```

###### Request:
```protobuf
message ListQuizzesRequest {
    ListQuizzesOption option = 1;
}
```
###### Response:
```protobuf
message ListQuizzesResponse {
    repeated QuizStat quizzes = 1;
}
```
###### Example:

```bash
curl "http://localhost:8080/v1/stats/quizzes/AVG_RATE"
```
#### Player Statistics

##### `GetPlayerStat`

Retrieves statistics for a specific player.

**HTTP Method**: `GET /v1/stats/player/{user_id}`

**gRPC Method**: `Statistics.GetPlayerStat`

###### Request:
```protobuf
message GetPlayerStatRequest {
    string user_id = 1; // Player ID
}
```
###### Response:
```protobuf
message PlayerStat {
    string user_id = 1;
    float total_score = 2;
    float best_score = 3;
    float avg_score = 4;
    int32 num_sessions = 5;
}
```
###### Example:
```bash
curl "http://localhost:8080/v1/stats/player/player_123"
```
##### `ListPlayers`

Lists players sorted by a specified option.

**HTTP Method**: `GET /v1/stats/players`

**gRPC Method**: `Statistics.ListPlayers`

###### Options:
```protobuf
enum ListPlayersOption {
    TOTAL_SCORE = 0; // Sort by total score (desc)
    BEST_SCORE = 1;  // Sort by best score (desc)
    AVG_SCORE = 2;   // Sort by average score (desc)
}
```
###### Request:
```protobuf
message ListPlayersRequest {
    ListPlayersOption option = 1;
}
```
###### Response:
```protobuf
message ListPlayersResponse {
    repeated PlayerStat players = 1;
}
```
###### Example:
```bash
curl "http://localhost:8080/v1/stats/players?option=BEST_SCORE"
```
#### Author Statistics

##### `GetAuthorStat`

Retrieves statistics for a specific author.

**HTTP Method**: `GET /v1/stats/author/{user_id}`

**gRPC Method**: `Statistics.GetAuthorStat`

###### Request:
```protobuf
message GetAuthorStatRequest {
    string user_id = 1; // Author ID
}
```
###### Response:
```protobuf
message AuthorStat {
    string user_id = 1;
    int32 num_quizzes = 2;       // Total quizzes created
    float avg_quiz_rate = 3;     // Average quiz rating (0.0–5.0)
    float best_quiz_rate = 4;    // Highest quiz rating (0.0–5.0)
}
```
###### Example:
```bash
curl "http://localhost:8080/v1/stats/author/author_123"
```

##### `ListAuthors`

Lists authors sorted by a specified option.

**HTTP Method**: `GET /v1/stats/authors`

**gRPC Method**: `Statistics.ListAuthors`

###### Options:
```protobuf
enum ListAuthorsOption {
    NUM_QUIZZES = 0;    // Sort by quiz count (desc)
    AVG_QUIZ_RATE = 1;  // Sort by average rating (desc)
    BEST_QUIZ_RATE = 2; // Sort by best rating (desc)
}
```

###### Request:
```protobuf
message ListAuthorsRequest {
    ListAuthorsOption option = 1;
}
```

###### Response:
```protobuf
message ListAuthorsResponse {
    repeated AuthorStat authors = 1;
}
```
###### Example:
```bash
curl "http://localhost:8080/v1/stats/authors?option=AVG_QUIZ_RATE"
```
