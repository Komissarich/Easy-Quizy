 AuthService API Documentation

## Overview
AuthService provides API for authentication, user management, friends, and favorite quizzes.

Base URL: `/v1/`

All requests except registration and login require an authentication token.

## Table of Contents
1. [Authentication](#authentication)
   - [Register](#register)
   - [Login](#login)
   - [Logout](#logout)
   - [Validate Token](#validate-token)
2. [Users](#users)
   - [Get Current User](#get-current-user)
   - [Get User by ID](#get-user-by-id)
   - [Update Current User](#update-current-user)
3. [Friends](#friends)
   - [Add Friend](#add-friend)
   - [Remove Friend](#remove-friend)
   - [Get Friends List](#get-friends-list)
4. [Favorite Quizzes](#favorite-quizzes)
   - [Add Favorite Quiz](#add-favorite-quiz)
   - [Remove Favorite Quiz](#remove-favorite-quiz)
   - [Get Favorite Quizzes](#get-favorite-quizzes)
5. [Data Models](#data-models)

## Authentication

### Register
Registers a new user.

**Endpoint**: `POST /v1/users/register`

**Request**:
```json
{
  "email": "string",
  "password": "string",
  "username": "string"
}
```
**Response**:
```json
{
  "user_id": "string"
}
```

### Login
Authenticates a user and returns a token.

**Endpoint**: `POST /v1/users/login`

**Request**:
```json
{
  "email": "string",
  "password": "string"
}
```
**Response**:
```json
{
  "token": "string",
  "user": {
    "id": "string",
    "username": "string",
    "email": "string",
    "created_at": "timestamp",
    "updated_at": "timestamp"
  }
}
```
### Logout
Ends the user session.

**Endpoint**: `POST /v1/users/logout`

**Request**:
```json
{
  "token": "string"
}
```
**Response**:
```json
{
  "success": "boolean"
}
```

### Validate-token
Validates the authentication token.

**Endpoint**: `POST /v1/users/validate`

**Request**:
```json
{
  "token": "string"
}
```
**Response**:
```json
{
  "valid": "boolean",
  "user": {
    "id": "string",
    "username": "string",
    "email": "string",
    "created_at": "timestamp",
    "updated_at": "timestamp"
  }
}
```

## Users

### Get Current User

Returns information about the authenticated user.

**Endpoint**: `GET /v1/users/me`

**Headers**: Authorization: Bearer <token>

**Response**:
```json
{
  "id": "string",
  "username": "string",
  "email": "string",
  "created_at": "timestamp",
  "updated_at": "timestamp"
}
```


### Get User By ID

Returns information about a specific user.

**Endpoint**: `GET /v1/users/{user_id}`

**Headers**: Authorization: Bearer <token>

**Response**:
```json
{
  "id": "string",
  "username": "string",
  "email": "string",
  "created_at": "timestamp",
  "updated_at": "timestamp"
}
```

### Update Current User

Updates the current user's information.

**Endpoint**: `PATCH /v1/users/me`

**Headers**: Authorization: Bearer <token>

**Request**:
```json
{
  "username": "string",
  "email": "string",
  "password": "string"
}
```

**Response**:
```json
{
  "id": "string",
  "username": "string",
  "email": "string",
  "created_at": "timestamp",
  "updated_at": "timestamp"
}
```

## Friends


### Add Friend

Adds a user to friends.

**Endpoint**: `POST /v1/users/friends/add`

**Headers**: Authorization: Bearer <token>

**Request**:
```json
{
  "friend_id": "string"
}
```

**Response**:
```json
{
  "success": "boolean",
  "message": "string"
}
```

### Remove Friend

Removes a user from friends.


**Endpoint**: `POST /v1/users/friends/remove`

**Headers**: Authorization: Bearer <token>

**Request**:
```json
{
  "friend_id": "string"
}
```

**Response**:
```json
{
  "success": "boolean",
  "message": "string"
}
```

### Get Friends List

Returns the current user's friends list.


**Endpoint**: `POST /v1/user/friends`

**Headers**: Authorization: Bearer <token>


**Response**:
```json
{
  "friends": [
    {
      "id": "string",
      "username": "string",
      "email": "string",
      "created_at": "timestamp",
      "updated_at": "timestamp"
    }
  ]
}
```
## Favorite Quizzes


### Add Favorite Quiz

Adds a quiz to favorites.


**Endpoint**: `POST /v1/users/favorites/quizzes/add`

**Headers**: Authorization: Bearer <token>

**Request**:
```json
{
  "quiz_id": "string"
}
```

**Response**:
```json
{
  "success": "boolean",
  "message": "string"
}
```


### Remove Favorite Quiz

Removes a quiz from favorites.


**Endpoint**: `POST /v1/users/favorites/quizzes/remove`

**Headers**: Authorization: Bearer <token>

**Request**:
```json
{
  "quiz_id": "string"
}

**Response**:
```json
{
  "success": "boolean",
  "message": "string"
}
```


### Get Favorite Quizzes

Returns a list of favorite quiz IDs.


**Endpoint**: `GET /v1/users/favorites/quizzes`

**Headers**: Authorization: Bearer <token>


**Response**:
```json
{
  "quiz_ids": ["string"]
}
```


