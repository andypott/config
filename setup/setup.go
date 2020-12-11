package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/sys/unix"
)

type system struct {
	disks int
	setup func([]string)
}

type checkResult struct {
	check   string
	success bool
	msg     string
}

type strSliceArgs []string

func (ssa *strSliceArgs) String() string {
	return fmt.Sprintf("%v", *ssa)
}

func (ssa *strSliceArgs) Set(argStr string) error {
	*ssa = strings.Split(argStr, ",")
	return nil
}

func printColor(color, msg string, eol bool) {
	end := ""
	if eol {
		end = "\n"
	}

	fmt.Printf("\033%s%s\033[0m%s", color, msg, end)
}

func printFailure(msg string, eol bool) {
	printColor("[31m", msg, eol)
}

func printSuccess(msg string, eol bool) {
	printColor("[32m", msg, eol)
}

func checkIsRoot(disks []string) checkResult {
	check := "Checking root user"
	msg := "Must be run as root user!"
	success := os.Getuid() == 0
	if success {
		msg = "OK"
	}
	return checkResult{check, success, msg}
}

func checkInstallDisks(disks []string) checkResult {
	check := "Checking devices exist"
	for _, disk := range disks {
		_, err := os.Stat("/sys/block/" + disk)
		if err != nil {
			msg := fmt.Sprintf("%s does not exist or is not a disk", disk)
			return checkResult{check, false, msg}
		}
	}
	return checkResult{check, true, "OK"}
}

func checkInstallDisksForPartitions(disks []string) checkResult {
	check := "Checking install device for partitions"
	for _, disk := range disks {
		dir := fmt.Sprintf("/sys/block/%s", disk)
		files, err := ioutil.ReadDir(dir)
		if err != nil {
			return checkResult{check, false, fmt.Sprintf("Error reading %s!", dir)}
		}
		partitions := 0
		for _, file := range files {
			if file.IsDir() {
				_, err := os.Stat(fmt.Sprintf("%s/%s/partition", dir, file.Name()))
				if err == nil {
					partitions++
				}
			}
		}
		if partitions != 0 {
			return checkResult{check, false, fmt.Sprintf("Found %d partitions on %s!", partitions, disk)}
		}
	}
	return checkResult{check, true, "OK"}
}

func runOrDie(name string, args ...string) {
	cmd := exec.Command(name, args...)
	err := cmd.Run()

	if err != nil {
		printFailure(fmt.Sprintf("Unable to run command: %s %v!", name, args), true)
		os.Exit(1)
	}
}

func strToMountOpts(opts string) (uintptr, string) {
	optSlice := strings.Split(opts, ",")
	var mountOpts uintptr = 0
	fsOpts := ""

	for _, opt := range optSlice {
		switch opt {
		case "rw":
			// Ignore, this is default
		case "relatime":
			mountOpts |= unix.MS_RELATIME
		default:
			fsOpts += opt + ","
		}
	}

	return mountOpts, strings.Trim(fsOpts, ",")
}

func mountOrDie(fs, partition, mountpoint, opts string) {
	mountOpts, fsOpts := strToMountOpts(opts)

	err := unix.Mount(partition, mountpoint, fs, mountOpts, fsOpts)
	if err != nil {
		printFailure(fmt.Sprintf("Unable to mount %s on %s!", mountpoint, partition), true)
		fmt.Println(err)
		os.Exit(1)
	}
}

func mountBtrfsOrDie(partition, mountpoint, opts, subvol string) {
	opts += ",subvol=" + subvol
	mountOrDie("btrfs", partition, mountpoint, opts)
}

func unmountOrDie(mountpoint string) {
	err := unix.Unmount(mountpoint, 0)

	if err != nil {
		printFailure(fmt.Sprintf("Unable to unmount %s!", mountpoint), true)
		fmt.Println(err)
		os.Exit(1)
	}
}

func mkdirOrDie(path string, perms os.FileMode) {
	err := os.MkdirAll(path, perms)
	if err != nil {
		printFailure(fmt.Sprintf("Unable create %s!", path), true)
		os.Exit(1)
	}
}

func getUuidOrDie(part string) string {
	// lsblk -n -o UUID /dev/nvme0n1p2
	cmd := exec.Command("lsblk", "-n", "-o", "UUID", part)
	out, err := cmd.Output()

	if err != nil {
		printFailure(fmt.Sprintf("Unable to get uuid for %s!", part), true)
		os.Exit(1)
	}

	return strings.TrimSpace(string(out))
}

func partName(disk string, num uint) string {
	prefix := ""
	if strings.Contains(disk, "nvme") {
		prefix = "p"
	}
	return fmt.Sprintf("/dev/%s%s%d", disk, prefix, num)
}

