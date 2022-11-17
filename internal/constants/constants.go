package constants

const (
	WorkSpace                  = "/home/codio/workspace/"
	GuidesFolder               = WorkSpace + ".guides"
	AssessmentsDescriptionFile = WorkSpace + GuidesFolder + "/assessments.json"
	AssessmentsFolder          = WorkSpace + GuidesFolder + "/assessments"

	GuidesContentFolder   = WorkSpace + GuidesFolder + "/content"
	GuidesDescriptionFile = GuidesFolder + "/metadata.json"
	GuidesBookFile        = WorkSpace + GuidesFolder + "/book.json"

	TmpContentFolder = WorkSpace + GuidesFolder + "/content_v3"
	IndexJsonFile    = "index.json"
	IndexFile        = "index"
	ContentFile      = "content-file"

	ContentHeaderFile = "header.html"
	ContentFooterFile = "footer.html"

	AlreadyInProgressFlag = WorkSpace + "converterV3AlreadyInProgress"
)
