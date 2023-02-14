package migrations

import (
	"io"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func TestMigrateDB(t *testing.T) {
	l := zerolog.New(io.Discard)

	t.Run("returns error empty connection string", func(t *testing.T) {
		db_url := ""
		source := "file://db/migrations."
		err := MigrateDB(db_url, source, &l)
		require.NotNil(t, err)
	})
	t.Run("returns no db scheme available", func(t *testing.T) {
		db_url := "postgres@example"
		source := "file://."
		err := MigrateDB(db_url, source, &l)
		require.ErrorContains(t, err, "no scheme")
	})

}
