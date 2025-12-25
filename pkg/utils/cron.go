package utils

import (
	"bytes"
	"fmt"
	"os/exec"
)

func GetCronTab() (string, error) {
	cmd := exec.Command("crontab", "-l")
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			return "", nil
		}
		return "", fmt.Errorf("failed to read crontab: %w, stderr: %s", err, stderr.String())
	}

	return out.String(), nil
}

func WriteCrontab(content string) error {
	cmd := exec.Command("crontab", "-")
	cmd.Stdin = bytes.NewBufferString(content)
	return cmd.Run()
}

func AddCronJob(name, schedule, command string) error {
	curr, err := GetCronTab()
	if err != nil {
		return err
	}

	entry := fmt.Sprintf(
		"# backup-tool:%s\n%s %s\n",
		name,
		schedule,
		command,
	)

	return WriteCrontab(curr + entry)
}

func RemoveCronJob(name string) error {
	curr, err := GetCronTab()
	if err != nil {
		return err
	}

	lines := bytes.Split([]byte(curr), []byte("\n"))
	var result []byte

	skip := false
	for i, line := range lines {
		if bytes.Contains(line, []byte("# backup-tool:"+name)) {
			skip = true
			continue
		}
		if skip {
			skip = false
			continue
		}
		result = append(result, line...)
		if i != len(lines)-1 {
			result = append(result, '\n')
		}
	}

	return WriteCrontab(string(result))
}
