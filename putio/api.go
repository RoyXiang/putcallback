package putio

import (
	"context"
)

func (put *Put) GetFileInfo(id int64) (name string, isDir bool) {
	file, err := put.Client.Files.Get(context.Background(), id)
	if err != nil {
		return
	}
	name = file.Name
	isDir = file.IsDir()
	return
}
