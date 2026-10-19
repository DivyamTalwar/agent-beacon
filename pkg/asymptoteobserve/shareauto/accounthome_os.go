//go:build !windows && !((!cgo && !darwin && !plan9) || android || (osusergo && !plan9))

package shareauto

import "os/user"

// accountHome asks the system user database: libSystem on macOS, getpwuid_r with cgo. Neither
// consults the environment.
func accountHome() string {
	u, err := user.Current()
	if err != nil {
		return ""
	}
	return u.HomeDir
}
