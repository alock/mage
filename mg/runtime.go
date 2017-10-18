package mg

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

// CacheEnv is the environment variable that users may set to change the
// location where mage stores its compiled binaries.
const CacheEnv = "MAGEFILE_CACHE"

// verboseEnv is the environment variable that indicates the user requested
// verbose mode when running a magefile.
const verboseEnv = "MAGEFILE_VERBOSE"

// Verbose reports whether a magefile was run with the verbose flag.
func Verbose() bool {
	return os.Getenv(verboseEnv) != ""
}

// CacheDir returns the directory where mage caches compiled binaries.  It
// defaults to $HOME/.magefile, but may be overridden by the MAGEFILE_CACHE
// environment variable.
func CacheDir() string {
	d := os.Getenv(CacheEnv)
	if d != "" {
		return d
	}
	switch runtime.GOOS {
	case "windows":
		return filepath.Join(os.Getenv("HOMEDRIVE"), os.Getenv("HOMEPATH"), "magefile")
	default:
		return filepath.Join(os.Getenv("HOME"), ".magefile")
	}
}

var ctx context.Context
var ctxCancel func()

func Context() (context.Context, func()) {
	if ctx != nil {
		return ctx, ctxCancel
	}

	if os.Getenv("MAGEFILE_TIMEOUT") != "" {
		timeout, err := time.ParseDuration(os.Getenv("MAGEFILE_TIMEOUT"))
		if err != nil {
			fmt.Printf("timeout error: %v\n", err)
			os.Exit(1)
		}

		ctx, ctxCancel = context.WithTimeout(context.Background(), timeout)
	} else {
		ctx = context.Background()
		ctxCancel = func() {}
	}
	return ctx, ctxCancel
}
