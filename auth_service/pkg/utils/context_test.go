package utils

import (
	"context"
	"eazy-quizy-auth/internal/entity"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithUserAndGetUser(t *testing.T) {
	testUser := &entity.User{ID: 1, Email: "user@mail.ru", Username: "Test User", Password: "password"}

	ctx := context.Background()
	ctxWithUser := WithUser(ctx, testUser)

	user, ok := GetUser(ctxWithUser)
	assert.True(t, ok, "Должен вернуть true, так как пользователь есть в ctx")
	assert.Equal(t, testUser, user, "Должен вернуть того же пользователя, который был добавлен")
}

func TestGetUserNoUserInContext(t *testing.T) {
	ctx := context.Background()

	user, ok := GetUser(ctx)

	assert.False(t, ok, "Должен вернуть false, так как пользователь не в ctx")
	assert.Nil(t, user, "Должен вернуть nil, так как пользователь не в ctx")
}

func TestGetUserWrongTypeInContext(t *testing.T) {
	ctx := context.WithValue(context.Background(), UserKey, "not a user")

	user, ok := GetUser(ctx)

	assert.False(t, ok, "Должен вернуть false, так как значение не является структурой User")
	assert.Nil(t, user, "Должен вернуть nil, так как значение не является структурой User")
}

func TestWithUserNilUser(t *testing.T) {
	ctx := context.Background()
	ctxWithNilUser := WithUser(ctx, nil)

	user, ok := GetUser(ctxWithNilUser)

	assert.True(t, ok, "Должен вернуть true, так как значение (nil) было установлено")
	assert.Nil(t, user, "Должен вернуть nil, так как nil был установлен как значение")
}

func TestWithUserOverwriteExistingUser(t *testing.T) {
	user1 := &entity.User{ID: 1, Email: "user1@mail.ru", Username: "Test User1", Password: "password123"}
	user2 := &entity.User{ID: 2, Email: "user2@mail.ru", Username: "Test User2", Password: "password321"}

	ctx := WithUser(context.Background(), user1)

	ctx = WithUser(ctx, user2)

	user, ok := GetUser(ctx)
	assert.True(t, ok)
	assert.Equal(t, user2, user, "Должен вернуть последнего добавленного пользователя")
}

func TestGetUserWithDifferentKey(t *testing.T) {
	testUser := &entity.User{ID: 1, Email: "user@mail.ru", Username: "Test User", Password: "password456"}

	otherKey := struct{}{}
	ctx := context.WithValue(context.Background(), otherKey, testUser)

	user, ok := GetUser(ctx)

	assert.False(t, ok, "Должен вернуть false, так как ключ другой")
	assert.Nil(t, user, "Должен вернуть nil, так как ключ другой")
}
