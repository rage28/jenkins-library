// Code generated by piper's step-generator. DO NOT EDIT.

package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/SAP/jenkins-library/pkg/config"
	"github.com/SAP/jenkins-library/pkg/log"
	"github.com/SAP/jenkins-library/pkg/piperenv"
	"github.com/SAP/jenkins-library/pkg/splunk"
	"github.com/SAP/jenkins-library/pkg/telemetry"
	"github.com/SAP/jenkins-library/pkg/validation"
	"github.com/spf13/cobra"
)

type transportRequestUploadCTSOptions struct {
	Description            string   `json:"description,omitempty"`
	Endpoint               string   `json:"endpoint,omitempty"`
	Client                 string   `json:"client,omitempty"`
	Username               string   `json:"username,omitempty"`
	Password               string   `json:"password,omitempty"`
	ApplicationName        string   `json:"applicationName,omitempty"`
	AbapPackage            string   `json:"abapPackage,omitempty"`
	OsDeployUser           string   `json:"osDeployUser,omitempty"`
	DeployConfigFile       string   `json:"deployConfigFile,omitempty"`
	TransportRequestID     string   `json:"transportRequestId,omitempty"`
	DeployToolDependencies []string `json:"deployToolDependencies,omitempty"`
	NpmInstallOpts         []string `json:"npmInstallOpts,omitempty"`
}

type transportRequestUploadCTSCommonPipelineEnvironment struct {
	custom struct {
		transportRequestID string
	}
}

func (p *transportRequestUploadCTSCommonPipelineEnvironment) persist(path, resourceName string) {
	content := []struct {
		category string
		name     string
		value    interface{}
	}{
		{category: "custom", name: "transportRequestId", value: p.custom.transportRequestID},
	}

	errCount := 0
	for _, param := range content {
		err := piperenv.SetResourceParameter(path, resourceName, filepath.Join(param.category, param.name), param.value)
		if err != nil {
			log.Entry().WithError(err).Error("Error persisting piper environment.")
			errCount++
		}
	}
	if errCount > 0 {
		log.Entry().Error("failed to persist Piper environment")
	}
}

// TransportRequestUploadCTSCommand This step uploads an UI5 application to the SAPUI5 ABAP repository.
func TransportRequestUploadCTSCommand() *cobra.Command {
	const STEP_NAME = "transportRequestUploadCTS"

	metadata := transportRequestUploadCTSMetadata()
	var stepConfig transportRequestUploadCTSOptions
	var startTime time.Time
	var commonPipelineEnvironment transportRequestUploadCTSCommonPipelineEnvironment
	var logCollector *log.CollectorHook
	var splunkClient *splunk.Splunk
	telemetryClient := &telemetry.Telemetry{}

	var createTransportRequestUploadCTSCmd = &cobra.Command{
		Use:   STEP_NAME,
		Short: "This step uploads an UI5 application to the SAPUI5 ABAP repository.",
		Long: `This step uploads an UI5 application from your project folder to the SAPUI5 ABAP repository of the SAPUI5 ABAP back-end infrastructure using the SAPUI5 Repository OData service.
It processes the results of the ` + "`" + `ui5 build` + "`" + ` command of the SAPUI5 toolset.`,
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			startTime = time.Now()
			log.SetStepName(STEP_NAME)
			log.SetVerbose(GeneralConfig.Verbose)

			GeneralConfig.GitHubAccessTokens = ResolveAccessTokens(GeneralConfig.GitHubTokens)

			path, _ := os.Getwd()
			fatalHook := &log.FatalHook{CorrelationID: GeneralConfig.CorrelationID, Path: path}
			log.RegisterHook(fatalHook)

			err := PrepareConfig(cmd, &metadata, STEP_NAME, &stepConfig, config.OpenPiperFile)
			if err != nil {
				log.SetErrorCategory(log.ErrorConfiguration)
				return err
			}
			log.RegisterSecret(stepConfig.Username)
			log.RegisterSecret(stepConfig.Password)

			if len(GeneralConfig.HookConfig.SentryConfig.Dsn) > 0 {
				sentryHook := log.NewSentryHook(GeneralConfig.HookConfig.SentryConfig.Dsn, GeneralConfig.CorrelationID)
				log.RegisterHook(&sentryHook)
			}

			if len(GeneralConfig.HookConfig.SplunkConfig.Dsn) > 0 {
				splunkClient = &splunk.Splunk{}
				logCollector = &log.CollectorHook{CorrelationID: GeneralConfig.CorrelationID}
				log.RegisterHook(logCollector)
			}

			if len(GeneralConfig.ANSServiceKey) > 0 {
				log.RegisterSecret(GeneralConfig.ANSServiceKey)
				ansHook := log.NewANSHook(GeneralConfig.ANSServiceKey, GeneralConfig.CorrelationID, GeneralConfig.ANSEventTemplateFilePath)
				log.RegisterHook(&ansHook)
			}

			validation, err := validation.New(validation.WithJSONNamesForStructFields(), validation.WithPredefinedErrorMessages())
			if err != nil {
				return err
			}
			if err = validation.ValidateStruct(stepConfig); err != nil {
				log.SetErrorCategory(log.ErrorConfiguration)
				return err
			}

			return nil
		},
		Run: func(_ *cobra.Command, _ []string) {
			stepTelemetryData := telemetry.CustomData{}
			stepTelemetryData.ErrorCode = "1"
			handler := func() {
				commonPipelineEnvironment.persist(GeneralConfig.EnvRootPath, "commonPipelineEnvironment")
				config.RemoveVaultSecretFiles()
				stepTelemetryData.Duration = fmt.Sprintf("%v", time.Since(startTime).Milliseconds())
				stepTelemetryData.ErrorCategory = log.GetErrorCategory().String()
				stepTelemetryData.PiperCommitHash = GitCommit
				telemetryClient.SetData(&stepTelemetryData)
				telemetryClient.Send()
				if len(GeneralConfig.HookConfig.SplunkConfig.Dsn) > 0 {
					splunkClient.Send(telemetryClient.GetData(), logCollector)
				}
			}
			log.DeferExitHandler(handler)
			defer handler()
			telemetryClient.Initialize(GeneralConfig.NoTelemetry, STEP_NAME)
			if len(GeneralConfig.HookConfig.SplunkConfig.Dsn) > 0 {
				splunkClient.Initialize(GeneralConfig.CorrelationID,
					GeneralConfig.HookConfig.SplunkConfig.Dsn,
					GeneralConfig.HookConfig.SplunkConfig.Token,
					GeneralConfig.HookConfig.SplunkConfig.Index,
					GeneralConfig.HookConfig.SplunkConfig.SendLogs)
			}
			transportRequestUploadCTS(stepConfig, &stepTelemetryData, &commonPipelineEnvironment)
			stepTelemetryData.ErrorCode = "0"
			log.Entry().Info("SUCCESS")
		},
	}

	addTransportRequestUploadCTSFlags(createTransportRequestUploadCTSCmd, &stepConfig)
	return createTransportRequestUploadCTSCmd
}

