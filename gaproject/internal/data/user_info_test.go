package data

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"time"
)

func setupDatabase() (*sql.DB, error) {
	connStr := "user=postgres password=2005 dbname=gaproject sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to open a DB connection: ", err)
		return nil, err
	}

	return db, nil
}

// UNIT TEST
func TestUserInfoInsert(t *testing.T) {
	// Create a new mock database.
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Create a UserInfoModel instance with the mock database.
	m := UserInfoModel{DB: db}

	// Define mock user.
	user := &User{
		Name:      "Test1",
		Surname:   "Test1",
		Email:     "Test1@example.com",
		Password:  password{hash: []byte("hashedpassword")},
		Role:      "user",
		Activated: true,
	}

	// Define mock rows.
	rows := sqlmock.NewRows([]string{"id", "created_at", "version"}).
		AddRow(1, time.Now(), 1)

	// Set up expectations for the mock database query.
	mock.ExpectQuery("^INSERT INTO user_info").
		WillReturnRows(rows)

	// Call the Insert method with the mock user.
	err = m.Insert(user)

	// Check for errors.
	assert.NoError(t, err)

	// Verify that all expectations were met.
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserInfoGetByEmail(t *testing.T) {
	// Create a new mock database.
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Create a UserInfoModel instance with the mock database.
	m := UserInfoModel{DB: db}

	// Define mock email.
	email := "Test1@example.com"

	// Define mock rows.
	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "fname", "lname", "email", "password_hash", "user_role", "activated", "version"}).
		AddRow(1, time.Now(), time.Now(), "John", "Doe", "john.doe@example.com", []byte("hashedpassword"), "user", true, 1)

	// Set up expectations for the mock database query.
	mock.ExpectQuery("^SELECT id, created_at, updated_at, fname, lname, email, password_hash, user_role, activated, version FROM user_info WHERE email").
		WillReturnRows(rows)

	// Call the GetByEmail method with the mock email.
	user, err := m.GetByEmail(email)

	// Check for errors.
	assert.NoError(t, err)
	assert.NotNil(t, user)

	// Verify that all expectations were met.
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserInfoUpdate(t *testing.T) {
	// Create a new mock database.
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Create a UserInfoModel instance with the mock database.
	m := UserInfoModel{DB: db}

	// Define mock user.
	user := &User{
		ID:        1,
		Name:      "Test",
		Surname:   "User",
		Email:     "test@example.com",
		Password:  password{hash: []byte("hashedpassword")},
		Activated: true,
		Version:   1,
	}

	// Define mock rows for UPDATE query.
	rows := sqlmock.NewRows([]string{"version"}).AddRow(2)

	// Set up expectations for the mock database query.
	mock.ExpectQuery("^UPDATE user_info").
		WithArgs(user.Name, user.Surname, user.Email, user.Password.hash, user.Activated, time.Now(), user.ID, user.Version).
		WillReturnRows(rows)

	// Call the Update method with the mock user.
	err = m.Update(user)

	// Check for errors.
	assert.NoError(t, err)

	// Verify that all expectations were met.
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserInfoDelete(t *testing.T) {
	// Create a new mock database.
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Create a UserInfoModel instance with the mock database.
	m := UserInfoModel{DB: db}

	// Define mock user ID.
	userID := int64(1)

	// Set up expectations for the mock database query.
	mock.ExpectExec("^DELETE FROM user_info").
		WithArgs(userID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Call the Delete method with the mock user ID.
	err = m.Delete(userID)

	// Check for errors.
	assert.NoError(t, err)

	// Verify that all expectations were met.
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// INTEGRATION TEST
func TestUserInfoUpdateIntegration(t *testing.T) {
	db, err := setupDatabase()
	if err != nil {
		t.Fatalf("Failed to set up database: %v", err)
	}
	defer db.Close()

	m := UserInfoModel{DB: db}

	user := &User{
		ID:        1,
		Name:      "test1",
		Surname:   "test1",
		Email:     "Zhanassetkazy@example.com",
		Password:  password{hash: []byte("updatedhashedpassword")},
		Activated: true,
		Version:   16,
	}

	err = m.Update(user)

	assert.NoError(t, err, "Failed to update user")

	updatedUser, err := m.Get(user.ID)
	assert.NoError(t, err, "Failed to get updated user")

	assert.Equal(t, user.Name, updatedUser.Name)
	assert.Equal(t, user.Surname, updatedUser.Surname)
	assert.Equal(t, user.Email, updatedUser.Email)
	assert.Equal(t, user.Password.hash, updatedUser.Password.hash)
	assert.Equal(t, user.Activated, updatedUser.Activated)
}

func TestUserInfoDeleteIntegration(t *testing.T) {
	db, err := setupDatabase()
	if err != nil {
		t.Fatalf("Failed to set up database: %v", err)
	}
	defer db.Close()

	m := UserInfoModel{DB: db}

	err = m.Delete(int64(3))

	assert.NoError(t, err, "Failed to delete user")

	deletedUser, err := m.Get(3)

	assert.Error(t, err, "Expected an error as the user should be deleted")
	assert.Equal(t, ErrRecordNotFound, err)
	assert.Nil(t, deletedUser, "Deleted user should be nil")
}
