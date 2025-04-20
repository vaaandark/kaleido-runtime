package kmsglog

import (
	"fmt"
	"os"
	"sync"
)

const (
	kaleidoRuntimeTag = "kaleido-runtime"
	kmsgPath          = "/dev/kmsg"
	// Facility for user-level messages
	facility = 1
	// Default priority (info level)
	defaultPriority = 6
)

var (
	once   sync.Once
	kmsgFD *os.File
)

// initKmsg opens /dev/kmsg with sync.Once
func initKmsg() error {
	var err error
	once.Do(func() {
		kmsgFD, err = os.OpenFile(kmsgPath, os.O_WRONLY|os.O_APPEND, 0)
	})
	return err
}

// Log writes message to kernel log with specified priority
func Log(priority int, format string, args ...interface{}) error {
	if err := initKmsg(); err != nil {
		return fmt.Errorf("failed to open %s: %v", kmsgPath, err)
	}

	msg := fmt.Sprintf(format, args...)
	// Format: <priority>tag: message\n
	_, err := fmt.Fprintf(kmsgFD, "<%d>%s: %s\n", facility*8+priority, kaleidoRuntimeTag, msg)
	return err
}

// InfoF logs at info level (priority 6)
func InfoF(format string, args ...interface{}) {
	_ = Log(defaultPriority, format, args...)
}

// ErrorF logs at error level (priority 3)
func ErrorF(format string, args ...interface{}) {
	_ = Log(3, format, args...)
}

// WarnF logs at warning level (priority 4)
func WarnF(format string, args ...interface{}) {
	_ = Log(4, format, args...)
}

// DebugF logs at debug level (priority 7)
func DebugF(format string, args ...interface{}) {
	_ = Log(7, format, args...)
}
