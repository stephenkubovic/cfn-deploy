package orchestrator

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/fatih/color"
	"github.com/stephenkubovic/cfn-deploy/internal/deploy"
	"github.com/stephenkubovic/cfn-deploy/internal/stackevents"
)

// Command represents is the main orchestration operation.
type Command struct {
	StackName    string
	TemplateFile string
	Profile      string
	Params       []string
	Debug        bool
	KmsKeyID     string
	S3Bucket     string
	S3Prefix     string
	ForceUpload  bool
}

// Execute triggers the cloudformation deploy operation and starts polling
// for stack events.
func (c *Command) Execute() error {
	deploy := c.CreateDeployCommand()

	if c.Debug {
		c.Debugf(strings.Join(deploy.AwsCliArgs(), " "))
	}

	cmd, stdout, stderr, err := deploy.Execute()
	if err != nil {
		return err
	}

	// receives progress updates from the deploy command
	progress := make(chan int)

	// receives `true` when stdout is done processing
	stdoutChan := make(chan bool)

	// receives `true` when stderr is done processing
	stderrChan := make(chan bool)

	// receives `true` when no more stack events are being read
	eventsChan := make(chan bool)

	stdoutReader := bufio.NewScanner(stdout)
	stderrReader := io.TeeReader(stderr, os.Stderr)

	go c.stdoutHandler(stdoutReader, progress, stdoutChan)
	go c.stderrHandler(stderrReader, stderrChan)
	go c.readStackEvents(progress, eventsChan)

	if err := cmd.Start(); err != nil {
		return err
	}

	<-stdoutChan

	if err := cmd.Wait(); err != nil {
		return err
	}

	<-stderrChan
	<-eventsChan

	return nil
}

func (c *Command) readStackEvents(progress <-chan int, done chan<- bool) {
	defer func() {
		done <- true
	}()

	msg := <-progress
	if msg == deploy.ProgressNoChangeset {
		return
	}

	opts := session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}
	if c.Profile != "" {
		opts.Profile = c.Profile
	}

	sess, err := session.NewSessionWithOptions(opts)
	if err != nil {
		log.Println("Could not create aws session", err)
		return
	}

	cfn := cloudformation.New(sess)
	t := time.Now().Add(-5 * time.Second)
	write := func(se []stackevents.Event) {
		for _, event := range se {
			writeStackEvent(event, os.Stdout)
		}
	}

	for {
		select {
		case msg := <-progress:
			if msg == deploy.ProgressEOF {
				events, err := stackevents.Read(cfn, t, c.StackName)
				if err != nil {
					log.Println("error describing stack events", err)
					break
				}
				write(events)
				return
			}
		default:
			stackEvents, err := stackevents.Read(cfn, t, c.StackName)
			if err != nil {
				log.Println("error describing stack events", err)
				continue
			}
			write(stackEvents)
			if len(stackEvents) > 0 {
				t = stackEvents[len(stackEvents)-1].Timestamp
			}
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func (c *Command) stdoutHandler(stdout *bufio.Scanner, progress chan<- int, done chan<- bool) {
	defer func() {
		c.Debugf("stdout handler finished")
		done <- true
	}()

	var line string

	for stdout.Scan() {
		line = stdout.Text()
		fmt.Println(line)
		p := deploy.Progress(line)
		c.Debugf("deploy lined mapped to progress: %d", p)
		progress <- p
	}

	if err := stdout.Err(); err != nil {
		log.Fatal(err)
	}

	progress <- deploy.ProgressEOF
}

func (c *Command) stderrHandler(r io.Reader, done chan<- bool) {
	defer func() {
		c.Debugf("stderr handler finished")
		done <- true
	}()

	ioutil.ReadAll(r)
}

func writeStackEvent(e stackevents.Event, w io.Writer) {
	var status string = e.ResourceStatus
	if e.IsOk() {
		status = color.New(color.FgGreen).Sprint(status)
	} else if e.IsFailure() {
		status = color.New(color.FgRed).Sprint(status)
	} else if e.IsProgress() {
		status = color.New(color.FgBlue).Sprint(status)
	}

	id := color.New(color.Underline).Sprint(e.LogicalResourceID)
	t := color.New(color.Faint).Sprint(e.Timestamp.Format("15:04:05"))
	pad := strings.Repeat(" ", 8)

	var b strings.Builder
	b.WriteString(fmt.Sprintf("%s %s %s %s", t, id, e.PhysicalResourceID, status))
	b.WriteString("\n")
	if e.ResourceStatusReason != "" {
		b.WriteString(fmt.Sprintf("%s %s\n", pad, color.New(color.Italic).Sprint(e.ResourceStatusReason)))
	}
	b.WriteString("\n")
	w.Write([]byte(b.String()))
}

// CreateDeployCommand creates the CloudFormation stack deployment command.
func (c *Command) CreateDeployCommand() *deploy.DeployCommand {
	cmd := deploy.
		New().
		SetStackName(c.StackName).
		SetTemplateFile(c.TemplateFile)

	if c.Profile != "" {
		cmd.SetProfile(c.Profile)
	}
	if len(c.Params) > 0 {
		cmd.SetParameterOverrides(c.Params)
	}
	if c.Debug {
		cmd.SetDebug()
	}
	if c.ForceUpload {
		cmd.SetForceUpload()
	}
	if c.KmsKeyID != "" {
		cmd.SetKmsKeyID(c.KmsKeyID)
	}
	if c.S3Bucket != "" {
		cmd.SetS3Bucket(c.S3Bucket)
	}
	if c.S3Prefix != "" {
		cmd.SetS3Prefix(c.S3Prefix)
	}
	return cmd
}

// Debugf logs the given formatted message if command debugging is turned on.
func (c *Command) Debugf(s string, args ...interface{}) {
	if c.Debug {
		log.Printf(s, args...)
	}
}
