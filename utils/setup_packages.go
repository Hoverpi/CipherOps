package utils

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

// ----- Distro detection -----
type DistroInfo struct {
	ID     string
	Like   string
	Name   string
	Version string
}

func DetectDistro() (DistroInfo, error) {
	f, err := os.Open("/etc/os-release")
	if err != nil {
		return DistroInfo{}, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	info := DistroInfo{}
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "ID=") {
			info.ID = strings.Trim(strings.SplitN(line, "=", 2)[1], "\"")
		}
		if strings.HasPrefix(line, "ID_LIKE=") {
			info.Like = strings.Trim(strings.SplitN(line, "=", 2)[1], "\"")
		}
		if strings.HasPrefix(line, "PRETTY_NAME=") {
			info.Name = strings.Trim(strings.SplitN(line, "=", 2)[1], "\"")
		}
		if strings.HasPrefix(line, "VERSION_ID=") {
			info.Version = strings.Trim(strings.SplitN(line, "=", 2)[1], "\"")
		}
	}
	if info.ID == "" {
		return info, errors.New("The distribution could not be detected (no /etc/os-release).")
	}
	return info, nil
}

// ----- Installer interface and some implemetations -----
type Installer interface {
	Install(packages []string) error
	Remove(packages []string) error
	IsInstalled(pkg string) (bool, error)
}

func NewInstallerFor(d DistroInfo) Installer {
	id := strings.ToLower(d.ID)
	like := strings.ToLower(d.Like)
	// simple heuristic
	if id == "ubuntu" || id == "debian" || strings.Contains(like, "debian") {
		return &AptInstaller{}
	}
	if id == "fedora" || id == "centos" || id == "rhel" || strings.Contains(like, "rhel") {
		return &DnfInstaller{}
	}
	if id == "arch" || id == "manjaro" || strings.Contains(like, "arch") {
		return &PacmanInstaller{}
	}
	if id == "opensuse" || strings.Contains(like, "suse") {
		return &ZypperInstaller{}
	}
	// fallback: try sh and let commands fall back
	return &GenericShellInstaller{}
}

// implement the methods defined on the installer interface
// ----- APT -----
type AptInstaller struct {}

func (a *AptInstaller) Install(pkgs []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()

	var output string
	var err error
	
	// Update
	if output, err = runCmd(ctx, true, "apt-get", "update", "-y"); err != nil {
		fmt.Printf("Warning: apt update failed: %v\nOutput: %s\n", err, output)
	}
	fmt.Println(output)

	args := append([]string{"install", "-y"}, pkgs...)
	if output, err = runCmd(ctx, true, "apt-get", args...); err != nil {
		return fmt.Errorf("failed to install packages %v: %w\nOutput: %s", pkgs, err, output)
	}
	fmt.Printf("apt-get install output: %s\n", output)
	return nil
}

func (a *AptInstaller) Remove(pkgs []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()

	var output string
	var err error

	args := append([]string{"remove", "-y"}, pkgs...)
	if output, err = runCmd(ctx, true, "apt-get", args...); err != nil {
		return fmt.Errorf("failed to remove packages: %v\nOutput: %s", pkgs, output)
	}
	fmt.Println(output)
	return nil
}

func (a *AptInstaller) IsInstalled(pkg string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if _, err := runCmd(ctx, false, "dpkg", "-s", pkg); err != nil {
		// dpkg -s returns an error if it is not installed
		return false, nil
	}
	return true, nil
}

// ----- DNF / YUM -----
type DnfInstaller struct {}

func (d *DnfInstaller) Install(pkgs []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()

	var output string
	var err error
	
	// Update
	if output, err = runCmd(ctx, true, "dnf", "update", "-y"); err != nil {
		fmt.Printf("Warning: dnf update failed: %v\nOutput: %s\n", err, output)
	}

	args := append([]string{"install", "-y"}, pkgs...)
	if output, err = runCmd(ctx, true, "dnf", args...); err != nil {
		if output, err = runCmd(ctx, true, "yum", args...); err != nil {
			return fmt.Errorf("failed to install packages: %v\nOutput: %s", pkgs, output)
		}
	}
	fmt.Printf("dnf/yum install output: %s\n", output)
	return nil
}

func (d *DnfInstaller) Remove(pkgs []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()

	var output string
	var err error

	args := append([]string{"remove", "-y"}, pkgs...)
	if output, err = runCmd(ctx, true, "dnf", args...); err != nil {
		if output, err = runCmd(ctx, true, "yum", args...); err != nil {
			return fmt.Errorf("failed to remove packages: %v\nOutput: %s", pkgs, output)
		}
	}
	fmt.Println(output)
	return nil
}

func (d *DnfInstaller) IsInstalled(pkg string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if _, err := runCmd(ctx, false, "rpm", "-q", pkg); err != nil {
		// rpm returns an error if the package is not installed
		return false, nil
	}
	return true, nil
}

// ----- PACMAN -----
type PacmanInstaller struct {}

