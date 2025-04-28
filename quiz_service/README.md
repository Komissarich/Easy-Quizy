# QuizService API Documentation

## Overview
QuizService provides CRUD operations for quizzes with questions and answers.

Base URL: `/v1/quiz`

## Endpoints

### Create Quiz
`POST /v1/quiz`

**Request**:
```json
{
  "name": "string",
  "author": "string",
  "image_id": "string (optional)",
  "description": "string (optional)",
  "question": [
    {
      "question_text": "string",
      "image_id": "string (optional)",
      "answer": [
        {
          "answer_text": "string",
          "is_correct": "boolean"
        }
      ]
    }
  ]
}
```

**Response**:
```json
{
  "quiz_id": "string",
  "short_id": "string",
  "message": "string"
}
```

### Get Quiz by ID
`GET /v1/quiz/{quiz_id}`

**Parameters**:
    quiz_id (string, required) - ID of the quiz to retrieve

**Response**:
```json
{
  "shortID": "string",
  "name": "string",
  "author": "string",
  "image_id": "string (optional)",
  "description": "string (optional)",
  "question": [
    {
      "question_text": "string",
      "image_id": "string (optional)",
      "answer": [
        {
          "answer_text": "string",
          "is_correct": "boolean"
        }
      ]
    }
  ]
}
```

### Get Quiz by Author
`GET /v1/quiz/author/{author}`

**Parameters**:
    author (string, required) - User ID of the quiz author

**Response**:
```json
{
  "author_quizzes": [
    {
      "quizzes": [
        {
          "shortID": "geo123",
          "name": "Geography Quiz",
          "author": "5f8d0d55b54764421b7156c3",
          "image_id": "img_12345",
          "description": "Test your geography knowledge",
          "question": [
            {
              "question_text": "What is the capital of France?",
              "image_id": "img_67890",
              "answer": [
                {
                  "answer_text": "Paris",
                  "is_correct": true
                }
              ]
            }
          ]
        }
      ]
    }
  ]
}
```

### List All Quizzes
`GET /v1/quiz/orderby`


**Response**:
```json
{
  "quizzes": [
    {
      "shortID": "geo123",
      "name": "Geography Quiz",
      "author": "5f8d0d55b54764421b7156c3",
      "image_id": "img_12345",
      "description": "Test your geography knowledge",
      "question": [
        {
          "question_text": "What is the capital of France?",
          "image_id": "img_67890",
          "answer": [
            {
              "answer_text": "Paris",
              "is_correct": true
            }
          ]
        }
      ]
    },
    {
      "shortID": "math456",
      "name": "Math Quiz",
      "author": "6h9e2e77d76986643d9378e5",
      "description": "Basic math test",
      "question": [
        {
          "question_text": "2 + 2 = ?",
          "answer": [
            {
              "answer_text": "4",
              "is_correct": true
            }
          ]
        }
      ]
    }
  ]
}
```






