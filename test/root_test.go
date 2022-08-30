package test_test

import (
	"log"
	"os"
	"os/exec"
	"runtime"
	"testing"

	"github.com/linzeyan/ops-cli/cmd/common"
)

const (
	mainCommand = "go"
	mainGo      = "../main.go"
	runCommand  = "run"
	testHost    = "google.com"
)

var binaryCommand = "../" + common.RepoName

func isWindows() bool {
	return runtime.GOOS == "windows"
}

func TestMain(m *testing.M) {
	if isWindows() {
		binaryCommand += ".exe"
	}
	if err := exec.Command(mainCommand, "build", "-trimpath", "-ldflags", "-s -w", "-o", binaryCommand, mainGo).Run(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	exitCode := m.Run()
	_ = os.Remove(binaryCommand)
	os.Exit(exitCode)
}
