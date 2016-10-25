package config

import (
	"bytes"
	"github.com/fsnotify/fsnotify"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type Protocol interface {
	Read(location *FileLocation) ([]byte, error)
	Watch(location *FileLocation, fileConfigGroup *FileConfigGroup) error
	Close() error
}

type LocalFileProtocol struct {
	watcher *fsnotify.Watcher
}

type FileChangeEventListener struct {
	watch       *fsnotify.Watcher
	configGroup *FileConfigGroup
	watchedFile string
}

func NewLocalFileProtocol() (*LocalFileProtocol, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Println("add watch error ")
		return nil, err
	}

	return &LocalFileProtocol{
		watcher: watcher,
	}, nil
}

//监听文件变化,重新加载属性列表
func (this *LocalFileProtocol) Watch(location *FileLocation, fileConfigGroup *FileConfigGroup) error {

	path, err := filepath.Abs(location.file)
	if err != nil {
		log.Printf("file path error:%v file path:%s\n", err, path)
		return err
	}

	err = this.watcher.Add(path)
	if err != nil {
		log.Printf("add file watch error:%v file path:%s\n", err, path)
		return err
	}

	go func() {
		for {
			select {
			case event := <-this.watcher.Events:
				log.Printf("accept watch event:%v", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Printf("modified file:%s\n", event.Name)
					fileConfigGroup.initConfig()
				}
			case err := <-this.watcher.Errors:
				log.Printf("accept watch error:%v\n", err)
			}
		}
	}()

	return nil
}

// 检查文件或目录是否存在
// 如果由 filename 指定的文件或目录存在则返回 true，否则返回 false
func exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

//获取文件里面的属性内容
func (this *LocalFileProtocol) Read(location *FileLocation) ([]byte, error) {
	path, err := filepath.Abs(location.file)
	if err != nil {
		log.Printf("file path error:%v file path:%s\n", err, path)
		return nil, err
	}

	_, err = os.Stat(path)
	if err != nil && !os.IsExist(err) {
		return nil, err
	}

	result, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("read file error:%v file path:%s\n", err, path)
		return nil, err
	}

	return result, nil
}

func (this *LocalFileProtocol) Close() error {
	return this.watcher.Close()
}

type HttpProtocol struct{}

//从http接口获取属性内容
func (this *HttpProtocol) Read(location *FileLocation) ([]byte, error) {
	url := location.protocol + `://` + string(bytes.TrimLeft([]byte(location.file), `/`))
	resp, err := http.Get(url)
	if err != nil {
		log.Printf(`request http properties error:%v %v`, location.protocol+`:`+location.file, err)
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf(`request http properties error:%v %v`, location.protocol+`:`+location.file, err)
		return nil, err
	}

	return body, nil
}

func (this *HttpProtocol) Watch(location *FileLocation, fileConfigGroup *FileConfigGroup) error {
	return nil
}

func (this *HttpProtocol) Close() error {
	return nil
}
