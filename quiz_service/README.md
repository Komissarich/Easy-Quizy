### Quiz service

This document describes the gRPC/REST API endpoints for the Quiz service. All endpoints are versioned under `/v1/quiz`.

---

##### Table of Contents
1. [Quiz Management](#quiz-management)
2. [Data Model Specifications](#data-model-specifications)

---

#### Quiz Management

##### `CreateQuiz`
Creates a new quiz.

**HTTP Method**: `POST /v1/quiz`    
**gRPC Method**: `QuizService.CreateQuiz`

###### Request:
```protobuf
message CreateQuizRequest {
    string name = 1;            // Quiz name (required)
    string author = 2;          // Author ID (required)
    optional string image_id = 3; // Image identifier (optional)
    optional string description = 4; // Quiz description (optional)
    repeated CreateQuestion questions = 5; // List of questions (min 1)
}

message CreateQuestion {
    string question_text = 1;     // Question text (required)
    optional string image_id = 2; // Image identifier (optional)
    repeated CreateAnswer answers = 3; // List of answers (min 2)
}

message CreateAnswer {
    string answer_text = 1; // Answer text (required)
    bool is_correct = 2;    // Marks correct answer (required)
}
```
###### Response:
```protobuf
message CreateQuizResponse {
    string quiz_id = 1; // Created quiz ID
    string message = 2; // Status message
}
```
###### Example:
```bash
curl -X POST "http://localhost:8080/v1/quiz" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "World Capitals",
    "author": "user_123",
    "description": "Test your geography knowledge",
    "questions": [
      {
        "question_text": "Capital of France",
        "answers": [
          {"answer_text": "Paris", "is_correct": true},
          {"answer_text": "London", "is_correct": false}
        ]
      }
    ]
  }'
```
##### `GetQuiz`

Retrieves quiz details by ID.

**HTTP Method**: `GET /v1/quiz/{quiz_id}`   
**gRPC Method**: `QuizService.GetQuiz`

###### Request:
```protobuf
message GetQuizRequest {
    string quiz_id = 1; // Quiz ID (required)
}
```
###### Response:
```protobuf
message GetQuizResponse {
    string name = 1;
    string author = 2;
    optional string image_id = 3;
    optional string description = 4;
    repeated CreateQuestion questions = 5; // Full question list
}
```
###### Example:
```bash
curl "http://localhost:8080/v1/quiz/quiz_abc123"
```
##### `GetQuizByAuthor`

Lists all quizzes by author.

**HTTP Method**: `GET /v1/quiz/{author}`    
**gRPC Method**: `QuizService.GetQuizByAuthor`

###### Request:
```protobuf
message GetQuizByAuthorRequest {
    string author = 1; // Author ID (required)
}
```
###### Response:
```protobuf
message GetQuizByAuthorResponse {
    repeated GetQuizResponse quizzes = 1; // All author's quizzes
}
```
###### Example:
```bash
curl "http://localhost:8080/v1/quiz/user_123"
```

#### Data Model Specifications

###### Quiz Validation:
- Minimum 1 question per quiz
- Each question must have ≥2 answers
- Exactly 1 correct answer per question

###### Image Handling:
- image_id references pre-uploaded images

Supported formats: PNG, JPEG, WebP

###### Error Codes:

- 400 Bad Request: Missing required fields

- 404 Not Found: Quiz/author not found

- 422 Unprocessable Entity: Validation errors