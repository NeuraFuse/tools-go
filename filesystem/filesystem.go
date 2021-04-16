package filesystem

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"

	//"os/user"
	"path"
	"path/filepath"

	"../errors"
	"../objects/strings"
	ostrings "../objects/strings"
	osTools "../os"
	"../runtime"
	copy "github.com/otiai10/copy"
)

func GiveProgramPermissions(path, user string) {
	ChangeOwnership(path, user, true)
	ChangeFilePermission(path, 0775)
}

func GetFileMode() os.FileMode {
	defaultChmodFilePermissions := 0777
	return os.FileMode(defaultChmodFilePermissions)
}

func FileContainsString(filePath, str string) bool {
	file := FileToString(filePath)
	return strings.Contains(file, str)
}

func AppendStringToFile(filePath, input string) {
	// If the file doesn't exist, create it, or append to the file
	CreateEmptyDir(filePath)
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, GetFileMode())
	errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to open file: "+filePath+"!", false, false, true)
	_, err = f.WriteString(input + "\n")
	errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to write string to file: "+filePath+"!", false, false, true)
	err = f.Close()
	errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to close file operation: "+filePath+"!", false, false, true)
}

func GetSize(path, unit string) float64 {
	var sizeInt int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			sizeInt += info.Size()
		}
		return err
	})
	errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to get directory size: "+path+"!", false, false, true)
	size := float64(sizeInt)
	switch unit {
	case "mb":
		size = size / (1024 * 1024)
	case "gb":
		size = size / (1024 * 1024) / 1024
	}
	return size
}

func GetUserHomeDir() string {
	usr, _ := user.Current()
	return usr.HomeDir
}

func Copy(sourcePath, destPath string, sudo bool) error {
	var err error
	var errMsg string
	if sudo {
		CreateDir(destPath, true)
		program := "sudo cp"
		args := "-r " + sourcePath + " " + destPath
		c := exec.Command("/bin/sh", "-c", program+" "+args)
		err = c.Run()
		errMsg = "Unable to sudo copy path: " + program + " " + args
	} else {
		err = copy.Copy(sourcePath, destPath)
		errMsg = "Unable to copy path: " + sourcePath + " --> " + destPath + " !"
	}
	errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), errMsg, false, false, true)
	return err
}

func Delete(path string, sudo bool) {
	if Exists(path) {
		absPath := GetAbsolutePathToFile(path)
		if sudo {
			program := "sudo rm"
			args := "-rf " + path
			c := exec.Command("/bin/sh", "-c", program+" "+args)
			err := c.Run()
			errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to sudo delete path: "+program+" "+args, false, false, true)
		} else {
			err := os.RemoveAll(absPath)
			errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to delete path: "+path+"!", false, true, true)
		}
	}
}

func Exists(filePath string) bool {
	absPath := GetAbsolutePathToFile(filePath)
	exists := false
	if _, err := os.Stat(absPath); !(os.IsNotExist(err)) {
		exists = true
	}
	return exists
}

func FileToString(filePath string) string {
	byteArr, err := ioutil.ReadFile(filePath)
	errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to load file into string: "+filePath+"!", false, false, true)
	return string(byteArr)
}

func ChangeFilePermission(filePath string, permission int) {
	err := os.Chmod(filePath, os.FileMode(permission))
	errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to change file permission: "+filePath+" / "+ostrings.ToString(permission), false, false, true)
}

func ChangeOwnership(path, username string, sudo bool) {
	var err error
	if sudo {
		program := "sudo chown"
		args := "-R " + username + " " + path
		c := exec.Command("/bin/sh", "-c", program+" "+args)
		err = c.Run()
	} else {
		uid := osTools.GetHostUID(username)
		gid := osTools.GetHostGID(username)
		err = os.Chown(path, strings.ToInt(uid), strings.ToInt(gid))
	}
	errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to change ownership: "+username+" | "+path, false, false, true)
}

func RemoveFile(filePath string) {
	err := os.Remove(filePath)
	errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to remove file: "+filePath+"!", false, false, true)
}

