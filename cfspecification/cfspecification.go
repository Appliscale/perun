package cfspecification

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"io"
	"net/http"
	"os/user"
	"github.com/Appliscale/cftool/cfconfiguration"
	"strings"
)

type Specification struct {
	PropertyTypes map[string]PropertyType
	ResourceSpecificationVersion string
	ResourceTypes map[string]Resource
}

type PropertyType struct {
	Documentation string
	Properties map[string]Property
}

type Property struct {
	Documentation string
	DuplicatesAllowed bool
	ItemType string
	PrimitiveItemType string
	PrimitiveType string
	Required bool
	Type string
	UpdateType string
}

type Resource struct {
	Documentation string
	Attributes map[string]Attribute
	Properties map[string]Property
}

type Attribute struct {
	ItemType string
	PrimitiveItemType string
	PrimitiveType string
	Type string
}

func GetSpecification(configurationFilePath string) (specification Specification, err error) {
	filePath, err := downloadSpecification(configurationFilePath)
	if err != nil  {
		return specification, err
	}

	return GetSpecificationFromFile(filePath)
}

func GetSpecificationFromFile(specificationFilePath string) (specification Specification, err error) {
	specificationFile, err := ioutil.ReadFile(specificationFilePath)
	if err != nil  {
		return specification, err
	}

	return parseSpecificationFile(specificationFile)
}

func downloadSpecification(configurationFilePath string) (filePath string, err error) {
	user, err := user.Current()
	if err != nil  {
		return
	}

	specificationDir := user.HomeDir + "/.Appliscale/cftool/specification"
	specificationFileUrl, err := cfconfiguration.GetSpecificationFileURL(configurationFilePath)
	if err != nil  {
		return
	}
	fileName := strings.Replace(specificationFileUrl, "https://", "", -1)
	fileName = strings.Replace(fileName, ".cloudfront.net/latest/gzip/CloudFormationResourceSpecification.json", "", -1)
	specificationFilePath := specificationDir + "/" + fileName + ".json"

	if _, err := os.Stat(specificationFilePath); err == nil {
		return specificationFilePath, nil
	}

	if _, err := os.Stat(specificationDir); os.IsNotExist(err) {
		os.MkdirAll(specificationDir, os.ModePerm)
	}
	out, err := os.Create(specificationFilePath)
	if err != nil  {
		return
	}
	defer out.Close()

	resp, err := http.Get(specificationFileUrl)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil  {
		return
	}

	return specificationFilePath, nil
}

func parseSpecificationFile(specificationFile []byte) (specification Specification, err error) {
	err = json.Unmarshal(specificationFile, &specification)
	if err != nil {
		return specification, err
	}

	return specification, nil
}