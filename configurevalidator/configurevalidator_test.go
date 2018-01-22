package configurevalidator

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"os/user"
	"fmt"
)

func TestGetPath(t *testing.T) {
	var a string
	a=GetPath()
	assert.Containsf(t,a,"/.config/perun/main.yaml","Incorrect path")

}

func TestMakeMapRegion(t *testing.T) {
	var a map[string]string
	a=MakeMapRegion()
	assert.NotEmptyf(t,a,"Map of regions isn't empty")
}

func TestMakePath(t *testing.T) {
	var a string
	a=MakePath()
	assert.NotEmptyf(t,a,"Path doesn't exist")
}

func TestMakeYaml(t *testing.T) {
	usr, err:=user.Current()
	path:=usr.HomeDir
	if err!=nil{
		fmt.Println("Error-path")
	}
	assert.FileExistsf(t,path+"/.config/perun/main.yaml","File yaml doesn't exist")
}

func TestSetProfile(t *testing.T) {
	var a bool
	a=SetProfile()
	assert.Falsef(t,a,"Empty profile")
}
func TestSetRegions(t *testing.T) {
	var a bool
	a=SetRegions()
	assert.Falsef(t,a,"Empty region")
}
func TestMakeArrayRegions(t *testing.T) {
	var a [14]string
	a=MakeArrayRegions()
	var b [14]string
	b[0]="us-east-2"
	b[1]="us-east-1"
	b[2]="us-west-1"
	b[3]="us-west-2"
	b[4]="ap-south-1"
	b[5]="ap-northeast-2"
	b[6]="ap-southeast-1"
	b[7]="ap-southeast-2"
	b[8]="ap-northeast-1"
	b[9]="ca-central-1"
	b[10]="eu-central-1"
	b[11]="eu-west-1"
	b[12]="eu-west-2"
	b[13]="sa-east-1"

	assert.ElementsMatchf(t,a,b,"Incorrect list of regions")
}
