package main

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
)

func pivotRoot(newroot string) error {
	putold := filepath.Join(newroot, "/.pivot_root")
	if err := syscall.Mount(
		newroot,
		newroot,
		"",
		syscall.MS_BIND|syscall.MS_REC,
		"",
	); err != nil {
		return fmt.Errorf("Error while mounting root on itself: %s", err)
	}
	if err := os.MkdirAll(putold, 0700); err != nil {
		return fmt.Errorf("Error while creating old folder: %s", err)
	}
	if err := syscall.PivotRoot(newroot, putold); err != nil {
		return fmt.Errorf("Error while pivoting: %s", err)
	}
	if err := os.Chdir("/"); err != nil {
		return fmt.Errorf("Error while checking /: %s", err)
	}
	putold = "/.pivot_root"
	if err := syscall.Unmount(
		putold,
		syscall.MNT_DETACH,
	); err != nil {
		return fmt.Errorf("Error while unmounting: %s", err)
	}
	if err := os.RemoveAll(putold); err != nil {
		return fmt.Errorf("Error while removing old folder: %s", err)
	}

	return nil
}

func checkRootFS(rootfsPath string) {
	if _, err := os.Stat(rootfsPath); os.IsNotExist(err) {
		fmt.Printf("Error: '%s' does not exist.", rootfsPath)
		os.Exit(1)
	}
}
