package core

type ModrinthJSON struct {
	Name          string         `json:"name"`
	VersionNumber string         `json:"version_number"`
	GameVersions  []string       `json:"game_versions"`
	Files         []ModrinthFile `json:"files"`
}

type ModrinthFile struct {
	Url      string `json:"url"`
	Filename string `json:"filename"`
	Primary  bool   `json:"primary"`
}
