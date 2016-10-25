package config

import (
	"fmt"
	"github.com/go-xweb/log"
	"strconv"
	"sync"
)

type ConfigGroup interface {
	Get(key string) string
}

type IObserver func(key, value string)

type GeneralConfigGroup struct {
	configMap           map[string]string
	lock                *sync.RWMutex
	internalConfigGroup ConfigGroup
	watchs              []IObserver
}

func NewGeneralConfigGroup(internalConfigGroup ConfigGroup) *GeneralConfigGroup {
	return &GeneralConfigGroup{
		lock:                &sync.RWMutex{},
		internalConfigGroup: internalConfigGroup,
		configMap:           make(map[string]string),
		watchs:              make([]IObserver, 0),
	}
}

func (this *GeneralConfigGroup) get(key string) string {
	this.lock.RLock()
	defer this.lock.RUnlock()

	value, ok := this.configMap[key]
	if !ok {
		return ``
	}

	return value
}

func (this *GeneralConfigGroup) Get(key string) string {

	value := this.get(key)

	if len(value) > 0 {
		return value
	}

	if this.internalConfigGroup == nil {
		return ``
	}

	value = this.internalConfigGroup.Get(key)

	this.put(key, value)

	return value
}

func (this *GeneralConfigGroup) GetInt(key string) (int, error) {
	value := this.Get(key)
	result, err := strconv.Atoi(value)
	if err != nil {
		log.Printf(`config error: name:%s, value:%s`, key, value)
		return -1, err
	}

	return result, nil

}

// ParseBool returns the boolean value represented by the string.
//
// It accepts 1, 1.0, t, T, TRUE, true, True, YES, yes, Yes,Y, y, ON, on, On,
// 0, 0.0, f, F, FALSE, false, False, NO, no, No, N,n, OFF, off, Off.
// Any other value returns an error.
func (this *GeneralConfigGroup) GetBool(key string) (bool, error) {
	val := this.Get(key)
	if len(val) > 0 {
		switch val {
		case "1", "t", "T", "true", "TRUE", "True", "YES", "yes", "Yes", "Y", "y", "ON", "on", "On":
			return true, nil
		case "0", "f", "F", "false", "FALSE", "False", "NO", "no", "No", "N", "n", "OFF", "off", "Off":
			return false, nil
		}
	}

	return false, fmt.Errorf("parsing %q: invalid syntax", val)
}

func (this *GeneralConfigGroup) PutAll(configs map[string]string) {
	if configs != nil && len(configs) > 0 {
		for key, value := range configs {
			this.Put(key, value)
		}
	}
}

func (this *GeneralConfigGroup) put(key, value string) string {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.configMap[key] = value

	return value
}

func (this *GeneralConfigGroup) Put(key, value string) string {
	if len(key) == 0 {
		return ``
	}

	preValue := this.Get(key)

	if preValue == value {
		return value
	}

	value = this.put(key, value)

	this.notify(key, value)

	return value
}

func (this *GeneralConfigGroup) Size() int {
	this.lock.RLock()
	defer this.lock.RUnlock()
	return len(this.configMap)
}

func (this *GeneralConfigGroup) clone() map[string]string {
	this.lock.RLock()
	defer this.lock.RUnlock()
	clone := make(map[string]string, len(this.configMap))
	for key, value := range this.configMap {
		clone[key] = value
	}
	return clone
}

func (this *GeneralConfigGroup) ForEach(callback func(key, value string)) {
	for key, value := range this.clone() {
		callback(key, value)
	}
}

func (this *GeneralConfigGroup) AddWatcher(watch IObserver) {
	this.watchs = append(this.watchs, watch)
}

func (this *GeneralConfigGroup) notify(key, value string) {
	for _, observer := range this.watchs {
		go func() {
			observer(key, value)
		}()
	}
}
