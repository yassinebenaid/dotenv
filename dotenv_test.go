package dotenv

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func tmpFile(t *testing.T, name string, content string) {
	t.Helper()
	err := os.WriteFile(name, []byte(content), 0666)
	if !assert.NoError(t, err) {
		return
	}
	t.Cleanup(func() {
		assert.NoError(t, os.Remove(name))
	})
}

func TestReadFallsBackToDotEnvFileInWorkingDirectory(t *testing.T) {
	tmpFile(t, ".env", "KEY=value")

	env, err := Read()
	if !assert.NoError(t, err) {
		return
	}

	assert.Equal(t, "value", env["KEY"])
}

func TestCanReadFile(t *testing.T) {
	tmpFile(t, "./.env-file", "KEY=value")

	env, err := Read("./.env-file")
	if !assert.NoError(t, err) {
		return
	}

	assert.Equal(t, "value", env["KEY"])
}

func TestReadReturnsErrorIfFileDoesntExist(t *testing.T) {
	_, err := Read("./.missing-env-file")
	assert.Error(t, err)
}

func TestCanReadManyFilesAndReferenceVariablesInBetweenFiles(t *testing.T) {
	tmpFile(t, "./.env-1", "KEY_1=value")
	tmpFile(t, "./.env-2", "KEY_2=$KEY_1")
	tmpFile(t, "./.env-3", "KEY_3=$KEY_2")

	env, err := Read("./.env-1", "./.env-2", "./.env-3")
	if !assert.NoError(t, err) {
		return
	}

	assert.Equal(t, "value", env["KEY_1"])
	assert.Equal(t, "value", env["KEY_2"])
	assert.Equal(t, "value", env["KEY_3"])
}

func TestReadReturnsParserErrorsIfThereAreAny(t *testing.T) {
	tmpFile(t, "./.env-file", "KEY_3 value")

	_, err := Read("./.env-file")
	if !assert.Error(t, err) {
		return
	}
	assert.Equal(t, `error parsing file [./.env-file], expected "=", found "v", line 1:7`, err.Error())
}

func TestCanLoad(t *testing.T) {
	tmpFile(t, "./.env-file", "KEY=value")

	os.Setenv("KEY", "something")
	assert.Equal(t, `something`, os.Getenv("KEY"))

	err := Load(false, "./.env-file")
	if !assert.NoError(t, err) {
		return
	}

	assert.NotEqual(t, `value`, os.Getenv("KEY"))
}

func TestLoadCanOverride(t *testing.T) {
	tmpFile(t, "./.env-file", "KEY=value")

	os.Setenv("KEY", "something")
	assert.Equal(t, `something`, os.Getenv("KEY"))

	err := Load(true, "./.env-file")
	if !assert.NoError(t, err) {
		return
	}

	assert.Equal(t, `value`, os.Getenv("KEY"))
}

func TestCanUnmarshal(t *testing.T) {
	input := []byte(`
		PORT=8080
		HOST=example.com:$PORT
	`)

	env, err := Unmarshal(input)
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, `8080`, env["PORT"])
	assert.Equal(t, `example.com:8080`, env["HOST"])
}
