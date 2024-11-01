package uploader

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/thiagozs/go-sftpc"
	"github.com/thiagozs/go-sftpsync/internal/domain"
	"github.com/thiagozs/go-sftpsync/pkg/csize"
	"github.com/thiagozs/go-sftpsync/pkg/utils"
	"github.com/thiagozs/go-xutils"
)

type Uploader struct {
	content        []domain.SftpContent
	sftp           *sftpc.SFTPClient
	cdirs          int
	cfiles         int
	utils          *xutils.XUtils
	localBasePath  string
	remoteBasePath string
}

func NewUploader(sftp *sftpc.SFTPClient, remoteBasePath,
	localBasePath string) *Uploader {
	return &Uploader{
		utils:          xutils.New(),
		sftp:           sftp,
		localBasePath:  localBasePath,
		remoteBasePath: remoteBasePath,
	}
}

func (u *Uploader) FileWaker() error {
	err := filepath.Walk(u.localBasePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			log.Printf("Scan Directory: %s\n", path)
			u.content = append(u.content, domain.SftpContent{Path: path, Info: info, Dir: true})
			u.cdirs++
		} else {
			log.Printf("Scan File: %s (Size: %s)\n", path, csize.FormatSize(info.Size()))

			if strings.HasSuffix(path, ".go") {
				log.Printf("Skipping file: %s\n", path)
				return nil
			}

			u.content = append(u.content, domain.SftpContent{Path: path, Info: info, File: true})
			u.cfiles++
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to process path: %v", err)
	}

	return nil
}

func (u *Uploader) UploadFiles() error {
	if err := u.sftp.CreateRemoteDirRecursive(u.remoteBasePath); err != nil {
		return fmt.Errorf("failed to create remote directory: %v", err)
	}

	for _, cont := range u.content {
		if cont.Dir {
			log.Printf("Local Dir: %s\n", cont.Path)

			fmt.Println(cont.Path)
			fmt.Println(u.remoteBasePath)

			var remotePath string

			if u.remoteBasePath != "" &&
				strings.Contains(cont.Path, u.remoteBasePath) {
				log.Printf("CONTAINS Remote Base Path: %s\n", u.remoteBasePath)
				_, sec, err := utils.ParsePath(cont.Path, u.remoteBasePath)
				if err != nil {
					return fmt.Errorf("failed to parse path: %v", err)
				}
				remotePath = filepath.Join(u.remoteBasePath, sec)

			} else {
				log.Printf("DOES NOT CONTAIN Remote Base Path: %s\n", u.remoteBasePath)
				paths := strings.Split(cont.Path, u.localBasePath)
				remotePath = filepath.Join(u.remoteBasePath, paths[1])
			}

			log.Printf("Remote Dir: %s\n", remotePath)

			if err := u.sftp.CreateRemoteDirRecursive(remotePath); err != nil {
				return fmt.Errorf("failed to create remote directory: %v", err)
			}

		} else if cont.File {
			log.Printf("Local File: %s Size:(%s)\n", cont.Path, csize.FormatSize(cont.Info.Size()))

			fmt.Println(cont.Path)
			fmt.Println(u.remoteBasePath)

			var remotePath string

			if u.remoteBasePath != "" &&
				strings.Contains(cont.Path, u.remoteBasePath) {
				_, sec, err := utils.ParsePath(cont.Path, u.remoteBasePath)
				if err != nil {
					return fmt.Errorf("failed to parse path: %v", err)
				}

				remotePath = filepath.Join(u.remoteBasePath, sec)

			} else {
				//fileName := filepath.Base(cont.Path)
				paths := strings.Split(cont.Path, u.localBasePath)
				remotePath = filepath.Join(u.remoteBasePath, paths[1])
				//remotePath = filepath.Join(u.remoteBasePath, fileName)
			}

			log.Printf("Remote File: %s\n", remotePath)

			remoteInfo, err := u.sftp.FileInfo(remotePath)
			if err != nil {
				if !strings.Contains(err.Error(), "file does not exist") {
					log.Printf("File not exist, Upload: %s (Size: %s)\n", cont.Path, csize.FormatSize(cont.Info.Size()))
					if err := u.sftp.UploadFile(cont.Path, remotePath); err != nil {
						return fmt.Errorf("failed to upload file: %v", err)
					}
					continue
				}
			}

			if remoteInfo == nil {
				log.Printf("Upload: %s (Size: %s)\n", cont.Path, csize.FormatSize(cont.Info.Size()))
				if err := u.sftp.UploadFile(cont.Path, remotePath); err != nil {
					// return fmt.Errorf("failed to upload file: %v", err)

					remotePath = remotePath[1:]

					remoteInfo, err = u.sftp.FileInfo(remotePath)
					if err != nil {
						log.Printf("Retry Upload: %s (Size: %s)\n", cont.Path, csize.FormatSize(cont.Info.Size()))

						if err := u.sftp.UploadFile(cont.Path, remotePath); err != nil {
							return fmt.Errorf("failed to upload file: %v", err)
						}
						continue

					}

					log.Println("Remotepath changed :", remotePath)

					if remoteInfo.Size() == cont.Info.Size() {
						log.Printf("File exists and is up-to-date: %s (Size: %s)\n", cont.Path, csize.FormatSize(cont.Info.Size()))
						continue
					}

					log.Printf("Retry Upload: %s (Size: %s)\n", cont.Path, csize.FormatSize(cont.Info.Size()))

					if err := u.sftp.UploadFile(cont.Path, remotePath); err != nil {
						return fmt.Errorf("failed to upload file: %v", err)
					}
				}

			} else {
				if remoteInfo.Size() == cont.Info.Size() {
					log.Printf("File exists and is up-to-date: %s (Size: %s)\n", cont.Path, csize.FormatSize(cont.Info.Size()))
					continue

				} else if remoteInfo.Size() != cont.Info.Size() {
					if err := u.sftp.RemoveFile(remotePath); err != nil {
						remotePath = remotePath[1:]
						if err := u.sftp.RemoveFile(remotePath); err != nil {
							return fmt.Errorf("failed to remove file: %v", err)
						}
					}
				}

				log.Printf("File exists but is outdated: %s\n", remotePath)
				log.Printf("Upload: %s (Size: %s)\n", remotePath, csize.FormatSize(cont.Info.Size()))
				if err := u.sftp.UploadFile(cont.Path, remotePath); err != nil {
					return fmt.Errorf("failed to upload file: %v", err)
				}
			}
		}
	}

	return nil
}

func (u *Uploader) Report() {
	fmt.Printf("Directories: %d\n", u.cdirs)
	fmt.Printf("Files: %d\n", u.cfiles)
}
