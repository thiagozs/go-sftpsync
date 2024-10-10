package cmd

import (
	"log"
	"slices"

	"github.com/spf13/cobra"
	"github.com/thiagozs/go-sftpsync/internal/domain"
	"github.com/thiagozs/go-sftpsync/internal/sftpsync"
	"github.com/thiagozs/go-sftpsync/pkg/utils"
)

var (
	host           string
	port           string
	user           string
	password       string
	operation      string
	localBasePath  string
	remoteBasePath string
	publicKey      string
	privateKey     string
	configFile     string
	profileName    string
	listProfiles   bool
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Syncronise files between local and remote directories",
	// 	Long: `A longer description that spans multiple lines and likely contains examples
	// and usage of using your command. For example:

	// Cobra is a CLI library for Go that empowers applications.
	// This application is a tool to generate the needed files
	// to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		opts := []sftpsync.Options{}

		// Config file settings
		if configFile != "" {
			config, err := utils.ReadConfig(configFile)
			if err != nil {
				log.Fatalf("failed to read config file: %v", err)
				return
			}

			// Profile settings list
			if listProfiles {
				profiles := []string{}
				for profile := range config.Profiles {
					profiles = append(profiles, profile)
				}

				slices.Sort(profiles)

				for _, profile := range profiles {
					log.Println(profile)
				}

				return
			}

			// Executing profile settings
			if profileName == "" {
				log.Println("profile name is required")
				return
			}

			pcfg := domain.ProfileConfig{}
			for profile, cfg := range config.Profiles {
				if profile == profileName {
					pcfg = cfg
					break
				}
			}

			opts = AddOptionIfNotEmpty(opts, sftpsync.WithHost, pcfg.Host)
			opts = AddOptionIfNotEmpty(opts, sftpsync.WithPort, pcfg.Port)
			opts = AddOptionIfNotEmpty(opts, sftpsync.WithUser, pcfg.User)
			opts = AddOptionIfNotEmpty(opts, sftpsync.WithPassword, pcfg.Password)
			opts = AddOptionIfNotEmpty(opts, sftpsync.WithPrivateKey, pcfg.PrivateKey)
			opts = AddOptionIfNotEmpty(opts, sftpsync.WithLocalBasePath, pcfg.LocalBasePath)
			opts = AddOptionIfNotEmpty(opts, sftpsync.WithRemoteBasePath, pcfg.RemoteBasePath)

			operation = pcfg.Operation
		}

		if host != "" && port != "" && user != "" {
			opts = AddOptionIfNotEmpty(opts, sftpsync.WithHost, host)
			opts = AddOptionIfNotEmpty(opts, sftpsync.WithPort, port)
			opts = AddOptionIfNotEmpty(opts, sftpsync.WithUser, user)
			opts = AddOptionIfNotEmpty(opts, sftpsync.WithPassword, password)
			opts = AddOptionIfNotEmpty(opts, sftpsync.WithPrivateKey, privateKey)
			opts = AddOptionIfNotEmpty(opts, sftpsync.WithLocalBasePath, localBasePath)
			opts = AddOptionIfNotEmpty(opts, sftpsync.WithRemoteBasePath, remoteBasePath)
		}

		sync, err := sftpsync.NewSftpSync(opts...)
		if err != nil {
			log.Fatalf("failed to create SFTP sync: %v", err)
			return
		}

		opkind := sftpsync.NewOperationKind()

		switch opkind.GetFromString(operation) {
		case sftpsync.UPLOAD:
			if err := sync.SyncUpload(); err != nil {
				log.Printf("failed to sync upload: %v", err)
				return
			}
		case sftpsync.DOWNLOAD:
			if err := sync.SyncDownload(); err != nil {
				log.Printf("failed to sync download: %v", err)
				return
			}
		default:
			log.Println("operation is required")
			return
		}

	},
}

func init() {
	rootCmd.AddCommand(syncCmd)

	syncCmd.PersistentFlags().StringVarP(&operation, "operation", "o", "", "Operation to perform: upload or download")
	syncCmd.PersistentFlags().StringVarP(&host, "host", "H", "", "SFTP host")
	syncCmd.PersistentFlags().StringVarP(&port, "port", "P", "", "SFTP port")
	syncCmd.PersistentFlags().StringVarP(&user, "user", "u", "", "SFTP user")
	syncCmd.PersistentFlags().StringVarP(&password, "password", "p", "", "SFTP password")
	syncCmd.PersistentFlags().StringVarP(&privateKey, "private-key", "k", "", "SFTP private key")
	syncCmd.PersistentFlags().StringVarP(&localBasePath, "local-base-path", "l", "", "Local base path")
	syncCmd.PersistentFlags().StringVarP(&remoteBasePath, "remote-base-path", "r", "", "Remote base path")
	syncCmd.PersistentFlags().StringVarP(&configFile, "config-file", "c", "", "Config file")
	syncCmd.PersistentFlags().StringVarP(&profileName, "profile", "f", "", "Profile name")
	syncCmd.PersistentFlags().BoolVarP(&listProfiles, "list-profiles", "L", false, "List profiles")
}

func AddOptionIfNotEmpty(opts []sftpsync.Options,
	option func(string) sftpsync.Options,
	value string) []sftpsync.Options {
	if value != "" {
		opts = append(opts, option(value))
	}
	return opts
}
