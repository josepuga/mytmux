package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/josepuga/goini"
)

var VERSION string = "..."
var tmuxPath string // Safe use with `which`
const usage = `MyTmux. %s
Automates the creation of sessions and windows in Tmux
(c)José Puga 2025. Under GPL 3 License.
Usage: %s <tmux workspaces ini>
`

func main() {
	code := realMain()
	os.Exit(code)
}

func realMain() int {
	args := os.Args
	if len(args) == 1 {
		fmt.Printf(usage, VERSION, filepath.Base(args[0]))
		return 0
	}
	workspacesFile := args[1]

	// Check if tmux is in the system
	var err error
	tmuxPath, err = exec.LookPath("tmux")
	if err != nil {
		fmt.Fprint(os.Stderr, "Error, tmux is not installed")
		return 1
	}

	ini := goini.NewIni()
	err = ini.LoadFromFile(workspacesFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening workspace file %s\n", workspacesFile)
		return 1
	}

	tmuxSessionList := []TmuxSession{}
	for _, section := range ini.GetSectionValues() {
		if section == "" { //Section "" is always available
			continue
		}
		tms := TmuxSession{}

		// [section] of the INI is the name of the session
		tms.Name = section

		// Get all windows
		// INI key=value;[command]. key -> Title, value -> Default path
		for _, k := range ini.GetSectionKeys(section) {
			tmw := TmuxWindow{}
			if k == "" {
				fmt.Fprintf(os.Stderr, "Error, empty title at [%s]\n", tms.Name)
				return 1
			}
			tmw.Title = k
			if after, ok := strings.CutPrefix(tmw.Title, "*"); ok {
				tmw.Title = after
				tms.DefaultWindow = after
			}
			v := ini.GetStringSlice(section, k, "", ";")
			if len(v) > 1 {
				tmw.Command = v[1]
			}
			if v[0] == "" {
				tmw.Path = "~"
			} else {
				tmw.Path = v[0]
			}
			tms.Windows = append(tms.Windows, tmw)
		}
		tmuxSessionList = append(tmuxSessionList, tms)
	}

	// Now we have the []TmuxSession filled. Time to exec tmux command.
	for _, tms := range tmuxSessionList {
		if tmuxSessionExists(tms.Name) {
			fmt.Printf("Session %s already exists. To kill it, use `tmux kill-session -t %s`\n", tms.Name, tms.Name)
			continue
		}
		err = tmuxCreateSession(tms)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating %s session: %s.\n", tms.Name, err)
		} else {
			fmt.Printf("Session %s created. Use `tmux attach-session -t %s`\n", tms.Name, tms.Name)
		}
	}

	return 0
}

func tmuxCreateSession(tms TmuxSession) error {
	// Session
	cmd := exec.Command(tmuxPath, "new-session", "-d",
		"-s", tms.Name)
	err := cmd.Run()
	if err != nil {
		return err
	}

	// Windows
	for _, tmw := range tms.Windows {
		// Create and define the window
		path := expandTilde(tmw.Path)
		cmd := exec.Command(tmuxPath, "new-window",
			"-t", tms.Name,
			"-n", tmw.Title,
			"-c", path)
		err := cmd.Run()
		if err != nil {
			return err
		}
		// Run shell commmand
		if tmw.Command != "" {
			sw := fmt.Sprintf("%s:%s", tms.Name, tmw.Title)
			cmd := exec.Command(tmuxPath, "send-keys",
				"-t", sw,
				tmw.Command, "C-m")
			err := cmd.Run()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func tmuxSessionExists(sessionName string) bool {
	cmd := exec.Command(tmuxPath, "has-session", "-t", sessionName)
	err := cmd.Run()
	return err == nil
}

func expandTilde(path string) string {
	if strings.HasPrefix(path, "~") {
		home, err := os.UserHomeDir()
		//WARNING: No error check!
		if err == nil && home != "" {
			return filepath.Join(home, strings.TrimPrefix(path, "~"))
		}
	}
	return path
}
