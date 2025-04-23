### Authentification service

This document describes the gRPC/REST API endpoints for the Auth service. All endpoints are versioned under `/v1/users`.

---

##### Table of Contents
1. [Authentication](#authentication)
2. [User Management](#user-management)
3. [Friends](#friends)
4. [Favorite Quizzes](#favorite-quizzes)

---

#### Authentication

##### `Register`
Creates a new user account.

**HTTP Method**: `POST /v1/users/register`  
**gRPC Method**: `AuthService.Register`

###### Request:
```protobuf
message RegisterRequest {
    string email = 1;
    string password = 2;
    optional string username = 3; // Optional display name
}
```

###### Response:
```protobuf
message RegisterResponse {
    string user_id = 1; // Newly created user ID
}
```

###### Example:
```bash
curl -X POST "http://localhost:8080/v1/users/register" \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "password": "secret", "username": "john"}'
```

##### `Login`
Authenticates a user and returns JWT token.

**HTTP Method**: `POST /v1/users/login`  
**gRPC Method**: `AuthService.Login`

###### Request:
```protobuf
message LoginRequest {
    string email = 1;
    string password = 2;
}
```

###### Response:
```protobuf
message LoginResponse {
    string token = 1;       // JWT for authenticated requests
    UserResponse user = 2;  // User profile data
}
```

###### Example:
```bash
curl -X POST "http://localhost:8080/v1/users/login" \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "password": "secret"}'
```

##### `Logout`
Invalidates user session token.

**HTTP Method**: `POST /v1/users/logout`  
**gRPC Method**: `AuthService.Logout`

###### Request:
```protobuf
message LogoutRequest {
    string token = 1; // JWT to invalidate
}
```

###### Response:
```protobuf
message LogoutResponse {
    bool success = 1;
}
```

###### Example:
```bash
curl -X POST "http://localhost:8080/v1/users/logout" \
  -H "Content-Type: application/json" \
  -d '{"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."}'
```

##### `ValidateToken`
Verifies JWT token validity.

**HTTP Method**: `POST /v1/users/validate`  
**gRPC Method**: `AuthService.ValidateToken`

###### Request:
```protobuf
message ValidateTokenRequest {
    string token = 1;
}
```

###### Response:
```protobuf
message ValidateTokenResponse {
    bool valid = 1;        // Token validity status
    UserResponse user = 2; // Decoded user data
}
```

---

#### User Management

##### `GetMe`
Retrieves current user profile.

**HTTP Method**: `POST /v1/users/me`  
**gRPC Method**: `AuthService.GetMe`

###### Request:
```protobuf
message GetMeRequest {
    string token = 1;
}
```

###### Response:
```protobuf
message UserResponse {
    string id = 1;
    string username = 2;
    string email = 3;
    google.protobuf.Timestamp created_at = 8; // Registration date
    google.protobuf.Timestamp updated_at = 9; // Last profile update
}
```

##### `UpdateMe`
Updates user profile information.

**HTTP Method**: `PATCH /v1/users/me`  
**gRPC Method**: `AuthService.UpdateMe`

###### Request:
```protobuf
message UpdateMeRequest {
    string token = 1;
    optional string username = 2; // New display name
    optional string email = 3;    // New email
    optional string password = 4; // New password
}
```

---

#### Friends

##### `AddFriend`
Adds another user to friends list.

**HTTP Method**: `POST /v1/users/friends/add`  
**gRPC Method**: `AuthService.AddFriend`

###### Request:
```protobuf
message AddFriendRequest {
    string token = 1;
    string friend_id = 2; // ID of user to add
}
```

###### Response:
```protobuf
message FriendResponse {
    bool success = 1;
    string message = 2; // Status description
}
```

##### `RemoveFriend`
Removes user from friends list.

**HTTP Method**: `POST /v1/users/friends/remove`  
**gRPC Method**: `AuthService.RemoveFriend`

###### Request:
```protobuf
message RemoveFriendRequest {
    string token = 1;
    string friend_id = 2; // ID of user to remove
}
```

##### `GetFriends`
Retrieves friends list.

**HTTP Method**: `POST /v1/user/friends`  
**gRPC Method**: `AuthService.GetFriends`

###### Request:
```protobuf
message GetFriendsRequest {
    string token = 1;
}
```

###### Response:
```protobuf
message FriendsListResponse {
    repeated UserResponse friends = 1; // List of friend profiles
}
```

---

#### Favorite Quizzes

##### `AddFavoriteQuiz`
Adds quiz to user favorites.

**HTTP Method**: `POST /v1/users/favorites/quizzes/add`  
**gRPC Method**: `AuthService.AddFavoriteQuiz`

###### Request:
```protobuf
message AddFavoriteQuizRequest {
    string token = 1;
    string quiz_id = 2; // ID of quiz to add
}
```

###### Response:
```protobuf
message FavoriteQuizResponse {
    bool success = 1;
    string message = 2;
}
```

##### `GetFavoriteQuizzes`
Lists user's favorite quizzes.

**HTTP Method**: `POST /v1/users/favorites/quizzes`  
**gRPC Method**: `AuthService.GetFavoriteQuizzes`

###### Response:
```protobuf
message FavoriteQuizzesResponse {
    repeated string quiz_ids = 1; // List of favorite quiz IDs
}
```

##### `RemoveFavoriteQuiz`
Removes quiz from favorites.

**HTTP Method**: `POST /v1/users/favorites/quizzes/remove`  
**gRPC Method**: `AuthService.RemoveFavoriteQuiz`

###### Request:
```protobuf
message RemoveFavoriteQuizRequest {
    string token = 1;
    string quiz_id = 2; // ID of quiz to remove
}
```