func (p *PacmanInstaller) Install(pkgs []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()

	var output string
	var err error
	
		// Update
	if output, err = runCmd(ctx, true, "pacman", "-Sy", "--noconfirm"); err != nil {
		fmt.Printf("Warning: pacman update failed: %v\nOutput: %s\n", err, output)
	}

	args := append([]string{"-S", "--noconfirm"}, pkgs...)
	if output, err = runCmd(ctx, true, "pacman", args...); err != nil {
		return fmt.Errorf("failed to install packages: %v\nOutput: %s", pkgs, output)
	}
	fmt.Printf("pacman install output: %s\n", output)
	return nil
}

func (p *PacmanInstaller) Remove(pkgs []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()

	var output string
	var err error

	args := append([]string{"-R", "--noconfirm"}, pkgs...)
	if output, err = runCmd(ctx, true, "pacman", args...); err != nil {
		return fmt.Errorf("failed to remove packages: %v\nOutput: %s", pkgs, output)
	}
	fmt.Println(output)
	return nil
}

func (p *PacmanInstaller) IsInstalled(pkg string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if _, err := runCmd(ctx, false, "pacman", "-Qi", pkg); err != nil {
		// pacman -Qi returns an error if the package is not installed
		return false, nil
	}
	return true, nil
}

// ----- ZYPPER -----
type ZypperInstaller struct{}

func (z *ZypperInstaller) Install(pkgs []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()

	var output string
	var err error
	
	// Update
	if output, err = runCmd(ctx, true, "zypper", "update", "--non-interactive"); err != nil {
		fmt.Printf("Warning: zypper update failed: %v\nOutput: %s\n", err, output)
	}

	args := append([]string{"--non-interactive", "install"}, pkgs...)
	if output, err = runCmd(ctx, true, "zypper", args...); err != nil {
		return fmt.Errorf("failed to install packages: %v\nOutput: %s", pkgs, output)
	}
	fmt.Printf("zypper install output: %s\n", output)
	return nil
}

func (z *ZypperInstaller) Remove(pkgs []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()

	var output string
	var err error

	args := append([]string{"--non-interactive", "remove"}, pkgs...)
	if output, err = runCmd(ctx, true, "zypper", args...); err != nil {
		return fmt.Errorf("failed to remove packages: %v\nOutput: %s", pkgs, output)
	}
	fmt.Println(output)
	return nil
}

func (z *ZypperInstaller) IsInstalled(pkg string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if _, err := runCmd(ctx, false, "rpm", "-q", pkg); err != nil {
		// rpm returns an error if the package is not installed
		return false, nil
	}
	return true, nil
}


// --- Generic fallback ---
type GenericShellInstaller struct{}

func (g *GenericShellInstaller) Install(pkgs []string) error {
	// Try apt, dns, and pacman in order to maximize compatibility.
	if _, err := exec.LookPath("apt-get"); err == nil {
		return (&AptInstaller{}).Install(pkgs)
	}
	if _, err := exec.LookPath("dnf"); err == nil {
		return (&DnfInstaller{}).Install(pkgs)
	}
	if _, err := exec.LookPath("pacman"); err == nil {
		return (&PacmanInstaller{}).Install(pkgs)
	}
	return errors.New("No known package manager available")
}
func (g *GenericShellInstaller) Remove(pkgs []string) error {
	// Try apt, dns, and pacman in order to maximize compatibility.
	if _, err := exec.LookPath("apt-get"); err == nil {
		return (&AptInstaller{}).Remove(pkgs)
	}
	if _, err := exec.LookPath("dnf"); err == nil {
		return (&DnfInstaller{}).Remove(pkgs)
	}
	if _, err := exec.LookPath("pacman"); err == nil {
		return (&PacmanInstaller{}).Remove(pkgs)
	}
	return errors.New("No known package manager available")
}
func (g *GenericShellInstaller) IsInstalled(pkg string) (bool, error) {
	// Try apt, dns, and pacman in order to maximize compatibility.
	if _, err := exec.LookPath("dpkg"); err == nil {
		return (&AptInstaller{}).IsInstalled(pkg)
	}
	if _, err := exec.LookPath("rpm"); err == nil {
		return (&DnfInstaller{}).IsInstalled(pkg)
	}
	if _, err := exec.LookPath("pacman"); err == nil {
		return (&PacmanInstaller{}).IsInstalled(pkg)
	}
	return false, errors.New("Cannot determine if the package is installed")
}

// ----- Service Manager (systemd, sysv) -----
type ServiceManager interface {
	Start(name string) error
	Restart(name string) error
	Stop(name string) error
	Enable(name string) error
	Disable(name string) error
	Status(name string) (string, error)
}

func NewServiceManager() ServiceManager {
	// If exists systemctl -> systemd
	if _, err := exec.LookPath("systemctl"); err == nil {
		return &SystemdService{}
	}
	// fallback a service
	return &SysVService{}
}

