// 07 CGroup
package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
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
	cmd.Env = []string{"PS1=hello:$(pwd) #", "PATH=/bin:/usr/bin/", "TERM=xterm"}
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

	// 100Mo de mémoire pour le CGroup demo (créé dans bash pour garder les droits)
	memoryCGroup := "/sys/fs/cgroup/memory/"
	if err := ioutil.WriteFile(filepath.Join(memoryCGroup, "demo/memory.limit_in_bytes"), []byte("100000000"), 0700); err != nil {
		fmt.Printf("Error while creating memory limit: %s\n", err)
	}
	if err := ioutil.WriteFile(filepath.Join(memoryCGroup, "demo/memory.swappiness"), []byte("0"), 0700); err != nil {
		fmt.Printf("Error while disabling memory swapiness: %s\n", err)
	}
	// Ajoute le PID courant dans le control group demo	if err := ioutil.WriteFile(filepath.Join(memoryCGroup, "demo/cgroup.procs"), []byte(strconv.Itoa(os.Getpid())), 0700); err != nil {
		fmt.Printf("Error while adding pid (%d) to demo Cgroup: %s\n", os.Getpid(), err)
	}
	// Expérimental: Supprime le nouveau CGroup après la suppression du container
	if err := ioutil.WriteFile(filepath.Join(memoryCGroup, "demo/notify_on_release"), []byte("1"), 0700); err != nil {
		fmt.Printf("Error while adding notify on release: %s\n", err)
	}

	// Mount
	rootfsPath := "/home/stoakes/go/src/recreez-votre-docker/steps/07_cgroups/centos"
	checkRootFS(rootfsPath)
	if err := mountProc(rootfsPath); err != nil {
		fmt.Printf("Error mounting /proc - %s\n", err)
		os.Exit(1)
	}
	if err := bindMountDeviceNode("/dev/urandom", rootfsPath+"/dev/urandom"); err != nil {
		fmt.Printf("Error running bind mount urandom: %s", err)
	}
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

	syscall.Unmount("proc", 0)
	syscall.Unmount("dev/urandom", 0)
}
