package config

import "strings"

const seprator = `/`

/**
 * Given a parent path and a child node, create a combined full path
 *
 * @param parent the parent
 * @param child  the child
 * @return full path
 */
func MakePath(parent, child string) string {

	if len(parent) == 0 {
		return seprator
	}

	path := parent
	if strings.Index(path, seprator) != 0 {
		path = seprator + path
	}

	if len(child) == 0 {
		return path
	}

	if strings.LastIndex(parent, seprator) != len(parent)-1 {
		path += seprator
	}

	if strings.LastIndex(child, seprator) == len(child)-1 {
		return path + child[1:]
	} else {
		return path + child
	}
}

/**
 * Given a full path, return the node name. i.e. "/one/two/three" will return "three"
 *
 * @param path the path
 * @return the node
 */
func getNodeFromPath(path string) string {
	index := strings.LastIndex(path, seprator)
	if index < 0 {
		return path
	}

	if index+1 >= len(path) {
		return ""
	}

	return path[index+1:]
}
