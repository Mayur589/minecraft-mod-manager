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

type Mod struct {
	ID          string
	Version     string
	Path        string
	Hash        string
	IsModrinth  bool
	Files       []ModrinthFile
	GameVersion []string
}

type jsonModFile struct {
	ID      string `json:"id"`
	Version string `json:"version"`
	Name    string `json:"name"`
}