type SystemdService struct{}

func (s *SystemdService) Start(name string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	_, err := runCmd(ctx, true, "systemctl", "start", name)
	return err
}
func (s *SystemdService) Restart(name string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	_, err := runCmd(ctx, true, "systemctl", "restart", name)
	return err
}
func (s *SystemdService) Stop(name string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	_, err := runCmd(ctx, true, "systemctl", "stop", name)
	return err
}
func (s *SystemdService) Enable(name string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	_, err := runCmd(ctx, true, "systemctl", "enable", "--now", name)
	return err
}
func (s *SystemdService) Disable(name string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	_, err := runCmd(ctx, true, "systemctl", "disable", "--now", name)
	return err
}
func (s *SystemdService) Status(name string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	out, err := runCmd(ctx, false, "systemctl", "status", name, "--no-pager")
	return out, err
}

type SysVService struct{}

func (s *SysVService) Start(name string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	_, err := runCmd(ctx, true, "service", name, "start")
	return err
}
func (s *SysVService) Restart(name string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	_, err := runCmd(ctx, true, "service", name, "restart")
	return err
}
func (s *SysVService) Stop(name string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	_, err := runCmd(ctx, true, "service", name, "stop")
	return err
}
func (s *SysVService) Enable(name string) error {
	// Not all sysv have enable; we will try chkconfig or update-rc.d
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	if _, err := exec.LookPath("chkconfig"); err == nil {
		_, err2 := runCmd(ctx, true, "chkconfig", name, "on")
		return err2
	}
	if _, err := exec.LookPath("update-rc.d"); err == nil {
		_, err2 := runCmd(ctx, true, "update-rc.d", name, "defaults")
		return err2
	}
	return errors.New("habilitar servicio no soportado en este sistema")
}
func (s *SysVService) Disable(name string) error {
	// Not all sysv have enable; we will try chkconfig or update-rc.d
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	if _, err := exec.LookPath("chkconfig"); err == nil {
		_, err2 := runCmd(ctx, true, "chkconfig", name, "off")
		return err2
	}
	if _, err := exec.LookPath("update-rc.d"); err == nil {
		_, err2 := runCmd(ctx, true, "update-rc.d", name, "defaults")
		return err2
	}
	return errors.New("habilitar servicio no soportado en este sistema")
}
func (s *SysVService) Status(name string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return runCmd(ctx, false, "service", name, "status")
}

// ----- Package name mapping -----
var PackageMap = map[string]map[string][]string{
	"docker": {
		"debian": {"docker.io"},
		"ubuntu": {"docker.io"},
		"fedora": {"docker"},
		"rhel":   {"docker"},
		"centos": {"docker"},
		"arch":   {"docker"},
		"opensuse": {"docker"},
	},
	"postgres": {
		"debian": {"postgresql"},
		"ubuntu": {"postgresql"},
		"fedora": {"postgresql-server"},
		"rhel":   {"postgresql-server"},
		"centos": {"postgresql-server"},
		"arch":   {"postgresql"},
		"opensuse": {"postgresql-server"},
	},
}

// --- High-level helpers that your code can call ---
// InstallPackages receives a list of “abstract names” (docker, postgres, etc.)
func InstallPackages(programs []string) error {
	distro, err := DetectDistro()
	if err != nil {
		return err
	}
	installer := NewInstallerFor(distro)
	for _, prog := range programs {
		pmap, ok := PackageMap[prog]
		if !ok {
			return fmt.Errorf("There is no packet mapping for '%s'", prog)
		}
		// try exact ID first, then ID_LIKE heuristic, then fallback to “any”
		var toInstall []string
		if v, exists := pmap[strings.ToLower(distro.ID)]; exists {
			toInstall = v
		} else if v, exists := pmap[strings.ToLower(distro.Like)]; exists {
			toInstall = v
		} else {
			// TODO
			// if there is a “default” key, you could use it; for now, error if it does not match
			return fmt.Errorf("There is no known package for %s in %s", prog, distro.ID)
		}
		// comprobar si ya está instalado
		filtered := make([]string, 0, len(toInstall))
		for _, pkg := range toInstall {
			ok, _ := installer.IsInstalled(pkg)
			if ok {
				continue
			}
			filtered = append(filtered, pkg)
		}
		if len(filtered) == 0 {
			continue
		}
		if err := installer.Install(filtered); err != nil {
			return fmt.Errorf("Error installing %s: %w", prog, err)
		}
	}
	return nil
}

// Service helper simple:
func Service(action, serviceName string) error {
	sm := NewServiceManager()
	switch action {
	case "start":
		return sm.Start(serviceName)
	case "restart":
		return sm.Restart(serviceName)
	case "stop":
		return sm.Stop(serviceName)
	case "enable":
		return sm.Enable(serviceName)
	case "status":
		_, err := sm.Status(serviceName)
		return err
	default:
		return fmt.Errorf("Unknown service action: %s", action)
	}
}