func runInteractiveOrDie(name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to run command: %s %v!", name, args)
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func ultra24(disks []string) {
	disk := disks[0]
	hostname := "ultra24"
	btrfsOpts := "rw,relatime,compress=zstd,ssd,space_cache"
	fat32Opts := "rw,relatime,fmask=0022,dmask=0022,codepage=437,iocharset=iso8859-1,shortname=mixed,utf8,errors=remount-ro"

	fmt.Print("Creating partitions...")

	dev := fmt.Sprintf("/dev/%s", disk)
	runOrDie("parted", "-s", dev, "mklabel", "gpt")
	runOrDie("parted", "-s", dev, "mkpart", "BOOT", "fat32", "1MiB", "513MiB")
	runOrDie("parted", "-s", dev, "set", "1", "esp", "on")
	runOrDie("parted", "-s", dev, "mkpart", "ROOT", "btrfs", "513Mib", "100%")

	printSuccess("OK", true)

	fmt.Print("Formatting partitions...")

	bootPart := partName(disk, 1)
	rootPart := partName(disk, 2)

	runOrDie("mkfs.fat", "-F", "32", bootPart)
	runOrDie("mkfs.btrfs", "-f", rootPart)

	printSuccess("OK", true)

	fmt.Print("Creating btrfs subvolumes...")

	mountBtrfsOrDie(rootPart, "/mnt", btrfsOpts, "/")

	runOrDie("btrfs", "subvolume", "create", "/mnt/_active")
	runOrDie("btrfs", "subvolume", "create", "/mnt/_active/root")
	runOrDie("btrfs", "subvolume", "create", "/mnt/_active/home")
	runOrDie("btrfs", "subvolume", "create", "/mnt/_active/tmp")
	runOrDie("btrfs", "subvolume", "create", "/mnt/_snapshots")

	unmountOrDie("/mnt")

	printSuccess("OK", true)

	fmt.Print("Mounting partitions for install...")

	var perms os.FileMode = 0777

	mountBtrfsOrDie(rootPart, "/mnt", btrfsOpts, "_active/root")
	mkdirOrDie("/mnt/home", perms)
	mountBtrfsOrDie(rootPart, "/mnt/home", btrfsOpts, "_active/home")
	mkdirOrDie("/mnt/tmp", perms)
	mountBtrfsOrDie(rootPart, "/mnt/tmp", btrfsOpts, "_active/tmp")
	mkdirOrDie("/mnt/mnt/defvol", perms)
	mountBtrfsOrDie(rootPart, "/mnt/mnt/defvol", btrfsOpts, "/")
	mkdirOrDie("/mnt/boot/efi", perms)
	mountOrDie("vfat", bootPart, "/mnt/boot/efi", fat32Opts)

	printSuccess("OK", true)

	fmt.Print("Running pacstrap...")
	runOrDie("pacstrap", "/mnt", "base", "linux", "linux-firmware", "git")
	printSuccess("OK", true)

	fmt.Print("Creating fstab...")
	espUuid := getUuidOrDie(partName(disk, 1))
	btrfsUuid := getUuidOrDie(partName(disk, 2))

	file, err := os.OpenFile("/mnt/etc/fstab", os.O_WRONLY|os.O_CREATE, 0664)
	if err != nil {
		printFailure("Unable to create fstab!", true)
		os.Exit(1)
	}
	file.WriteString("# Generated automatically - remember to update the setup script if updating this file!\n")
	file.WriteString(fmt.Sprintf("UUID=%s / btrfs %s,subvol=_active/root 0 1\n", btrfsUuid, btrfsOpts))
	file.WriteString(fmt.Sprintf("UUID=%s /boot/efi vfat %s 0 2\n", espUuid, fat32Opts))
	file.WriteString(fmt.Sprintf("UUID=%s /home btrfs %s,subvol=_active/home 0 2\n", btrfsUuid, btrfsOpts))
	file.WriteString(fmt.Sprintf("UUID=%s /tmp btrfs %s,subvol=_active/tmp 0 2\n", btrfsUuid, btrfsOpts))
	file.WriteString(fmt.Sprintf("UUID=%s /mnt/defvol btrfs %s,subvol=/ 0 2\n", btrfsUuid, btrfsOpts))

	if err = file.Close(); err != nil {
		printFailure("Unable to close fstab! For some reason!", true)
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	printSuccess("OK", true)

	fmt.Print("Setting timezone...")
	runOrDie("arch-chroot", "/mnt", "ln", "-sf", "/usr/share/zoneinfo/Europe/London", "/etc/localtime")
	runOrDie("arch-chroot", "/mnt", "hwclock", "--systohc")
	printSuccess("OK", true)

	fmt.Print("Creating User...")
	runOrDie("arch-chroot", "/mnt", "useradd", "-m", "-G", "wheel", "andy")
	printSuccess("OK", true)

	fmt.Print("Cloning Config Repo...")
	runOrDie("arch-chroot", "-u", "andy", "/mnt", "git", "clone", "https://github.com/andypott/config", "/home/andy/config")
	printSuccess("OK", true)

	fmt.Print("Configuring installed system...")
	runOrDie("arch-chroot", "/mnt", "/home/andy/config/bin/sysconf", "-system", hostname, "-installgrub")
	printSuccess("OK", true)

	fmt.Println("Please set your password...")
	runInteractiveOrDie("arch-chroot", "/mnt", "passwd", "andy")
	printSuccess("...OK", true)

}

func main() {
	systems := map[string]system{
		"ultra24": {disks: 1, setup: ultra24},
	}

	var system string
	var disks strSliceArgs
	flag.StringVar(&system, "system", "", "Required. The hostname of the system to setup")
	flag.Var(&disks, "disks", "Required. Comma seperated list of disks to use. Order matters!!")

	flag.Parse()

	if system == "" {
		flag.Usage()
		os.Exit(1)
	}

	if len(disks) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	if _, ok := systems[system]; !ok {
		fmt.Fprintf(os.Stderr, "Unknown system %s\n", system)
		os.Exit(1)
	}

	if len(disks) != systems[system].disks {
		fmt.Fprintf(os.Stderr, "%s requires exactly %d disks\n", system, systems[system].disks)
		os.Exit(1)
	}

	checks := []func([]string) checkResult{
		checkIsRoot,
		checkInstallDisks,
		checkInstallDisksForPartitions,
	}

	failures := 0
	for _, c := range checks {
		res := c(disks)
		fmt.Print(res.check)
		fmt.Print("...")
		if res.success {
			printSuccess(res.msg, true)
		} else {
			printFailure(res.msg, true)
			failures++
		}
	}

	if failures > 0 {
		printFailure("All checks must pass to continue. Exiting.", true)
		os.Exit(1)
	} else {
		printSuccess("All checks passed", true)
		systems[system].setup(disks)
	}
}
