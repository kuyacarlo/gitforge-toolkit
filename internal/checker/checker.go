package checker

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/kuyacarlo/gitforge-toolkit/internal/config"
)

type CheckResult struct {
	Binary    string
	Installed bool
	Version   string
	Required  string
	Passed    bool
	Error     string
}

func CheckTool(req config.ToolRequirement) CheckResult {
	res := CheckResult{
		Binary:   req.Binary,
		Required: req.MinVersion,
	}

	path, err := exec.LookPath(req.Binary)
	if err != nil {
		res.Installed = false
		res.Passed = false
		res.Error = fmt.Sprintf("binary '%s' not found in PATH", req.Binary)
		return res
	}

	res.Installed = true
	res.Version = detectVersion(path, req.Binary)

	if req.MinVersion == "" {
		res.Passed = true
		return res
	}

	if res.Version == "" {
		res.Passed = false
		res.Error = "could not detect version"
		return res
	}

	if compareSemVer(res.Version, req.MinVersion) >= 0 {
		res.Passed = true
	} else {
		res.Passed = false
		res.Error = fmt.Sprintf("version %s is lower than required %s", res.Version, req.MinVersion)
	}

	return res
}

func detectVersion(path string, binary string) string {
	var cmd *exec.Cmd
	switch binary {
	case "java":
		cmd = exec.Command(path, "-version")
	case "go":
		cmd = exec.Command(path, "version")
	default:
		cmd = exec.Command(path, "--version")
	}

	out, err := cmd.CombinedOutput()
	if err != nil && len(out) == 0 {
		return ""
	}

	text := string(out)
	re := regexp.MustCompile(`(?:version\s+(?:go)?|v)?([0-9]+\.[0-9]+(?:\.[0-9]+)?)`)
	matches := re.FindStringSubmatch(text)
	if len(matches) > 1 {
		return matches[1]
	}

	return ""
}

func compareSemVer(v1, v2 string) int {
	parts1 := parseVer(v1)
	parts2 := parseVer(v2)

	for i := 0; i < 3; i++ {
		p1 := 0
		if i < len(parts1) {
			p1 = parts1[i]
		}
		p2 := 0
		if i < len(parts2) {
			p2 = parts2[i]
		}
		if p1 > p2 {
			return 1
		}
		if p1 < p2 {
			return -1
		}
	}
	return 0
}

func parseVer(v string) []int {
	v = strings.TrimPrefix(v, "v")
	var res []int
	parts := strings.Split(v, ".")
	for _, p := range parts {
		var n int
		fmt.Sscanf(p, "%d", &n)
		res = append(res, n)
	}
	return res
}

func CheckProject(cfg *config.ProjectConfig) []CheckResult {
	var results []CheckResult
	for _, tool := range cfg.Tools {
		results = append(results, CheckTool(tool))
	}
	return results
}
