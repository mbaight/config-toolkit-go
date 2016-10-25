package config

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"github.com/cyfdecyf/bufio"
	"github.com/go-xweb/log"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
	"io"
)

type ContentTypeResolve interface {
	resolve(data []byte, encoding string) (map[string]string, error)
}

type (
	PropertiesContentType struct{}
	XmlContentType        struct{}
	JsonContentType       struct{}
)

var (
	bNumComment = []byte{'#'} // number signal
	bEqual      = []byte{'='} // equal signal
)

const (
	ContentType_properties = iota
	ContentType_xml
)

func bufferReader(data []byte, encoding string) io.Reader {
	return transform.NewReader(bytes.NewReader([]byte(data)), unicode.UTF8.NewEncoder())
}

//从XML内容解析出属性列表
func (this *PropertiesContentType) resolve(data []byte, encoding string) (map[string]string, error) {
	result := make(map[string]string, 10)
	reader := bufio.NewReader(bufferReader(data, encoding))
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}

		if len(line) == 0 {
			continue
		}

		line = bytes.TrimSpace(line)
		if bytes.HasPrefix(line, bNumComment) {
			continue
		}

		keyValue := bytes.SplitN(line, bEqual, 2)
		if len(keyValue) != 2 {
			return result, fmt.Errorf("parse ini property error:%v, content:%v", err, line)
		}

		key := string(bytes.TrimSpace(keyValue[0]))   // key name case insensitive
		value := string(bytes.TrimSpace(keyValue[1])) // key name case insensitive
		result[key] = value
	}

	return result, nil
}

//从ini内容解析出属性列表
func (this *XmlContentType) resolve(data []byte, encoding string) (map[string]string, error) {
	result := make(map[string]string, 10)
	decoder := xml.NewDecoder(bufferReader(data, encoding))
	var key string
	var err error
	for t, err := decoder.Token(); err == nil; t, err = decoder.Token() {
		switch token := t.(type) {
		case xml.StartElement:
			if "entry" == token.Name.Local {
				for _, attr := range token.Attr {
					if "key" == attr.Name.Local {
						key = attr.Value
						break
					}
				}
			}
		case xml.EndElement:
		case xml.CharData:
			result[key] = string([]byte(token))
		default:
			log.Errorf("parse xml fail:%v", token)
			return nil, fmt.Errorf("parse xml fail:%v", token)
		}
	}

	return result, err
}

func (this *JsonContentType) resolve(data []byte, encoding string) (map[string]string, error) {
	return nil, fmt.Errorf("properties with json not implement:%v", encoding)
}

//根据ContentType选择解析器
func selectContentTypeResolve(contentType int) (ContentTypeResolve, error) {
	switch contentType {
	case ContentType_properties:
		return &PropertiesContentType{}, nil
	case ContentType_xml:
		return &XmlContentType{}, nil
	default:
		return nil, fmt.Errorf("contentType not implement with:%d", contentType)
	}
}
