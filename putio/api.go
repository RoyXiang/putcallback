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
	if file.ParentID != RootFolderId && file.ParentID == put.DefaultDownloadFolderId {
		fullPath = fmt.Sprintf("%s/%s", put.DefaultDownloadFolder, fullPath)
	} else {
		for folderId := file.ParentID; folderId != RootFolderId; {
			folder, err := put.Client.Files.Get(ctx, folderId)
			if err != nil {
				return nil
			}
			fullPath = fmt.Sprintf("%s/%s", folder.Name, fullPath)
			folderId = folder.ParentID
		}
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

func (put *Put) DeleteFile(id int64) bool {
	ctx := context.Background()
	if _, err := put.Client.Files.Get(ctx, id); err != nil {
		// file may be deleted by user
		return true
	}
	err := put.Client.Files.Delete(ctx, id)
	return err == nil
}
