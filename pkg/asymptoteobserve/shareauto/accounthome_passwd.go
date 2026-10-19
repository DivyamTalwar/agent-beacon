//go:build (!cgo && !darwin && !windows && !plan9) || android || (osusergo && !windows && !plan9)

package shareauto

import (
	"os"
	"strconv"
	"sync"
)

var accountHomeOnce = sync.OnceValue(func() string {
	uid := strconv.Itoa(os.Getuid())
	if data, err := os.ReadFile("/etc/passwd"); err == nil {
		if home := passwdHome(string(data), uid); home != "" {
			return home
		}
	}
	return getentHome(uid)
})

// accountHome reads the user database itself: in this build os/user's Current and LookupId fall
// back to $HOME when the uid is not in /etc/passwd. getent reaches directory-service accounts
// through NSS; without it an unlisted uid has no home, so auto-share is off.
func accountHome() string { return accountHomeOnce() }
