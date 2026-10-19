//go:build !((!cgo && !darwin && !windows && !plan9) || android || (osusergo && !windows && !plan9))

package shareauto

import "os/user"

// accountHome asks the system user database: libSystem on macOS, getpwuid_r with cgo, the
// process token on Windows. None of them consult the environment.
func accountHome() string {
	u, err := user.Current()
	if err != nil {
		return ""
	}
	return u.HomeDir
}
