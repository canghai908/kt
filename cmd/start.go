package cmd

import (
	"bufio"
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

var (
	Start = &cli.Command{
		Name:   "st",
		Usage:  "start vm",
		Action: start,
	}
)

func start(*cli.Context) error {
	filepath := "./kvm.list"
	file, err := os.Open(filepath)
	if err != nil {
		log.Println(err)
		return err
	}
	fileScanner := bufio.NewScanner(file)
	i := 0
	for fileScanner.Scan() {
		line := strings.TrimSpace(fileScanner.Text())
		newline := strings.Replace(line, "xml", "qcow2", -1)
		b, err1 := PathExists(newline)
		if b && err1 == nil {
			VmT := strings.Split(newline, "/")
			VmID := strings.Split(VmT[5], ".")
			var cmd1, cmd2, cmd3, cmd4, cmd5 *exec.Cmd
			var imginfo, define, status, newstatus, tlist []byte
			log.Println("--------------------------------")
			//检查vm是否启动
			cmd4 = exec.Command("virsh", "list")
			if newstatus, err = cmd4.Output(); err != nil {
				log.Println(err)
				continue
			}
			if strings.Contains(strings.Trim(string(newstatus), "\n"), VmID[0]) {
				log.Println("VM", VmID[0], "is running")
				continue
			}
			//检查虚拟机磁盘文件大小
			cmd1 = exec.Command("qemu-img", "info", newline)
			if imginfo, err = cmd1.Output(); err != nil {
				fmt.Println(err)
			}
			info := strings.Split(strings.Trim(string(imginfo), "\n"), "\n")
			disksize := strings.Split(info[2], ":")
			if !strings.Contains(disksize[1], "50G") {
				log.Println("VM磁盘文件大小错误", info[2])
				continue
			}

			log.Println("VM ID", VmID[0])
			log.Println("VM Disk file", newline)
			log.Println("VM Disk Size", strings.Trim(disksize[1], " "))
			//Define VM
			log.Println("Define VM", VmID[0])
			cmd2 = exec.Command("virsh", "define", line)
			if define, err = cmd2.Output(); err != nil {
				log.Println(err)
				break
			}
			log.Println(strings.Trim(string(define), "\n"))
			time.Sleep(1 * time.Second)
			//Start VM
			log.Println("Start VM", VmID[0])
			cmd3 = exec.Command("virsh", "start", VmID[0])
			if status, err = cmd3.Output(); err != nil {
				log.Println(err)
				continue
			}
			log.Println(strings.Trim(string(status), "\n"))
			time.Sleep(1 * time.Second)
			//virsh list
			cmd5 = exec.Command("virsh", "list")
			if tlist, err = cmd5.Output(); err != nil {
				log.Println(err)
				continue
			}
			log.Println(strings.Trim(string(tlist), "\n"))
		} else {
			log.Println("VM xml not found")
		}
	}
	// handle first encountered error while reading
	if err := fileScanner.Err(); err != nil {
		log.Println(err)
		return err
	}

	defer file.Close()
	log.Println("VM Count:", i)
	var cmd6 *exec.Cmd
	var plist []byte
	cmd6 = exec.Command("virsh", "list")
	if plist, err = cmd6.Output(); err != nil {
		log.Println(err)
		return err
	}
	fmt.Println(strings.Trim(string(plist), "\n"))
	return nil
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
