package config_test

import (
	"config"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewZookeeperConfigGroupWithCache(t *testing.T) {
	version := `0.0.1`
	configProfile := config.NewZkConfigProfile(
		`localhost:2181`,
		`/config`,
		version,
	)

	zkConfigGroup, err := config.NewZookeeperConfigGroupWithCache(
		configProfile,
		`application`,
		`/etc/application.cache`,
	)

	assert.NoError(t, err, `NewZookeeperConfigGroupWithCache failed:%v`, configProfile)

	zkConfigGroup.ForEach(func(key, value string) {
		fmt.Println(key + `:` + value)
	})
}

func TestNewFileConfigGroup(t *testing.T) {
	version := `0.0.1`
	configProfile := config.NewZkConfigProfile(
		`localhost:2181`,
		`/config`,
		version,
	)

	zkConfigGroup, err := config.NewZookeeperConfigGroupWithCache(
		configProfile,
		`application`,
		`/etc/application.cache`,
	)

	assert.NoError(t, err, `NewZookeeperConfigGroupWithCache failed:%v`, configProfile)

	fileConfigGroup, err := config.NewFileConfigGroup(
		zkConfigGroup,
		config.NewFileConfigProfileWithVersion(`UTF8`, config.ContentType_properties, version),
		`/etc/application.properties`,
	)
	assert.NoError(t, err, `NewFileConfigGroup failed:%v`, configProfile)

	fileConfigGroup.ForEach(func(key, value string) {
		fmt.Println(key + `:` + value)
	})
}
