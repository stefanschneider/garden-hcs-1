package container

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"code.cloudfoundry.org/cflager"
	"code.cloudfoundry.org/garden"
	"github.com/pborman/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/hpcloud/garden-hcs/windows_containers"
)

var windowsServerCodeImagePath string

func TestMain(m *testing.M) {
	cmd := exec.Command("docker", "inspect", "microsoft/windowsservercore")
	imageJson, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	windowsImage := []map[string]interface{}{}
	err = json.Unmarshal(imageJson, &windowsImage)
	if err != nil {
		panic(err)
	}

	graphDriver := windowsImage[0]["GraphDriver"].(map[string]interface{})
	graphDriverData := graphDriver["Data"].(map[string]interface{})
	graphDriverDataDir := graphDriverData["dir"].(string)
	windowsServerCodeImagePath = graphDriverDataDir

	cflager.AddFlags(flag.CommandLine)

	retCode := m.Run()

	os.Exit(retCode)
}

func TestCreateContainer(t *testing.T) {
	assert := assert.New(t)

	logger, _ := cflager.New("windows-garden-tests")

	id := uuid.New()
	handle := id
	rootPath := "WindowsServerCore:dummy"
	hostIP := "127.0.0.1"
	driverInfo := windows_containers.NewDriverInfo("c:\\garden-windows\\tests")
	properties := garden.Properties{}

	containerSpec := garden.ContainerSpec{
		Handle:     handle,
		Properties: properties,
		RootFSPath: rootPath,
	}

	container, err := NewContainer(id, handle, containerSpec, logger, hostIP, driverInfo, windowsServerCodeImagePath)
	defer container.Stop(true)

	assert.Nil(err)
}

func TestRunInContainer(t *testing.T) {
	assert := assert.New(t)

	logger, _ := cflager.New("windows-garden-tests")

	id := uuid.New()
	handle := id
	rootPath := "WindowsServerCore:dummy"
	hostIP := "127.0.0.1"
	driverInfo := windows_containers.NewDriverInfo("c:\\garden-windows\\tests")
	properties := garden.Properties{}

	containerSpec := garden.ContainerSpec{
		Handle:     handle,
		Properties: properties,
		RootFSPath: rootPath,
	}

	container, err := NewContainer(id, handle, containerSpec, logger, hostIP, driverInfo, windowsServerCodeImagePath)
	defer container.Stop(true)

	assert.Nil(err)

	processSpec := garden.ProcessSpec{
		Path: "cmd.exe",
		Args: []string{"/c", "ver"},
		Env:  []string{},
		Dir:  "c:\\",
	}

	pio := garden.ProcessIO{
		Stdin:  nil,
		Stdout: nil,
		Stderr: nil,
	}

	pt, err := container.Run(processSpec, pio)
	assert.NoError(err)

	exitCode, err := pt.Wait()

	assert.Nil(err)
	assert.Equal(0, exitCode)
}

func TestRunInContainerLinuxPaths(t *testing.T) {
	assert := assert.New(t)

	logger, _ := cflager.New("windows-garden-tests")

	id := uuid.New()
	handle := id
	rootPath := "WindowsServerCore:dummy"
	hostIP := "127.0.0.1"
	driverInfo := windows_containers.NewDriverInfo("c:\\garden-windows\\tests")
	properties := garden.Properties{}

	containerSpec := garden.ContainerSpec{
		Handle:     handle,
		Properties: properties,
		RootFSPath: rootPath,
	}

	container, err := NewContainer(id, handle, containerSpec, logger, hostIP, driverInfo, windowsServerCodeImagePath)
	defer container.Stop(true)

	assert.Nil(err)

	processSpec := garden.ProcessSpec{
		Path: "c:\\Windows\\System32\\cmd.exe",
		Args: []string{"/c", "ver"},
		Env:  []string{},
		Dir:  "c:\\",
	}

	pio := garden.ProcessIO{
		Stdin:  nil,
		Stdout: nil,
		Stderr: nil,
	}

	pt, err := container.Run(processSpec, pio)
	assert.Nil(err)

	exitCode, err := pt.Wait()

	assert.Nil(err)
	assert.Equal(0, exitCode)
}

func TestRunInContainerWithOutput(t *testing.T) {
	assert := assert.New(t)

	logger, _ := cflager.New("windows-garden-tests")

	id := uuid.New()
	handle := id
	rootPath := "WindowsServerCore:dummy"
	hostIP := "127.0.0.1"
	driverInfo := windows_containers.NewDriverInfo("c:\\garden-windows\\tests")
	properties := garden.Properties{}

	containerSpec := garden.ContainerSpec{
		Handle:     handle,
		Properties: properties,
		RootFSPath: rootPath,
	}

	container, err := NewContainer(id, handle, containerSpec, logger, hostIP, driverInfo, windowsServerCodeImagePath)
	defer container.Stop(true)

	assert.Nil(err)

	processSpec := garden.ProcessSpec{
		Path: "cmd.exe",
		Args: []string{"/c", "ver"},
		Env:  []string{},
		Dir:  "c:\\",
	}

	stdout := bytes.NewBufferString("")

	pio := garden.ProcessIO{
		Stdin:  nil,
		Stdout: stdout,
		Stderr: nil,
	}

	pt, err := container.Run(processSpec, pio)
	assert.Nil(err)

	exitCode, err := pt.Wait()

	assert.Nil(err)
	assert.Equal(0, exitCode)

	output := stdout.String()
	assert.Contains(output, "Windows")
}

