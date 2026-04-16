package testhelpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDatabaseNameFromFileUsesInternalPackagePath(t *testing.T) {
	t.Parallel()

	name := databaseNameFromFile("C:/Code/do-your-dailies/server/internal/domain/chores/test_helpers_test.go")

	assert.Equal(t, "dailies_test_internal_domain_chores", name)
}

func TestDatabaseNameFromFileNormalizesWindowsPath(t *testing.T) {
	t.Parallel()

	name := databaseNameFromFile("C:\\Code\\do-your-dailies\\server\\internal\\api\\json\\json_test.go")

	assert.Equal(t, "dailies_test_internal_api_json", name)
}

func TestDatabaseNameFromFileFallsBackWhenInternalMissing(t *testing.T) {
	t.Parallel()

	name := databaseNameFromFile("C:/Code/do-your-dailies/server/main_test.go")

	assert.Equal(t, "dailies_test", name)
}

func TestDatabaseNameFromFileLimitsLength(t *testing.T) {
	t.Parallel()

	name := databaseNameFromFile("C:/Code/do-your-dailies/server/internal/really/long/package/name/that/keeps/going/and/going/and/going/test_file_test.go")

	assert.LessOrEqual(t, len(name), 63)
}