func Move(sourcePath, destPath string, sudo bool) {
	Copy(sourcePath, destPath, sudo)
	Delete(sourcePath, sudo)
}

func CreateEmptyDir(filePath string) {
	dir, _ := filepath.Split(filePath)
	if !Exists(dir) {
		CreateDir(dir, false)
	}
}

func CreateEmptyFile(filePath string) {
	CreateEmptyDir(filePath)
	err := ioutil.WriteFile(filePath, []byte(""), GetFileMode())
	errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to write empty file!", false, false, true)
}

func SaveByteArrayToFile(byteArray []byte, filePath string) {
	createFolderFromFilePath(filePath)
	err := ioutil.WriteFile(filePath, byteArray, GetFileMode())
	errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to save byte array to file: "+filePath+" !", false, false, true)
}

func createFolderFromFilePath(filePath string) {
	if !Exists(filePath) {
		CreateDir(GetDirPathFromFilePath(filePath), false)
	}
}

func GetDirPathFromFilePath(filePath string) string {
	return path.Dir(filePath)
}

func GetFileNameFromDirPath(dirPath string) string {
	return filepath.Base(dirPath)
}

func StreamToBytes(stream io.Reader) []byte {
	buf := new(bytes.Buffer)
	buf.ReadFrom(stream)
	return buf.Bytes()
}

func GetAbsolutePathToFile(filePath string) string {
	absPath, _ := filepath.Abs(filePath)
	return absPath
}

func FileToBytes(filePath string) []byte {
	filename, _ := filepath.Abs(filePath)
	byteArray, err := ioutil.ReadFile(filename)
	errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to load file to byte array: "+filePath+" !", false, false, true)
	return byteArray
}

func RenameFile(filePathOld string, filePathNew string) {
	err := os.Rename(filePathOld, filePathNew)
	errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to rename file: "+filePathOld+" --> "+filePathNew+" !", false, false, true)
}

func Explorer(objectType, path string, search, exclude []string) []string {
	var objects []string
	if objectType == "files" {
		files, err := os.Open(path)
		errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to open directory: "+path+" !", false, false, true)
		defer files.Close()
		list, _ := files.Readdirnames(0)
		for _, object := range list {
			objects = append(objects, object)
		}
	} else if objectType == "directories" {
		files, err := ioutil.ReadDir(path)
		errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to read directory: "+path+" !", false, false, true)
		for _, object := range files {
			objects = append(objects, object.Name())
		}
	}
	for _, exc := range exclude {
		objects = filterFiles(objects, exc, false)
	}
	for _, se := range search {
		objects = filterFiles(objects, se, true)
	}
	return objects
}

func filterFiles(input []string, filter string, filterOnly bool) []string {
	var filtered []string
	for _, entry := range input {
		var filterCond bool
		if filterOnly {
			filterCond = strings.HasSuffix(entry, filter)
		} else {
			filterCond = !strings.HasSuffix(entry, filter)
			if filter == "." {
				filterCond = !strings.HasPrefix(entry, filter)
			}
		}
		if filterCond {
			filtered = append(filtered, entry)
		}
	}
	return filtered
}

func CreateDir(dirPath string, sudo bool) error {
	var err error
	var errMsg string
	if sudo {
		program := "sudo mkdir"
		args := "-p " + dirPath
		c := exec.Command("/bin/sh", "-c", program+" "+args)
		err = c.Run()
		errMsg = "Unable to sudo create dir: " + dirPath + "!"
	} else {
		err = os.MkdirAll(dirPath, GetFileMode())
		errMsg = "Unable to create dir: " + dirPath + "!"
	}
	errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), errMsg, false, false, true)
	return err
}

func DirIsEmpty(path string) bool {
	f, err := os.Open(path)
	errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to open directory: "+path+" !", false, false, true)
	defer f.Close()
	_, err = f.Readdirnames(1) // Or f.Readdir(1)
	if err == io.EOF {
		return true
	}
	return false // Either not empty or error, suits both cases
}

func GetWorkingDir() string {
	dir, err := os.Getwd()
	errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "", false, false, true)
	return dir + "/"
}
