package config

import (
	"fmt"
	"github.com/go-xweb/log"
)

type FileConfigGroup struct {
	configProfile *FileConfigProfile
	location      *FileLocation
	protocol      Protocol
	*GeneralConfigGroup
}

func NewFileConfigGroup(internalConfigGroup ConfigGroup, configProfile *FileConfigProfile, location string) (*FileConfigGroup, error) {
	group := &FileConfigGroup{
		configProfile:      configProfile,
		location:           newFileLocation(location),
		GeneralConfigGroup: NewGeneralConfigGroup(internalConfigGroup),
	}

	if err := group.initConfig(); err != nil {
		return group, err
	}

	return group, nil
}

func (this *FileConfigGroup) initConfig() error {
	protocol := this.location.selectProtocol()
	if protocol == nil {
		return fmt.Errorf("can't resolve protocol:%v", this.location)
	}

	this.protocol = protocol

	contentTypeResolve, err := selectContentTypeResolve(this.configProfile.contentType)
	if err != nil {
		log.Printf("fileConfigGroup init failed :%v, contentType%s", err, this.configProfile.contentType)
		return err
	}

	data, err := this.protocol.Read(this.location)
	if err != nil {
		log.Printf("fileConfigGroup init failed :%v, read file error:%s\n", err, this.location.file)
		return err
	}

	properties, err := contentTypeResolve.resolve(data, this.location.protocol)
	if err != nil {
		log.Printf("fileConfigGroup init failed :%v, contentType%s", err, this.configProfile.contentType)
		return err
	}

	for key, value := range properties {
		this.Put(key, value)
	}

	return nil
}
