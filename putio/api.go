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
		folder, err := put.Client.Files.Get(ctx, folderId)
		if err != nil {
			return nil
		}
		fullPath = fmt.Sprintf("%s/%s", folder.Name, fullPath)
		folderId = folder.ParentID
	}
	return &FileInfo{
		ID:          file.ID,
		Name:        file.Name,
		IsDir:       file.IsDir(),
		Size:        file.Size,
		FullPath:    fullPath,
		ContentType: file.ContentType,
	}
}

func (put *Put) DeleteFolder(id int64, force bool) bool {
	ctx := context.Background()
	if !force {
		file, err := put.Client.Files.Get(ctx, id)
		if err != nil || !file.IsDir() || file.Size > 0 {
			return false
		}
	}
	err := put.Client.Files.Delete(ctx, id)
	return err == nil
}
