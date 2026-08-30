package serializer

import (
	"encoding/gob"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/zafir0101/ssienv/internal/domain"
)

const appDirName = "ssienv"

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

func Serialize(label string, obj domain.Controller) error {
	path, err := binPath(label)
	if err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	enc := gob.NewEncoder(file)
	return enc.Encode(obj)
}

func Deserialize(label string) (domain.Controller, bool, error) {
	path, err := binPath(label)
	if err != nil {
		return nil, false, err
	}

	file1, err := os.Open(path)
	file2, err := os.Open(path)
	if err != nil {
		return nil, false, err
	}
	defer file1.Close()
	defer file2.Close()

	dec1 := gob.NewDecoder(file1)
	dec2 := gob.NewDecoder(file2)

	ins_controller := &domain.InstitutionController{}
	ind_controller := &domain.IndividualController{}

	if err := dec1.Decode(ins_controller); err != nil {
		return nil, false, err
	}
	if err := dec2.Decode(ind_controller); err != nil {
		return nil, false, err
	}

	IsInstitutional := ins_controller.CloudAgentAPI != nil
	if IsInstitutional {
		return ins_controller, IsInstitutional, nil
	}
	return ind_controller, IsInstitutional, nil
}

func WithMutateCommand(cmd *cobra.Command, fn func(coData ControllerData) error) error {
	controllerLabel, err := cmd.Flags().GetString("controller")
	if err != nil {
		return err
	}

	controller, IsInstitutional, err := Deserialize(controllerLabel)
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

	if err := Serialize(controllerLabel, controller); err != nil {
		return err
	}

	return nil
}

func WithPureCommand(cmd *cobra.Command, fn func(cmdData ControllerData) error) error {
	controllerLabel, err := cmd.Flags().GetString("controller")
	if err != nil {
		return err
	}

	controller, IsInstitutional, err := Deserialize(controllerLabel)
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

	return nil
}

type ControllerData struct {
	Controller      domain.Controller
	IsInstitutional bool
}
