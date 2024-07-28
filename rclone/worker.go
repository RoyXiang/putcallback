package rclone

import (
	"context"
	"fmt"
	"log"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/RoyXiang/putcallback/notification"
	"github.com/RoyXiang/putcallback/putio"
	"github.com/RoyXiang/putcallback/utils"
	"github.com/samber/lo"
)

func (r *Remote) IsValid(p string) bool {
	if r.Path == "" {
		return true
	}
	return strings.HasPrefix(p, r.Path)
}

func (r *Remote) FullPath(p string, ignorePath bool) string {
	var remotePath string
	remoteName := r.Name
	if ignorePath {
		remotePath = p
	} else {
		remotePath = path.Join(r.Path, p)
	}
	if remoteName != "" {
		remoteName += ":"
		if remotePath == "." {
			remotePath = ""
		}
	}
	return remoteName + remotePath
}

func SendFileIdToWorker(fileId int64) {
	callbackMu.Lock()
	defer callbackMu.Unlock()

	fileInfo := Put.GetFileInfo(fileId)
	if fileInfo == nil {
		return
	}
	go Put.CleanupTransfers()

	if delayBeforeTransfer > 0 {
		go func() {
			notification.Send(fmt.Sprintf("%s downloaded, transfer will begin in %s", fileInfo.Name, delayBeforeTransfer))
			time.Sleep(delayBeforeTransfer)
			// Get file info from Put.io again in case the file was trashed
			fileInfo = Put.GetFileInfo(fileId)
			if fileInfo == nil {
				return
			}
			taskChan <- fileInfo
		}()
	} else {
		taskChan <- fileInfo
	}
}

func worker() {
	defer workerWg.Done()
	for fileInfo := range taskChan {
		workerWg.Add(1)
		if fileInfo.IsDir {
			go moveFolder(fileInfo)
		} else {
			go moveFile(fileInfo)
		}
	}
}

func checkBeforeTransfer(info *putio.FileInfo) bool {
	if !remoteSrc.IsValid(info.FullPath) {
		notification.Send(fmt.Sprintf("%s downloaded", info.Name))
		return false
	}
	return true
}

func moveFolder(folder *putio.FileInfo) {
	defer workerWg.Done()

	if !checkBeforeTransfer(folder) {
		log.Printf("Folder %s skipped", folder.Name)
		return
	}

	if err := transferSem.Acquire(context.Background(), argSmallFileTransfers); err != nil {
		log.Printf("Failed acquiring semaphore while moving folder %s", folder.Name)
		return
	}
	defer transferSem.Release(argSmallFileTransfers)

	if folder.Size > 0 {
		log.Printf("Moving folder %s...", folder.Name)

		src := remoteSrc.FullPath(folder.FullPath, true)
		dest := remoteDest.FullPath(folder.Name, false)
		if rcCopyDir(src, dest) {
			Put.DeleteFile(folder.ID)
			notification.Send(fmt.Sprintf("%s moved", folder.Name))
		} else {
			notification.Send(fmt.Sprintf("error occurred, %s was not moved", folder.Name))
		}
	} else {
		Put.DeleteFile(folder.ID)
	}
}

func moveFile(file *putio.FileInfo) {
	defer workerWg.Done()

	if checkBeforeTransfer(file) {
		if ext := filepath.Ext(file.Name); ext != "" && lo.Contains[string](excludeFileTypes, ext[1:]) {
			Put.DeleteFile(file.ID)
			log.Printf("File %s filtered", file.Name)
			return
		}
		log.Printf("Moving file %s...", file.Name)
	} else {
		log.Printf("File %s skipped", file.Name)
		return
	}

	var weight int64
	if file.Size < argMultiThreadCutoff {
		weight = 1
	} else {
		weight = 2
	}
	if err := transferSem.Acquire(context.Background(), weight); err != nil {
		log.Printf("Failed acquiring semaphore while moving file %s", file.Name)
		return
	}
	defer transferSem.Release(weight)

	newFilename := file.Name
	if strings.HasPrefix(file.ContentType, putio.ContentTypeVideo) {
		switch renamingStyle {
		case RenamingStyleAnime:
			newFilename = utils.RenameFileInAnimeStyle(file.Name)
		case RenamingStyleTv:
			newFilename = utils.RenameFileInTvStyle(file.Name)
		}
	}

	src := remoteSrc.FullPath(file.FullPath, true)
	dest := remoteDest.FullPath(newFilename, false)
	if rcCopyFile(src, dest) {
		Put.DeleteFile(file.ID)
		if file.Name == newFilename {
			notification.Send(fmt.Sprintf("%s moved", file.Name))
		} else {
			notification.Send(fmt.Sprintf("%s moved and renamed", file.Name))
		}
	} else {
		notification.Send(fmt.Sprintf("error occurred, %s was not moved", file.Name))
	}
}
