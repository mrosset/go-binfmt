package main

import "errors"
import "flag"
import "fmt"
import "github.com/str1ngs/util"
import "log"
import "os"
import "os/exec"
import "path/filepath"
import "runtime"
import "syscall"

const (
	REG_NAME  = "GO"
	REG_DONE  = "/proc/sys/fs/binfmt_misc/GO"
	REG_FILE  = "/proc/sys/fs/binfmt_misc/register"
	REG_MOUNT = "/proc/sys/fs/binfmt_misc"
)

var (
	// flags
	fregister   = flag.Bool("register", false, "register .go extenstions with binfmt")
	funregister = flag.Bool("unregister", false, "unregister .go extenstions with binfmt")

	// errors
	ErrorOsNotSupported = errors.New("binfmt is only supported on linux")
	ErrorPermissions    = errors.New("need to be root to register or unregister with binfmt")
	ErrorRegPath        = errors.New("run go-binfmt with its absolute path in order to register")
	ErrorMount          = errors.New("kernel needs binfmt support to mount /proc/sys/fs/binfmt_misc")
)

func init() {
	if runtime.GOOS != "linux" {
		log.Fatal(ErrorOsNotSupported)
	}
	log.SetFlags(log.Lshortfile)
}

func main() {
	flag.Parse()
	switch {
	case *fregister:
		register()
	case *funregister:
		unregister()
	default:
		run()
	}
}

func run() {
	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}
	dir, file := filepath.Split(flag.Arg(0))
	goRun := exec.Command("go", "run", file)
	goRun.Stdout = os.Stdout
	goRun.Stderr = os.Stderr
	goRun.Dir = dir
	err := goRun.Run()
	util.CheckFatal(err)
}

func register() {
	if os.Geteuid() != 0 {
		log.Fatal(ErrorPermissions)
	}
	if !util.FileExists(REG_FILE) {
		fmt.Println("binfmt_misc not mounted. mounting...")
		err := syscall.Mount("binfmt_misc", REG_MOUNT, "binfmt_misc", uintptr(0), "")
		util.CheckFatal(err)
	}
	unregister()
	bin := os.Args[0]
	if !filepath.IsAbs(bin) {
		log.Fatal(ErrorRegPath)
	}
	binFmt := fmt.Sprintf(":%s:E::go::%s:", REG_NAME, bin)
	fd, err := os.OpenFile(REG_FILE, os.O_WRONLY, 0)
	util.CheckFatal(err)
	defer fd.Close()
	_, err = fd.WriteString(binFmt)
	util.CheckFatal(err)
	fmt.Printf("Go interpreter %s registered with binfmt\n", bin)
}

func unregister() {
	if !util.FileExists(REG_DONE) {
		fmt.Println(".go extenstions are not registered with binfmt, skipping")
		return
	}
	if os.Geteuid() != 0 {
		log.Fatal(ErrorPermissions)
	}
	fd, err := os.OpenFile(REG_DONE, os.O_WRONLY, 0)
	util.CheckFatal(err)
	defer fd.Close()
	_, err = fd.WriteString("-1")
	util.CheckFatal(err)
	fmt.Println("Go interpreter unregistered with binfmt")
}
