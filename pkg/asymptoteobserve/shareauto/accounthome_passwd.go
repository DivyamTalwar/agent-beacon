//go:build (!cgo && !darwin && !windows && !plan9) || android || (osusergo && !windows && !plan9)

package shareauto

import (
	"os"
	"strconv"
	"strings"
)

// accountHome reads /etc/passwd itself: in this build os/user's Current and LookupId fall back
// to $HOME when the uid is not listed there. An unlisted uid has no home, so auto-share is off.
func accountHome() string {
	data, err := os.ReadFile("/etc/passwd")
	if err != nil {
		return ""
	}
	uid := strconv.Itoa(os.Getuid())
	for _, line := range strings.Split(string(data), "\n") {
		fields := strings.Split(strings.TrimSpace(line), ":")
		if len(fields) >= 6 && fields[2] == uid && !strings.HasPrefix(fields[0], "#") {
			return fields[5]
		}
	}
	return ""
}
