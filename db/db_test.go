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
			name:          "Multiple databases",
			filters:       nil,
			databases:     []interface{}{"db1", "db2", "db3"},
			expectedQuery: `SELECT * FROM INFORMATION_SCHEMA.PROCESSLIST WHERE COMMAND != 'Sleep' AND USER != 'system user' AND DB IN (?,?,?) ORDER BY time DESC`,
			mockArgs:      []interface{}{"db1", "db2", "db3"},
			expectErr:     false,
			expectedCount: 3,
		},
		{
			name:          "Potential injection (not blocked yet)",
			filters:       []string{"state='Sending query'; DROP TABLE users--"},
			databases:     nil,
			expectedQuery: `SELECT * FROM INFORMATION_SCHEMA.PROCESSLIST WHERE COMMAND != 'Sleep' AND USER != 'system user' AND state='Sending query'; DROP TABLE users-- ORDER BY time DESC`,
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
			name:          "Many databases",
			filters:       nil,
			databases:     []interface{}{"db1", "db2", "db3", "db4", "db5"},
			expectedQuery: `SELECT * FROM INFORMATION_SCHEMA.PROCESSLIST WHERE COMMAND != 'Sleep' AND USER != 'system user' AND DB IN (?,?,?,?,?) ORDER BY time DESC`,
			mockArgs:      []interface{}{"db1", "db2", "db3", "db4", "db5"},
			expectErr:     false,
			expectedCount: 5,
		},
		{
			name:          "Complex combined filters",
			filters:       []string{"state='Sending query'", "user_id=100", "priority='high'"},
			databases:     []interface{}{"db1", "db2", "db3", "db4"},
			expectedQuery: `SELECT * FROM INFORMATION_SCHEMA.PROCESSLIST WHERE COMMAND != 'Sleep' AND USER != 'system user' AND DB IN (?,?,?,?) AND state='Sending query' AND user_id=100 AND priority='high' ORDER BY time DESC`,
			mockArgs:      []interface{}{"db1", "db2", "db3", "db4"},
			expectErr:     false,
			expectedCount: 12,
		},
		{
			name:          "Injection Baseline - multiple statements",
			filters:       []string{"state='Sending query'; DROP TABLE users;--", "time > 1000"},
			databases:     nil,
			expectedQuery: `SELECT * FROM INFORMATION_SCHEMA.PROCESSLIST WHERE COMMAND != 'Sleep' AND USER != 'system user' AND state='Sending query'; DROP TABLE users;-- AND time > 1000 ORDER BY time DESC`,
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
