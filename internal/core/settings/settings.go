package settings

import (
	"gopkg.in/yaml.v3"
	"slv.sh/slv/internal/core/commons"
)

type Settings struct {
	path *string
	*manifest
}

type manifest struct {
	AllowChanges      bool `yaml:"allow_changes"`
	AllowCreateEnv    bool `yaml:"allow_create_env"`
	AllowCreateGroup  bool `yaml:"allow_create_group"`
	SyncInterval      int  `yaml:"sync_interval"`
	AllowGroups       bool `yaml:"allow_groups"`
	AllowVaultSharing bool `yaml:"allow_vault_sharing"`
}

func (settings Settings) MarshalYAML() (any, error) {
	return settings.manifest, nil
}

func (settings *Settings) UnmarshalYAML(value *yaml.Node) (err error) {
	return value.Decode(&settings.manifest)
}

func NewManifest(path string) (settings *Settings, err error) {
	if commons.FileExists(path) {
		return nil, errManifestPathExistsAlready
	}
	settings = &Settings{
		path:     &path,
		manifest: &manifest{},
	}
	return
}

func GetManifest(path string) (settings *Settings, err error) {
	if !commons.FileExists(path) {
		return nil, errManifestNotFound
	}
	settings = &Settings{}
	if err = commons.ReadFromYAML(path, settings); err != nil {
		return nil, err
	}
	settings.path = &path
	return
}
