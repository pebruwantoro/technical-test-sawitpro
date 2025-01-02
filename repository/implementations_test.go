package repository

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

// var (
// 	mockRepo *MockRepositoryInterface
// )

type testCase struct {
	name     string
	request  interface{}
	response interface{}
	mockFunc func(m sqlmock.Sqlmock)
	err      error
}

func TestCreateEstate(t *testing.T) {
	testCases := []testCase{
		{
			name: "Test Create Estate - Success",
			request: Estate{
				Id:     "1",
				Width:  10,
				Length: 10,
			},
			mockFunc: func(m sqlmock.Sqlmock) {
				m.ExpectQuery(regexp.QuoteMeta(`INSERT INTO estates (id, width, length) VALUES ($1, $2, $3) returning id;`)).
					WithArgs("1", 10, 10).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("1"))
			},
			response: Estate{
				Id: "1",
			},
			err: nil,
		},
		{
			name: "Test Create Estate - Error",
			request: Estate{
				Id:     "1",
				Width:  10,
				Length: 10,
			},
			mockFunc: func(m sqlmock.Sqlmock) {
				m.ExpectQuery(regexp.QuoteMeta(`INSERT INTO estates (id, width, length) VALUES ($1, $2, $3) returning id;`)).
					WithArgs("1", 10, 10).
					WillReturnError(fmt.Errorf("error"))
			},
			response: Estate{},
			err:      fmt.Errorf("error"),
		},
	}

	for _, tc := range testCases {
		db, mock, _ := sqlmock.New()
		defer db.Close()

		repo := &Repository{
			Db: db,
		}

		tc.mockFunc(mock)

		res, err := repo.CreateEstate(context.Background(), tc.request.(Estate))
		assert.Equal(t, res, tc.response)
		assert.Equal(t, err, tc.err)
	}
}

func TestCreateEstateTree(t *testing.T) {
	testCases := []testCase{
		{
			name: "Test Create Estate Tree - Success",
			request: EstateTree{
				Id:       "1",
				EstateId: "1",
				X:        10,
				Y:        10,
				Height:   10,
			},
			mockFunc: func(m sqlmock.Sqlmock) {
				m.ExpectQuery(regexp.QuoteMeta(`INSERT INTO trees (id, estate_id, x, y, height) VALUES ($1, $2, $3, $4, $5) returning id;`)).
					WithArgs("1", "1", 10, 10, 10).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("1"))
			},
			response: EstateTree{
				Id:       "1",
				EstateId: "1",
				X:        10,
				Y:        10,
				Height:   10,
			},
			err: nil,
		},
		{
			name: "Test Create Estate Tree - Error",
			request: EstateTree{
				Id:       "1",
				EstateId: "1",
				X:        10,
				Y:        10,
				Height:   10,
			},
			mockFunc: func(m sqlmock.Sqlmock) {
				m.ExpectQuery(regexp.QuoteMeta(`INSERT INTO trees (id, estate_id, x, y, height) VALUES ($1, $2, $3, $4, $5) returning id;`)).
					WithArgs("1", "1", 10, 10, 10).
					WillReturnError(fmt.Errorf("error"))
			},
			response: EstateTree{},
			err:      fmt.Errorf("error"),
		},
	}

	for _, tc := range testCases {
		db, mock, _ := sqlmock.New()
		defer db.Close()

		repo := &Repository{
			Db: db,
		}

		tc.mockFunc(mock)

		res, err := repo.CreateEstateTree(context.Background(), tc.request.(EstateTree))
		assert.Equal(t, res, tc.response)
		assert.Equal(t, err, tc.err)
	}
}

