package utils

type MenuEntry struct {
	Id        string         `yaml:"Id"`
	Name      string         `yaml:"Name"`
	Submenues []SubmenuEntry `yaml:"Submenues"`
}

type SubmenuEntry struct {
	Id       string `yaml:"Id"`
	Name     string `yaml:"Name"`
	Function string `yaml:"Function"`
}

type Component struct {
	Name   string
	Parent string
	Code   []string
}

type GitHubApiRef struct {
	Ref    string            `json:"ref"`
	NodeId string            `json:"node_id"`
	Url    string            `json:"url"`
	Object map[string]string `json:"object"`
}
