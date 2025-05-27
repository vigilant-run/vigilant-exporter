package tail

import (
	"fmt"
	"os"

	"github.com/nxadm/tail"
)

type TailConfig struct {
	Path        string
	StartOffset int64
}

func NewTail(cfg TailConfig) (*tail.Tail, error) {
	if _, err := os.Stat(cfg.Path); err != nil {
		if os.IsNotExist(err) {
			return nil, TailFileDoesNotExistError{path: cfg.Path}
		} else if os.IsPermission(err) {
			return nil, TailFileInvalidPermissionError{path: cfg.Path}
		}
		return nil, TailFileError{path: cfg.Path, err: err}
	}

	file, err := os.Open(cfg.Path)
	if err != nil {
		if os.IsPermission(err) {
			return nil, TailFileInvalidPermissionError{path: cfg.Path}
		}
		return nil, TailFileError{path: cfg.Path, err: err}
	}
	file.Close()

	tcfg := tail.Config{
		Follow:    true,
		ReOpen:    true,
		Poll:      true,
		MustExist: true,
		Logger:    tail.DiscardingLogger,
		Location: &tail.SeekInfo{
			Offset: cfg.StartOffset,
			Whence: 0,
		},
	}

	t, err := tail.TailFile(cfg.Path, tcfg)
	if err != nil {
		return nil, err
	}

	return t, nil
}

type TailFileDoesNotExistError struct {
	path string
}

func (e TailFileDoesNotExistError) Error() string {
	return fmt.Sprintf("file does not exist: %s", e.path)
}

type TailFileInvalidPermissionError struct {
	path string
}

func (e TailFileInvalidPermissionError) Error() string {
	return fmt.Sprintf("file is not readable due to invalid permissions: %s", e.path)
}

type TailFileError struct {
	path string
	err  error
}

func (e TailFileError) Error() string {
	return fmt.Sprintf("file cannot be read: %s: %v", e.path, e.err)
}
