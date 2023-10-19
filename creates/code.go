/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package creates

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/eximiait/builder-cli/gitlab"
)

func cloneRepo(repoURL, token, tmpDir string) error {
	// Inserta el token en la URL
	parsedURL := strings.Replace(repoURL, "https://", "https://oauth2:"+token+"@", 1)
	os.RemoveAll(tmpDir) // Limpiar antes de su uso

	cmd := exec.Command("git", "clone", parsedURL, tmpDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func CreateCodeRepository(gitlabHost, originRepoUrl, token, targetGroupName, targetRepoName string) {

	// Limpia los espacios y saltos de línea de las cadenas
	originRepoUrl = strings.TrimSpace(originRepoUrl)
	fmt.Println()
	tmpDir := "tmp"
	if err := cloneRepo(originRepoUrl, token, tmpDir); err != nil {
		fmt.Printf("\nError al clonar el repositorio: %v\n", err)
	}

	gitlab.PushToNewCleanRepo(tmpDir, gitlabHost, targetGroupName, targetRepoName, token)

}
