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
	"strings"

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

func readTargetRepositoryData() (map[string]string, error) {
	var groupName string
	var repoName string

	fmt.Print("Nombre del grupo destino: ")
	_, err := fmt.Scan(&groupName)
	if err != nil {
		fmt.Println("Error al leer el nombre del grupo destino:", err)
		os.Exit(1)
	}

	fmt.Print("Nombre del repositorio destino: ")
	_, err = fmt.Scan(&repoName)
	if err != nil {
		fmt.Println("Error al leer el nombre del repositorio destino:", err)
		os.Exit(1)
	}

	// Limpiamos los caracteres de nueva línea de los strings leídos
	groupName = strings.TrimSpace(groupName)
	repoName = strings.TrimSpace(repoName)

	// Se pide el token de acceso para el repo destino
	fmt.Print("Access Token destino: ")
	targetToken := getAccessToken()

	return map[string]string{
		"groupName":   groupName,
		"repoName":    repoName,
		"targetToken": targetToken,
	}, nil
}

var createCodeRepository = &cobra.Command{
	Use:   "createCodeRepository",
	Short: "A code repository is created with the scaffolding of the application",
	Long:  `A code repository is created with the scaffolding of the application`,
	Run: func(cmd *cobra.Command, args []string) {

		fmt.Print("Introduzca el Access Token provisto para clonar un repo starter: ")
		token := getAccessToken()
		groupsApiURL := GitlabHostOrigin + "/api/v4/groups"
		repoMap := getReposForGroup(groupsApiURL, startersGroupName, token)

		originRepoUrl := selectApplicationType(repoMap)

		fmt.Println()
		result, err := readTargetRepositoryData()
		if err != nil {
			fmt.Printf("Hubo un error: %v\n", err)
			return
		}

		creates.CreateCodeRepository(GitlabHostTarget, originRepoUrl, result["targetToken"], result["groupName"], result["repoName"])
	},
}

var createGitOpsRepository = &cobra.Command{
	Use:   "createGitOpsRepository",
	Short: "A GitOps repository is created with the scaffolding of the environment",
	Long:  `A GitOps repository is created with the scaffolding of the environment`,
	Run: func(cmd *cobra.Command, args []string) {
		creates.CreateGitOpsRepository()
	},
}

func init() {
	rootCmd.AddCommand(createCodeRepository)
	rootCmd.AddCommand(createGitOpsRepository)
}
