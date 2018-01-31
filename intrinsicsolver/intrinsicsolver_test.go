// Copyright 2018 Appliscale
//
// Maintainers and contributors are listed in README file inside repository.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package intrinsicsolver

import (
	"os"
	"testing"

	"github.com/Appliscale/perun/logger"
	"github.com/stretchr/testify/assert"
)

var sink logger.Logger

func setup() {
	sink = logger.Logger{}
}

func TestMain(m *testing.M) {
	setup()
	retCode := m.Run()
	os.Exit(retCode)
}

func TestDeleteFromSlice(t *testing.T) {
	slice := []string{"a", "b", "c", "d"}
	compareSlice := []string{"a", "b", "%%%", "%%%"}
	deleteFromSlice(slice, 1, 1)
	deleteFromSlice(slice, 1, 1)
	assert.Equal(t, compareSlice, slice, "MSG")
}

func TestFixRef(t *testing.T) {
	/*
		from Doc's:
		YAML:
		!Ref "logicalName"
		JSON:
		{ "Ref" : "logicalName" }
	*/
	line := "!Ref \"logicalName\""
	expected := "{ \"Ref\" : \"logicalName\" }"
	fixRef(&line)
	assert.Equal(t, expected, line, "MSG")
}

func TestFixImportValue(t *testing.T) {
	/*
		from Doc's:
		YAML:
		!ImportValue "sharedValueToImport"
		JSON:
		{ "Fn::ImportValue" : "sharedValueToImport" }
	*/
	line := "Key: !ImportValue \"sharedValueToImport\""
	expected := "Key:  { \"Fn::ImportValue\" :  \"sharedValueToImport\"}"
	fixImportValue(&line)
	assert.Equal(t, expected, line, "MSG")
}

func TestSub(t *testing.T) {
	/*
				from Doc's:
				YAML:
				!Sub
		  			- String
		  			- { Var1Name: Var1Value, Var2Name: Var2Value }
				JSON:
				{ "Fn::Sub" : [ String, { Var1Name: Var1Value, Var2Name: Var2Value } ] }
	*/
	line1 := "        - Key: !Sub"
	line2 := "            - \"String\""
	line3 := "            - KeyInside: something"
	linesArray := make([]string, 0)
	linesArray = append(linesArray, line1)
	linesArray = append(linesArray, line2)
	linesArray = append(linesArray, line3)

	expected := "        - Key: { \"Fn::Sub\" : [ \"String\", {\"KeyInside\" : something}]}"

	fixSub(&line1, linesArray, 0)

	assert.Equal(t, expected, line1, "MSG")

}

func TestFixFindInMap(t *testing.T) {
	/*
		from Doc's:
		YAML:
		!FindInMap [ MapName, TopLevelKey, SecondLevelKey ]
		JSON:
		{ "Fn::FindInMap" : [ "MapName", "TopLevelKey", "SecondLevelKey"] }
	*/
	line := "!FindInMap [ MapName, TopLevelKey, SecondLevelKey ]"
	expected := " { \"Fn::FindInMap\" : [\"MapName\", \"TopLevelKey\", \"SecondLevelKey\"] }"

	fixFindInMap(&line)

	assert.Equal(t, expected, line)

}

func TestFixGetAtt(t *testing.T) {
	/*
		from Doc's:
		YAML:
		!GetAtt logicalNameOfResource.attributeName
		JSON:
		{ "Fn::GetAtt" : [ "logicalNameOfResource", "attributeName" ] }
	*/
	line := "Key: !GetAtt logicalNameOfResource.attributeName"
	expected := "Key: { \"Fn::GetAtt\" : [\"logicalNameOfResource\", \"attributeName\" ] }"
	fixGetAtt(&line)

	assert.Equal(t, expected, line, "MSG")
}

func TestFixEquals(t *testing.T) {
	/*
		from Doc's:
		YAML:
		!Equals [value_1, value_2]
		JSON:
		"Fn::Equals" : ["value_1", "value_2"]
	*/
	line := "Key: !Equals [value_1, value_2]"
	expected := "Key: {\"Fn::Equals\" : [\"value_1\", \"value_2\"]}"
	fixEquals(&line)

	assert.Equal(t, expected, line, "MSG")
}

