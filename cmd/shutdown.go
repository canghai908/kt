package cmd

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
	"os/exec"
	"strings"
	"time"
)

var (
	Shutdown = &cli.Command{
		Name:   "sh",
		Usage:  "shutdown vm",
		Action: shutdown,
	}
)

func shutdown(*cli.Context) error {
	var cmd, cmd2, cmd3 *exec.Cmd
	var resize, size, status []byte
	var err error
	//查看vm列表
	cmd = exec.Command("virsh", "list")
	if resize, err = cmd.Output(); err != nil {
		log.Println(err)
		return err
	}
	StrList := strings.Trim(string(resize), "\n")
	list := strings.Split(StrList, "\n")
	//小于2无vm running推出
	if len(list) == 2 {
		log.Println("No VM Running")
		return err
	}
	//显示vm
	fmt.Println("VM Running Count:", len(list)-2)
	for i := 0; i < len(list); {
		fmt.Println(list[i])
		i++
	}
	//关闭vm
	fmt.Println("Shutdown VM")
	p := 1
	for i := 2; i < len(list); {
		VmName := strings.Split(list[i], " ")
		log.Println("--------------------------------")
		log.Println("Action", p, ",Shutdown VM ", strings.Trim(VmName[3], " "))
		cmd2 = exec.Command("virsh", "shutdown", strings.Trim(VmName[3], " "))
		if size, err = cmd2.Output(); err != nil {
			log.Println(err)
			continue
		}
		log.Println(strings.Trim(string(size), "\n"))
		time.Sleep(1 * time.Second)
		p++
		i++
	}
	cmd3 = exec.Command("virsh", "list")
	if status, err = cmd3.Output(); err != nil {
		log.Println(err)
		return err
	}
	log.Println(strings.Trim(string(status), "\n"))
	return nil
}
