// 06 Pivot root
package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	switch os.Args[1] {
	case "run":
		run()
	case "child":
		child()
	default:
		panic(fmt.Sprintf("Unknow command %s", os.Args[1]))
	}
}

func run() {

	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = []string{"PS1=marcel:$(pwd) #", "PATH=/bin:/usr/bin/", "TERM=xterm"}
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS |
			syscall.CLONE_NEWNS |
			syscall.CLONE_NEWIPC |
			syscall.CLONE_NEWPID |
			syscall.CLONE_NEWNET |
			syscall.CLONE_NEWUSER,
		UidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      os.Getuid(),
				Size:        1,
			},
		},
		GidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      os.Getgid(),
				Size:        1,
			},
		},
	}

	if err := cmd.Run(); err != nil {
		panic(err)
	}
}

func child() {
	fmt.Printf("Namespace initialization\n")

	rootfsPath := "/home/stoakes/go/src/recreez-votre-docker/steps/06_pivot-root/centos"
	checkRootFS(rootfsPath)

	// Monte /proc
	if err := mountProc(rootfsPath); err != nil {
		fmt.Printf("Error mounting /proc - %s\n", err)
		os.Exit(1)
	}

	// Monte /dev/urandom depuis l'hôte
	if err := bindMountDeviceNode("/dev/urandom", rootfsPath+"/dev/urandom"); err != nil {
		fmt.Printf("Error running bind mount urandom: %s", err)
	}

	// Pivote
	if err := pivotRoot(rootfsPath); err != nil {
		fmt.Printf("Error running pivot_root - %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Running %v as PID %d\n", os.Args[2:], os.Getpid())

	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		panic(err)
	}

	// Une fois la commande exécutée: unmount
	syscall.Unmount("proc", 0)
	syscall.Unmount("dev/urandom", 0)
}
