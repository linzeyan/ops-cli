package test_test

import (
	"log"
	"os"
	"os/exec"
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

func TestMain(m *testing.M) {
	if common.IsWindows() {
		binaryCommand += ".exe"
	}
	cmd := exec.Command(mainCommand, "build", "-trimpath", "-ldflags", "-s -w", "-o", binaryCommand, mainGo)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "CGO_ENABLED=0")
	if err := cmd.Run(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	exitCode := m.Run()
	_ = os.Remove(binaryCommand)
	os.Exit(exitCode)
}
