package content

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"

	copyF "github.com/otiai10/copy"

	"github.com/codio/guides-converter-v3/internal/guidespaths"
	"github.com/codio/guides-converter-v3/internal/utils"
)

var metadataSections = make(map[string]map[string]interface{})
var escapeNameRe = regexp.MustCompile("[^a-zA-Z0-9]")

func Convert() error {
	utils.MakeDir(guidespaths.GetGuidesPaths().TmpContentFolder)
	var metadata map[string]interface{}
	var structure = make(map[string]interface{})
	var rootChildrenWithoutBook []interface{}
	bookJsonIsExists, _ := utils.FileIsExists(guidespaths.GetGuidesPaths().GuidesBookFile)
	if err := utils.GetParsedJson(guidespaths.GetGuidesPaths().GuidesDescriptionFile, &metadata); err != nil {
		return err
	}
	if bookJsonIsExists {
		if err := utils.GetParsedJson(guidespaths.GetGuidesPaths().GuidesBookFile, &structure); err != nil {
			return err
		}
	}

	if sections, exists := metadata["sections"]; exists {
		for _, item := range sections.([]interface{}) {
			section := item.(map[string]interface{})
			if id, err := convertInterfaceToString(section["id"]); err == nil {
				metadataSections[id] = section
				if !bookJsonIsExists {
					addToRootChildrenWithoutBook(&rootChildrenWithoutBook, section)
				}
			}
		}
	}
	if !bookJsonIsExists {
		structure["children"] = rootChildrenWithoutBook
		structure["name"] = ""
	}

	if err := createRootMetadata(metadata, structure); err != nil {
		return err
	}
	if err := copyHtmlHeaderFooter(); err != nil {
		return err
	}
	if children, exists := structure["children"]; exists {
		for _, item := range children.([]interface{}) {
			node := item.((map[string]interface{}))
			err := convertNodeToNewFormat(guidespaths.GetGuidesPaths().TmpContentFolder, &metadata, &node)
			if err != nil {
				log.Println(err.Error())
			}
		}
	}
	return nil
}

func addToRootChildrenWithoutBook(structurePtr *[]interface{}, section map[string]interface{}) {
	var item = make(map[string]interface{})
	item["title"] = section["title"]
	item["id"] = section["id"]
	item["type"] = "page"
	item["pageId"] = section["id"]

	*structurePtr = append(*structurePtr, item)
}

func createRootMetadata(metadata, structure map[string]interface{}) error {
	newMetadata := make(map[string]interface{})
	for k, v := range metadata {
		if k != "sections" {
			newMetadata[k] = v
		}
	}
	newMetadata["title"] = structure["name"]
	if children, exists := structure["children"]; exists {
		order, err := getOrder(children)
		if err != nil {
			return err
		}
		newMetadata["order"] = order
	}
	pathToRootJson := filepath.Join(guidespaths.GetGuidesPaths().TmpContentFolder, guidespaths.GetGuidesPaths().IndexJsonFile)
	return utils.WriteJson(pathToRootJson, newMetadata)
}

func getOrder(children interface{}) ([]string, error) {
	nodes, ok := children.([]interface{})
	if !ok {
		return nil, fmt.Errorf("get order error")
	}
	order := make([]string, 0)
	for _, node := range nodes {
		node := node.(map[string]interface{})
		nodeIsNotExists, err := nodeNotExists(node)
		if err != nil {
			return nil, err
		}
		if nodeIsNotExists {
			continue
		}
		nodeFileName, err := getNodeFileName(node)
		if err != nil {
			return nil, err
		}
		order = append(order, nodeFileName)
	}
	return order, nil
}

func nodeNotExists(node map[string]interface{}) (bool, error) {
	nodeType, err := convertInterfaceToString(node["type"])
	if err != nil {
		return true, err
	}
	if nodeType != "page" {
		return false, nil
	}
	pageId, err := convertInterfaceToString(node["pageId"])
	if err != nil {
		return true, err
	}
	section, sectionExists := metadataSections[pageId]
	if !sectionExists {
		return true, nil
	}
	contentFilePath, err := convertInterfaceToString(section[guidespaths.GetGuidesPaths().ContentFile])
	if err != nil {
		return true, err
	}
	fullFilePath := filepath.Join(guidespaths.GetGuidesPaths().WorkSpace, contentFilePath)
	fileExists, err := utils.FileIsExists(fullFilePath)
	if err != nil {
		return true, err
	}
	return !fileExists, nil
}

func getNodeFileName(node map[string]interface{}) (string, error) {
	title, err := convertInterfaceToString(node["title"])
	if err != nil {
		return "", err
	}
	id, err := convertInterfaceToString(node["id"])
	if err != nil {
		return "", err
	}
	suffix := id[0:4]
	name := fmt.Sprintf("%s-%s", escapeName(title), suffix)
	return name, nil
}

func escapeName(name string) string {
	return escapeNameRe.ReplaceAllString(name, "-")
}

