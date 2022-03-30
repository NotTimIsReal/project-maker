/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new [project_name] [project_url]",
	Short: "Create New Project",
	Long:  `Creates a new project. Requires a project name and a project URL with github or gitlab.`,
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		exeLocation, err := os.Executable()
		exeLocation = strings.Join(strings.Split(exeLocation, "/")[:len(strings.Split(exeLocation, "/"))-1], "/")
		projects, err := os.ReadFile(fmt.Sprintf("%s/projects.json", exeLocation))
		if err != nil {
			fmt.Println("Error reading file:", err)
			os.Exit(1)
		}
		var projectsMap map[string]string
		err = json.Unmarshal(projects, &projectsMap)
		if err != nil {
			fmt.Println("Error parsing file:", err)
			os.Exit(1)
		}
		if !IsUrl(args[1]) {
			fmt.Print("Please enter a valid URL, this uses git")
			os.Exit(1)
		}
		projectsMap[args[0]] = args[1]
		projects, err = json.Marshal(projectsMap)
		if err != nil {
			fmt.Println("Error making file:", err)
			os.Exit(1)
		}
		os.WriteFile(fmt.Sprintf("%s/projects.json", exeLocation), projects, 0644)
		fmt.Print(color.BlueString("Project Added To DB"))

	},
}

func init() {
	rootCmd.AddCommand(newCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// newCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// newCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
func IsUrl(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}
