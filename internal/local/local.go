package local

import (
	"os"
	"path/filepath"
	"runtime"
)

type Info struct {
	Found   bool
	Bundled bool
	Exe     string
	DataDir string
	Hint    string
}

func exeDir() string {
	exe, err := os.Executable()
	if err != nil {
		return ""
	}
	return filepath.Dir(exe)
}

func Detect() Info {
	var candidates []string
	dataDir := ""
	switch runtime.GOOS {
	case "windows":
		candidates = []string{
			filepath.Join(exeDir(), "boinc.exe"),
			`C:\Program Files\BOINC\boinc.exe`,
		}
		dataDir = os.Getenv("ProgramData") + `\BOINC`
	case "darwin":
		candidates = []string{
			filepath.Join(exeDir(), "boinc"),
			"/Applications/BOINCManager.app/Contents/Resources/boinc",
		}
		home, _ := os.UserHomeDir()
		dataDir = home + "/Library/Application Support/BOINC"
	default:
		candidates = []string{
			filepath.Join(exeDir(), "boinc_client"),
			filepath.Join(exeDir(), "boinc"),
			"/usr/bin/boinc_client",
			"/usr/bin/boinc",
			"/usr/local/bin/boinc_client",
		}
		home, _ := os.UserHomeDir()
		dataDir = home + "/.local/share/BOINC"
		if runtime.GOOS == "freebsd" {
			candidates = append(candidates, "/usr/local/bin/boinc_client")
		}
	}
	for i, p := range candidates {
		if st, err := os.Stat(p); err == nil && !st.IsDir() {
			return Info{
				Found:   true,
				Bundled: i <= 1 && runtime.GOOS != "linux" || filepath.Dir(p) == exeDir(),
				Exe:     p,
				DataDir: dataDir,
			}
		}
	}
	return Info{
		Hint: "Install the BOINC client, or connect to remote hosts from Servers.",
		DataDir: dataDir,
	}
}
