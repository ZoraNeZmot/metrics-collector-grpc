package main

import (
	"log"
	"os"
	"os/exec"
	"time"
)

func main() {
	if os.Getenv("IS_DAEMON") != "1" {
		cmd := exec.Command(os.Args[0], os.Args[:1]...)

		cmd.Env = append(os.Environ(), "IS_DAEMON=1")

		cmd.Stdin = nil

		f, err := os.OpenFile("daemon.log", os.O_CREATE|os.O_TRUNC, os.ModePerm)
		if err != nil {
			log.Fatal("Failed to open log file for daemon")
		}
		cmd.Stderr = f
		cmd.Stdout = f

		if err := cmd.Start(); err != nil {
			log.Fatalln("Failed to start a daemon.")
		}
		os.Exit(0)
	}

	log.Printf("DAEMON PID: %v\n", os.Getpid())
	for i := 0; ; i++ {
		time.Sleep(2 * time.Second)
		log.Printf("DO SOME WORK %d", i)
	}
}
