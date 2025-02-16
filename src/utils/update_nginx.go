package utils

import (
	"log"
	"os/exec"
)

func ReloadNginxConfig() error {
	cmd := exec.Command("sh", "/app/update_nginx.sh")
	err := cmd.Run()
	if err != nil {
		log.Printf("Error recargando Nginx: %v", err)
		return err
	}
	return nil
}
