package guidespaths

var (
	WorkSpace                  = "/home/codio/workspace/"
	GuidesFolder               = WorkSpace + ".guides"
	AssessmentsDescriptionFile = GuidesFolder + "/assessments.json"
	AssessmentsFolder          = GuidesFolder + "/assessments"

	GuidesContentFolder   = GuidesFolder + "/content"
	GuidesDescriptionFile = GuidesFolder + "/metadata.json"
	GuidesBookFile        = GuidesFolder + "/book.json"

	TmpContentFolder = GuidesFolder + "/content_v3"
	IndexJsonFile    = "index.json"
	IndexFile        = "index"
	ContentFile      = "content-file"

	ContentHeaderFile = "header.html"
	ContentFooterFile = "footer.html"

	AlreadyInProgressFlag = WorkSpace + "converterV3AlreadyInProgress"
)