func TestGetStatsByEstateId(t *testing.T) {
	testCases := []testCase{
		{
			name:    "Test Get Stats By Estate Id - Success",
			request: "1",
			mockFunc: func(m sqlmock.Sqlmock) {
				m.ExpectQuery(regexp.QuoteMeta(`SELECT COALESCE(COUNT(*), 0) AS count, COALESCE(MAX(height), 0) AS max_height, COALESCE(MIN(height), 0) AS min_height,	COALESCE(PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY height), 0) AS median_height FROM trees WHERE estate_id = $1;`)).
					WithArgs("1").
					WillReturnRows(sqlmock.NewRows([]string{"count", "max_height", "min_height", "median_height"}).AddRow(2, 25, 21, 23))

			},
			response: StatsEstate{
				Count:  2,
				Min:    21,
				Max:    25,
				Median: 23,
			},
			err: nil,
		},
		{
			name:    "Test Get Stats By Estate Id - Error",
			request: "1",
			mockFunc: func(m sqlmock.Sqlmock) {
				m.ExpectQuery(regexp.QuoteMeta(`SELECT COALESCE(COUNT(*), 0) AS count, COALESCE(MAX(height), 0) AS max_height, COALESCE(MIN(height), 0) AS min_height,	COALESCE(PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY height), 0) AS median_height FROM trees WHERE estate_id = $1;`)).
					WithArgs("1").
					WillReturnError(fmt.Errorf("error"))
			},
			response: StatsEstate{},
			err:      fmt.Errorf("error"),
		},
	}

	for _, tc := range testCases {
		db, mock, _ := sqlmock.New()
		defer db.Close()

		repo := &Repository{
			Db: db,
		}

		tc.mockFunc(mock)

		res, err := repo.GetStatsByEstateId(context.Background(), tc.request.(string))
		assert.Equal(t, res, tc.response)
		assert.Equal(t, err, tc.err)
	}
}

func TestGetEstateById(t *testing.T) {
	testCases := []testCase{
		{
			name:    "Test Get Stats By Estate Id - Success",
			request: "1",
			mockFunc: func(m sqlmock.Sqlmock) {
				m.ExpectQuery(regexp.QuoteMeta(`SELECT id, width, length FROM estates WHERE id = $1;`)).
					WithArgs("1").
					WillReturnRows(sqlmock.NewRows([]string{"id", "width", "length"}).AddRow("1", 10, 10))

			},
			response: Estate{
				Id:     "1",
				Width:  10,
				Length: 10,
			},
			err: nil,
		},
		{
			name:    "Test Get Stats By Estate Id - Error",
			request: "1",
			mockFunc: func(m sqlmock.Sqlmock) {
				m.ExpectQuery(regexp.QuoteMeta(`SELECT id, width, length FROM estates WHERE id = $1;`)).
					WithArgs("1").
					WillReturnError(fmt.Errorf("error"))
			},
			response: Estate{},
			err:      fmt.Errorf("error"),
		},
	}

	for _, tc := range testCases {
		db, mock, _ := sqlmock.New()
		defer db.Close()

		repo := &Repository{
			Db: db,
		}

		tc.mockFunc(mock)

		res, err := repo.GetEstateById(context.Background(), tc.request.(string))
		assert.Equal(t, res, tc.response)
		assert.Equal(t, err, tc.err)
	}
}

func TestGetTreesByEstateId(t *testing.T) {
	testCases := []testCase{
		{
			name:    "Test Get Stats By Estate Id - Success",
			request: "1",
			mockFunc: func(m sqlmock.Sqlmock) {
				m.ExpectQuery(regexp.QuoteMeta(`SELECT id, estate_id, x, y, height FROM trees WHERE estate_id = $1;`)).
					WithArgs("1").
					WillReturnRows(sqlmock.NewRows([]string{"id", "estate_id", "x", "y", "height"}).
						AddRow("1", "1", 10, 10, 10).
						AddRow("2", "1", 11, 11, 10))

			},
			response: []EstateTree{
				{
					Id:       "1",
					EstateId: "1",
					X:        10,
					Y:        10,
					Height:   10,
				},
				{
					Id:       "2",
					EstateId: "1",
					X:        11,
					Y:        11,
					Height:   10,
				},
			},
			err: nil,
		},
	}

	for _, tc := range testCases {
		db, mock, _ := sqlmock.New()
		defer db.Close()

		repo := &Repository{
			Db: db,
		}

		tc.mockFunc(mock)

		res, err := repo.GetTreesByEstateId(context.Background(), tc.request.(string))
		assert.Equal(t, res, tc.response)
		assert.Equal(t, err, tc.err)
	}
}
