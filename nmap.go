package nmap

import (
	"fmt"
	"os/exec"
	"bytes"
	"github.com/pkg/errors"
	"encoding/xml"
)

type Nmap struct {
	SystemPath string
	Args       []string
	Ports      string
	Hosts      string
	Exclude    string
	Result     []byte
}

func (n *Nmap) SetSystemPath(systemPath string) {
	if systemPath != "" {
		n.SystemPath = systemPath
	}
}
func (n *Nmap) SetArgs(arg ...string) {
	n.Args = arg
}
func (n *Nmap) SetPorts(ports string) {
	n.Ports = ports
}
func (n *Nmap) SetHosts(hosts string) {
	n.Hosts = hosts
}

// 排除扫描IP/IP段
func (n *Nmap) SetExclude(exclude string) {
	n.Exclude = exclude
}
func (n *Nmap) Run() error {
	var (
		cmd        *exec.Cmd
		outb, errs bytes.Buffer
	)

	if n.Hosts != "" {
		n.Args = append(n.Args, n.Hosts)
	}

	if n.Ports != "" {
		n.Args = append(n.Args, "-p")
		n.Args = append(n.Args, n.Ports)
	}

	if n.Exclude != "" {
		n.Args = append(n.Args, "--exclude")
		n.Args = append(n.Args, n.Exclude)
	}

	n.Args = append(n.Args, "-oX")
	n.Args = append(n.Args, "-")

	cmd = exec.Command(n.SystemPath, n.Args ...)
	fmt.Println(cmd.Args)
	cmd.Stdout = &outb
	cmd.Stderr = &errs
	err := cmd.Run()

	if errs.Len() > 0 {
		return errors.New(errs.String())
	}
	if err != nil {
		return err
	}
	n.Result = outb.Bytes()
	return nil
}

// Parse takes a byte array of nmap xml data and unmarshals it into an
// NmapRun struct. All elements are returned as strings, it is up to the caller
// to check and cast them to the proper type.
func (n *Nmap) Parse() (*NmapRun, error) {
	r := &NmapRun{}
	err := xml.Unmarshal(n.Result, r)
	return r, err
}

func New() *Nmap {
	return &Nmap{
		SystemPath: "nmap",
	}
}
