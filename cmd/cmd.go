/*
Copyright © 2022 Henry Huang <hhh@rutcode.com>
This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.
This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.
You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/

package cmd

import (
	"crypto/sha256"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"trellis.tech/trellis.v1/pkg/router"
	"trellis.tech/trellis.v1/pkg/server"
	"trellis.tech/trellis.v1/pkg/server/grpc_server"
	"trellis.tech/trellis.v1/pkg/server/http_server"
	"trellis.tech/trellis.v1/pkg/tracing"
	"trellis.tech/trellis.v1/pkg/trellis"

	"github.com/prometheus/client_golang/prometheus"
	"gopkg.in/yaml.v3"
	"trellis.tech/trellis/common.v1/errcode"
	"trellis.tech/trellis/common.v1/flagext"
)

const (
	configFileOption = "config.file"
	configExpandENV  = "config.expand-env"
)

// configHash exposes information about the loaded config
var configHash = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "trellis_config_hash",
		Help: "Hash of the currently active config file.",
	},
	[]string{"sha256"},
)

func Run() {
	var cfg = trellis.ServerConfig{}

	configFile, expandENV := parseConfigFileParameter(os.Args[1:])

	// This sets default values from flags to the config.
	// It needs to be called before parsing the config file!
	flagext.ParseFlags(&cfg)

	if configFile != "" {
		if err := LoadConfig(configFile, expandENV, &cfg); err != nil {
			fmt.Fprintf(os.Stderr, "error loading config from %s: %v\n", configFile, err)
			os.Exit(1)
		}
	}

	flagext.IgnoredFlag(flag.CommandLine, configFileOption, "Configuration file to load.")
	_ = flag.CommandLine.Bool(configExpandENV, false, "Expands ${var} or $var in config according to the values of the environment variables.")

	usage := flag.CommandLine.Usage
	flag.CommandLine.Usage = func() { /* don't do anything by default, we will print usage ourselves, but only when requested. */ }
	flag.CommandLine.Init(flag.CommandLine.Name(), flag.ContinueOnError)

	if err := flag.CommandLine.Parse(os.Args[1:]); err == flag.ErrHelp {
		// Print available parameters to stdout, so that users can grep/less it easily.
		flag.CommandLine.SetOutput(os.Stdout)
		usage()
		return
	} else if err != nil {
		fmt.Fprintln(flag.CommandLine.Output(), "Run with -help to get list of available parameters")
	}

	tracingCloser, err := tracing.InitTracer(cfg.ServerName, &cfg.TracingConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error init tracer %s: %v\n", configFile, err)
		os.Exit(1)
	}

	var (
		svr server.Server
		r   router.Router
	)

	r, err = router.NewRouter(cfg.RouterConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error new router %s: %v\n", configFile, err)
		os.Exit(1)
	}

	switch t := server.ServerType_name[int32(cfg.ServerType)]; t {
	case "ALL":
		svr, err = http_server.NewServer(
			http_server.ServerName(cfg.ServerName),
			http_server.Config(&cfg.HTTPServerConfig),
			http_server.Router(r),
			http_server.Tracing(cfg.TracingConfig.Enable))
		if err != nil {
			os.Exit(1)
		}
		svr, err = grpc_server.NewServer(
			grpc_server.ServerName(cfg.ServerName),
			grpc_server.Config(&cfg.GrpcServerConfig),
			grpc_server.Router(r),
			grpc_server.Tracing(cfg.TracingConfig.Enable))
		if err != nil {
			os.Exit(1)
		}
	case "HTTP":
		svr, err = http_server.NewServer(
			http_server.ServerName(cfg.ServerName),
			http_server.Config(&cfg.HTTPServerConfig),
			http_server.Router(r),
			http_server.Tracing(cfg.TracingConfig.Enable))
		if err != nil {
			os.Exit(1)
		}
	case "GRPC":
		svr, err = grpc_server.NewServer(
			grpc_server.ServerName(cfg.ServerName),
			grpc_server.Config(&cfg.GrpcServerConfig),
			grpc_server.Router(r),
			grpc_server.Tracing(cfg.TracingConfig.Enable))
		if err != nil {
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "error server type %d\n", cfg.ServerType)
		os.Exit(1)
	}

	if err = svr.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "error start server: %v\n", err)
		os.Exit(1)
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGQUIT)
	<-ch

	if tracingCloser != nil {
		err = tracingCloser.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed close tracing closer: %v\n", err)
		}
	}

	if err := svr.Stop(); err != nil {
		fmt.Fprintf(os.Stderr, "error stop server: %v\n", err)
		os.Exit(1)
	}
}

// Parse -config.file and -config.expand-env option via separate flag set, to avoid polluting default one and calling flag.Parse on it twice.
func parseConfigFileParameter(args []string) (configFile string, expandEnv bool) {
	// ignore errors and any output here. Any flag errors will be reported by main flag.Parse() call.
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.SetOutput(ioutil.Discard)

	// usage not used in these functions.
	fs.StringVar(&configFile, configFileOption, "", "")
	fs.BoolVar(&expandEnv, configExpandENV, false, "")

	// Try to find -config.file and -config.expand-env option in the flags. As Parsing stops on the first error, eg. unknown flag, we simply
	// try remaining parameters until we find config flag, or there are no params left.
	// (ContinueOnError just means that flag.Parse doesn't call panic or os.Exit, but it returns error, which we ignore)
	for len(args) > 0 {
		_ = fs.Parse(args)
		args = args[1:]
	}

	return
}

// LoadConfig read YAML-formatted config from filename into cfg.
func LoadConfig(filename string, expandENV bool, cfg *trellis.ServerConfig) error {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return errcode.NewErrors(err, errcode.New("Error reading config file"))
	}

	// create a sha256 hash of the config before expansion and expose it via
	// the config_info metric
	hash := sha256.Sum256(buf)
	configHash.Reset()
	configHash.WithLabelValues("sha256", fmt.Sprintf("%x", hash)).Set(1)

	if expandENV {
		buf = expandEnv(buf)
	}

	err = yaml.Unmarshal(buf, cfg)
	if err != nil {
		return errcode.NewErrors(err, errcode.New("Error parsing config file"))
	}
	return nil
}

// expandEnv replaces ${var} or $var in config according to the values of the current environment variables.
// The replacement is case-sensitive. References to undefined variables are replaced by the empty string.
// A default value can be given by using the form ${var:default value}.
func expandEnv(config []byte) []byte {
	return []byte(os.Expand(string(config), func(key string) string {
		keyAndDefault := strings.SplitN(key, ":", 2)
		key = keyAndDefault[0]

		v := os.Getenv(key)
		if v == "" && len(keyAndDefault) == 2 {
			v = keyAndDefault[1] // Set value to the default.
		}
		return v
	}))
}
