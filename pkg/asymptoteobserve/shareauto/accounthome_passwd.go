//go:build (!cgo && !darwin && !windows && !plan9) || android || (osusergo && !windows && !plan9)

package shareauto

import (
	"os"
	"strconv"
	"sync"
)

var (
	accountUID = strconv.Itoa(os.Getuid())
	listedHome = sync.OnceValue(func() string {
		data, err := os.ReadFile("/etc/passwd")
		if err != nil {
			return ""
		}
		return passwdHome(string(data), accountUID)
	})
	directoryHome = sync.OnceValue(func() string { return getentHome(accountUID) })
)

// accountHome reads the user database itself: in this build os/user's Current and LookupId fall
// back to $HOME when the uid is not in /etc/passwd. getent reaches directory-service accounts
// through NSS; without it an unlisted uid has no home, so auto-share is off.
func accountHome() string {
	if home := listedHome(); home != "" {
		return home
	}
	return directoryHome()
}
