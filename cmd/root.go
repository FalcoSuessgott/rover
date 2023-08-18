package cmd

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/caarlos0/env/v6"
	"github.com/im2nguyen/rover/pkg/parser"
	"github.com/spf13/cobra"
)

const envVarPrefix = "ROVER_"

type opts struct {
	// General
	Name          string `env:"NAME"`
	ZipFileName   string `env:"ZIP_FILE_NAME"`
	WorkingDir    string `env:"WORKING_DIRECTORY"`
	ShowSensitive bool   `env:"SHOW_SENSITIVE"`
	GenImage      bool   `env:"GEN_IMAGE"`
	Listener      string `env:"LISTENER"`
	Standalone    bool   `env:"STANDALONE"`

	// Plan Options
	PlanPath     string `env:"PLAN_JSON"`
	PlanJSONPath string `env:"PLAN_JSON_PATH"`

	// TF
	TFPath           string   `env:"TF_PATH"`
	TFWorkspace      string   `env:"TF_WORKSPACE"`
	TFVarFiles       []string `env:"TF_VAR_FILES"`
	TFVars           []string `env:"TF_VARS"`
	TFBackendConfigs []string `env:"TF_BACKEND_CONFIGS"`

	// TFE & TFC
	TFCHostName      string `env:"TFC_HOSTNAME"`
	TFCOrgName       string `env:"TFC_ORG_NAME"`
	TFCWorkspaceName string `env:"TFC_WORKSPACE_NAME"`
	TFCNewRun        bool   `env:"TFC_NEW_RUN"`
}

func newDefaultOpts() *opts {
	return &opts{
		Name:          "rover",
		ZipFileName:   "rover",
		WorkingDir:    ".",
		Listener:      "0.0.0.0:9000",
		ShowSensitive: false,
		Standalone:    false,
		GenImage:      false,
		TFPath:        "/bin/terraform",
		TFCHostName:   "app.terraform.io",
		TFCNewRun:     false,
	}
}

// NewRootCmd vkv root command.
func NewRootCmd(v string, writer io.Writer) *cobra.Command {
	o := newDefaultOpts()

	if err := o.parseEnvs(); err != nil {
		log.Fatal(err)
	}

	cmd := &cobra.Command{
		Use:           "rover",
		Short:         "Interactive Terraform visualization, State and configuration explorer",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.validateFlags(); err != nil {
				return err
			}

			log.Println("Starting Rover...")

			err := parser.GenerateAssets()
			if err != nil {
				log.Fatal(err.Error())
			}
			return cmd.Help()
		},
	}

	// General
	cmd.Flags().StringVarP(&o.Name, "name", "n", o.Name, "Configuration name")
	cmd.Flags().StringVarP(&o.ZipFileName, "zip-file", "z", o.Name, "Standalone zip file name")
	cmd.Flags().StringVarP(&o.WorkingDir, "working-dir", "w", o.Name, "Path to Terraform configuration")
	cmd.Flags().BoolVarP(&o.ShowSensitive, "show-sensitive", "s", o.ShowSensitive, "Display sensitive values")
	cmd.Flags().BoolVarP(&o.GenImage, "gen-image", "g", o.ShowSensitive, "Generate graph image")
	cmd.Flags().BoolVar(&o.Standalone, "standalone", o.Standalone, "Generate standalone HTML files")

	// Plan
	cmd.Flags().StringVarP(&o.PlanPath, "plan-path", "p", o.PlanPath, "Plan file path")
	cmd.Flags().StringVarP(&o.PlanJSONPath, "plan-json-path", "j", o.PlanJSONPath, "Plan JSON file path")

	// TF
	cmd.Flags().StringVar(&o.TFPath, "tf-path", o.TFPath, "Path to Terraform binary")
	cmd.Flags().StringVar(&o.TFCWorkspaceName, "tf-workspace", o.TFPath, "Terraform Cloud Workspace name")
	cmd.Flags().StringSliceVar(&o.TFVarFiles, "tf-var-files", o.TFVarFiles, "Path to *.tfvars files")
	cmd.Flags().StringSliceVar(&o.TFVars, "tf-vars", o.TFVars, "Terraform variable (key=value)")
	cmd.Flags().StringSliceVar(&o.TFBackendConfigs, "tf-backend-configs", o.TFBackendConfigs, "Path to *.tfbackend files")

	// TFC & TFE
	cmd.Flags().StringVar(&o.TFCHostName, "tfc-hostname", o.TFCHostName, "Terraform Cloud/Enterprise Hostname")
	cmd.Flags().StringVar(&o.TFCOrgName, "tfc-org", o.TFCOrgName, "Terraform Cloud Organization name")
	cmd.Flags().StringVar(&o.TFCWorkspaceName, "tfc-workspace", o.TFCWorkspaceName, "Terraform Cloud Workspace name")
	cmd.Flags().BoolVar(&o.TFCNewRun, "tfc-run", o.TFCNewRun, "Create new Terraform Cloud run")

	return cmd
}

func (o *opts) validateFlags() error {
	switch {
	}

	return nil
}

// Execute invokes the command.
func Execute(version string) error {
	if err := NewRootCmd(version, os.Stdout).Execute(); err != nil {
		return fmt.Errorf("[ERROR] %w", err)
	}

	return nil
}

func (o *opts) parseEnvs() error {
	if err := env.Parse(o, env.Options{
		Prefix: envVarPrefix,
	}); err != nil {
		return err
	}

	return nil
}
