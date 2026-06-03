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

type SkillDetail struct {
	Name            string
	Version         string
	Description     string
	Tags            string
	Versions        []string
	Downloads       int
	Docs            string
	Changelog       string
	Acknowledgments string
	Files           []string
	FileContents    map[string]string
}

type RuleDetail struct {
	Name            string
	Version         string
	Description     string
	Tags            string
	Versions        []string
	Downloads       int
	Docs            string
	Changelog       string
	Acknowledgments string
	Files           []string
	FileContents    map[string]string
}
