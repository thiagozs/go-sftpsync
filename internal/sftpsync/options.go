package sftpsync

type Options func(*SftpSyncParams) error

type SftpSyncParams struct {
	RemoteBasePath string
	LocalBasePath  string
	Operation      OperationKind
	User           string
	Password       string
	Port           string
	Host           string
	PrivateKey     string
	PublicKey      string
}

func newSftpSyncParams(opts ...Options) (*SftpSyncParams, error) {
	params := &SftpSyncParams{}
	for _, opt := range opts {
		if err := opt(params); err != nil {
			return nil, err
		}
	}
	return params, nil
}

func WithRemoteBasePath(remoteBasePath string) Options {
	return func(params *SftpSyncParams) error {
		params.RemoteBasePath = remoteBasePath
		return nil
	}
}

func WithLocalBasePath(localBasePath string) Options {
	return func(params *SftpSyncParams) error {
		params.LocalBasePath = localBasePath
		return nil
	}
}

func WithOperation(operation OperationKind) Options {
	return func(params *SftpSyncParams) error {
		params.Operation = operation
		return nil
	}
}

func WithUser(user string) Options {
	return func(params *SftpSyncParams) error {
		params.User = user
		return nil
	}
}

func WithPassword(password string) Options {
	return func(params *SftpSyncParams) error {
		params.Password = password
		return nil
	}
}

func WithPort(port string) Options {
	return func(params *SftpSyncParams) error {
		params.Port = port
		return nil
	}
}

func WithHost(host string) Options {
	return func(params *SftpSyncParams) error {
		params.Host = host
		return nil
	}
}

func WithPrivateKey(privateKey string) Options {
	return func(params *SftpSyncParams) error {
		params.PrivateKey = privateKey
		return nil
	}
}

func WithPublicKey(publicKey string) Options {
	return func(params *SftpSyncParams) error {
		params.PublicKey = publicKey
		return nil
	}
}

// Getters -----

func (p *SftpSyncParams) GetRemoteBasePath() string {
	return p.RemoteBasePath
}

func (p *SftpSyncParams) GetLocalBasePath() string {
	return p.LocalBasePath
}

func (p *SftpSyncParams) GetOperation() OperationKind {
	return p.Operation
}

func (p *SftpSyncParams) GetUser() string {
	return p.User
}

func (p *SftpSyncParams) GetPassword() string {
	return p.Password
}

func (p *SftpSyncParams) GetPort() string {
	return p.Port
}

func (p *SftpSyncParams) GetHost() string {
	return p.Host
}

func (p *SftpSyncParams) GetPrivateKey() string {
	return p.PrivateKey
}

func (p *SftpSyncParams) GetPublicKey() string {
	return p.PublicKey
}

// Setters -----

func (p *SftpSyncParams) SetRemoteBasePath(remoteBasePath string) {
	p.RemoteBasePath = remoteBasePath
}

func (p *SftpSyncParams) SetLocalBasePath(localBasePath string) {
	p.LocalBasePath = localBasePath
}

func (p *SftpSyncParams) SetOperation(operation OperationKind) {
	p.Operation = operation
}

func (p *SftpSyncParams) SetUser(user string) {
	p.User = user
}

func (p *SftpSyncParams) SetPassword(password string) {
	p.Password = password
}

func (p *SftpSyncParams) SetPort(port string) {
	p.Port = port
}

func (p *SftpSyncParams) SetHost(host string) {
	p.Host = host
}

func (p *SftpSyncParams) SetPrivateKey(privateKey string) {
	p.PrivateKey = privateKey
}

func (p *SftpSyncParams) SetPublicKey(publicKey string) {
	p.PublicKey = publicKey
}
