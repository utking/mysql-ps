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
		filters       []Filter
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
			name:          "Injection prevented",
			filters:       []Filter{{Column: "STATE", Operator: "=", Value: "Sending query'; DROP TABLE users--"}},
			databases:     nil,
			expectedQuery: `SELECT * FROM INFORMATION_SCHEMA.PROCESSLIST WHERE COMMAND != 'Sleep' AND USER != 'system user' AND STATE = ? ORDER BY time DESC`,
			mockArgs:      []interface{}{"Sending query'; DROP TABLE users--"},
			expectErr:     false,
			expectedCount: 1,
		},
		{
			name:          "Combined filters",
			filters:       []Filter{{Column: "STATE", Operator: "=", Value: "Sending query"}},
			databases:     []interface{}{"db1", "db2"},
			expectedQuery: `SELECT * FROM INFORMATION_SCHEMA.PROCESSLIST WHERE COMMAND != 'Sleep' AND USER != 'system user' AND DB IN (?,?) AND STATE = ? ORDER BY time DESC`,
			mockArgs:      []interface{}{"db1", "db2", "Sending query"},
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
			filters:       []Filter{{Column: "STATE", Operator: "=", Value: "Sending query"}, {Column: "USER", Operator: "=", Value: "admin"}, {Column: "TIME", Operator: ">", Value: "100"}},
			databases:     []interface{}{"db1", "db2", "db3", "db4"},
			expectedQuery: `SELECT * FROM INFORMATION_SCHEMA.PROCESSLIST WHERE COMMAND != 'Sleep' AND USER != 'system user' AND DB IN (?,?,?,?) AND STATE = ? AND USER = ? AND TIME > ? ORDER BY time DESC`,
			mockArgs:      []interface{}{"db1", "db2", "db3", "db4", "Sending query", "admin", "100"},
			expectErr:     false,
			expectedCount: 12,
		},
		{
			name:          "Injection prevented - multiple filters",
			filters:       []Filter{
				{Column: "STATE", Operator: "=", Value: "Sending query'; DROP TABLE users;--"},
				{Column: "TIME", Operator: ">", Value: "1000"},
			},
			databases:     nil,
			expectedQuery: `SELECT * FROM INFORMATION_SCHEMA.PROCESSLIST WHERE COMMAND != 'Sleep' AND USER != 'system user' AND STATE = ? AND TIME > ? ORDER BY time DESC`,
			mockArgs:      []interface{}{"Sending query'; DROP TABLE users;--", "1000"},
			expectErr:     false,
			expectedCount: 2,
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
