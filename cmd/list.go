package cmd

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"os/exec"
	"strings"
)

var (
	List = &cli.Command{
		Name:   "ls",
		Usage:  "list vm",
		Action: list,
	}
)

func list(*cli.Context) error {
	var cmd *exec.Cmd
	var resize []byte
	var err error
	cmd = exec.Command("virsh", "list")
	if resize, err = cmd.Output(); err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(strings.Trim(string(resize), "\n"))
	return nil
}
