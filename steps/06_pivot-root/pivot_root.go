package main

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
)

// Pivote le dossier cible
func pivotRoot(newroot string) error {
	putold := filepath.Join(newroot, "/.pivot_root")

	// Triche pour valider les contraintes de pivot_root:
	// bind mount newroot sur lui même
	if err := syscall.Mount(
		newroot,
		newroot,
		"",
		syscall.MS_BIND|syscall.MS_REC,
		"",
	); err != nil {
		return fmt.Errorf("Error while mounting root on itself: %s", err)
	}

	// Créé l'ancien dossier
	if err := os.MkdirAll(putold, 0700); err != nil {
		return fmt.Errorf("Error while creating old folder: %s", err)
	}

	// Pivote
	if err := syscall.PivotRoot(newroot, putold); err != nil {
		return fmt.Errorf("Error while pivoting: %s", err)
	}

	//
	if err := os.Chdir("/"); err != nil {
		return fmt.Errorf("Error while checking /: %s", err)
	}

	// Unmount ancien dossier (maintenant disponible à /.pivot_root)
	putold = "/.pivot_root"
	if err := syscall.Unmount(
		putold,
		syscall.MNT_DETACH,
	); err != nil {
		return fmt.Errorf("Error while unmounting: %s", err)
	}

	// remove .pivot_root
	if err := os.RemoveAll(putold); err != nil {
		return fmt.Errorf("Error while removing old folder: %s", err)
	}

	return nil
}

// Vérifie l'existence du futur root du container
func checkRootFS(rootfsPath string) {
	if _, err := os.Stat(rootfsPath); os.IsNotExist(err) {
		fmt.Printf("Error: '%s' does not exist.", rootfsPath)
		os.Exit(1)
	}
}
