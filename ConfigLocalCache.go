package config

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type ConfigLocalCache struct {
	localCachePath string
}

const (
	SUFFIX = ".cache"
)

func newConfigLocalCache(localCachePath string) *ConfigLocalCache {
	if len(localCachePath) == 0 {
		localCachePath = os.TempDir()
	}

	if strings.LastIndex(localCachePath, string(filepath.Separator)) != len(localCachePath)-1 {
		localCachePath = localCachePath + string(filepath.Separator)
	}
	return &ConfigLocalCache{
		localCachePath: localCachePath,
	}
}

/**
 * 缓存配置到本地
 *
 * @param configNode
 * @param node
 */
func (this *ConfigLocalCache) saveLocalCache(configNode *ZookeeperConfigGroup, node string) (string, error) {
	properties := &bytes.Buffer{}
	for key, value := range configNode.exportProperties() {
		_, err := properties.WriteString(key)
		if err != nil {
			return ``, err
		}
		properties.WriteString(`=`)
		_, err = properties.WriteString(value)
		if err != nil {
			return ``, err
		}
		_, err = properties.WriteString("\n")
		if err != nil {
			return ``, err
		}
	}

	localPath := this.genLocalPath(node)
	index := strings.LastIndex(localPath, string(os.PathSeparator))
	os.MkdirAll(localPath[:index], 0644)

	err := ioutil.WriteFile(localPath, properties.Bytes(), 0644)
	if err != nil {
		return ``, err
	}

	return localPath, nil
}

func (this *ConfigLocalCache) genLocalPath(node string) string {
	return this.localCachePath + node + SUFFIX
}
