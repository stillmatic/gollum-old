package tools

import (
	"database/sql"
	"strings"
)

type SQLToolRunner struct {
	db *sql.DB
}

func QueryResultToString(resp *sql.Rows) (string, error) {
	defer resp.Close()
	sb := strings.Builder{}
	// append columns
	cols, err := resp.Columns()
	if err != nil {
		return "", err
	}
	columnNames := strings.Join(cols, ",")
	sb.WriteString(columnNames)
	sb.WriteString("\n")
	sb.WriteString(strings.Repeat("-", len(columnNames)))
	sb.WriteString("\n")

	values := make([]sql.RawBytes, len(cols))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for resp.Next() {
		err = resp.Scan(scanArgs...)
		if err != nil {
			return "", err
		}

		for _, val := range values {
			if val == nil {
				sb.WriteString("NULL")
			} else {
				sb.WriteString(string(val))
			}
			sb.WriteString(",")
		}
		sb.WriteString("\n")
	}
	if err = resp.Err(); err != nil {
		return "", err
	}
	return sb.String(), nil
}

func (r *SQLToolRunner) Run(arg string) (string, error) {
	resp, err := r.db.Query(arg)
	if err != nil {
		return "", err
	}
	return QueryResultToString(resp)
}

func NewSQLToolRunner(db *sql.DB) *SQLToolRunner {
	return &SQLToolRunner{
		db: db,
	}
}
