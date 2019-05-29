package cfnyaml

import (
	"fmt"
	"gopkg.in/yaml.v3"
)

type packageablePropertyDefinition struct {
	ResourceType string
	Path         []string
	Rewritten    func(bucket, key, versionId string) *yaml.Node
}

func standardS3Uri(bucket, key, versionId string) *yaml.Node {
	val := fmt.Sprintf("s3://%s/%s", bucket, key)
	if versionId != "" {
		val += fmt.Sprintf("?versionId=%s", versionId)
	}
	return &yaml.Node{
		Kind:  yaml.ScalarNode,
		Tag:   "!!str",
		Value: val,
	}
}

func toYamlNode(in interface{}) *yaml.Node {
	b, _ := yaml.Marshal(in)
	n := yaml.Node{}
	_ = yaml.Unmarshal(b, &n)
	return n.Content[0]
}

var packageablePropertyDefinitions = []packageablePropertyDefinition{
	{
		ResourceType: "AWS::ApiGateway::RestApi",
		Path:         []string{"Properties", "BodyS3Location"},
		Rewritten: func(bucket, key, versionId string) *yaml.Node {
			return toYamlNode(struct {
				Bucket  string `yaml:"Bucket"`
				Key     string `yaml:"Key"`
				Version string `yaml:"Version,omitempty"`
			}{Bucket: bucket, Key: key, Version: versionId})
		},
	},
	{
		ResourceType: "AWS::Lambda::Function",
		Path:         []string{"Properties", "Code"},
		Rewritten: func(bucket, key, versionId string) *yaml.Node {
			return toYamlNode(struct {
				Bucket  string `yaml:"S3Bucket"`
				Key     string `yaml:"S3Key"`
				Version string `yaml:"S3ObjectVersion,omitempty"`
			}{Bucket: bucket, Key: key, Version: versionId})
		},
	},
	{
		ResourceType: "AWS::Serverless::Function",
		Path:         []string{"Properties", "CodeUri"},
		Rewritten: func(bucket, key, versionId string) *yaml.Node {
			return toYamlNode(struct {
				Bucket  string `yaml:"Bucket"`
				Key     string `yaml:"Key"`
				Version string `yaml:"Version,omitempty"`
			}{Bucket: bucket, Key: key, Version: versionId})
		},
	},
	{
		ResourceType: "AWS::AppSync::GraphQLSchema",
		Path:         []string{"Properties", "DefinitionS3Location"},
		Rewritten:    standardS3Uri,
	},
	{
		ResourceType: "AWS::AppSync::Resolver",
		Path:         []string{"Properties", "RequestMappingTemplateS3Location"},
		Rewritten:    standardS3Uri,
	},
	{
		ResourceType: "AWS::AppSync::Resolver",
		Path:         []string{"Properties", "ResponseMappingTemplateS3Location"},
		Rewritten:    standardS3Uri,
	},
	{
		ResourceType: "AWS::Serverless::Api",
		Path:         []string{"Properties", "DefinitionUri"},
		Rewritten: func(bucket, key, versionId string) *yaml.Node {
			return toYamlNode(struct {
				Bucket  string `yaml:"Bucket"`
				Key     string `yaml:"Key"`
				Version string `yaml:"Version,omitempty"`
			}{Bucket: bucket, Key: key, Version: versionId})
		},
	},
	{
		ResourceType: "AWS::Include",
		Path:         []string{"Properties", "Location"},
		Rewritten:    standardS3Uri,
	},
	{
		ResourceType: "AWS::ElasticBeanstalk::ApplicationVersion",
		Path:         []string{"Properties", "SourceBundle"},
		Rewritten: func(bucket, key, versionId string) *yaml.Node {
			return toYamlNode(struct {
				Bucket string `yaml:"Bucket"`
				Key    string `yaml:"Key"`
			}{Bucket: bucket, Key: key})
		},
	},
	{
		ResourceType: "AWS::CloudFormation::Stack",
		Path:         []string{"Properties", "TemplateURL"},
		Rewritten:    standardS3Uri,
	},
}
