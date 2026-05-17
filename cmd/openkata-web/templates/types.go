package templates

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
	Downloads       int
	Docs            string
	Changelog       string
	Acknowledgments string
	Files           []string
}
