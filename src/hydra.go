/*
Code responsible for the hydra CLI.

Author: Shravan Asati
Originally Written: 27 March 2021
Last edited: 4 June 2021
*/

package main

import (
	"fmt"
	"github.com/thatisuday/commando"
	"regexp"
	"strings"
)

const (
	NAME    string = "hydra"
	VERSION string = "2.2.0"
)

var (
	supportedLangs    []string          = []string{"go", "python", "web", "flask", "ruby", "c", "c++"}
	supportedLicenses map[string]string = map[string]string{
		"APACHE": "Apache License",
		"BSD":    "Berkeley Software Distribution 3-Clause",
		"EPL":    "Eclipse Public License",
		"GPL":    "GNU General Public License v3",
		"MIT":    "Massachusetts Institute of Technology License",
		"MPL":    "Mozilla Public License",
		"UNI":    "The Unilicense"}
)

func stringInSlice(s string, slice []string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}
	return false
}

func wrongProjectName(projectName string) bool {
	match, _ := regexp.MatchString(`\.|\?|\*|\:|\,|\'|\"|\||\|\|/<|>`, projectName)
	return match
}

func main() {
	fmt.Println(NAME, VERSION)

	// * basic configuration
	commando.
		SetExecutableName(NAME).
		SetVersion(VERSION).
		SetDescription("hydra is command line utility used to generate language-specific project structure. \nFor more detailed information and documentation, visit https://github.com/shravanasati/hydra . \n")

	commando.
		Register(nil).
		SetAction(func(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {
			fmt.Println("\nExecute `hydra -h` for help.")
		})

	// * the list command
	commando.
		Register("list").
		SetShortDescription("Lists supported languages, licenses and user configurations.").
		SetDescription("Lists supported languages, licenses and user configurations.").
		AddArgument("item", "The item to list. Valid options are `langs`, `licenses` and `configs`.", "").
		SetAction(func(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {

			if args["item"].Value == "langs" {
				fmt.Println(list("langs"))
			} else if args["item"].Value == "licenses" {
				fmt.Println(list("licenses"))
			} else if args["item"].Value == "configs" {
				fmt.Println(list("configs"))
			} else {
				fmt.Println(list(args["item"].Value))
			}
		})

	commando.
		Register("config").
		SetShortDescription("Alter or set the hydra user configuration.").
		SetDescription("Alter or set the hydra user configuration.").
		AddFlag("name", "The user's full name.", commando.String, "default").
		AddFlag("github-username", "The user's GitHub username.", commando.String, "default").
		AddFlag("default-lang", "The user's default language for project initialisation.", commando.String, "default").
		AddFlag("default-license", "The user's default license for project initialisation.", commando.String, "default").
		SetAction(func(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {
			config(
				flags["name"].Value.(string), 
				flags["github-username"].Value.(string), 
				flags["default-lang"].Value.(string), 
				strings.ToUpper(flags["default-license"].Value.(string)))
		})

	// * the init command
	commando.
		Register("init").
		SetShortDescription("Intialises the project structure.").
		SetDescription("Intialises the project structure.\n\nUsage: \n name : project name \n lang : programming language in which the project is being built.").
		AddArgument("name", "Name of the project", "").
		AddArgument("lang", "Language/framework of the project. To view valid options for the this parameter, execute `hydra list langs`.", "default").
		AddFlag("license", "The license to initialise the project with.", commando.String, "default").
		SetAction(func(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {

			// * checking if user has properly configured hydra (full name and github username)
			if !checkForCorrectConfig() {
				fmt.Println("Error: You've not set your hydra configuration. You cannot proceed without setting the necessary configuration.\nTo set configuration, execute `hydra config --name \"YOUR NAME\" --github-username \"YOUR GITHUB USERNAME\"` .\nFor further assistance regarding hydra configuration, type in `hydra config -h` .")
				return
			}

			// * checking for correct license
			license := strings.ToUpper(flags["license"].Value.(string))
			if license == "DEFAULT" {
				license = getConfig("defaultLicense")
			}
			if !stringInSlice(license, []string{"MIT", "BSD", "MPL", "EPL", "GPL", "APACHE", "UNI"}) {
				fmt.Printf("Invalid value for flag license: '%v'.\n", license)
				fmt.Println("You've either provided invalid license flag in the init command, or you've set wrong license in your hydra configuration.\nTo see your hydra configuration, execute `hydra list configs`.")
				return
			}

			// * checking for correct project language
			projectLang := strings.ToLower(args["lang"].Value)
			if projectLang == "default" {
				projectLang = getConfig("defaultLang")
			}

			projectName := args["name"].Value

			// * checking the project name
			if wrongProjectName(projectName) {
				fmt.Printf(`Error: Invalid project name: '%v'. Characters like (, " | \ ? / : ; < >) are not allowed in filenames.`+"\n", projectName)
				return
			}

			init := Initializer{
				projectName: projectName,
				license:     license,
				lang:        projectLang,
			}
			switch projectLang {
			case "python":
				init.pythonInit()
			case "go":
				init.goInit()
			case "web":
				init.webInit()
			case "flask":
				init.flaskInit()
			case "c":
				init.cInit()
			case "c++":
				init.cppInit()
			case "ruby":
				init.rubyInit()
			default:
				fmt.Printf("Unsupported language type: '%v'. Cannot initiate the project. \nHint: You've either a typo at the language name, or the hydra default language configuration is wrong.", projectLang)
			}
		})

	commando.
		Register("update").
		SetShortDescription("The update command updates hydra to the latest release.").
		SetDescription("The update command downloads and installs the latest hydra release.").
		SetAction(func(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {
			update()
		})

	commando.Parse(nil)
	deletePreviousInstallation()
}
