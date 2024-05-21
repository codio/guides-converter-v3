package guidespaths

import (
	"os"
	"path/filepath"
	"sync"
)

type GuidesPaths struct {
	WorkSpace                  string
	GuidesFolder               string
	AssessmentsDescriptionFile string
	AssessmentsFolder          string

	GuidesContentFolder   string
	GuidesDescriptionFile string
	GuidesBookFile        string

	TmpContentFolder string
	IndexJsonFile    string
	IndexFile        string
	ContentFile      string

	ContentHeaderFile string
	ContentFooterFile string

	AlreadyInProgressFlag string
}

var (
	once       sync.Once
	guidesPath *GuidesPaths
)

func GetGuidesPaths() *GuidesPaths {
	if guidesPath == nil {
		once.Do(
			func() {
				workspace := "/home/codio/workspace"
				args := os.Args[1:]
				if len(args) > 0 {
					workspace = args[0]
				}
				guidesFolder := filepath.Join(workspace, ".guides")
				guidesPath = &GuidesPaths{
					WorkSpace:                  workspace,
					GuidesFolder:               guidesFolder,
					AssessmentsDescriptionFile: filepath.Join(guidesFolder, "assessments.json"),
					AssessmentsFolder:          filepath.Join(guidesFolder, "assessments"),

					GuidesContentFolder:   filepath.Join(guidesFolder, "content"),
					GuidesDescriptionFile: filepath.Join(guidesFolder, "metadata.json"),
					GuidesBookFile:        filepath.Join(guidesFolder, "book.json"),

					TmpContentFolder: filepath.Join(guidesFolder, "content_v3"),
					IndexJsonFile:    "index.json",
					IndexFile:        "index",
					ContentFile:      "content-file",

					ContentHeaderFile: "header.html",
					ContentFooterFile: "footer.html",

					AlreadyInProgressFlag: filepath.Join(workspace, "converterV3AlreadyInProgress"),
				}
			})
	}
	return guidesPath
}
