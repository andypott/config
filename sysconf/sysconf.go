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

func runOrDie(name string, args ...string) {
	cmd := exec.Command(name, args...)
	err := cmd.Run()

	if err != nil {
		fmt.Fprintf("Unable to run command: %s %v!", name, args)
		os.Exit(1)
	}
}

func enableServices(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf("Unable to open %s!", filename)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		runOrDie("systemctl", "enable", scanner.Text())
	}

	if err = file.Close(); err != nil {
		fmt.Fprintf("Unable to open %s!", filename)
		os.Exit(1)
	}
}

func main() {
	var system string
	if os.Getuid() != 0 {
		fmt.Fprintln(os.Stderr, "Must be run as root user!")
		os.Exit(1)
	}

	flag.StringVar(&system, "system", "", "Required. The hostname of the system to configure")
	flag.Parse()
	if system == "" {
		flag.Usage()
		os.Exit(1)
	}

	srcPath, err := filepath.Abs(filepath.Dir(os.Args[0]) + "/../sysfiles")
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

	systemDir := srcPath + "/" + system
	sharedDir := srcPath + "/shared"

	runOrDie("pacman", "-S", "--needed", "--noconfirm", "-", "<", systemDir+"/pkgs")
	runOrDie("pacman", "-S", "--needed", "--noconfirm", "-", "<", sharedDir+"/pkgs")
	copyDir(systemDir+"/files", "/")
	copyDir(sharedDir+"/files", "/")
	enableServices(systemDir + "/services")
	enableServices(sharedDir + "/services")

}
