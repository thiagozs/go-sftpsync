package sftpsync

import (
	"github.com/thiagozs/go-sftpc"
	"github.com/thiagozs/go-sftpsync/internal/downloader"
	"github.com/thiagozs/go-sftpsync/internal/uploader"
)

type SftpSync struct {
	params *SftpSyncParams
	sftp   *sftpc.SFTPClient
}

func NewSftpSync(opts ...Options) (*SftpSync, error) {
	params, err := newSftpSyncParams(opts...)
	if err != nil {
		return nil, err
	}

	optsftp := []sftpc.Options{
		sftpc.WithHost(params.GetHost()),
		sftpc.WithPort(params.GetPort()),
		sftpc.WithUser(params.GetUser()),
		sftpc.WithPassword(params.GetPassword()),
		sftpc.WithPrivateKeyB64(params.GetPrivateKey()),
	}

	sftp, err := sftpc.NewSFTPClient(optsftp...)
	if err != nil {
		return nil, err
	}

	return &SftpSync{params: params, sftp: sftp}, nil
}

func (s *SftpSync) SyncDownload() error {

	downloader := downloader.NewDownloader(
		s.sftp, s.params.GetRemoteBasePath(),
		s.params.GetLocalBasePath())

	err := downloader.FileWaker()
	if err != nil {
		return err
	}

	err = downloader.DownloadFiles()
	if err != nil {
		return err
	}

	downloader.Report()

	return nil
}

func (s *SftpSync) SyncUpload() error {

	uploader := uploader.NewUploader(
		s.sftp, s.params.GetRemoteBasePath(),
		s.params.GetLocalBasePath())

	err := uploader.FileWaker()
	if err != nil {
		return err
	}

	err = uploader.UploadFiles()
	if err != nil {
		return err
	}

	uploader.Report()

	return nil
}
