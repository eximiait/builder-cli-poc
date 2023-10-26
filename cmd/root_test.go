package cmd

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestRootCmd(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("builder-cli", rootCmd.Use, "rootCmd.Use should be 'builder-cli'")
	assert.Equal("CLI para poder crear y gestionar proyectos GitOps", rootCmd.Short, "rootCmd.Short should be 'CLI para poder crear y gestionar proyectos GitOps'")
	assert.Equal("CLI para poder crear y gestionar proyectos GitOps", rootCmd.Long, "rootCmd.Long should be 'CLI para poder crear y gestionar proyectos GitOps'")
}

func TestRootCmdFlags(t *testing.T) {
	assert := assert.New(t)

	// Check verbose flag
	assert.False(viper.GetBool("verbose"), "verbose flag default value should be false")

	// Check debug flag
	assert.False(viper.GetBool("debug"), "debug flag default value should be false")

	// Check gitlab-host-origin flag
	assert.Equal("https://gitlab.com", viper.GetString("gitlabHostOrigin"), "gitlabHostOrigin flag default value should be 'https://gitlab.com'")

	// Check gitlab-host flag
	assert.Equal("https://gitlab.com", viper.GetString("gitlabHost"), "gitlabHost flag default value should be 'https://gitlab.com'")
}
