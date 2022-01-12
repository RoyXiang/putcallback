package putio

import (
	"context"
	"fmt"
)

func (put *Put) GetFileInfo(id int64) *FileInfo {
	ctx := context.Background()
	file, err := put.Client.Files.Get(ctx, id)
	if err != nil {
		return nil
	}
	fullPath := file.Name
	folderId := file.ParentID
	for folderId != RootFolderId {
		folder, _ := put.Client.Files.Get(ctx, folderId)
		fullPath = fmt.Sprintf("%s/%s", folder.Name, fullPath)
		folderId = folder.ParentID
	}
	return &FileInfo{
		Name:     file.Name,
		IsDir:    file.IsDir(),
		FullPath: fullPath,
	}
}
