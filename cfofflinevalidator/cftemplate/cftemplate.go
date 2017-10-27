package cftemplate

type Template struct {
	AWSTemplateFormatVersion string
	Description string
	Metadata map[string]interface{}
	Parameters map[string]interface{}
	Mappings map[string]interface{}
	Conditions map[string]interface{}
	Transform map[string]interface{}
	Resources map[string]Resource
	Outputs map[string]interface{}
}

type Resource struct {
	Type string
	Properties map[string]interface{}
}