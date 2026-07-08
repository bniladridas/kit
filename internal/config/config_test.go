package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestInitDefaults(t *testing.T) {
	_ = os.Setenv("KIT_TEST_KEY", "test-value")

	v, err := Init()
	assert.NoError(t, err)
	assert.NotNil(t, v)
}

func TestConfigPath(t *testing.T) {
	tmpDir := t.TempDir()
	oldConfigPath := ConfigPath
	ConfigPath = tmpDir
	defer func() { ConfigPath = oldConfigPath }()

	v, err := Init()
	assert.NoError(t, err)
	assert.NotNil(t, v)
}

func TestSaveAndRead(t *testing.T) {
	tmpDir := t.TempDir()
	oldConfigPath := ConfigPath
	oldConfigName := ConfigName
	oldConfigType := ConfigType
	ConfigPath = tmpDir
	ConfigName = "testconfig"
	ConfigType = "yaml"
	defer func() {
		ConfigPath = oldConfigPath
		ConfigName = oldConfigName
		ConfigType = oldConfigType
	}()

	v := viper.New()
	v.Set("key", "value")

	err := Save(v)
	assert.NoError(t, err)

	assert.FileExists(t, filepath.Join(tmpDir, "testconfig.yaml"))
}
