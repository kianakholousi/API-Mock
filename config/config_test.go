package config_test

import (
	"flight-data-api/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitFileNotExisting(t *testing.T) {
	testParam := config.Params{FileType: "xml", FilePath: "."}

	_, err := config.Init(testParam)

	assert.NotEqual(t, nil, err)
}
