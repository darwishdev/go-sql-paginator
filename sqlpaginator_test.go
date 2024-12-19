package sqlpaginator

import (
	"errors"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Helper function to normalize queries by removing extra spaces and formatting.
func normalizeQuery(query string) string {
	re := regexp.MustCompile(`\s+`)

	q := strings.ReplaceAll(strings.ToLower(query), "( ", "(")
	q = strings.ReplaceAll(q, " )", ")")
	return strings.TrimSpace(re.ReplaceAllString(q, " "))
}

func TestAppendCount(t *testing.T) {
	paginator := NewSqlpaginator()

	tests := []struct {
		name        string
		queryBase   string
		expected    string
		errorString string
	}{
		// Basic scenarios
		{"with WHERE clause", "SELECT * FROM table_name WHERE condition", "select * ,c.count count from table_name cross join count c where condition", ""},
		{"without WHERE clause", "SELECT * FROM table_name", "select * ,c.count count from table_name cross join count c", ""},
		{"with GROUP BY", "SELECT * FROM table_name GROUP BY column", "select * ,c.count count from table_name cross join count c group by column", ""},
		{"with ORDER BY", "SELECT * FROM table_name ORDER BY column ASC", "", ""},
		{"with JOIN and WHERE", "SELECT * FROM table_name t JOIN other_table o ON t.id = o.id WHERE condition", "select * ,c.count count from table_name t cross join count c join other_table o on t.id = o.id where condition", ""},

		// Combined cases
		{"with WHERE and GROUP BY", "SELECT * FROM table_name WHERE condition GROUP BY column", "select * ,c.count count from table_name cross join count c where condition group by column", ""},
		{"with JOIN, GROUP BY, and ORDER BY", "SELECT * FROM table_name t JOIN other_table o ON t.id = o.id GROUP BY column ORDER BY column DESC", "", ""},
		{"complex query", "SELECT * FROM table_name t cross join count c JOIN other_table o ON t.id = o.id WHERE condition GROUP BY column ORDER BY column DESC", "", ""},
		{"empty queryBase", "", "", "queryBase cannot be empty"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := paginator.AppendCount(test.queryBase)
			assert.Equal(t, normalizeQuery(test.expected), normalizeQuery(result))
		})
	}
}

func TestNewPaginatedQuery(t *testing.T) {
	paginator := NewSqlpaginator()

	tests := []struct {
		name     string
		params   PaginatedQueryParams
		expected string
		err      error
	}{
		// Basic scenarios
		{
			"with WHERE",
			PaginatedQueryParams{
				QueryBase:    "SELECT * FROM table_name WHERE condition",
				SortColumn:   "id",
				SortFunction: "asc",
				SortValue:    "42",
				PrimaryKey:   "id",
			},
			"with count as ( select count(*) from table_name where condition ) select * ,c.count count from table_name cross join count c where condition and id > 42 order by id asc , id asc limit 2",
			nil,
		},
		{
			"with GROUP BY",
			PaginatedQueryParams{
				QueryBase:    "SELECT * FROM table_name GROUP BY column",
				SortColumn:   "id",
				SortFunction: "desc",
				SortValue:    "",
				PrimaryKey:   "id",
			},
			"with count as ( select count(*) from table_name group by column ) select * ,c.count count from table_name cross join count c group by column order by id desc , id desc limit 2",
			nil,
		},
		{
			"with ORDER BY",
			PaginatedQueryParams{
				QueryBase:    "SELECT * FROM table_name ORDER BY column ASC",
				SortColumn:   "id",
				SortFunction: "asc",
				SortValue:    "10",
				PrimaryKey:   "id",
			},
			"",
			nil,
		},
		{
			"complex query",
			PaginatedQueryParams{
				QueryBase:    "SELECT * FROM table_name t JOIN other_table o ON t.id = o.id WHERE condition GROUP BY column ORDER BY column DESC",
				SortColumn:   "id",
				SortFunction: "desc",
				SortValue:    "100",
				PrimaryKey:   "id",
			},
			"",
			nil,
		},
		// Edge case
		{
			"missing QueryBase",
			PaginatedQueryParams{
				QueryBase:    "",
				SortColumn:   "id",
				SortFunction: "asc",
				SortValue:    "42",
				PrimaryKey:   "id",
			},
			"",
			errors.New("QueryBase cannot be empty"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := paginator.NewPaginatedQuery(test.params)
			assert.Equal(t, normalizeQuery(test.expected), normalizeQuery(result))
			if test.err != nil {
				assert.EqualError(t, err, test.err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
