package git

import (
	"os/exec"

	"github.com/Salamandra2402/git_pswitcher/profile"
)

//SwitchToProfile switch to chosen git profile
func SwitchToProfile(p profile.Profile) error {
	var err error
	if _, err = exec.Command("git", "config", "--global", "user.name", p.Name).Output(); err != nil {
		return err
	}

	if _, err = exec.Command("git", "config", "--global", "user.email", p.Email).Output(); err != nil {
		return err
	}
	return nil
}
