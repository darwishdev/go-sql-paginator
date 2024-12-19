package main

import (
	"fmt"

	sqlpaginator "github.com/darwishdev/go-sql-paginator"
)

func main() {
	paginator := sqlpaginator.NewSqlpaginator()
	newQ, err := paginator.NewPaginatedQuery(sqlpaginator.PaginatedQueryParams{QueryBase: "slecct * from accounts_schema.role", PrimaryKey: "role_id", SortColumn: "role_id", SortValue: ""})

	fmt.Println("vim-go", newQ, err)
}
