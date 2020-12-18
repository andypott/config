package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
)

var STDOUT io.Writer
var STDERR io.Writer

func linkDirContents(src string, dest string) {
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
			linkDirContents(srcFilename, destFilename)
		} else {
			if destFile, err := os.Lstat(destFilename); err != nil {
				// File doesn't exist
				if err = os.Symlink(srcFilename, destFilename); err != nil {
					fmt.Fprintf(os.Stderr, "Unable to link %s to %s\n", srcFilename, destFilename)
					fmt.Fprintln(os.Stderr, err)
					os.Exit(1)
				}
			} else {
				remove := destFile.Mode().IsRegular()
				if !remove && (destFile.Mode()&os.ModeSymlink != 0) {
					// Symlink is already there, need to check if it's correct
					linkDest, err := os.Readlink(destFilename)
					if err != nil {
						fmt.Fprintf(os.Stderr, "%s is already a link, but it is unreadable!\n", destFilename)
						fmt.Fprintln(os.Stderr, err)
						os.Exit(1)
					}
					if linkDest != srcFilename {
						// symlink points somewhere else, remove it
						remove = true
					}
				}
				if remove {
					if err = os.Remove(destFilename); err != nil {
						fmt.Fprintf(os.Stderr, "Unable to link %s to %s; couldn't remove existing file\n", srcFilename, destFilename)
						fmt.Fprintln(os.Stderr, err)
						os.Exit(1)
					}
					if err = os.Symlink(srcFilename, destFilename); err != nil {
						fmt.Fprintf(os.Stderr, "Unable to link %s to %s\n", srcFilename, destFilename)
						fmt.Fprintln(os.Stderr, err)
						os.Exit(1)
					}
				}
			}
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
		fmt.Fprintf(os.Stderr, "Unable to run command: %s %v!\n", name, args)
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
		fmt.Fprintf(os.Stderr, "Unable to run command: %s %v!\n", name, args)
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func enableServices(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to open %s!\n", filename)
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		runOrDie("systemctl", "enable", scanner.Text())
	}

	if err = file.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "Unable to open %s!\n", filename)
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func downloadFileOrDie(src, dest string) {
	response, err := http.Get(src)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to download %s!\n", src)
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err = os.MkdirAll(filepath.Dir(dest), 0700); err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create path for %s!\n", dest)
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	destFile, err := os.Create(dest)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create file %s!\n", dest)
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	_, err = io.Copy(destFile, response.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to copy to file!\n", dest)
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err = response.Body.Close(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if err = destFile.Close(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func fileExists(fullpath string) bool {
	_, err := os.Stat(fullpath)
	return err == nil
}

func setupNeovim(home string) {
	vimPluggedLocation := home + "/.local/share/nvim/site/autoload/plug.vim"
	if !fileExists(vimPluggedLocation) {
		downloadFileOrDie("https://raw.githubusercontent.com/junegunn/vim-plug/master/plug.vim", vimPluggedLocation)
	}

	runOrDie("nvim", "--headless", "+PlugInstall", "+qall")
	runOrDie("nvim", "--headless", "+GoInstallBinaries", "+qall")
}

func main() {
	if os.Getuid() == 0 {
		fmt.Fprintln(os.Stderr, "Must NOT be run as root user!")
		os.Exit(1)
	}

	var system string
	var withOutput bool

	flag.StringVar(&system, "system", "", "Optional. The hostname of the system to configure.")
	flag.BoolVar(&withOutput, "output", false, "Optional. Display the output of the commands run.")
	flag.Parse()
	if system == "" {
		var err error
		if system, err = os.Hostname(); err != nil {
			fmt.Fprintln(os.Stderr, "No hostname provided and unable to get current hostname")
			os.Exit(1)
		}
	}

	srcPath, err := filepath.Abs(filepath.Dir(os.Args[0]) + "/../dotfiles")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	_, err = os.Stat(srcPath + "/" + system)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to find files for %s.\n", system)
		flag.Usage()
		os.Exit(1)
	}

	if withOutput {
		STDOUT = os.Stdout
		STDERR = os.Stderr
	}

	systemDir := srcPath + "/" + system
	sharedDir := srcPath + "/shared"
	user, err := user.Current()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Unable to get current user!")
		os.Exit(1)
	}
	homeDir := user.HomeDir

	linkDirContents(systemDir+"/files", homeDir)
	linkDirContents(sharedDir+"/files", homeDir)

	setupNeovim(homeDir)

	//enableServices(systemDir + "/services")
	//enableServices(sharedDir + "/services")
}
