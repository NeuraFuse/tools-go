package hash

import (
	"os"
	"crypto/sha256"
	"io"
	"../../errors"
	"../../runtime"
	//"../../filesystem"
	"../../objects/strings"
	"golang.org/x/mod/sumdb/dirhash"
)

func SHA256Folder(dirPath string) string {
	hash, err := dirhash.HashDir(dirPath, "", dirhash.DefaultHash)
	errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to generate SHA-256 hash for dir: "+dirPath+"!", false, true, true)
	return hash
}

func SHA256File(filePath string) string { // TODO: Bugfix
	f, err := os.Open(filePath)
	errMsg := "Unable to generate SHA-256 hash for file: "+filePath+"!"
	errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), errMsg, false, true, true)
	defer f.Close()
	h := sha256.New()
	_, err = io.Copy(h, f)
	errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), errMsg, false, true, true)
	return strings.BytesToString(h.Sum(nil))
}

/*func SHA256File(filePath string) string {
	var hash string
	var file []byte = filesystem.FileToBytes(filePath)
	sha256 := sha256.Sum256(file)
	hash = strings.BytesToString([]byte(sha256[:]))
	return hash
}*/