/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/eximiait/builder-cli/creates"
	"github.com/spf13/cobra"
)

// Mapping de lenguajes a sus respectivas URLs
var languageURLMap = map[string]string{
	"java":   "https://gitlab.com/demo-cicd1473038/app-code",
	"dotnet": "https://gitlab.com/demo-cicd1473038/app-code-dotnet",
	"vuejs":  "https://gitlab.com/demo-cicd1473038/app-code-vuejs", // Añade aquí la URL para vuejs si la tienes.
}

var createCodeRepository = &cobra.Command{
	Use:   "createCodeRepository",
	Short: "Se crea un repositorio de código con el scaffolding de la aplicación",
	Long:  `Se crea un repositorio de código con el scaffolding de la aplicación`,
	Run: func(cmd *cobra.Command, args []string) {

		// Obtener lenguajes desde el mapa
		languages := make([]string, 0, len(languageURLMap))
		for lang := range languageURLMap {
			languages = append(languages, lang)
		}

		// Mostrar opciones de lenguaje
		fmt.Println("Por favor, selecciona un lenguaje:")
		for i, lang := range languages {
			fmt.Printf("%d. %s\n", i+1, lang)
		}

		// Leer selección del usuario
		var choice int
		_, err := fmt.Scan(&choice)
		if err != nil || choice < 1 || choice > len(languages) {
			fmt.Println("Selección inválida")
			os.Exit(1)
		}

		selectedLanguage := languages[choice-1]
		fmt.Printf("Has seleccionado: %s\n", selectedLanguage)
		fmt.Printf("URL asociada: %s\n", languageURLMap[selectedLanguage])
		creates.CreateCodeRepository(GitlabHost, selectedLanguage, languageURLMap[selectedLanguage])
	},
}

var createGitOpsRepository = &cobra.Command{
	Use:   "createGitOpsepository",
	Short: "Se crea un repositorio GitOps con el scaffolding del ambiente",
	Long:  `Se crea un repositorio GitOps con el scaffolding del ambiente`,
	Run: func(cmd *cobra.Command, args []string) {
		creates.CreateGitOpsRepository()
	},
}

func init() {
	rootCmd.AddCommand(createCodeRepository)
	rootCmd.AddCommand(createGitOpsRepository)
}

func getLanguageByIndex(index int) string {
	var i int = 1
	for lang := range languageURLMap {
		if i == index {
			return lang
		}
		i++
	}
	return ""
}
