package sqlpaginator

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"text/template"
)

type PaginatedQueryParams struct {
	QueryBase    string
	SortColumn   string
	SortFunction string
	SortValue    string
	PrimaryKey   string
}

type SqlpaginatorInterface interface {
	GetFrom(queryBase string) string
	AppendCount(queryBase string) string
	ConstructPaginator(sortColumn string, sortFunction string, sortValue string, primaryKey string) string
	NormalizeQuery(query string) string
	NewPaginatedQuery(PaginatedQueryParams) (string, error)
}

type Sqlpaginator struct{}

func NewSqlpaginator() SqlpaginatorInterface {
	return &Sqlpaginator{}
}

func (s *Sqlpaginator) GetFrom(queryBase string) string {
	if queryBase == "" {
		fmt.Print(errors.New("queryBase cannot be empty"))
		return ""
	}
	queryParts := strings.Split(strings.ToLower(queryBase), "from")
	if len(queryParts) != 2 {
		fmt.Print(errors.New("invalid queryBase format, missing 'from'"))
		return ""
	}
	return queryParts[1]
}

func (s *Sqlpaginator) AppendCount(queryBase string) string {
	if queryBase == "" {
		fmt.Print(errors.New("queryBase cannot be empty"))
		return ""
	}

	queryBaseLower := strings.ToLower(queryBase)
	if !strings.Contains(queryBaseLower, "from") {
		fmt.Print(errors.New("queryBase does not contain 'from'"))
		return ""
	}

	if strings.Contains(queryBaseLower, "order by") {
		fmt.Print(errors.New("queryBase cannot contain 'order by'"))
		return ""
	}

	// Extract the part of the query after "FROM"
	fromIndex := strings.Index(queryBaseLower, "from") + len("from")
	restOfQuery := queryBase[fromIndex:]

	// Match possible table definitions (including schema and alias cases)
	tableRegex := `\s*([\w.]+)(\s+as\s+\w+|\s+\w+)?`
	re := regexp.MustCompile(tableRegex)
	matches := re.FindStringSubmatch(restOfQuery)

	if len(matches) == 0 {
		fmt.Print(errors.New("unable to identify table name after 'FROM'"))
		return ""
	}

	// Build the updated query by injecting the CROSS JOIN clause
	tableDef := matches[0] // Full match of table definition
	// Remove reserved keywords if they appear within the captured tableDef
	reservedKeywords := []string{"where", "group", "join", "having", "limit"}
	for _, keyword := range reservedKeywords {
		tableDef = strings.TrimSpace(strings.ReplaceAll(strings.ToLower(tableDef), keyword, ""))
	}

	queryBase = strings.Replace(strings.ToLower(queryBase), "from", ",c.count count from", 1)
	queryWithCrossJoin := strings.Replace(queryBase, tableDef, tableDef+" cross join count c", 1)

	return queryWithCrossJoin
}

func (s Sqlpaginator) ConstructPaginator(sortColumn string, sortFunction string, sortValue string, primaryKey string) string {
	if sortColumn == "" {
		fmt.Print(errors.New("SortColumn cannot be empty"))
		return ""
	}

	if sortFunction == "" {
		sortFunction = "asc"
	}
	if primaryKey == "" {
		fmt.Print(errors.New("PrimaryKey cannot be empty"))
		return ""
	}

	operator := ">"
	if strings.ToLower(sortFunction) == "desc" {
		operator = "<"
	}

	fmt.Println("sort value pased", sortValue)
	if sortValue == "" {
		fmt.Println("sort value not pased")
		paginator := fmt.Sprintf(" order by %s %s , %s %s limit %d", sortColumn, sortFunction, primaryKey, sortFunction, 2)
		return paginator
	}
	paginator := fmt.Sprintf("%s %s %s order by %s %s , %s %s limit %d", sortColumn, operator, sortValue, sortColumn, sortFunction, primaryKey, sortFunction, 2)
	return paginator
}

func (s Sqlpaginator) WhereOperator(querBase string, sortValue string) string {
	if sortValue == "" {
		return ""
	}
	if strings.Contains(strings.ToLower(querBase), "where") {
		return " and "
	}
	return "where"

}
func (s Sqlpaginator) NormalizeQuery(query string) string {
	// Replace multiple spaces with a single space
	re := regexp.MustCompile(`\s+`)
	query = re.ReplaceAllString(query, " ")
	// Trim leading and trailing spaces
	return strings.TrimSpace(strings.ToLower(query))
}
func (s Sqlpaginator) NewPaginatedQuery(params PaginatedQueryParams) (string, error) {
	if params.QueryBase == "" {
		return "", errors.New("QueryBase cannot be empty")
	}
	if params.SortColumn == "" || params.PrimaryKey == "" {
		return "", errors.New("all parameters must be provided")
	}
	if params.SortFunction == "" {
		params.SortFunction = "asc"
	}
	params.QueryBase = s.NormalizeQuery(params.QueryBase)
	const templateContent = `with count as (select count(*) from {{ GetFrom .QueryBase }}) {{ AppendCount .QueryBase }} {{WhereOperator .QueryBase .SortValue }} {{ ConstructPaginator .SortColumn .SortFunction .SortValue .PrimaryKey }}`

	funcMap := template.FuncMap{
		"WhereOperator":      s.WhereOperator,
		"GetFrom":            s.GetFrom,
		"AppendCount":        s.AppendCount,
		"ConstructPaginator": s.ConstructPaginator,
	}
	// Parse the SQL template.
	tmpl, err := template.New("sql").Funcs(funcMap).Parse(templateContent)
	if err != nil {
		return "", err
	}
	// Use a buffer to capture the generated SQL output.
	var sqlBuffer bytes.Buffer
	err = tmpl.Execute(&sqlBuffer, params)
	if err != nil {
		return "", err
	}
	return s.NormalizeQuery(sqlBuffer.String()), nil

}
