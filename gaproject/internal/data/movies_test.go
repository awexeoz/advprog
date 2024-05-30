package data

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// UNIT TEST
func TestMovieInsert(t *testing.T) {
	// Create a new mock database.
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Create a MovieModel instance with the mock database.
	m := MovieModel{DB: db}

	// Define a mock movie.
	mockMovie := &Movie{
		Title:   "Test Movie 1",
		Year:    2024,
		Runtime: 120,
		Genres:  []string{"Action", "Adventure"},
	}

	// Set up expectations for the mock database query.
	mock.ExpectQuery("^INSERT INTO movies (.+) RETURNING id, created_at, version").
		WithArgs(mockMovie.Title, mockMovie.Year, mockMovie.Runtime, pq.Array(mockMovie.Genres)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "version"}).
			AddRow(1, time.Now(), 1))

	// Call the Insert method with the mock movie.
	err = m.Insert(mockMovie)

	// Check for errors.
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Verify that all expectations were met.
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestMovieGet(t *testing.T) {
	// Create a new mock database.
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Create a MovieModel instance with the mock database.
	m := MovieModel{DB: db}

	// Define a mock movie ID.
	mockID := int64(1)

	// Define a mock movie.
	mockMovie := &Movie{
		ID:        1,
		CreatedAt: time.Now(),
		Title:     "Mock Movie",
		Year:      2022,
		Runtime:   120,
		Genres:    []string{"Action", "Adventure"},
		Version:   1,
	}

	// Set up expectations for the mock database query.
	mock.ExpectQuery("^SELECT (.+) FROM movies WHERE id = \\$1").
		WithArgs(mockID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "title", "year", "runtime", "genres", "version"}).
			AddRow(mockMovie.ID, mockMovie.CreatedAt, mockMovie.Title, mockMovie.Year, mockMovie.Runtime, pq.Array(mockMovie.Genres), mockMovie.Version))

	// Call the Get method with the mock movie ID.
	movie, err := m.Get(mockID)

	// Check for errors and compare the retrieved movie with the mock movie.
	assert.NoError(t, err)
	assert.Equal(t, mockMovie, movie)

	// Verify that all expectations were met.
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestMovieUpdate(t *testing.T) {
	// Create a new mock database.
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Create a MovieModel instance with the mock database.
	m := MovieModel{DB: db}

	// Define mock movie.
	movie := &Movie{
		ID:      1,
		Title:   "Test Movie",
		Year:    2023,
		Runtime: 120,
		Genres:  []string{"Action", "Adventure"},
		Version: 1,
	}

	// Define mock rows for UPDATE query.
	rows := sqlmock.NewRows([]string{"version"}).AddRow(2)

	// Set up expectations for the mock database query.
	mock.ExpectQuery("^UPDATE movies").
		WithArgs(movie.Title, movie.Year, movie.Runtime, pq.Array(movie.Genres), movie.ID, movie.Version).
		WillReturnRows(rows)

	// Call the Update method with the mock movie.
	err = m.Update(movie)

	// Check for errors.
	assert.NoError(t, err)

	// Verify that all expectations were met.
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestMovieDelete(t *testing.T) {
	// Create a new mock database.
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Create a MovieModel instance with the mock database.
	m := MovieModel{DB: db}

	// Define mock movie ID.
	movieID := int64(1)

	// Set up expectations for the mock database query.
	mock.ExpectExec("^DELETE FROM movies").
		WithArgs(movieID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Call the Delete method with the mock movie ID.
	err = m.Delete(movieID)

	// Check for errors.
	assert.NoError(t, err)

	// Verify that all expectations were met.
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// INTEGRATION TEST
func TestMovieInsertIntegration(t *testing.T) {
	db, err := setupDatabase()
	if err != nil {
		t.Fatalf("Failed to set up database: %v", err)
	}
	defer db.Close()

	m := MovieModel{DB: db}

	movie := &Movie{
		Title:   "Test Movie 1",
		Year:    2021,
		Runtime: 120,
		Genres:  []string{"Action", "Adventure"},
	}

	err = m.Insert(movie)
	assert.NoError(t, err, "Failed to insert movie")
}

func TestMovieGetIntegration(t *testing.T) {
	db, err := setupDatabase()
	if err != nil {
		t.Fatalf("Failed to set up database: %v", err)
	}
	defer db.Close()

	m := MovieModel{DB: db}

	movie := &Movie{
		Title:   "Test Movie 1",
		Year:    2021,
		Runtime: 120,
		Genres:  []string{"Action", "Adventure"},
	}

	err = m.Insert(movie)
	assert.NoError(t, err, "Failed to insert movie")

	insertedMovie, err := m.Get(movie.ID)
	assert.NoError(t, err, "Failed to get movie by ID")

	assert.Equal(t, movie.ID, insertedMovie.ID)
	assert.Equal(t, movie.Title, insertedMovie.Title)
	assert.Equal(t, movie.Year, insertedMovie.Year)
	assert.Equal(t, movie.Runtime, insertedMovie.Runtime)
	assert.Equal(t, movie.Genres, insertedMovie.Genres)
}