func TestFixUserData(t *testing.T) {
	/*
	   		from Doc's:
	   		YAML:
	   		UserData:
	     		  Fn::Base64:
	               !Sub |
	                 #!/bin/bash -xe
	                 yum update -y aws-cfn-bootstrap
	                 /opt/aws/bin/cfn-init -v --stack ${AWS::StackName} --resource LaunchConfig --configsets wordpress_install --region ${AWS::Region}
	                 /opt/aws/bin/cfn-signal -e $? --stack ${AWS::StackName} --resource WebServerGroup --region ${AWS::Region}

	   		JSON:
	   		UserData: { "Fn::Base64" : { "Fn::Join": ["\n", ["#!/bin/bash -xe","yum update -y aws-cfn-bootstrap",{ "Fn::Sub" : "/opt/aws/bin/cfn-init -v --stack ${AWS::StackName} --resource LaunchConfig --configsets wordpress_install --region ${AWS::Region}" },{ "Fn::Sub" : "/opt/aws/bin/cfn-signal -e $? --stack ${AWS::StackName} --resource WebServerGroup --region ${AWS::Region}" }]]}}

	*/
	line1 := "  Fn::Base64:"
	line2 := "    !Sub |"
	line3 := "      #!/bin/bash -xe"
	line4 := "      yum update -y aws-cfn-bootstrap"
	line5 := "      /opt/aws/bin/cfn-init -v --stack ${AWS::StackName} --resource LaunchConfig --configsets wordpress_install --region ${AWS::Region}"
	line6 := "      /opt/aws/bin/cfn-signal -e $? --stack ${AWS::StackName} --resource WebServerGroup --region ${AWS::Region}"
	line7 := ""
	linesArray := make([]string, 0)
	linesArray = append(linesArray, line1)
	linesArray = append(linesArray, line2)
	linesArray = append(linesArray, line3)
	linesArray = append(linesArray, line4)
	linesArray = append(linesArray, line5)
	linesArray = append(linesArray, line6)
	linesArray = append(linesArray, line7)

	expected := "  {\"Fn::Base64\" : { \"Fn::Join\" : [\"\n\", [\"#!/bin/bash -xe\",\"yum update -y aws-cfn-bootstrap\",{ \"Fn::Sub\" : \"/opt/aws/bin/cfn-init -v --stack ${AWS::StackName} --resource LaunchConfig --configsets wordpress_install --region ${AWS::Region}\"},{ \"Fn::Sub\" : \"/opt/aws/bin/cfn-signal -e $? --stack ${AWS::StackName} --resource WebServerGroup --region ${AWS::Region}\"}]]}}"
	fixUserData(&line1, linesArray, 0)

	assert.Equal(t, expected, line1, "MSG")

}

func TestFixIf(t *testing.T) {
	/*
		   		from Doc's:
		   		YAML:
				   !If [condition_name, value_if_true, value_if_false]

		   		JSON:
		   		{ "Fn::If": [condition_name, value_if_true, value_if_false]}

	*/
	line1 := "!If [condition_name, value_if_true, value_if_false]"
	line2 := ""
	linesArray := make([]string, 0)
	linesArray = append(linesArray, line1)
	linesArray = append(linesArray, line2)

	expected := "{ \"Fn::If\" : [\"condition_name\", \" value_if_true\", \" value_if_false\"]}"
	fixIf(&line1, linesArray, 0)

	assert.Equal(t, expected, line1, "MSG")
}

func TestFixSplit(t *testing.T) {
	/*
		from Doc's:
		YAML:
		!Split [ "|" , "a|b|c" ]
		JSON:
		{ "Fn::Split" : [ "|" , "a|b|c" ] }
	*/
	line := "!Split [ \"|\" , \"a|b|c\" ]"
	expected := "{ \"Fn::Split\" : [ \"|\" , \"a|b|c\" ] }"
	fixSplit(&line)

	assert.Equal(t, expected, line, "MSG")

}

func TestFixGetAZs(t *testing.T) {
	/*
		from Doc's:
		YAML:
		!GetAZs region
		JSON:
		{ "Fn::GetAZs" : "region" }
	*/
	line := "!GetAZs region"
	expected := "{ \"Fn::GetAZs\" : \"region\" }"
	fixGetAZs(&line)

	assert.Equal(t, expected, line, "MSG")
}

func TestFixSelect(t *testing.T) {
	/*
		from Doc's:
		YAML:
		!Select [ "1", [ "apples", "grapes", "oranges", "mangoes" ] ]
		JSON:
		{ "Fn::Select" : [ "1", [ "apples", "grapes", "oranges", "mangoes" ] ] }
	*/
	line1 := "!Select [ \"1\", [ \"apples\", \"grapes\", \"oranges\", \"mangoes\" ] ]"
	line2 := ""
	linesArray := make([]string, 0)
	linesArray = append(linesArray, line1)
	linesArray = append(linesArray, line2)
	expected := "{ \"Fn::Select\" : [ \"1\", [ \"apples\", \"grapes\", \"oranges\", \"mangoes\" ] ]}"

	fixSelect(&line1, linesArray, 0)

	assert.Equal(t, expected, line1, "MSG")
}
