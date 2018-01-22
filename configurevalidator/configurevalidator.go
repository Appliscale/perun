package configurevalidator

import (
	"fmt"
	"os"
	"os/user"

)

var profile string
var NumberRegion int
var regions map[string]string
var numbers[14]string
var file os.File
var usr *user.User
var path string
var exist bool

func GetExist() bool{
	return exist
}

func MakeArrayRegions () [14]string {
	numbers[0]="us-east-2"
	numbers[1]="us-east-1"
	numbers[2]="us-west-1"
	numbers[3]="us-west-2"
	numbers[4]="ap-south-1"
	numbers[5]="ap-northeast-2"
	numbers[6]="ap-southeast-1"
	numbers[7]="ap-southeast-2"
	numbers[8]="ap-northeast-1"
	numbers[9]="ca-central-1"
	numbers[10]="eu-central-1"
	numbers[11]="eu-west-1"
	numbers[12]="eu-west-2"
	numbers[13]="sa-east-1"

	return numbers

}

func ShowRegions( ) {
numbers=MakeArrayRegions()
	fmt.Println("New configuration file. Choose region:")
	for i:=0; i<len(numbers);i++{
		fmt.Println("Number ",i," region ",numbers[i])
	}
	for !SetRegions(){
		SetRegions()
	}

	fmt.Println("Profile: ")
	for !SetProfile(){
		SetProfile()
	}

	}

func SetRegions( ) (ok bool){
	//get number from the console
	fmt.Scan(&NumberRegion)
	if NumberRegion>=0 && NumberRegion<14 {
	fmt.Println("Your region is: " + numbers[NumberRegion])
	ok=true
	return ok
} else {
	ok=false
	fmt.Println("Invalid region, try again")
	return ok
}
}

func SetProfile() (ok bool){
	//get name of the profile
	fmt. Scan(&profile)
	if profile!="" {
		fmt.Println("Your profile is: " + profile)
		ok=true
		return ok
	} else {
		ok=false
		fmt.Println("Invalid profile, try again")
		return ok
	}
}

func MakePath() string{
	usr, err := user.Current()
	if err!=nil {
		fmt.Println("Error-path")
		return ""
	}
	 path:= usr.HomeDir
	 return path
}

func MakeYaml() {
	//make file with region and profile
	file, err:=os.Create(MakePath()+"/.config/perun/main.yaml")

	if err !=nil {
		fmt.Println("Error- file")
		return
	}

	defer file.Close()
	_,err=file.WriteString("Default profiles: "+profile +"\n")
	_,err=file.WriteString("Default region: "+numbers[NumberRegion]+"\n")
	_,err=file.WriteString("SpecificationURL: \n")

	for j:=0; j<len(numbers);j++{
		_,err=file.WriteString("  "+numbers[j]+":"+" "+regions[numbers[j]]+"\n")
	}
	if err!=nil {
		fmt.Println("Error- write to file")
	}
}

func IsFileExist(a bool ){
exist=a
}

func GetPath()string{
	configurePath:=MakePath()+"/.config/perun/main.yaml"
	return configurePath
}

func MakeMapRegion ()  map[string]string{
	regions =make(map[string]string)
	regions["us-east-2"]="https://dnwj8swjjbsbt.cloudfront.net"
	regions["us-east-1"]="https://d1uauaxba7bl26.cloudfront.net"
	regions["us-west-1"]="https://d1uauaxba7bl26.cloudfront.net"
	regions["us-west-2"]="https://d201a2mn26r7lk.cloudfront.net"
	regions["ap-south-1"]="https://d2senuesg1djtx.cloudfront.net"
	regions["ap-northeast-2"]="https://d1ane3fvebulky.cloudfront.net"
	regions["ap-southeast-1"]="https://doigdx0kgq9el.cloudfront.net"
	regions["ap-southeast-2"]="https://d2stg8d246z9di.cloudfront.net"
	regions["ap-northeast-1"]="https://d33vqc0rt9ld30.cloudfront.net"
	regions["ca-central-1"]="https://d2s8ygphhesbe7.cloudfront.net"
	regions["eu-central-1"]="https://d1mta8qj7i28i2.cloudfront.net"
	regions["eu-west-1"]="https://d3teyb21fexa9r.cloudfront.net"
	regions["eu-west-2"]="https://d1742qcu2c1ncx.cloudfront.net"
	regions["sa-east-1"]="https://d3c9jyj3w509b0.cloudfront.net"
	return regions
}






