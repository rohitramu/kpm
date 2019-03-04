package types

import (
	"os"
	"time"
)

// MockFileInfo implements os.FileInfo to satisfy the tar writer interface when adding streamed files
type MockFileInfo struct {
	MockName    string
	MockSize    int64
	MockMode    os.FileMode
	MockModTime time.Time
	MockIsDir   bool
	MockSys     interface{}
}

// Name returns the base name of the file
func (f MockFileInfo) Name() string { return f.MockName }

// Size returns the length in bytes for regular files; system-dependent for others
func (f MockFileInfo) Size() int64 { return f.MockSize }

// Mode returns the file mode bits
func (f MockFileInfo) Mode() os.FileMode { return f.MockMode }

// ModTime returns the modification time
func (f MockFileInfo) ModTime() time.Time { return f.MockModTime }

// IsDir is an abbreviation for Mode().IsDir()
func (f MockFileInfo) IsDir() bool { return f.MockIsDir }

// Sys is an underlying data source (can return nil)
func (f MockFileInfo) Sys() interface{} { return f.MockSys }
