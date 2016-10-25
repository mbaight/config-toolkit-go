package config

import (
	"regexp"
	"strings"
)

type FileLocation struct {
	file     string
	protocol string
}

const (
	FILE  = "file"
	HTTP  = "http"
	HTTPS = "https"
)

func newFileLocation(location string) *FileLocation {
	if strings.Index(location, HTTPS) == 0 {
		return &FileLocation{
			file:     strings.TrimLeft(location[len(HTTPS)+1:], `/`),
			protocol: HTTPS,
		}
	}

	if strings.Index(location, HTTP) == 0 {
		return &FileLocation{
			file:     strings.TrimLeft(location[len(HTTP)+1:], `/`),
			protocol: HTTP,
		}
	}

	regex, _ := regexp.Compile(`\w:\\`)

	if strings.Index(location, FILE) == 0 {
		if regex.MatchString(location) {
			return &FileLocation{
				file:     strings.TrimLeft(location[len(FILE)+1:], `/`),
				protocol: FILE,
			}
		}

		return &FileLocation{
			file:     `/` + strings.TrimLeft(location[len(FILE)+1:], `/`),
			protocol: FILE,
		}
	}

	if regex.MatchString(location[:3]) {
		return &FileLocation{
			file:     location,
			protocol: FILE,
		}
	}

	index := strings.Index(location, `:`)
	if index < 0 {
		return &FileLocation{
			file:     location,
			protocol: FILE,
		}
	}

	if index == 0 {
		return &FileLocation{
			file:     location[1:],
			protocol: FILE,
		}
	}

	return &FileLocation{
		file:     location[0:index],
		protocol: strings.ToLower(location[index+1:]),
	}
}

//
func (this *FileLocation) selectProtocol() Protocol {
	switch this.protocol {
	case FILE:
		protocol, err := NewLocalFileProtocol()
		if err != nil {
			return nil
		}
		return protocol
	case HTTP:
		return &HttpProtocol{}
	case HTTPS:
		return &HttpProtocol{}
	}

	return nil
}
