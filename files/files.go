package files

import (
	"io"
	"os"
	"path/filepath"
)

//FileToByte for change file in path to []byte
func FileToByte(path string) (string, []byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", nil, err
	}
	defer file.Close()
	filename := filepath.Base(path)
	fileInfo, _ := file.Stat()
	var size int64 = fileInfo.Size()
	p := make([]byte, size)
	for {
		_, err = file.Read(p)
		if err == io.EOF {
			break
		}
	}
	return filename, p, nil
}
