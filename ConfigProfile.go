package config

import (
	"github.com/curator-go/curator"
	"github.com/emirpasic/gods/sets/hashset"
	"time"
)

const (
	/**
	 * 加载所有属性
	 */
	KeyLoadingMode_ALL = 0
	/**
	 * 包含某些属性
	 */
	KeyLoadingMode_INCLUDE = 1
	/**
	 * 排除某些属性
	 */
	KeyLoadingMode_EXCLUDE = 2
)

type ConfigProfile struct {
	Version        string
	KeyLoadingMode uint8
	KeysSpecified  *hashset.Set
}

type ZookeeperConfigProfile struct {
	ConnectStr           string
	RootNode             string
	RetryPolicy          *curator.ExponentialBackoffRetry
	ConsistencyCheck     bool
	ConsistencyCheckRate time.Duration
	*ConfigProfile
}

type FileConfigProfile struct {
	contentType  int
	fileEncoding string
	*ConfigProfile
}

func NewZkConfigProfile(connectStr, rootNode, version string) *ZookeeperConfigProfile {
	return &ZookeeperConfigProfile{
		ConnectStr:           connectStr,
		RootNode:             rootNode,
		RetryPolicy:          curator.NewExponentialBackoffRetry(time.Second, 3, 15*time.Second),
		ConsistencyCheck:     true,
		ConsistencyCheckRate: 1 * time.Second,
		ConfigProfile: &ConfigProfile{
			KeyLoadingMode: KeyLoadingMode_ALL,
			Version:        version,
			KeysSpecified:  hashset.New(),
		},
	}
}

func (this *ZookeeperConfigProfile) versionedRootNode() string {
	if len(this.Version) == 0 {
		return this.RootNode
	}

	return MakePath(this.RootNode, this.Version)
}

func NewFileConfigProfile(fileEncoding string, contentType int) *FileConfigProfile {
	fileConfigProfile := &FileConfigProfile{
		contentType:  contentType,
		fileEncoding: fileEncoding,
		ConfigProfile: &ConfigProfile{
			KeyLoadingMode: KeyLoadingMode_ALL,
			KeysSpecified:  hashset.New(),
			Version:        ``,
		},
	}

	return fileConfigProfile
}

func NewFileConfigProfileWithVersion(fileEncoding string, contentType int, version string) *FileConfigProfile {
	return &FileConfigProfile{
		contentType:  contentType,
		fileEncoding: fileEncoding,
		ConfigProfile: &ConfigProfile{
			KeyLoadingMode: KeyLoadingMode_ALL,
			KeysSpecified:  hashset.New(),
			Version:        version,
		},
	}
}
