package custom_executable

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/codecrafters-io/shell-tester/internal/custom_executable/completer/completer_configuration"
)

type CompleterExecutableSpecification struct {
	Path                   string
	SecretValue            string
	CompleterConfiguration completer_configuration.CompleterConfiguration
}

func (s *CompleterExecutableSpecification) Create() error {
	if err := s.verify(); err != nil {
		return err
	}

	executableName := "completer"
	err := createExecutableForOSAndArch(executableName, s.Path)
	if err != nil {
		return fmt.Errorf("CodeCrafters Internal Error: copying executable failed: %w", err)
	}

	// Add secret to executable
	err = addSecretCodeToExecutable(s.Path, s.SecretValue)
	if err != nil {
		return fmt.Errorf("CodeCrafters Internal Error: adding secret code to executable failed: %w", err)
	}

	configBytes, err := json.Marshal(s.CompleterConfiguration)
	if err != nil {
		return fmt.Errorf("CodeCrafters Internal Error: marshal completer config failed: %w", err)
	}

	// Write the configuration file with the same secret name so
	// the completer can retrieve it later
	configPath := filepath.Join("/tmp", s.SecretValue)
	if err := os.WriteFile(configPath, configBytes, 0644); err != nil {
		return fmt.Errorf("CodeCrafters Internal Error: write completer config failed: %w", err)
	}

	if err := reSignExecutableDarwinArm64(s.Path); err != nil {
		return err
	}

	return nil
}

func (s *CompleterExecutableSpecification) verify() error {
	if len(s.SecretValue) != 10 {
		return fmt.Errorf("CodeCrafters Internal Error: CompleterExecutableSpecification.SecretValue length must be 10")
	}

	if !filepath.IsAbs(s.Path) {
		return fmt.Errorf("Codecrafters Internal Error: CompleterExecutableSpecification.Path must be absolute")
	}

	if err := s.CompleterConfiguration.Verify(); err != nil {
		return err
	}

	return nil
}
