package shareauto

import (
	"context"
	"os"
	"os/exec"
	"strings"
	"time"
)

// getentPaths are fixed: PATH can come from a repository.
var getentPaths = []string{"/usr/bin/getent", "/bin/getent"}

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
