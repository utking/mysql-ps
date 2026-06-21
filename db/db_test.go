package db

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/utking/mysql-ps/helpers"
)

// MockQuerier is a mock implementation of the Querier interface for testing purposes.
type MockQuerier struct {
	SelectFunc func(dest interface{}, query string, args ...interface{}) error
}

func (m *MockQuerier) Select(dest interface{}, query string, args ...interface{}) error {
	return m.SelectFunc(dest, query, args...)
}

func TestDBStore_GetProcessList(t *testing.T) {
	tests := []struct {
		name          string
		filters       []string
		databases     []interface{}
		expectedQuery string
		mockArgs      []interface{}
		expectErr     bool
		expectedCount int
	}{
		{
			name:          "No filters",
			filters:       nil,
			databases:     nil,
			expectedQuery: `SELECT * FROM INFORMATION_SCHEMA.PROCESSLIST WHERE COMMAND != 'Sleep' AND USER != 'system user' ORDER BY time DESC`,
			mockArgs:      nil,
			expectErr:     false,
			expectedCount: 1,
		},
		{
			name:          "Database filter only",
			filters:       nil,
			databases:     []interface{}{"db1", "db2"},
			expectedQuery: `SELECT * FROM INFORMATION_SCHEMA.PROCESSLIST WHERE COMMAND != 'Sleep' AND USER != 'system user' AND DB IN (?,?) ORDER BY time DESC`,
			mockArgs:      []interface{}{"db1", "db2"},
			expectErr:     false,
			expectedCount: 2,
		},
		{
			name:          "Arbitrary filters only",
			filters:       []string{"state='Sending query'"},
			databases:     nil,
			expectedQuery: `SELECT * FROM INFORMATION_SCHEMA.PROCESSLIST WHERE COMMAND != 'Sleep' AND USER != 'system user' AND state='Sending query' ORDER BY time DESC`,
			mockArgs:      nil,
			expectErr:     false,
			expectedCount: 1,
		},
		{
			name:          "Combined filters",
			filters:       []string{"state='Sending query'"},
			databases:     []interface{}{"db1", "db2"},
			expectedQuery: `SELECT * FROM INFORMATION_SCHEMA.PROCESSLIST WHERE COMMAND != 'Sleep' AND USER != 'system user' AND DB IN (?,?) AND state='Sending query' ORDER BY time DESC`,
			mockArgs:      []interface{}{"db1", "db2"},
			expectErr:     false,
			expectedCount: 3,
		},
		{
			name:          "Error handling",
			filters:       nil,
			databases:     nil,
			expectedQuery: `SELECT * FROM INFORMATION_SCHEMA.PROCESSLIST WHERE COMMAND != 'Sleep' AND USER != 'system user' ORDER BY time DESC`,
			mockArgs:      nil,
			expectErr:     true,
		},
		{
			name:          "Empty database list",
			filters:       nil,
			databases:     []interface{}{},
			expectedQuery: `SELECT * FROM INFORMATION_SCHEMA.PROCESSLIST WHERE COMMAND != 'Sleep' AND USER != 'system user' ORDER BY time DESC`,
			mockArgs:      nil,
			expectErr:     false,
			expectedCount: 1,
		},
		{
			name:          "Nil arbitrary filters",
			filters:       nil,
			databases:     nil,
			expectedQuery: `SELECT * FROM INFORMATION_SCHEMA.PROCESSLIST WHERE COMMAND != 'Sleep' AND USER != 'system user' ORDER BY time DESC`,
			mockArgs:      nil,
			expectErr:     false,
			expectedCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var selectFunc func(dest interface{}, query string, args ...interface{}) error
			if tt.expectErr {
				selectFunc = func(dest interface{}, query string, args ...interface{}) error {
					return errors.New("mock database error")
				}
			} else {
				selectFunc = func(dest interface{}, query string, args ...interface{}) error {
					// Verify the query
					assert.Equal(t, tt.expectedQuery, query)
					// Verify the arguments
					assert.ElementsMatch(t, tt.mockArgs, args)

					items := dest.(*[]helpers.ProcessItem)
					*items = make([]helpers.ProcessItem, tt.expectedCount)
					return nil
				}
			}

			mockQuerier := &MockQuerier{SelectFunc: selectFunc}
			store := &DBStore{db: mockQuerier}

			list, err := store.GetProcessList(tt.filters, tt.databases)

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedCount, len(list))
			}
		})
	}
}
