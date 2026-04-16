package tasks

import (
	"fmt"
)

func PreWeigh() error {
	configs, _, err := discoverBuilds(solutionsDir)
	if err != nil {
		return err
	}

	var weighLines []string
	for _, c := range configs {
		weighLines = append(weighLines, fmt.Sprintf("docker image inspect %s --format {{.Size}}", c.Tag))
	}

	if err := saveScript(weighScript, weighLines, false); err != nil {
		return fmt.Errorf("failed to save weigh script: %w", err)
	}

	fmt.Printf("Pre-weigh complete. Generated weigh.sh with %d targets.\n", len(configs))
	return nil
}
