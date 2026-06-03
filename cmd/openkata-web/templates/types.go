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

// ArtifactDetail is implemented by both SkillDetail and RuleDetail.
type ArtifactDetail interface {
	ArtifactName() string
	ArtifactVersion() string
}

func (s SkillDetail) ArtifactName() string { return s.Name }
func (s SkillDetail) ArtifactVersion() string { return s.Version }
func (r RuleDetail) ArtifactName() string  { return r.Name }
func (r RuleDetail) ArtifactVersion() string { return r.Version }