func addTransportRequestUploadCTSFlags(cmd *cobra.Command, stepConfig *transportRequestUploadCTSOptions) {
	cmd.Flags().StringVar(&stepConfig.Description, "description", `Deployed with Piper based on SAP Fiori tools`, "The description of the application. The description is only taken into account for a new upload. In case of an update the description will not be updated.")
	cmd.Flags().StringVar(&stepConfig.Endpoint, "endpoint", os.Getenv("PIPER_endpoint"), "The ODATA service endpoint: https://<host>:<port>")
	cmd.Flags().StringVar(&stepConfig.Client, "client", os.Getenv("PIPER_client"), "The ABAP client")
	cmd.Flags().StringVar(&stepConfig.Username, "username", os.Getenv("PIPER_username"), "Service user for uploading to the SAPUI5 ABAP repository")
	cmd.Flags().StringVar(&stepConfig.Password, "password", os.Getenv("PIPER_password"), "Service user password for uploading to the SAPUI5 ABAP repository")
	cmd.Flags().StringVar(&stepConfig.ApplicationName, "applicationName", os.Getenv("PIPER_applicationName"), "Name of the UI5 application")
	cmd.Flags().StringVar(&stepConfig.AbapPackage, "abapPackage", os.Getenv("PIPER_abapPackage"), "ABAP package name of the UI5 application")
	cmd.Flags().StringVar(&stepConfig.OsDeployUser, "osDeployUser", `node`, "Docker image user performing the deployment")
	cmd.Flags().StringVar(&stepConfig.DeployConfigFile, "deployConfigFile", `ui5-deploy.yaml`, "Configuration file for the fiori deployment")
	cmd.Flags().StringVar(&stepConfig.TransportRequestID, "transportRequestId", os.Getenv("PIPER_transportRequestId"), "ID of the transport request to which the UI5 application is uploaded")
	cmd.Flags().StringSliceVar(&stepConfig.DeployToolDependencies, "deployToolDependencies", []string{`@ui5/cli`, `@sap/ux-ui5-tooling`, `@ui5/logger`, `@ui5/fs`}, "List of additional dependencies to fiori related packages. By default a standard node docker image is used on which the dependencies are installed. Provide an empty list, in case your docker image already contains the required dependencies")
	cmd.Flags().StringSliceVar(&stepConfig.NpmInstallOpts, "npmInstallOpts", []string{}, "List of additional installation options for the npm install call. `-g`, `--global` is always assumed. Can be used for e.g. providing custom registries (`--registry https://your.registry.com`) or for providing the verbose flag (`--verbose`) for troubleshooting")

	cmd.MarkFlagRequired("username")
	cmd.MarkFlagRequired("password")
	cmd.MarkFlagRequired("applicationName")
	cmd.MarkFlagRequired("abapPackage")
	cmd.MarkFlagRequired("transportRequestId")
}

