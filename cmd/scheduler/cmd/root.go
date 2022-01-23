/*
 *     Copyright 2020 The Dragonfly Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package cmd

import (
	"context"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"d7y.io/dragonfly/v2/cmd/dependency"
	logger "d7y.io/dragonfly/v2/internal/dflog"
	"d7y.io/dragonfly/v2/internal/dflog/logcore"
	"d7y.io/dragonfly/v2/pkg/dfpath"
	"d7y.io/dragonfly/v2/scheduler"
	"d7y.io/dragonfly/v2/scheduler/config"
	"d7y.io/dragonfly/v2/version"
)

var (
	cfg *config.Config
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "scheduler",
	Short: "the scheduler of dragonfly",
	Long: `Scheduler is a long-running process which receives and manages download tasks from the client, notify the CDN to return to the source, 
generate and maintain a P2P network during the download process, and push suitable download nodes to the client`,
	Args:              cobra.NoArgs,
	DisableAutoGenTag: true,
	SilenceUsage:      true,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// Initialize dfpath
		d, err := initDfpath(cfg.Server)
		if err != nil {
			return err
		}

		// Initialize logger
		if err := logcore.InitScheduler(cfg.Console, d.LogDir()); err != nil {
			return errors.Wrap(err, "init scheduler logger")
		}

		// Validate config
		if err := cfg.Validate(); err != nil {
			return err
		}

		return runScheduler(ctx, d)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logger.Error(err)
		os.Exit(1)
	}
}

func init() {
	// Initialize default scheduler config
	cfg = config.New()
	// Initialize cobra
	dependency.InitCobra(rootCmd, true, cfg)
}

func initDfpath(cfg *config.ServerConfig) (dfpath.Dfpath, error) {
	var options []dfpath.Option
	if cfg.LogDir != "" {
		options = append(options, dfpath.WithLogDir(cfg.LogDir))
	}

	if cfg.CacheDir != "" {
		options = append(options, dfpath.WithCacheDir(cfg.CacheDir))
	}

	return dfpath.New(options...)
}

func runScheduler(ctx context.Context, d dfpath.Dfpath) error {
	logger.Infof("Version:\n%s", version.Version())

	// scheduler config values
	s, _ := yaml.Marshal(cfg)

	logger.Infof("scheduler configuration:\n%s", string(s))

	ff := dependency.InitMonitor(cfg.Verbose, cfg.PProfPort, cfg.Telemetry)
	defer ff()

	svr, err := scheduler.New(ctx, cfg, d)
	if err != nil {
		return err
	}

	dependency.SetupQuitSignalHandler(func() { svr.Stop() })
	return svr.Serve()
}
