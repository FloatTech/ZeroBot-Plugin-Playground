package midicreate

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"runtime"

	"github.com/sirupsen/logrus"
)

func command(arg ...string) (result string, err error) {
	name := "/bin/bash"
	c := "-c"
	if runtime.GOOS == "windows" {
		name = "cmd"
		c = "/C"
	}
	arg = append([]string{c}, arg...)
	logrus.Infoln("命令:", arg)
	cmd := exec.Command(name, arg...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", fmt.Errorf("Can not obtain stdout pipe for command, %s\n", err)
	}

	if err = cmd.Start(); err != nil {
		return "", fmt.Errorf("The command is err, %s\n", err)
	}

	bytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		return "", fmt.Errorf("ReadAll Stdout,%s\n", err)
	}

	// if err = cmd.Wait(); err != nil {
	// 	return "", fmt.Errorf("Wait,%s\n", err)
	// }

	result = string(bytes)
	return
}
