package migrations

import (
	"io"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func TestMigrateDB(t *testing.T) {
	logger := zerolog.New(io.Discard)

	tests := []struct {
		name   string
		dbURL  string
		source string
	}{
		{name: "fails if dbURL is empty", dbURL: "", source: "file://db/migrations"},
		{name: "fails if dbURL and source are empty", dbURL: "", source: ""},
		{name: "fails if source is empty", dbURL: "postgres@example:5432", source: ""},
		{name: "fails if source dir doesn't exist", dbURL: "postgres@example:5432", source: "file://dirdoesnt/exist"},
	}

	for _, tc := range tests {
		err := MigrateDB(tc.dbURL, tc.source, &logger)
		require.NotNil(t, err)
	}

	t.Run("returns no db scheme available", func(t *testing.T) {
		db_url := "postgres@example"
		source := "file://."
		err := MigrateDB(db_url, source, &logger)
		require.ErrorContains(t, err, "no scheme")
	})

}
