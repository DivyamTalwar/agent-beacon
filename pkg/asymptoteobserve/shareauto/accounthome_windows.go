//go:build windows

package shareauto

import "syscall"

// accountHome reads only the profile directory from the process token. user.Current also asks
// for the display name, which on a domain-joined machine can query a domain controller.
func accountHome() string {
	t, err := syscall.OpenCurrentProcessToken()
	if err != nil {
		return ""
	}
	defer t.Close()
	dir, err := t.GetUserProfileDirectory()
	if err != nil {
		return ""
	}
	return dir
}
