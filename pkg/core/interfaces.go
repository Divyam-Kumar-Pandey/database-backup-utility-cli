package core

type Config map[string]interface{}

type Database interface {
	Backup(config Config, outputPath string) (string, error)
	Restore(config Config, backupPath string) error
	TestConnection(config Config) error
}

type Storage interface {
	Upload(localPath, remotePath string) (string, error)
	Download(remotePath, localPath string) (string, error)
	ListFiles(prefix string) ([]string, error)
}
