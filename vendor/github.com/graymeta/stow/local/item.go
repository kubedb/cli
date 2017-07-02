package local

import (
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Metadata constants describe the metadata available
// for a local Item.
const (
	MetadataPath       = "path"
	MetadataIsDir      = "is_dir"
	MetadataDir        = "dir"
	MetadataName       = "name"
	MetadataMode       = "mode"
	MetadataModeD      = "mode_d"
	MetadataPerm       = "perm"
	MetadataINode      = "inode"
	MetadataSize       = "size"
	MetadataIsHardlink = "is_hardlink"
	MetadataIsSymlink  = "is_symlink"
	MetadataLink       = "link"
)

type item struct {
	path     string
	infoOnce sync.Once // protects info
	info     os.FileInfo
	infoErr  error
	metadata map[string]interface{}
}

func (i *item) ID() string {
	return i.path
}

func (i *item) Name() string {
	return filepath.Base(i.path)
}

func (i *item) Size() (int64, error) {
	err := i.ensureInfo()
	if err != nil {
		return 0, err
	}
	return i.info.Size(), nil
}

func (i *item) URL() *url.URL {
	return &url.URL{
		Scheme: "file",
		Path:   filepath.Clean(i.path),
	}
}

func (i *item) ETag() (string, error) {
	err := i.ensureInfo()
	if err != nil {
		return "", nil
	}
	return i.info.ModTime().String(), nil
}

// Open opens the file for reading.
func (i *item) Open() (io.ReadCloser, error) {
	return os.Open(i.path)
}

// LimitedReadCloser wraps io.LimitedReader and exposes the Close() method.
type LimitedReadCloser struct {
	io.ReadCloser
	io.Reader
}

// Read reads data from the limited reader.
func (l *LimitedReadCloser) Read(p []byte) (int, error) {
	return l.Reader.Read(p)
}

// LimitReadCloser returns a new reader wraps r in an io.LimitReader, but also
// exposes the Close() method.
func LimitReadCloser(r io.ReadCloser, n int64) *LimitedReadCloser {
	return &LimitedReadCloser{ReadCloser: r, Reader: io.LimitReader(r, n)}
}

func (i *item) Partial(length, offset int64) (io.ReadCloser, error) {
	if offset < 0 {
		return nil, errors.New("offset is negative")
	}
	if length < 0 {
		return nil, fmt.Errorf("invalid length %d", length)
	}

	f, err := os.Open(i.path)
	if err != nil {
		return nil, err
	}

	_, err = f.Seek(offset, 0)
	if err != nil {
		f.Close()
		return nil, err
	}
	if length > 0 {
		return LimitReadCloser(f, int64(length)), nil
	}

	return f, nil
}

func (i *item) LastMod() (time.Time, error) {
	err := i.ensureInfo()
	if err != nil {
		return time.Time{}, nil
	}

	return i.info.ModTime(), nil
}

func (i *item) ensureInfo() error {
	i.infoOnce.Do(func() {
		i.info, i.infoErr = os.Lstat(i.path) // retrieve item file info

		i.infoErr = i.setMetadata(i.info) // merge file and metadata maps
		if i.infoErr != nil {
			return
		}
	})
	return i.infoErr
}

func (i *item) setMetadata(info os.FileInfo) error {
	fileMetadata := getFileMetadata(i.path, info) // retrieve file metadata
	i.metadata = fileMetadata
	return nil
}

// Metadata gets stat information for the file.
func (i *item) Metadata() (map[string]interface{}, error) {
	err := i.ensureInfo()
	if err != nil {
		return nil, err
	}
	return i.metadata, nil
}
