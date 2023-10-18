package gitlab

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strconv"

	"github.com/eximiait/builder-cli/common"
)

// Crea un nuevo repositorio en GitLab usando la API
func createGitLabRepo(gitlabHost, namespaceName, repoName, token string) (string, error) {
	projectApiURL := gitlabHost + "/api/v4/projects"
	namespaceID, err := getNamespaceID(gitlabHost, namespaceName, token)
	if err != nil {
		log.Fatalf("Error al obtener ID del namespace: %v", err)
	}

	data := map[string]string{
		"namespace_id": namespaceID,
		"name":         repoName,
	}

	payload, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error al procesar el payload : %v", err)
		return "", err
	}

	req, err := http.NewRequest("POST", projectApiURL, bytes.NewBuffer(payload))
	if err != nil {
		log.Printf("Error al hacer POST : %v", err)
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("PRIVATE-TOKEN", token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error reading GitLab response body: %v", err)
		} else {
			log.Printf("GitLab error response: %s", string(bodyBytes))
		}
		return "", fmt.Errorf("failed to create repo, status: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	// Retorna la URL git para el repositorio recién creado
	return result["http_url_to_repo"].(string), nil
}

// Inicializa un nuevo repositorio, lo conecta a GitLab y realiza un push
func PushToNewCleanRepo(localRepoPath, gitlabHost, namespaceName, newRepoName, token string) error {

	// Crear un directorio temporal
	tmpDir, err := os.MkdirTemp("", "newrepo")
	if err != nil {
		log.Printf("Error al crear directorio temporal: %v", err)
		return err
	}
	defer os.RemoveAll(tmpDir) // Limpiar después de su uso

	// Copiar contenido (excluyendo .git) al directorio temporal
	err = common.CopyDir(localRepoPath, tmpDir)
	if err != nil {
		log.Printf("Error al copiar directorio : %v", err)
		return err
	}

	// Inicializar un nuevo repositorio git en el directorio temporal
	cmd := exec.Command("git", "-C", tmpDir, "init")
	if err := cmd.Run(); err != nil {
		log.Printf("Error al inicializar el nuevo repositorio git : %v", err)
		return err
	}

	// Añadir todos los archivos y realizar el commit inicial
	cmd = exec.Command("git", "-C", tmpDir, "add", ".")
	if err := cmd.Run(); err != nil {
		log.Printf("Error al añadir los archivos a git (git add) : %v", err)
		return err
	}
	cmd = exec.Command("git", "-C", tmpDir, "commit", "-m", "Initial commit")
	if err := cmd.Run(); err != nil {
		log.Printf("Error hacer commit (git commit) : %v", err)
		return err
	}

	// Conectar el repositorio local al nuevo repositorio de GitLab y hacer push
	var newRepoURL string
	newRepoURL, err = createGitLabRepo(gitlabHost, namespaceName, newRepoName, token)
	log.Printf("Nueva URL repo: %v", newRepoURL)
	if err != nil {
		log.Printf("Error al ejecutar la creación del repositorio en Gitlab : %v", err)
		return err
	}

	// Modificar la URL del repositorio para incluir el token
	parsedURL, err := url.Parse(newRepoURL)
	if err != nil {
		return fmt.Errorf("Error al parsear la URL: %v", err)
	}
	parsedURL.User = url.UserPassword("oauth2", token)
	modifiedRepoURL := parsedURL.String()

	cmd = exec.Command("git", "remote", "add", "origin", modifiedRepoURL)
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Error al establecer nuevo remoto: %v", err)
	}

	cmd = exec.Command("git", "push", "-u", "origin", "main")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		log.Printf("Error al ejecutar push : %v", err)
		return fmt.Errorf("Error al hacer push al nuevo repositorio: %v", err)
	}

	return nil // No hubo errores
}

func getNamespaceID(gitlabHost, namespaceName, token string) (string, error) {
	apiURL := gitlabHost + "/api/v4/groups/" + namespaceName

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		log.Printf("Error al hacer GET de los grupos de Gitlab: %v", err)
		return "", err
	}

	req.Header.Set("PRIVATE-TOKEN", token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("failed to fetch namespace, status: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	namespaceID, ok := result["id"].(float64) // IDs are typically represented as float64 by the JSON parser
	if !ok {
		return "", fmt.Errorf("could not get namespace ID")
	}

	return strconv.Itoa(int(namespaceID)), nil
}
