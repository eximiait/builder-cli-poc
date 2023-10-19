/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/eximiait/builder-cli/creates"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var startersGroupName = "starters4273342"

type Repository struct {
	Name   string `json:"name"`
	WebURL string `json:"web_url"`
}

func getReposForGroup(groupsApiURL, groupName, token string) map[string]string {
	// Mapa para almacenar los resultados
	repoMap := make(map[string]string)

	// Construyendo la URL
	url := fmt.Sprintf("%s/%s/projects", groupsApiURL, groupName)

	// Crear la solicitud HTTP
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("Error al crear la solicitud: %v", err)
	}

	req.Header.Set("PRIVATE-TOKEN", token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Error al hacer la solicitud: %v", err)
	}
	defer resp.Body.Close()

	// Leer el cuerpo de la respuesta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error al leer la respuesta: %v", err)
	}

	// Deserializar el JSON
	var repos []Repository
	if err := json.Unmarshal(body, &repos); err != nil {
		log.Fatalf("Error al deserializar el JSON: %v", err)
	}

	// Rellenar el mapa
	for _, repo := range repos {
		repoMap[repo.Name] = repo.WebURL
	}

	return repoMap
}

func selectApplicationType(repos map[string]string) string {
	// Mostrar todos los nombres de repositorios
	keys := make([]string, 0, len(repos))
	for k := range repos {
		keys = append(keys, k)
	}

	fmt.Println("Por favor, selecciona un repositorio:")
	for i, key := range keys {
		fmt.Printf("%d. %s\n", i+1, key)
	}

	// Leer selección del usuario
	var choice int
	_, err := fmt.Scan(&choice)
	if err != nil || choice < 1 || choice > len(keys) {
		log.Fatalf("Selección fuera de rango. Por favor, selecciona un número entre 1 y %d.", len(keys))
		os.Exit(1)
	}

	// Devolver la URL del repositorio elegido
	return repos[keys[choice-1]]
}

func getAccessToken() string {
	byteToken, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Println("\nError al leer el Access Token.")
		os.Exit(1)
	}
	return string(byteToken)
}

var createCodeRepository = &cobra.Command{
	Use:   "createCodeRepository",
	Short: "Se crea un repositorio de código con el scaffolding de la aplicación",
	Long:  `Se crea un repositorio de código con el scaffolding de la aplicación`,
	Run: func(cmd *cobra.Command, args []string) {

		fmt.Print("Introduzca el Access Token provisto para clonar un repo starter: ")
		token := getAccessToken()
		groupsApiURL := GitlabHost + "/api/v4/groups"
		repoMap := getReposForGroup(groupsApiURL, startersGroupName, token)

		repoUrl := selectApplicationType(repoMap)

		creates.CreateCodeRepository(GitlabHost, repoUrl, token)
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