func TestRunInContainerEnv(t *testing.T) {
	assert := assert.New(t)

	logger, _ := cflager.New("windows-garden-tests")

	id := uuid.New()
	handle := id
	rootPath := "WindowsServerCore:dummy"
	hostIP := "127.0.0.1"
	driverInfo := windows_containers.NewDriverInfo("c:\\garden-windows\\tests")
	properties := garden.Properties{}

	containerSpec := garden.ContainerSpec{
		Handle:     handle,
		Properties: properties,
		RootFSPath: rootPath,
		Env:        []string{"INSTANCE_INDEX=0"},
	}

	container, err := NewContainer(id, handle, containerSpec, logger, hostIP, driverInfo, windowsServerCodeImagePath)
	defer container.Stop(true)

	assert.Nil(err)

	processSpec := garden.ProcessSpec{
		Path: "cmd.exe",
		Args: []string{"/c", "set"},
		Env:  []string{},
		Dir:  "c:\\",
	}

	stdout := bytes.NewBufferString("")

	pio := garden.ProcessIO{
		Stdin:  nil,
		Stdout: stdout,
		Stderr: nil,
	}

	pt, err := container.Run(processSpec, pio)
	assert.Nil(err)

	exitCode, err := pt.Wait()

	assert.Nil(err)
	assert.Equal(0, exitCode)

	output := stdout.String()
	assert.Contains(output, "INSTANCE_INDEX")
}

func TestRunInContainerWithNetwork(t *testing.T) {
	assert := assert.New(t)

	logger, _ := cflager.New("windows-garden-tests")

	id := uuid.New()
	handle := id
	rootPath := "WindowsServerCore:dummy"
	hostIP := "127.0.0.1"
	driverInfo := windows_containers.NewDriverInfo("c:\\garden-windows\\tests")
	properties := garden.Properties{}

	containerSpec := garden.ContainerSpec{
		Handle:     handle,
		Properties: properties,
		RootFSPath: rootPath,
	}

	container, err := NewContainer(id, handle, containerSpec, logger, hostIP, driverInfo, windowsServerCodeImagePath)
	defer container.Stop(true)

	assert.Nil(err)

	processSpec := garden.ProcessSpec{
		Path: "powershell.exe",
		Args: []string{
			"-command",
			"$l = New-Object System.Net.HttpListener ; $l.Prefixes.Add('http://*:8080/'); $l.Start(); while ($l.IsListening) { $c = $l.GetContext() ; $q = $c.Request; Write-Output (date); $r = $c.Response ; $m = [System.Text.ASCIIEncoding]::ASCII.GetBytes(((gci -path env:*) | Out-String)); $r.ContentLength64 = $m.Length ; $r.OutputStream.Write($m, 0, $m.Length) ; $r.OutputStream.Dispose(); }",
		},
		Env: []string{},
		Dir: "c:\\",
	}

	pio := garden.ProcessIO{
		Stdin:  nil,
		Stdout: nil,
		Stderr: nil,
	}

	_, err = container.Run(processSpec, pio)
	assert.Nil(err)

	// Sleep a bit, let the server start
	time.Sleep(10000 * time.Millisecond)

	resp, err := http.Get(fmt.Sprintf("http://%s:8080/", container.containerIp))
	assert.Nil(err)

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	body := buf.String()

	assert.Contains(body, "windir")
}

func TestRunInContainerWithStreamIn(t *testing.T) {
	assert := assert.New(t)

	logger, _ := cflager.New("windows-garden-tests")

	id := uuid.New()
	handle := id
	rootPath := "WindowsServerCore:dummy"
	hostIP := "127.0.0.1"
	driverInfo := windows_containers.NewDriverInfo("c:\\garden-windows\\tests")
	properties := garden.Properties{}

	containerSpec := garden.ContainerSpec{
		Handle:     handle,
		Properties: properties,
		RootFSPath: rootPath,
	}

	container, err := NewContainer(id, handle, containerSpec, logger, hostIP, driverInfo, windowsServerCodeImagePath)
	defer container.Stop(true)
	assert.Nil(err)

	workDir, err := os.Getwd()
	assert.Nil(err)

	tarFile := filepath.Join(workDir, "../test-assets/files.tar")
	tarStream, err := os.Open(tarFile)
	assert.Nil(err)

	streamInSpec := garden.StreamInSpec{
		Path:      "\\testfiles",
		TarStream: tarStream,
	}

	err = container.StreamIn(streamInSpec)
	assert.Nil(err)

	processSpec := garden.ProcessSpec{
		Path: "cmd",
		Args: []string{"/c", "dir"},
		Env:  []string{},
		Dir:  "c:\\testfiles",
	}

	stdout := bytes.NewBufferString("")

	pio := garden.ProcessIO{
		Stdin:  nil,
		Stdout: stdout,
		Stderr: nil,
	}

	pt, _ := container.Run(processSpec, pio)
	pt.Wait()

	output := stdout.String()
	assert.Contains(output, "file1.txt")
}
