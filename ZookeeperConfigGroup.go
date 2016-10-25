package config

import (
	"github.com/curator-go/curator"
	"github.com/samuel/go-zookeeper/zk"
	"log"
	"time"
)

type ZookeeperConfigGroup struct {
	configProfile    *ZookeeperConfigProfile
	node             string //节点名字
	client           curator.CuratorFramework
	configLocalCache *ConfigLocalCache
	*GeneralConfigGroup
}

func NewZookeeperConfigGroup(configProfile *ZookeeperConfigProfile, node string) (*ZookeeperConfigGroup, error) {
	group := &ZookeeperConfigGroup{
		configProfile:      configProfile,
		node:               node,
		GeneralConfigGroup: NewGeneralConfigGroup(nil),
	}
	err := group.initConfigs()
	if err != nil {
		return nil, err
	}

	return group, nil
}

func NewZookeeperConfigGroupWithCache(configProfile *ZookeeperConfigProfile, node string, cachePath string) (*ZookeeperConfigGroup, error) {
	group := &ZookeeperConfigGroup{
		configProfile:      configProfile,
		node:               node,
		configLocalCache:   newConfigLocalCache(cachePath),
		GeneralConfigGroup: NewGeneralConfigGroup(nil),
	}

	err := group.initConfigs()
	if err != nil {
		return nil, err
	}

	return group, nil
}

/**
 * 初始化节点
 */
func (this *ZookeeperConfigGroup) initConfigs() error {

	builder := &curator.CuratorFrameworkBuilder{
		ConnectionTimeout: 1 * time.Second,
		SessionTimeout:    1 * time.Second,
		RetryPolicy:       this.configProfile.RetryPolicy,
	}

	this.client = builder.ConnectString(this.configProfile.ConnectStr).Build()

	// this is one method of getting event/async notifications
	err := this.client.Start()
	if err != nil {
		log.Printf("start zookeeper client error:%v", err)
		return err
	}

	this.client.CuratorListenable().AddListener(curator.NewCuratorListener(
		func(client curator.CuratorFramework, event curator.CuratorEvent) error {
			if event.Type() == curator.WATCHED {
				someChange := false

				switch event.WatchedEvent().Type {
				case zk.EventNodeChildrenChanged:
					this.loadNode()
					someChange = true
				case zk.EventNodeDataChanged:
					this.reloadKey(event.Path())
					someChange = true
				default:

				}

				if someChange {
					log.Printf("reload properties with %s", event.Path())
					if this.configLocalCache != nil {
						_, err = this.configLocalCache.saveLocalCache(this, this.node)
						if err != nil {
							log.Printf("save to local file error:%v %v", this.configLocalCache, err)
						}
					}

				}
			}

			return nil
		}))

	err = this.loadNode()
	if err != nil {
		log.Printf("load node error:%v", err)
		return err
	}

	if this.configLocalCache != nil {
		_, err = this.configLocalCache.saveLocalCache(this, this.node)
		if err != nil {
			log.Printf("save to local file error:%v %v", this.configLocalCache, err)
			return err
		}
	}

	// Consistency check
	if this.configProfile.ConsistencyCheck {
		go func() {
			time.Sleep(1 * time.Second)
			for {
				<-time.After(this.configProfile.ConsistencyCheckRate)
				this.loadNode()
			}
		}()
	}

	return nil
}

/**
 * 加载节点并监听节点变化
 */
func (this *ZookeeperConfigGroup) loadNode() error {
	nodePath := MakePath(this.configProfile.versionedRootNode(), this.node)
	childrenBuilder := this.client.GetChildren()
	children, err := childrenBuilder.Watched().ForPath(nodePath)
	if err != nil {
		return err
	}

	configs := make(map[string]string, len(children))
	for _, item := range children {
		key, value, err := this.loadKey(MakePath(nodePath, item))
		if err != nil {
			log.Printf("load property error:%s, %v", item, err)
			return err
		}

		if len(key) > 0 {
			configs[key] = value
		}
	}

	this.PutAll(configs)
	return nil
}

//重新加载某一子节点
func (this *ZookeeperConfigGroup) reloadKey(nodePath string) {
	key, value, _ := this.loadKey(nodePath)
	if len(key) > 0 {
		this.Put(key, value)
	}
}

//加载某一子节点
func (this *ZookeeperConfigGroup) loadKey(nodePath string) (string, string, error) {
	nodeName := getNodeFromPath(nodePath)

	keysSpecified := this.configProfile.KeysSpecified
	switch this.configProfile.KeyLoadingMode {
	case KeyLoadingMode_INCLUDE:
		if keysSpecified == nil || keysSpecified.Contains(nodeName) {
			return ``, ``, nil
		}
	case KeyLoadingMode_EXCLUDE:
		if keysSpecified.Contains(nodeName) {
			return ``, ``, nil
		}
	case KeyLoadingMode_ALL:
	default:

	}

	data := this.client.GetData()
	value, err := data.Watched().ForPath(nodePath)
	if err != nil {
		return ``, ``, err
	}

	return nodeName, string(value), nil
}

/**
 * 导出属性列表
 *
 */
func (this *ZookeeperConfigGroup) exportProperties() map[string]string {
	result := make(map[string]string, this.Size())

	this.ForEach(func(key, value string) {
		result[key] = value
	})

	return result
}
