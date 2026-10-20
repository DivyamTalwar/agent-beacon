package shareauto

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// getentPaths are fixed: PATH can come from a repository.
var getentPaths = []string{"/usr/bin/getent", "/bin/getent"}

// lookupAccountHome is the passwd build's accountHome: the /etc/passwd home (listed), else the
// directory-service home (directory, getent), asked for only when a settings file exists under
// $HOME. getent can make NSS query a directory server, and every hook process would otherwise pay
// that even when auto-share was never turned on. $HOME can come from a repository, but it only
// decides whether the lookup runs, never which file is read: that is under the home getent
// returns. The gate is one stat, checked on every call, so the CLI sees the hooks' answer once it
// has written the file.
func lookupAccountHome(listed, directory func() string) string {
	if home := listed(); home != "" {
		return home
	}
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		return ""
	}
	if _, err := os.Stat(filepath.Join(home, ".beacon", "endpoint", FileName)); err != nil {
		return ""
	}
	return directory()
}

// getentHome asks NSS for uid's home. getent runs with an empty environment, so LD_PRELOAD and
// the like from the hooks' environment cannot reach it.
func getentHome(uid string) string {
	for _, p := range getentPaths {
		if _, err := os.Stat(p); err != nil {
			continue
		}
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		cmd := exec.CommandContext(ctx, p, "passwd", uid)
		cmd.Env = []string{}
		out, err := cmd.Output()
		if err != nil {
			return ""
		}
		return passwdHome(string(out), uid)
	}
	return ""
}

// passwdHome is the home field of uid's entry in passwd-format data.
func passwdHome(data, uid string) string {
	for _, line := range strings.Split(data, "\n") {
		fields := strings.Split(strings.TrimSpace(line), ":")
		if len(fields) >= 6 && fields[2] == uid && !strings.HasPrefix(fields[0], "#") {
			return fields[5]
		}
	}
	return ""
}
