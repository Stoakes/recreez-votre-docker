// 01 un simple exécuteur de tâche
package main

import (
	"fmt"
	"os"
	"os/exec"
)

// go run main.go run cmd args
func main() {
	switch os.Args[1] {
	case "run":
		run()
	default:
		panic("Unknow command")
	}
}

// Execute la commande passée en paramètres
func run() {
	fmt.Printf("Running %v\n", os.Args[2:])

	// Initialise la commande
	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	// Lie les entrées et sorties standard à celles de l'OS
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	// Variable d'environement pour différencier l'intérieur de l'extérieur
	cmd.Env = []string{"PS1=hello :", "HELLO=hello", "TERM=xterm"}

	// Lancement de la commande
	if err := cmd.Run(); err != nil {
		panic(err)
	}
}
