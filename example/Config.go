package config

import (
	"entity"
	"github.com/mbaight/config-toolkit-go"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
)

var fileConfigGroup = func() *config.FileConfigGroup {
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

	if err != nil {
		panic(err)
	}

	zkConfigGroup.AddWatcher(func(propertyName, value string) {
		log.Printf("config change:%s,%s\n", propertyName, value)
	})

	fileConfigGroup, err := config.NewFileConfigGroup(
		zkConfigGroup,
		config.NewFileConfigProfileWithVersion(`UTF8`, config.ContentType_properties, version),
		`/etc/application.properties`,
	)

	if err != nil {
		log.Printf("fileConfigGroup init error:%v\n", err)
	}

	fileConfigGroup.AddWatcher(func(propertyName, value string) {
		log.Printf("config change:%s,%s\n", propertyName, value)
	})

	return fileConfigGroup
}()


//根据propertyName获取属性
func GetConfig(propertyName string) string {
	return fileConfigGroup.Get(propertyName)
}

//获取string类型的配置,如果获取失败,则返回defaultValue
func GetStringConfigWithDefault(propertyName string, defaultValue string) string {
	result := fileConfigGroup.Get(propertyName)
	if len(result) > 0 {
		return result
	}

	return defaultValue
}

//获取int类型的配置,如果获取失败,则返回defaultValue
func GetIntConfigWithDefault(propertyName string, defaultValue int) int {
	result, err := fileConfigGroup.GetInt(propertyName)
	if err != nil {
		return defaultValue
	}
	return result
}

//获取bool类型的配置,如果获取失败,则返回defaultValue
func GetBoolConfigWithDefault(propertyName string, defaultValue bool) bool {
	result, err := fileConfigGroup.GetBool(propertyName)
	if err != nil {
		return defaultValue
	}

	return result
}