// retrieve step metadata
func transportRequestUploadCTSMetadata() config.StepData {
	var theMetaData = config.StepData{
		Metadata: config.StepMetadata{
			Name:        "transportRequestUploadCTS",
			Aliases:     []config.Alias{{Name: "transportRequestUploadFile", Deprecated: false}},
			Description: "This step uploads an UI5 application to the SAPUI5 ABAP repository.",
		},
		Spec: config.StepSpec{
			Inputs: config.StepInputs{
				Secrets: []config.StepSecrets{
					{Name: "uploadCredentialsId", Description: "Jenkins 'Username with password' credentials ID containing user and password to authenticate against the ABAP system.", Type: "jenkins", Aliases: []config.Alias{{Name: "changeManagement/credentialsId", Deprecated: false}}},
				},
				Parameters: []config.StepParameters{
					{
						Name:        "description",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS", "GENERAL"},
						Type:        "string",
						Mandatory:   false,
						Aliases:     []config.Alias{{Name: "applicationDescription"}},
						Default:     `Deployed with Piper based on SAP Fiori tools`,
					},
					{
						Name:        "endpoint",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS", "GENERAL"},
						Type:        "string",
						Mandatory:   false,
						Aliases:     []config.Alias{{Name: "changeManagement/endpoint"}, {Name: "changeManagement/cts/endpoint"}},
						Default:     os.Getenv("PIPER_endpoint"),
					},
					{
						Name:        "client",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS", "GENERAL"},
						Type:        "string",
						Mandatory:   false,
						Aliases:     []config.Alias{{Name: "changeManagement/client"}, {Name: "changeManagement/cts/client"}},
						Default:     os.Getenv("PIPER_client"),
					},
					{
						Name:        "username",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS", "GENERAL"},
						Type:        "string",
						Mandatory:   true,
						Aliases:     []config.Alias{},
						Default:     os.Getenv("PIPER_username"),
					},
					{
						Name:        "password",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS"},
						Type:        "string",
						Mandatory:   true,
						Aliases:     []config.Alias{},
						Default:     os.Getenv("PIPER_password"),
					},
					{
						Name:        "applicationName",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS", "GENERAL"},
						Type:        "string",
						Mandatory:   true,
						Aliases:     []config.Alias{},
						Default:     os.Getenv("PIPER_applicationName"),
					},
					{
						Name:        "abapPackage",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS", "GENERAL"},
						Type:        "string",
						Mandatory:   true,
						Aliases:     []config.Alias{},
						Default:     os.Getenv("PIPER_abapPackage"),
					},
					{
						Name:        "osDeployUser",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS", "GENERAL"},
						Type:        "string",
						Mandatory:   false,
						Aliases:     []config.Alias{{Name: "changeManagement/osDeployUser"}, {Name: "changeManagement/cts/osDeployUser"}},
						Default:     `node`,
					},
					{
						Name:        "deployConfigFile",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS", "GENERAL"},
						Type:        "string",
						Mandatory:   false,
						Aliases:     []config.Alias{{Name: "changeManagement/deployConfigFile"}, {Name: "changeManagement/cts/deployConfigFile"}},
						Default:     `ui5-deploy.yaml`,
					},
					{
						Name: "transportRequestId",
						ResourceRef: []config.ResourceReference{
							{
								Name:  "commonPipelineEnvironment",
								Param: "custom/transportRequestId",
							},
						},
						Scope:     []string{"PARAMETERS"},
						Type:      "string",
						Mandatory: true,
						Aliases:   []config.Alias{},
						Default:   os.Getenv("PIPER_transportRequestId"),
					},
					{
						Name:        "deployToolDependencies",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS", "GENERAL"},
						Type:        "[]string",
						Mandatory:   false,
						Aliases:     []config.Alias{{Name: "changeManagement/deployToolDependencies"}, {Name: "changeManagement/cts/deployToolDependencies"}},
						Default:     []string{`@ui5/cli`, `@sap/ux-ui5-tooling`, `@ui5/logger`, `@ui5/fs`},
					},
					{
						Name:        "npmInstallOpts",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS", "GENERAL"},
						Type:        "[]string",
						Mandatory:   false,
						Aliases:     []config.Alias{{Name: "changeManagement/npmInstallOpts"}, {Name: "changeManagement/cts/npmInstallOpts"}},
						Default:     []string{},
					},
				},
			},
			Containers: []config.Container{
				{Name: "fiori-client", Image: "node"},
			},
			Outputs: config.StepOutputs{
				Resources: []config.StepResources{
					{
						Name: "commonPipelineEnvironment",
						Type: "piperEnvironment",
						Parameters: []map[string]interface{}{
							{"name": "custom/transportRequestId"},
						},
					},
				},
			},
		},
	}
	return theMetaData
}
