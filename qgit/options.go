package qgit

// Options required for setting up the git client.
type Options struct {
	RepoPath string
	RepoUrl  string
	Username string
	Token    string
	//TODO Logger   *logrus.Logger
}

type Option func(*Options) error

// GetDefaultOptions returns default configuration options for a git Client.
func GetDefaultOptions() Options {
	return Options{
		Username: "qlik-pipeline-cd-helper",
	}
}

// TODO WithLogger is an Option to set a logger to be used by the Client.
// func WithLogger(logger *Logger.Log) Option {
// 	return func(opt *Options) error {
// 		opt.Logger = logger
// 		return nil
// 	}
// }

// WithRepoPath is an Option to set the local repo path
func WithRepoPath(path string) Option {
	return func(opt *Options) error {
		opt.RepoPath = path
		return nil
	}
}

// WithRepoUrl is an Option to set the remote repo url
func WithRepoUrl(url string) Option {
	return func(opt *Options) error {
		opt.RepoUrl = url
		return nil
	}
}

// WithUsername is an Option to set token git auth
func WithUsername(username string) Option {
	return func(opt *Options) error {
		opt.Username = username
		return nil
	}
}

// WithToken is an Option to set token git auth
func WithToken(token string) Option {
	return func(opt *Options) error {
		opt.Token = token
		return nil
	}
}


func compileOptions(opts ...Option) (*Options, error) {
	options := GetDefaultOptions()
	for _, opt := range opts {
		if err := opt(&options); err != nil {
			return nil, err
		}
	}
	return &options, nil
}
