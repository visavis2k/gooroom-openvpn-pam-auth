package main

import (
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/hashicorp/go-hclog"
	"github.com/spf13/viper"
)

var logger hclog.Logger

func main() {
	logger = hclog.Default()
	logger.SetLevel(hclog.Debug)

	r := new(req)
	r.service = os.Getenv("PAM_SERVICE")
	r.username = os.Getenv("PAM_USER")

	logger.Debug("===============================")
	logger.Debug("Start: gooroom-pam-openvpn-auth")
	logger.Debug("===============================")

	viper.SetConfigFile("/etc/gooroom/gooroom-client-server-register/gcsr.conf")
	viper.SetConfigType("props")
	viper.SetDefault("vpn_config", "/etc/openvpn/client/gooroom.ovpn")

	if err := viper.ReadInConfig(); err != nil {
		logger.Error("Error loading config", "error", err)
		os.Exit(5)
	}

	configFile := viper.GetString("vpn_config")

	if err := r.getPassword(); err != nil {
		logger.Error("Error reading password", "error", err)
		os.Exit(2)
	}

	c := exec.Command("/usr/bin/bash", "-c", "/usr/sbin/openvpn --config "+configFile+" --auth-user-pass <(echo -e \""+r.username+"\n"+r.password+"\")")
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

	logger.Debug("=============================")
	logger.Debug("End: gooroom-pam-openvpn-auth")
	logger.Debug("=============================")

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
