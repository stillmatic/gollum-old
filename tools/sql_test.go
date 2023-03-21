package tools_test

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"testing"

	_ "embed"

	"github.com/stillmatic/gollum/tools"
	"github.com/stretchr/testify/assert"
	_ "modernc.org/sqlite"
)

func setup(t *testing.T) tools.Tool {
	t.Helper()
	db, err := sql.Open("sqlite", ":memory:")
	assert.NoError(t, err)
	migr, err := os.ReadFile("../testdata/migration.sql")
	assert.NoError(t, err)
	_, err = db.Exec(string(migr))
	assert.NoError(t, err)
	// check migration
	res, err := db.Query("SELECT * FROM customers")
	assert.NoError(t, err)
	defer res.Close()
	i := 0
	for res.Next() {
		i++
	}
	assert.Equal(t, 4, i)

	sb := &strings.Builder{}
	sb.WriteString("Run SQL queries. Available schemas:\n")

	sqliteToolRunner := tools.NewSQLToolRunner(db)
	resp, err := sqliteToolRunner.Run("SELECT sql FROM sqlite_master WHERE type='table' and name NOT LIKE 'sqlite_%' ORDER BY name")
	assert.NoError(t, err)
	sb.WriteString(resp + "\n")

	tableResp, err := db.Query("SELECT name FROM sqlite_master WHERE type='table' and name NOT LIKE 'sqlite_%' ORDER BY name")
	assert.NoError(t, err)
	tables := make([]string, 0)
	for tableResp.Next() {
		var name string
		err = tableResp.Scan(&name)
		assert.NoError(t, err)
		tables = append(tables, name)
	}
	sb.WriteString("Sample data for each table:\n")
	for _, table := range tables {
		// NB: do not use this in production, vulnerable to SQL injection
		tableResp, err := sqliteToolRunner.Run(fmt.Sprintf("SELECT * FROM %s LIMIT 3'", table))
		assert.NoError(t, err)
		sb.WriteString(tableResp + "\n")
	}

	SQLTool := tools.Tool{
		Name:        "sql",
		Description: sb.String(),
		Run:         sqliteToolRunner.Run,
	}

	return SQLTool
}

func TestSQLToolQuery(t *testing.T) {
	sqlTool := setup(t)
	resp, err := sqlTool.Run("SELECT * FROM customers")
	assert.NoError(t, err)
	assert.Equal(t, `id,first_name,last_name,email,phone
-----------------------------------
1,John,Doe,john.doe@example.com,555-123-4567,
2,Jane,Smith,jane.smith@example.com,555-987-6543,
3,Michael,Johnson,michael.johnson@example.com,555-234-5678,
4,Emily,Williams,emily.williams@example.com,555-876-5432,
`, resp)

	resp, err = sqlTool.Run("SELECT last_name, COUNT(*) FROM customers INNER JOIN orders ON customers.id = orders.customer_id GROUP BY last_name")
	assert.NoError(t, err)
	assert.Equal(t, `last_name,COUNT(*)
------------------
Doe,2,
Johnson,1,
Smith,2,
Williams,1,
`, resp)

	resp, err = sqlTool.Run("DROP TABLE customers")
	assert.NoError(t, err)
	assert.Equal(t, `Error: near "DROP": syntax error`, resp)
	resp, err = sqlTool.Run("SELECT * FROM customers")
	assert.NoError(t, err)
	assert.Contains(t, resp, `4`)
}
