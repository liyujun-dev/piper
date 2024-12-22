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
		fmt.Printf("Current Context: %s\n", cfg.CurrentContext)
		fmt.Println("Configured Contexts:")
		for _, ctx := range cfg.Contexts {
			fmt.Printf("- Name: %s, Provider: %s, Server: %s\n", ctx.Name, ctx.Provider, ctx.Server)
		}
	},
}

var useContextCmd = &cobra.Command{
	Use:   "use-context [context-name]",
	Short: "Set the current context",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		contextName := args[0]
		cfg, err := config.LoadConfig(configFile)
		if err != nil {
			fmt.Println("Error loading config:", err)
			return
		}

		// 检查上下文是否存在
		var contextFound bool
		for _, ctx := range cfg.Contexts {
			if ctx.Name == contextName {
				contextFound = true
				break
			}
		}

		if !contextFound {
			fmt.Printf("Context '%s' not found\n", contextName)
			return
		}

		// 更新 current-context
		cfg.CurrentContext = contextName
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

		fmt.Printf("Switched to context '%s'\n", contextName)
	},
}

var currentContextCmd = &cobra.Command{
	Use:   "current-context",
	Short: "Display the current context",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig(configFile)
		if err != nil {
			fmt.Println("Error loading config:", err)
			return
		}
		fmt.Printf("Current Context: %s\n", cfg.CurrentContext)
	},
}
var addContextCmd = &cobra.Command{
	Use:   "add-context [name] [provider] [token] [server]",
	Short: "Add a new context",
	Args:  cobra.ExactArgs(4),
	Run: func(cmd *cobra.Command, args []string) {
		context := config.Context{
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

		if err := config.AddContext(cfg, context); err != nil {
			fmt.Println(err)
			return
		}

		if err := config.SaveConfig(configFile, cfg); err != nil {
			fmt.Println("Error saving config:", err)
			return
		}

		fmt.Printf("Added context '%s'\n", context.Name)
	},
}

var removeContextCmd = &cobra.Command{
	Use:   "remove-context [context-name]",
	Short: "Remove a context",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig(configFile)
		if err != nil {
			fmt.Println("Error loading config:", err)
			return
		}

		if err := config.RemoveContext(cfg, args[0]); err != nil {
			fmt.Println(err)
			return
		}

		if err := config.SaveConfig(configFile, cfg); err != nil {
			fmt.Println("Error saving config:", err)
			return
		}

		fmt.Printf("Removed context '%s'\n", args[0])
	},
}

var listContextsCmd = &cobra.Command{
	Use:   "list-contexts",
	Short: "List all contexts",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig(configFile)
		if err != nil {
			fmt.Println("Error loading config:", err)
			return
		}

		contexts := config.ListContexts(cfg)
		fmt.Println("Configured Contexts:")
		for _, ctx := range contexts {
			fmt.Printf("- Name: %s, Provider: %s, Server: %s\n", ctx.Name, ctx.Provider, ctx.Server)
		}
	},
}

func init() {
	configCmd.PersistentFlags().StringVar(&configFile, "config", "config.yaml", "Path to the config file")
	configCmd.AddCommand(viewConfigCmd)
	configCmd.AddCommand(useContextCmd)
	configCmd.AddCommand(currentContextCmd)
	configCmd.AddCommand(addContextCmd)
	configCmd.AddCommand(removeContextCmd)
	configCmd.AddCommand(listContextsCmd)
	rootCmd.AddCommand(configCmd)
}
