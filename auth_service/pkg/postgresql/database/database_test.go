package database

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew_ConnectionError(t *testing.T) {
	ctx := context.Background()

	testCases := []struct {
		name     string
		username string
		password string
		host     string
		port     int
		dbname   string
	}{
		{
			name:     "invalid credentials",
			username: "invalid",
			password: "invalid",
			host:     "localhost",
			port:     5432,
			dbname:   "test_db",
		},
		{
			name:     "invalid host",
			username: "postgres",
			password: "postgres",
			host:     "invalid_host",
			port:     5432,
			dbname:   "test_db",
		},
		{
			name:     "invalid port",
			username: "postgres",
			password: "postgres",
			host:     "localhost",
			port:     9999,
			dbname:   "test_db",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, err := New(ctx, tc.username, tc.password, tc.host, tc.port, tc.dbname)
			require.Error(t, err)
			assert.Nil(t, db)
			assert.Contains(t, err.Error(), "failed to connect to database")
		})
	}
}

func TestClose_Success(t *testing.T) {
	ctx := context.Background()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	mock.ExpectClose().WillReturnError(nil)

	db := &DB{DB: mockDB}
	err = db.Close(ctx)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestClose_Error(t *testing.T) {
	ctx := context.Background()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	expectedErr := errors.New("close error")
	mock.ExpectClose().WillReturnError(expectedErr)

	db := &DB{DB: mockDB}
	err = db.Close(ctx)
	assert.EqualError(t, err, expectedErr.Error())

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestClose_Timeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	db := &DB{DB: mockDB}
	err = db.Close(ctx)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "all expectations were already fulfilled, call to database Close was not expected")

	assert.NoError(t, mock.ExpectationsWereMet())
}
