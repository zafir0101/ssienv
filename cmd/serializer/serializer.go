package serializer

import (
	"encoding/gob"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/zafir0101/ssienv/internal/domain"
)

const appDirName = "ssienv"

func init() {
	gob.Register(&domain.IndividualController{})
	gob.Register(&domain.InstitutionController{})
}

func dataDir() (string, error) {
	baseDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	dir := filepath.Join(baseDir, appDirName, "data")

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}

	return dir, nil
}

func binPath(label string) (string, error) {
	dir, err := dataDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, label+".bin"), nil
}

func serialize(label string, obj domain.Controller) error {
	path, err := binPath(label)
	if err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return gob.NewEncoder(file).Encode(&obj)
}

func deserialize(label string) (domain.Controller, bool, error) {
	path, err := binPath(label)
	if err != nil {
		return nil, false, err
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, false, err
	}
	defer file.Close()

	var controller domain.Controller
	if err := gob.NewDecoder(file).Decode(&controller); err != nil {
		return nil, false, err
	}

	_, IsInstitutional := controller.(*domain.InstitutionController)

	return controller, IsInstitutional, nil
}

func WithMutateCommand(cmd *cobra.Command, fn func(coData ControllerData) error) error {
	controllerLabel, err := cmd.Flags().GetString("controller")
	if err != nil {
		return err
	}

	controller, IsInstitutional, err := deserialize(controllerLabel)
	if err != nil {
		return err
	}

	cmdData := ControllerData{
		Controller:      controller,
		IsInstitutional: IsInstitutional,
	}

	if err := fn(cmdData); err != nil {
		return err
	}

	return serialize(controllerLabel, controller)
}

func WithPureCommand(cmd *cobra.Command, fn func(cmdData ControllerData) error) error {
	controllerLabel, err := cmd.Flags().GetString("controller")
	if err != nil {
		return err
	}

	controller, IsInstitutional, err := deserialize(controllerLabel)
	if err != nil {
		return err
	}

	cmdData := ControllerData{
		Controller:      controller,
		IsInstitutional: IsInstitutional,
	}

	return fn(cmdData)
}

type ControllerData struct {
	Controller      domain.Controller
	IsInstitutional bool
}
