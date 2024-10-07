package sftpsync

import "strings"

type OperationKind int

const (
	UNKNOWN OperationKind = iota
	DOWNLOAD
	UPLOAD
)

func NewOperationKind() OperationKind {
	return UNKNOWN
}

func (o OperationKind) String() string {
	return [...]string{"UNKNOWN", "DOWNLOAD", "UPLOAD"}[o]
}

func (o OperationKind) Int() int {
	return int(o)
}

func (o OperationKind) IsValid() bool {
	return o == DOWNLOAD || o == UPLOAD
}

func (o OperationKind) GetFromString(s string) OperationKind {
	switch strings.ToUpper(s) {
	case DOWNLOAD.String():
		return DOWNLOAD
	case UPLOAD.String():
		return UPLOAD
	default:
		return UNKNOWN
	}
}
