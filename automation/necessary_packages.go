package automate

import (
	"fmt"
	"os"
	"CipherOps/utils"
)

func SetupNecessaryPkgs() {
	fmt.Println("------------------------------------------------------------------------------")
	//utils.DryRun = true //testing flag

	if err := utils.InstallPackages([]string{"docker", "postgres"}); err != nil {
		fmt.Fprintf(os.Stderr, "Error instalando: %v\n", err)
		return
	}

	if err := utils.Service("enable", "docker"); err != nil {
		fmt.Fprintf(os.Stderr, "no se pudo habilitar docker: %v\n", err)
		return
	}
	if err := utils.Service("restart", "docker"); err != nil {
		fmt.Fprintf(os.Stderr, "no se pudo reiniciar docker: %v\n", err)
		return
	}

	if err := utils.Service("enable", "postgresql"); err != nil {
		fmt.Println("Advertencia: no se pudo habilitar postgresql directamente:", err)
	}
	fmt.Println("------------------------------------------------------------------------------")
}