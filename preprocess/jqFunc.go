package preprocess

import (
	"errors"
	"os"
	"os/exec"
	"strings"
)

func prepareJQ(jqDirs ...string) (jqWD, oriWD string, err error) {
	fn := "prepareJQ"
	jq := "jq.exe" // windows
	// jq := "jq" // linux

	oriWD, err = os.Getwd()
	FailOnErr("Getwd() 1 fatal @ %v: %w", fn, err)
	for _, jqWD = range jqDirs {
		if _, err = os.Stat(jqWD + jq); err == nil {
			FailOnErr("Chdir() fatal @ %v: %w", fn, os.Chdir(jqWD))
			jqWD, err = os.Getwd()
			FailOnErr("Getwd() 2 fatal @ %v: %w", fn, err)
			return jqWD, oriWD, nil
		}
	}
	FailOnErr("%s is not found @ %v", jq, errors.New(fn))
	return "", "", nil
}

// FmtJSONStr : <json string> must not have single quote <'>
func FmtJSONStr(json string, jqDirs ...string) string {
	_, oriWD, _ := prepareJQ(jqDirs...)
	defer func() { os.Chdir(oriWD) }()

	json = "'" + strings.ReplaceAll(json, "'", "\\'") + "'" // *** deal with <single quote> in "echo" ***

	// cmdstr := "echo " + json + ` | ./jq .`    // linux
	// cmd := exec.Command("bash", "-c", cmdstr) // linux

	cmdstr := "echo " + json + ` | ./jq.exe .`            // windows
	cmd := exec.Command("PowerShell", "-Command", cmdstr) // windows

	if output, err := cmd.Output(); err == nil {
		return string(output)
	}
	FailOnErr("cmd.Output() error @ %v", errors.New("FmtJSONStr"))
	return ""
}

// FmtJSONFile : <file> is the <relative path> to <jq>
func FmtJSONFile(file string, jqDirs ...string) string {
	_, oriWD, _ := prepareJQ(jqDirs...)
	defer func() { os.Chdir(oriWD) }()

	// cmdstr := "cat " + file + ` | ./jq .`     // linux
	// cmd := exec.Command("bash", "-c", cmdstr) // linux

	cmdstr := "cat " + file + ` | ./jq.exe .`             // windows
	cmd := exec.Command("PowerShell", "-Command", cmdstr) // windows

	if output, err := cmd.Output(); err == nil {
		return string(output)
	}
	FailOnErr("cmd.Output() error @ %v", errors.New("FmtJSONFile"))
	return ""
}
