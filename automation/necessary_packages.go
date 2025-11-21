package automate

import (
	"fmt"
	"os"
	"CipherOps/utils"
)

func SetupNecessaryPkgs() {
	fmt.Println("------------------------------------------------------------------------------")
	utils.DryRun = true //testing flag

	if err := utils.InstallPackages([]string{"docker", "postgres"}); err != nil {
		fmt.Fprintf(os.Stderr, "Failed installing: %v\n", err)
		return
	}

	if err := utils.Service("enable", "docker"); err != nil {
		fmt.Fprintf(os.Stderr, "Cannot enable docker: %v\n", err)
		return
	}
	if err := utils.Service("restart", "docker"); err != nil {
		fmt.Fprintf(os.Stderr, "Cannot restart docker: %v\n", err)
		return
	}

	if err := utils.Service("enable", "postgresql"); err != nil {
		fmt.Println("Cannot enable postgres:", err)
	}
	fmt.Println("------------------------------------------------------------------------------")
}