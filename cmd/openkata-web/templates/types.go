package templates

type TabDef struct {
	ID    string
	Label string
}

type SkillEntry struct {
	Name          string
	Version       string
	Description   string
	Tags          string
	Downloads     int
	BestScore     int
	BestModel     string
	HasBenchmarks bool
	Models        []BenchmarkModel
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
	Models          []BenchmarkModel
	TesslScore      int
	Published       bool
}

type BenchmarkModel struct {
	ID            string
	Label         string
	Effectiveness int
	Scenarios     []BenchmarkScenario
}

type BenchmarkScenario struct {
	Name        string
	Description string
	Pass        bool
}

type StatsData struct {
	Empty          bool
	TotalDownloads int
	Artifacts      []ArtifactStats
	Types          []TypeStats
	Clients        []ClientStats
	Countries      []CountryStats
	Sources        []SourceStats
	PageLoads      int
	HumanPageLoads int
	BotPageLoads   int
	PagePaths      []PathStats
	Events         []DownloadEvent
	PageMetrics    []DailyMetric
}

type StatsDetailData struct {
	Artifact  string
	Version   string
	Versions  []string
	Total     int
	Clients   []ClientStats
	Countries []CountryStats
	Events    []DownloadEvent
}

type ArtifactStats struct {
	Name      string
	Type      string
	Downloads int
}

type TypeStats struct {
	Type      string
	Downloads int
}

type ClientStats struct {
	Client    string
	Downloads int
}

type CountryStats struct {
	Country   string
	Downloads int
}

type SourceStats struct {
	Source    string
	Downloads int
}

type DailyMetric struct {
	Date        string `json:"date"`
	Invocations int    `json:"invocations"`
}

type PathStats struct {
	Path  string
	Type  string
	Count int
	Bot   bool
}

type DownloadEvent struct {
	Artifact  string `json:"artifact"`
	Version   string `json:"version"`
	Source    string `json:"source"`
	Client    string `json:"client"`
	Country   string `json:"country"`
	Timestamp string `json:"timestamp"`
}
