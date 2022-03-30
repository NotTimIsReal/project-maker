/*
Copyright © 2022 NotTimIsReal

*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
	git "gopkg.in/src-d/go-git.v4"

	"time"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
)

// buildCmd represents the build command
const defaultProjects string = `{
  "discord-bot-js": "https://github.com/YourBetterAssistant/yourbetterassistant",
  "python-chatbot": "https://github.com/YourBetterAssistant/chatbot"
}
`

var buildCmd = &cobra.Command{
	Use:       "build [project_name] [project-dir]",
	Aliases:   []string{"init", "b", "i"},
	Example:   "project-maker build js-cli -D project_dir",
	Args:      cobra.MinimumNArgs(2),
	ValidArgs: []string{"js-cli", "python-chatbot"},
	Short:     "Builds and initialised a new project",
	Long:      `Builds and insitialises a new project. Requrires dir name and project type`,
	Run: func(cmd *cobra.Command, args []string) {
		//print to console "Searching for %s in %s"
		exeLocation, err := os.Executable()
		exeLocation = strings.Join(strings.Split(exeLocation, "/")[:len(strings.Split(exeLocation, "/"))-1], "/")
		if err != nil {
			panic(err)
		}
		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond) // Build our new spinner
		s.Start()
		var res string = color.GreenString(fmt.Sprintf(" Searching for %s in the project folders...", args[0]))
		s.Suffix = res
		content, err := os.ReadFile(fmt.Sprintf("%s/projects.json", exeLocation))
		time.Sleep(1 * time.Second)
		if err != nil {
			s.Stop()
			os.WriteFile(fmt.Sprintf("%s/projects.json", exeLocation), []byte(defaultProjects), 0644)
			fmt.Print("A Project.json file had not been created, please re run this command now to fix this problem")
			os.Exit(1)
		}
		if len(content) == 0 {
			s.Stop()
			fmt.Printf("projects.json is malformed please delete it and try again. The Location is %s/projects.json", exeLocation)
			os.Exit(1)
		}
		s.Stop()
		//projects looks like {"something":"something"}
		var projects map[string]string
		json.Unmarshal(content, &projects)
		var repo = projects[args[0]]
		if repo == "" {
			fmt.Printf("%s is not a valid project Add it or use a valid project", args[0])
			os.Exit(1)
		}
		s = spinner.New(spinner.CharSets[9], 100*time.Millisecond) // Build our new spinner

		s.Start()
		s.Suffix = color.BlueString(fmt.Sprintf(" Found The Code Template At %s, Now Cloning...", repo))
		_, err = git.PlainClone(args[1], false, &git.CloneOptions{
			URL: repo,
		})
		if err != nil {
			s.Stop()
			fmt.Print(color.RedString(fmt.Sprintf("Error Cloning %s, this is likely due to a git repo already existing in %s", repo, args[1])))
			os.Exit(1)
		}
		s.Stop()
		fmt.Printf(color.GreenString("Project %s has been created at %s, Now Running Set Up Script", args[0], args[1]))
		s = spinner.New(spinner.CharSets[9], 100*time.Millisecond) // Build our new spinner
		s.Start()
		s.Suffix = color.GreenString(fmt.Sprintf(" Running Set Up Script... This Might Take A While"))
		//cd into k and run main.go
		_, err = os.ReadFile(fmt.Sprintf("%s/project-setup.py", args[1]))
		if err != nil {
			s.Stop()
			fmt.Print(color.RedString("Error: project-setup.py was not found in the project directory"))
			os.Exit(1)
		}
		//run and show output
		var command = exec.Command("python3", "project-setup.py")
		command.Dir = args[1]
		status, err := command.Output()
		if strings.Contains(string(status), "setup.project.status.sucess") {
			s.Stop()
			fmt.Print(color.GreenString("Project Setup Sucessful"))
			os.Exit(0)
		}
		if err != nil {
			s.Stop()
			fmt.Print(err)

			fmt.Print(color.RedString("Something happened while running setup-script, Make Sure You Have Made A Valid project-setup.py file"))
		}
		s.Stop()

	},
}

func init() {
	rootCmd.AddCommand(buildCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// buildCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// buildCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
