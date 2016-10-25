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
		`192.168.1.159:2181`,
		`/yidonglianjie-alimama`,
		version,
	)

	zkConfigGroup, err := config.NewZookeeperConfigGroupWithCache(configProfile, `alimama`, `e:\properties\alimama.properties`)
	assert.NoError(t, err, `NewZookeeperConfigGroupWithCache failed:%v`, configProfile)

	zkConfigGroup.ForEach(func(key, value string) {
		fmt.Println(key + `:` + value)
	})
}

func TestNewFileConfigGroup(t *testing.T) {
	version := `0.0.1`
	configProfile := config.NewZkConfigProfile(
		`192.168.1.159:2181`,
		`/yidonglianjie-alimama`,
		version,
	)

	zkConfigGroup, err := config.NewZookeeperConfigGroupWithCache(configProfile, `alimama`, `e:\\properties\alimama.properties`)
	assert.NoError(t, err, `NewZookeeperConfigGroupWithCache failed:%v`, configProfile)

	fileConfigGroup, err := config.NewFileConfigGroup(
		zkConfigGroup,
		config.NewFileConfigProfileWithVersion(`UTF8`, config.ContentType_properties, version),
		`e:\\properties\alimama.properties`,
	)
	assert.NoError(t, err, `NewFileConfigGroup failed:%v`, configProfile)

	fileConfigGroup.ForEach(func(key, value string) {
		fmt.Println(key + `:` + value)
	})
}
