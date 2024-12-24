package cmd

import (
	"fmt"
	"os"

	"github.com/liyujun-dev/piper/config"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var configFile string

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration",
}

var viewConfigCmd = &cobra.Command{
	Use:   "view",
	Short: "View current configuration",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig(configFile)
		if err != nil {
			fmt.Println("Error loading config:", err)
			return
		}
		fmt.Printf("Current Profile: %s\n", cfg.CurrentProfile)
		fmt.Println("Configured Profiles:")
		for _, profile := range cfg.Profiles {
			fmt.Printf("- Name: %s, Provider: %s, Server: %s\n", profile.Name, profile.Provider, profile.Server)
		}
	},
}

var useProfileCmd = &cobra.Command{
	Use:   "use-profile [profile-name]",
	Short: "Set the current profile",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		profileName := args[0]
		cfg, err := config.LoadConfig(configFile)
		if err != nil {
			fmt.Println("Error loading config:", err)
			return
		}

		// 检查 profile 是否存在
		var profileFound bool
		for _, profile := range cfg.Profiles {
			if profile.Name == profileName {
				profileFound = true
				break
			}
		}

		if !profileFound {
			fmt.Printf("Profile '%s' not found\n", profileName)
			return
		}

		// 更新 current-profile
		cfg.CurrentProfile = profileName
		data, err := yaml.Marshal(cfg)
		if err != nil {
			fmt.Println("Error marshaling config:", err)
			return
		}

		err = os.WriteFile(configFile, data, 0644)
		if err != nil {
			fmt.Println("Error writing config file:", err)
			return
		}

		fmt.Printf("Switched to profile '%s'\n", profileName)
	},
}

var currentProfileCmd = &cobra.Command{
	Use:   "current-profile",
	Short: "Display the current profile",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig(configFile)
		if err != nil {
			fmt.Println("Error loading config:", err)
			return
		}
		fmt.Printf("Current Profile: %s\n", cfg.CurrentProfile)
	},
}
var addProfileCmd = &cobra.Command{
	Use:   "add-profile [name] [provider] [token] [server]",
	Short: "Add a new profile",
	Args:  cobra.ExactArgs(4),
	Run: func(cmd *cobra.Command, args []string) {
		profile := config.Profile{
			Name:     args[0],
			Provider: args[1],
			Token:    args[2],
			Server:   args[3],
		}

		cfg, err := config.LoadConfig(configFile)
		if err != nil {
			fmt.Println("Error loading config:", err)
			return
		}

		if err := config.AddProfile(cfg, profile); err != nil {
			fmt.Println(err)
			return
		}

		if err := config.SaveConfig(configFile, cfg); err != nil {
			fmt.Println("Error saving config:", err)
			return
		}

		fmt.Printf("Added profile '%s'\n", profile.Name)
	},
}

var removeProfileCmd = &cobra.Command{
	Use:   "remove-profile [profile-name]",
	Short: "Remove a profile",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig(configFile)
		if err != nil {
			fmt.Println("Error loading config:", err)
			return
		}

		if err := config.RemoveProfile(cfg, args[0]); err != nil {
			fmt.Println(err)
			return
		}

		if err := config.SaveConfig(configFile, cfg); err != nil {
			fmt.Println("Error saving config:", err)
			return
		}

		fmt.Printf("Removed profile '%s'\n", args[0])
	},
}

var listProfilesCmd = &cobra.Command{
	Use:   "list-profiles",
	Short: "List all profiles",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig(configFile)
		if err != nil {
			fmt.Println("Error loading config:", err)
			return
		}

		profiles := config.ListProfiles(cfg)
		fmt.Println("Configured Profiles:")
		for _, profile := range profiles {
			fmt.Printf("- Name: %s, Provider: %s, Server: %s\n", profile.Name, profile.Provider, profile.Server)
		}
	},
}

func init() {
	configCmd.PersistentFlags().StringVar(&configFile, "config", "config.yaml", "Path to the config file")
	configCmd.AddCommand(viewConfigCmd)
	configCmd.AddCommand(useProfileCmd)
	configCmd.AddCommand(currentProfileCmd)
	configCmd.AddCommand(addProfileCmd)
	configCmd.AddCommand(removeProfileCmd)
	configCmd.AddCommand(listProfilesCmd)
	rootCmd.AddCommand(configCmd)
}
