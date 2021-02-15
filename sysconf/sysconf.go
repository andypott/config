package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

var STDOUT io.Writer
var STDERR io.Writer

func copyFile(src string, dest string) {
	destFile, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	srcFile, err := os.Open(src)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	destFile.Sync()

	if err = srcFile.Close(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if err = destFile.Close(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func copyDir(src string, dest string) {
	contents, err := ioutil.ReadDir(src)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for _, file := range contents {
		srcFilename := src + "/" + file.Name()
		destFilename := dest + "/" + file.Name()

		if file.IsDir() {
			os.MkdirAll(destFilename, 0755)
			copyDir(srcFilename, destFilename)
		} else if file.Mode().IsRegular() {
			copyFile(srcFilename, destFilename)
		}
	}
}

func runWithStdinOrDie(stdinFile string, name string, args ...string) {
	file, err := os.Open(stdinFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to open %s!\n", stdinFile)
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	cmd := exec.Command(name, args...)
	cmd.Stdin = file
	cmd.Stdout = STDOUT
	cmd.Stderr = STDERR

	err = cmd.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to run command: %s %v!", name, args)
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func runOrDie(name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Stdout = STDOUT
	cmd.Stderr = STDERR

	err := cmd.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to run command: %s %v!", name, args)
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func enableServices(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to open %s!", filename)
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		runOrDie("systemctl", "enable", scanner.Text())
	}

	if err = file.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "Unable to open %s!", filename)
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func main() {
	if os.Getuid() != 0 {
		fmt.Fprintln(os.Stderr, "Must be run as root user!")
		os.Exit(1)
	}

	var system string
	var withOutput bool
	var installgrub bool

	flag.StringVar(&system, "system", "", "Optional. The hostname of the system to configure.")
	flag.BoolVar(&withOutput, "output", false, "Optional. Display the output of the commands run.")
	flag.BoolVar(&installgrub, "installgrub", false, "Optional.  Install grub bootloader.")
	flag.Parse()
	if system == "" {
		var err error
		if system, err = os.Hostname(); err != nil {
			fmt.Fprintln(os.Stderr, "No hostname provided and unable to get current hostname")
			os.Exit(1)
		}
	}

	exePath, err := os.Executable()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	srcPath, err := filepath.Abs(filepath.Dir(exePath) + "/../sysfiles")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	_, err = os.Stat(srcPath + "/" + system)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to find config for %s.\n", system)
		flag.Usage()
		os.Exit(1)
	}

	if withOutput {
		STDOUT = os.Stdout
		STDERR = os.Stderr
	}

	systemDir := srcPath + "/" + system
	sharedDir := srcPath + "/shared"
	// Need to copy before running pacman to ensure that pacman.conf is there
	copyDir(systemDir+"/files", "/")
	copyDir(sharedDir+"/files", "/")

	runWithStdinOrDie(systemDir+"/pkgs", "pacman", "-Sy", "--needed", "--noconfirm", "-")
	runWithStdinOrDie(sharedDir+"/pkgs", "pacman", "-S", "--needed", "--noconfirm", "-")

	runOrDie("locale-gen")

	enableServices(systemDir + "/services")
	enableServices(sharedDir + "/services")

	if installgrub {
		runOrDie("grub-install", "--target=x86_64-efi", "--efi-directory=/boot/efi", "--bootloader-id=Arch")
	}
	runOrDie("grub-mkconfig", "-o", "/boot/grub/grub.cfg")

}
