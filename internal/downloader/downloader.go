package downloader

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/thiagozs/go-sftpc"
	"github.com/thiagozs/go-sftpsync/internal/domain"
	"github.com/thiagozs/go-sftpsync/pkg/csize"
	"github.com/thiagozs/go-sftpsync/pkg/utils"
	"github.com/thiagozs/go-xutils"
)

type Downloader struct {
	content        []domain.SftpContent
	sftp           *sftpc.SFTPClient
	cdirs          int
	cfiles         int
	utils          *xutils.XUtils
	localBasePath  string
	remoteBasePath string
}

func NewDownloader(sftp *sftpc.SFTPClient, remoteBasePath,
	localBasePath string) *Downloader {
	return &Downloader{
		utils:          xutils.New(),
		sftp:           sftp,
		localBasePath:  localBasePath,
		remoteBasePath: remoteBasePath,
	}
}

func (d *Downloader) FileWaker() error {
	err := d.sftp.WalkFile(d.remoteBasePath, func(path string, info os.FileInfo) error {
		if info.IsDir() {
			log.Printf("Scan Directory: %s\n", path)
			d.content = append(d.content, domain.SftpContent{Path: path, Info: info, Dir: true})

			d.cdirs++

		} else {
			log.Printf("Scan File: %s (Size: %s)\n", path, csize.FormatSize(info.Size()))
			d.content = append(d.content, domain.SftpContent{Path: path, Info: info, File: true})

			d.cfiles++
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to walk remote directory: %v", err)
	}

	return nil
}

func (d *Downloader) DownloadFiles() error {
	if d.utils.Files().DirectoryExist(d.localBasePath) {
		log.Printf("Local Directory exists: %s\n", d.localBasePath)
	} else {
		log.Printf("Local Directory does not exist, create: %s\n", d.localBasePath)
		if err := d.utils.Files().CreateDirAll(d.localBasePath); err != nil {
			return fmt.Errorf("failed to create directory: %v", err)
		}
	}

	for _, cont := range d.content {
		if cont.Dir {
			log.Printf("Remote Directory: %s\n", cont.Path)

			_, sec, err := utils.ParsePath(cont.Path, d.remoteBasePath)
			if err != nil {
				return fmt.Errorf("failed to parse path: %v", err)
			}

			localPath := filepath.Join(d.localBasePath, sec)
			log.Printf("Local Directory: %s\n", localPath)

			if d.utils.Files().DirectoryExist(localPath) {
				log.Printf("Directory exists: %s\n", localPath)
				continue
			} else {
				log.Printf("Directory does not exist, create: %s\n", localPath)
				if err := d.utils.Files().CreateDirAll(localPath); err != nil {
					return fmt.Errorf("failed to create directory: %v", err)
				}
			}

		} else if cont.File {
			log.Printf("Remote File: %s\n", cont.Path)

			_, sec, err := utils.ParsePath(cont.Path, d.remoteBasePath)
			if err != nil {
				return fmt.Errorf("failed to parse path: %v", err)
			}

			localPath := filepath.Join(d.localBasePath, sec)

			if d.utils.Files().IsFile(localPath) {
				infoLocal, err := utils.GetFileInfo(localPath)
				if err != nil {
					return fmt.Errorf("failed to get file info: %v", err)
				}

				log.Printf("Local File: %s (Size: %s)\n", localPath, csize.FormatSize(infoLocal.Size()))

				if cont.Info.Size() != infoLocal.Size() {
					log.Printf("File exists: %s (Size: %s)\n", localPath, csize.FormatSize(infoLocal.Size()))
					if err := d.utils.Files().RemoveFile(localPath); err != nil {
						return fmt.Errorf("failed to remove file: %v", err)
					}
					log.Printf("download: %s (Size: %s)\n", localPath, csize.FormatSize(infoLocal.Size()))
					if err := d.sftp.DownloadFile(cont.Path, localPath); err != nil {
						return fmt.Errorf("failed to download file: %v", err)
					}
				} else {
					log.Printf("File exists and is up-to-date: %s (Size: %s)\n", localPath, csize.FormatSize(infoLocal.Size()))
				}
			} else {
				log.Printf("File does not exist, download: %s (Size: %s)\n", localPath, csize.FormatSize(cont.Info.Size()))
				if err := d.sftp.DownloadFile(cont.Path, localPath); err != nil {
					return fmt.Errorf("failed to download file: %v", err)
				}
			}
		}

	}

	return nil
}

func (d *Downloader) Report() {
	fmt.Printf("Directories: %d\n", d.cdirs)
	fmt.Printf("Files: %d\n", d.cfiles)
}
