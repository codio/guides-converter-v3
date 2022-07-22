package content

import (
	// "github.com/codio/guides-converter-v3/internal/utils"
)

const (
	GUIDES_FOLDER = ".guides"
	GUIDES_DESCRIPTION_FILE = GUIDES_FOLDER + "/metadata.json"
	GUIDES_BOOK_FILE = GUIDES_FOLDER + "/book.json"

	TMP_CONTENT_FOLDER_NAME = "content_v3"
	TMP_CONTENT_FOLDER = GUIDES_FOLDER + "/" + TMP_CONTENT_FOLDER_NAME
	INDEX_FILE = "index"
)


func Convert() error {
	return nil
}
