package midicreate

import (
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
		return
	}

	if err = cmd.Start(); err != nil {
		return
	}

	bytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		return
	}

	if err = cmd.Wait(); err != nil {
		return
	}

	result = string(bytes)
	return
}
