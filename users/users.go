package users

import (
	"../filesystem"
	"../objects/strings"
)

var idActive string = ""
var BasePath string = "users"
var ProjectPathActive = ""

func SetIDActive(id string) {
	idActive = id
}

func GetIDActive() string {
	return idActive
}

var clusterRecentlyDeleted bool = false
func SetClusterRecentlyDeleted(deleted bool) {
	clusterRecentlyDeleted = deleted
}

func GetClusterRecentlyDeleted() bool {
	return clusterRecentlyDeleted
}

func Create(id string) {
	userPath := BasePath + "/" + id
	if !filesystem.Exists(userPath) {
		filesystem.CreateDir(userPath, false)
	}
}

func Exists(id string) bool {
	return strings.ArrayContains(GetAllIDs(), id)
}

func Existing() bool {
	if !filesystem.Exists(BasePath) {
		filesystem.CreateDir(BasePath, false)
		return false
	} else {
		if len(filesystem.Explorer("files", BasePath, []string{}, []string{"lost+found"})) == 0 {
			return false
		} else {
			return true
		}
	}
}

func GetAllIDs() []string {
	return filesystem.Explorer("files", BasePath, []string{}, []string{"hidden", ".yaml"})
}