func copyHtmlHeaderFooter() error {
	if err := copyF.Copy(
		filepath.Join(guidespaths.GetGuidesPaths().GuidesContentFolder, guidespaths.GetGuidesPaths().ContentHeaderFile),
		filepath.Join(guidespaths.GetGuidesPaths().TmpContentFolder, guidespaths.GetGuidesPaths().ContentHeaderFile)); err != nil && !os.IsNotExist(err) {
		return err
	}
	if err := copyF.Copy(
		filepath.Join(guidespaths.GetGuidesPaths().GuidesContentFolder, guidespaths.GetGuidesPaths().ContentFooterFile),
		filepath.Join(guidespaths.GetGuidesPaths().TmpContentFolder, guidespaths.GetGuidesPaths().ContentFooterFile)); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func convertNodeToNewFormat(parentPath string, metadataPtr, nodePtr *map[string]interface{}) error {
	err := convertNodeContent(parentPath, metadataPtr, nodePtr)
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	err = createNodeMetadata(parentPath, metadataPtr, nodePtr)
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	nodeFileName, err := getNodeFileName(*nodePtr)
	if err != nil {
		return err
	}
	nodePath := filepath.Join(parentPath, nodeFileName)
	node := *nodePtr
	if children, exists := node["children"]; exists {
		for _, item := range children.([]interface{}) {
			node := item.(map[string]interface{})
			err := convertNodeToNewFormat(nodePath, metadataPtr, &node)
			if err != nil {
				log.Println(err.Error())
			}
		}
	}
	return nil
}

func createNodeMetadata(parentPath string, metadataPtr, nodePtr *map[string]interface{}) error {
	node := *nodePtr
	nodeFileName, err := getNodeFileName(node)
	if err != nil {
		return err
	}
	filePath := parentPath
	metadataFileName := nodeFileName
	nType, err := convertInterfaceToString(node["type"])
	if err != nil {
		return err
	}
	if nType != "page" {
		metadataFileName = guidespaths.GetGuidesPaths().IndexFile
		filePath = filepath.Join(parentPath, nodeFileName)
	}
	content, err := getNodeMetadataContent(filePath, metadataFileName, metadataPtr, nodePtr)
	if err != nil {
		return err
	}
	utils.MakeDir(filePath)
	pathToFile := filepath.Join(filePath, metadataFileName) + ".json"
	if err := utils.WriteJson(pathToFile, content); err != nil {
		return err
	}
	return nil
}

func getNodeMetadataContent(filePath, metadataFileName string, metadataPtr, nodePtr *map[string]interface{}) (map[string]interface{}, error) {
	var nodeMetadata = make(map[string]interface{})
	node := *nodePtr
	if _, exists := node["pageId"]; exists {
		pageId, err := convertInterfaceToString(node["pageId"])
		if err != nil {
			return nil, err
		}
		if section, exists := metadataSections[pageId]; exists {
			nodeMetadata = section
		} else {
			return nil, fmt.Errorf("metadata for node %s not found", pageId)
		}
	}
	title, err := convertInterfaceToString(node["title"])
	if err == nil {
		nodeMetadata["title"] = title
	}
	contentType, err := convertInterfaceToString(nodeMetadata["type"])
	if err == nil {
		nodeMetadata["contentType"] = contentType
	}
	newMetadata := make(map[string]interface{})
	for k, v := range nodeMetadata {
		if k != "sections" {
			newMetadata[k] = v
		}
	}

	typeN, err := convertInterfaceToString(node["type"])
	if err == nil {
		newMetadata["type"] = typeN
	}
	id, err := convertInterfaceToString(node["id"])
	if err == nil {
		newMetadata["id"] = id
	}
	if _, exists := nodeMetadata[guidespaths.GetGuidesPaths().ContentFile]; exists {
		delete(newMetadata, guidespaths.GetGuidesPaths().ContentFile)
	}
	if children, exists := node["children"]; exists {
		order, err := getOrder(children)
		if err != nil {
			return nil, err
		}
		newMetadata["order"] = order
	}
	return newMetadata, nil
}

func convertNodeContent(parentPath string, metadataPtr, nodePtr *map[string]interface{}) error {
	var nodeMetadata = make(map[string]interface{})
	node := *nodePtr
	pageId, err := convertInterfaceToString(node["pageId"])
	if err != nil {
		return nil
	}
	if section, exists := metadataSections[pageId]; exists {
		nodeMetadata = section
	} else {
		return fmt.Errorf("metadata for node %s not found", pageId)
	}
	nodeContentPath, err := convertInterfaceToString(nodeMetadata[guidespaths.GetGuidesPaths().ContentFile])
	if err != nil {
		return fmt.Errorf("content for node %s not found", pageId)
	}
	if _, err := os.Stat(filepath.Join(guidespaths.GetGuidesPaths().WorkSpace, nodeContentPath)); os.IsNotExist(err) {
		return fmt.Errorf("content for node %s not found", pageId)
	}
	extension := filepath.Ext(nodeContentPath)
	nodeFileName, err := getNodeFileName(node)
	if err != nil {
		return err
	}

	fileName := guidespaths.GetGuidesPaths().IndexFile + extension
	if typeN, err := convertInterfaceToString(node["type"]); err == nil && typeN == "page" {
		fileName = nodeFileName + extension
		if err := copyF.Copy(
			filepath.Join(guidespaths.GetGuidesPaths().WorkSpace, nodeContentPath),
			filepath.Join(parentPath, fileName)); err != nil {
			log.Printf("error copy node content: %s\n", err)
		}
		return nil
	}

	if err := copyF.Copy(
		filepath.Join(guidespaths.GetGuidesPaths().WorkSpace, nodeContentPath),
		filepath.Join(parentPath, nodeFileName, fileName)); err != nil {
		log.Printf("error copy node content: %s\n", err)
	}
	return nil
}

func convertInterfaceToString(src interface{}) (string, error) {
	var out string
	var ok bool
	if out, ok = src.(string); !ok {
		return "", fmt.Errorf("convert interface to string error")
	}
	return out, nil
}
