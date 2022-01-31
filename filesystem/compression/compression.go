package compression

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"

	"github.com/neurafuse/tools-go/errors"
	"github.com/neurafuse/tools-go/filesystem"
	"github.com/neurafuse/tools-go/logging"
	"github.com/neurafuse/tools-go/objects/strings"
	"github.com/neurafuse/tools-go/runtime"
	"github.com/neurafuse/tools-go/vars"
)

func ExtractTarGz(filePath string) {
	logging.Log([]string{"", vars.EmojiDir, vars.EmojiProcess}, "Extracting .tar.gz archive: "+filePath, 0)
	gzipStream, err := os.Open(filePath)
	errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to open file: "+filePath, false, false, true)
	uncompressedStream, err := gzip.NewReader(gzipStream)
	errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to create gzip.NewReader!", false, false, true)
	tarReader := tar.NewReader(uncompressedStream)
	for true {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to tarReader.Next()!", false, false, true)
		switch header.Typeflag {
		case tar.TypeDir:
			err := os.Mkdir(header.Name, filesystem.GetFileMode())
			errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to os.MkDir!", false, false, true)
		case tar.TypeReg:
			outFile, err := os.Create(header.Name)
			errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to os.Create()!", false, false, true)
			_, err = io.Copy(outFile, tarReader)
			errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to io.Copy()!", false, false, true)
			outFile.Close()
		default:
			errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unknown type ("+strings.BytesToString([]byte{header.Typeflag})+"/"+header.Name+")!", false, false, true)
		}
	}
	logging.Log([]string{"\n", vars.EmojiDir, vars.EmojiSuccess}, "Extracted .tar.gz archive.", 0)
}
