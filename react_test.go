package gollum_test

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"github.com/PullRequestInc/go-gpt3"
	"github.com/stillmatic/gollum"
	"github.com/stillmatic/gollum/tools"
	"github.com/stretchr/testify/assert"
	_ "modernc.org/sqlite"
)

const (
	testConvoName = "test"
)

func setupAgent(t *testing.T) *gollum.ReactAgent {
	t.Helper()
	reg := tools.NewToolRegistry()
	apiKey := os.Getenv("OPENAI_KEY")
	if apiKey == "" {
		t.Fatal("OPENAI_KEY is not set")
	}
	aiClient := gpt3.NewClient(apiKey)
	r := gollum.NewReactAgent(aiClient, reg)
	return r
}

func setupSQL(t *testing.T) tools.Tool {
	t.Helper()
	db, err := sql.Open("sqlite", ":memory:")
	assert.NoError(t, err)
	migr, err := os.ReadFile("testdata/migration.sql")
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

	sqliteToolRunner := tools.NewSQLToolRunner(db)
	resp, err := sqliteToolRunner.Run("SELECT sql FROM sqlite_master WHERE type='table' and name NOT LIKE 'sqlite_%' ORDER BY name")
	assert.NoError(t, err)

	SQLTool := tools.Tool{
		Name:        "sql",
		Description: "Run SQL queries. Available schemas: " + resp,
		Run:         sqliteToolRunner.Run,
	}

	return SQLTool
}

func TestWikipediaEndToEnd(t *testing.T) {
	r := setupAgent(t)
	r.NewConversation(testConvoName)
	ctx := context.Background()
	err := r.Speak(ctx, testConvoName, "Question: What does England share borders with?")
	assert.NoError(t, err)
}

func TestSQL(t *testing.T) {
	r := setupAgent(t)
	sqlTool := setupSQL(t)
	r.Registry.Register(sqlTool)
	r.NewConversation(testConvoName)
	ctx := context.Background()
	err := r.Speak(ctx, testConvoName, "Question: How many customers do we have?")
	assert.NoError(t, err)
}
