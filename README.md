# sqlpaginator

The `sqlpaginator` package provides a set of utilities for paginating SQL queries. It simplifies the process of building paginated queries with sorting and filtering options, by dynamically constructing SQL queries with pagination support.

## Features

- **Base Query Extraction**: Extracts the base query after the `FROM` clause.
- **Count Query**: Appends a count query to calculate the total number of records.
- **Paginator Construction**: Constructs the SQL paginator query based on the provided sort column, sort direction, and primary key.
- **Query Normalization**: Normalizes queries by trimming excess spaces and converting to lowercase.
- **Dynamic Query Generation**: Generates a fully paginated SQL query with appropriate `ORDER BY`, `LIMIT`, and `WHERE` clauses.
  
## Installation

To use `sqlpaginator`, you need to import it into your Go project:

```go
import "path/to/sqlpaginator"
```

## Usage

Below is an example of how to use `sqlpaginator` to create a paginated SQL query.

### Example

```go
package main

import (
	"fmt"
	"log"

	"path/to/sqlpaginator"
)

func main() {
	// Initialize Sqlpaginator
	paginator := sqlpaginator.NewSqlpaginator()

	// Define the parameters for pagination
	params := sqlpaginator.PaginatedQueryParams{
		QueryBase:    "SELECT * FROM users",
		SortColumn:   "created_at",
		SortFunction: "asc",
		SortValue:    "2024-01-01",
		PrimaryKey:   "user_id",
	}

	// Generate the paginated query
	query, err := paginator.NewPaginatedQuery(params)
	if err != nil {
		log.Fatal(err)
	}

	// Print the generated query
	fmt.Println(query)
}
```

In this example:
- `QueryBase`: The base SQL query without pagination (`"SELECT * FROM users"`).
- `SortColumn`: The column used for sorting (`"created_at"`).
- `SortFunction`: The direction of sorting (`"asc"` or `"desc"`).
- `SortValue`: The value used for sorting comparisons.
- `PrimaryKey`: The primary key of the table (`"user_id"`).

### Methods

#### `NewSqlpaginator() SqlpaginatorInterface`

Creates a new instance of `Sqlpaginator`.

#### `GetFrom(queryBase string) string`

Extracts the part of the query after the `FROM` clause.

#### `AppendCount(queryBase string) string`

Appends a `CROSS JOIN` clause to the query, adding a count query.

#### `ConstructPaginator(sortColumn, sortFunction, sortValue, primaryKey string) string`

Constructs the paginator query, adding the `ORDER BY` and `LIMIT` clauses based on the provided parameters.

#### `NormalizeQuery(query string) string`

Normalizes a query by replacing multiple spaces with a single space and trimming leading/trailing spaces.

#### `NewPaginatedQuery(params PaginatedQueryParams) (string, error)`

Generates a full paginated query with `WITH COUNT` and pagination support.

