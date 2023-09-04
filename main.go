package main

import (
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/hashicorp/go-hclog"
)

var logger hclog.Logger

func main() {
	logger = hclog.Default()
	logger.SetLevel(hclog.Debug)

	r := new(req)
	r.service = os.Getenv("PAM_SERVICE")
	r.username = os.Getenv("PAM_USER")

	logger.Debug("Start: gooroom-pam-openvpn-auth")

	if err := r.getPassword(); err != nil {
		logger.Debug("Error reading password", "error", err)
		os.Exit(2)
	}

	logger.Debug("========================================")
	logger.Debug(r.service)
	logger.Debug(r.username)
	logger.Debug(r.password)

	os.Setenv("USERNAME", r.username)
	os.Setenv("PASSWORD", r.password)

	c := exec.Command("/usr/bin/bash", "-c", "/usr/sbin/openvpn --config /etc/openvpn/client/gooroom.ovpn --auth-user-pass <(echo -e \""+r.username+"\n"+r.password+"\")")
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	w, err := c.StdinPipe()
	if err != nil {
		logger.Error("Can't open Stdin", "error", err)
	}
	defer w.Close()

	if err := c.Start(); err != nil {
		logger.Error("Faild to start openvpn client", "error", err)
		os.Exit(1)
	}

	logger.Debug("End")

	os.Exit(0)
}

type req struct {
	service  string
	username string
	password string
}

func (r *req) getPassword() error {
	b, err := io.ReadAll(os.Stdin)
	if err != nil {
		return err
	}

	r.password = strings.TrimSuffix(string(b), string('\x00'))

	return nil
}
