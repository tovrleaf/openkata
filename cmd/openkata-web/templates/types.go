package templates

type TabDef struct {
	ID    string
	Label string
}

type SkillEntry struct {
	Name        string
	Version     string
	Description string
	Tags        string
	Downloads   int
}

type ArtifactDetail struct {
	Type            string // "skills", "rules", "profiles"
	Name            string
	Version         string
	Description     string
	Tags            string
	Versions        []string
	Downloads       int
	Docs            string
	Rationale       string
	Changelog       string
	Acknowledgments string
	Files           []string
	FileContents    map[string]string
	Prev            string // name of previous artifact (empty if first)
	Next            string // name of next artifact (empty if last)
}
