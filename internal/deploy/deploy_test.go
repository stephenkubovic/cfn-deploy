package deploy

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func contains(src []string, test string) bool {
	for _, v := range src {
		if v == test {
			return true
		}
	}
	return false
}

func mockExecCommand(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

const successOutput = "success output"
const failureOutput = "failure output"

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	if contains(os.Args, "success-stack") {
		fmt.Fprint(os.Stdout, successOutput)
	} else if contains(os.Args, "failure-stack") {
		fmt.Fprint(os.Stderr, failureOutput)
	}
	os.Exit(0)
}

func TestStderrPipe(t *testing.T) {
	execCommand = mockExecCommand
	defer func() {
		execCommand = exec.Command
	}()

	c := New().SetStackName("failure-stack")
	cmd, stdout, stderr, err := c.Execute()
	if err != nil {
		t.Error(err)
	}

	if err := cmd.Start(); err != nil {
		t.Error(err)
	}

	out, _ := ioutil.ReadAll(stdout)
	if string(out) != "" {
		t.Errorf("expected stdout to be empty but got '%s'", out)
	}

	e, _ := ioutil.ReadAll(stderr)
	if string(e) != failureOutput {
		t.Errorf("expected stderr to be '%s' but got '%s'", failureOutput, e)
	}

	err = cmd.Wait()
	if err != nil {
		t.Error(err)
	}
}

func TestStdoutPipe(t *testing.T) {
	execCommand = mockExecCommand
	defer func() {
		execCommand = exec.Command
	}()

	c := New().SetStackName("success-stack")
	cmd, stdout, stderr, err := c.Execute()
	if err != nil {
		t.Error(err)
	}

	if err := cmd.Start(); err != nil {
		t.Error(err)
	}

	out, _ := ioutil.ReadAll(stdout)
	if string(out) != successOutput {
		t.Errorf("expected stdout to be '%s' but got '%s'", successOutput, out)
	}

	e, _ := ioutil.ReadAll(stderr)
	if string(e) != "" {
		t.Errorf("expected stderr to be empty but got '%s'", e)
	}

	err = cmd.Wait()
	if err != nil {
		t.Error(err)
	}
}

func TestAwsCliArgs(t *testing.T) {
	c := New().SetStackName("test-stack").SetDebug()
	args := c.AwsCliArgs()
	s := strings.Join(args, " ")
	expected := "cloudformation deploy --stack-name test-stack --debug"
	if s != expected {
		t.Errorf("expected arg string to be %s but got %s", expected, s)
	}
}

func TestAwsCliArgsEmptyArgs(t *testing.T) {
	c := New()
	args := c.AwsCliArgs()
	s := strings.Join(args, " ")
	expected := "cloudformation deploy"
	if s != expected {
		t.Errorf("expected arg string to be %s but got %s", expected, s)
	}
}
