package metadata

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type ProjectMetadata struct {
	ProjectID    string `json:"project_id"`
	ProjectName  string `json:"project_name"`
	IsAutoRecord bool   `json:"is_auto_record"`
}

func SaveMetadata(debugoDir string, projectID string, projectname string, isAutoRecord bool) error {
	metaData := ProjectMetadata{
		ProjectID:    projectID,
		ProjectName:  projectname,
		IsAutoRecord: isAutoRecord,
	}

	data, err := json.MarshalIndent(metaData, "", " ")
	if err != nil {
		return err
	}

	metaDataPath := filepath.Join(debugoDir, "metadata.json")

	return os.WriteFile(metaDataPath, data, 0644)
}

func LoadMetadata(debugoDir string) (*ProjectMetadata, error) {
	metadataPath := filepath.Join(debugoDir, "metadata.json")

	data, err := os.ReadFile(metadataPath)
	if err != nil {
		return nil, err
	}

	var metadata ProjectMetadata
	err = json.Unmarshal(data, &metadata)
	if err != nil {
		return nil, err
	}

	return &metadata, nil
}
