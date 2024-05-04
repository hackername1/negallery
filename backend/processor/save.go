package processor

import (
	"os"
)

func (img *Image) Save(path string) error {
	err := os.WriteFile(path, img.ImageBytes, 0644)
	if err != nil {
		return err
	}
	return nil
}
