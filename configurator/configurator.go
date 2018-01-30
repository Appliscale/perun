package configurator

import (
	"github.com/Appliscale/perun/configuration"
	"github.com/Appliscale/perun/logger"
	"os"
	"os/user"
	"strconv"
)

var resourceSpecificationURL = map[string]string{
	"us-east-2":      "https://dnwj8swjjbsbt.cloudfront.net",
	"us-east-1":      "https://d1uauaxba7bl26.cloudfront.net",
	"us-west-1":      "https://d68hl49wbnanq.cloudfront.net",
	"us-west-2":      "https://d201a2mn26r7lk.cloudfront.net",
	"ap-south-1":     "https://d2senuesg1djtx.cloudfront.net",
	"ap-northeast-2": "https://d1ane3fvebulky.cloudfront.net",
	"ap-southeast-1": "https://doigdx0kgq9el.cloudfront.net",
	"ap-southeast-2": "https://d2stg8d246z9di.cloudfront.net",
	"ap-northeast-1": "https://d33vqc0rt9ld30.cloudfront.net",
	"ca-central-1":   "https://d2s8ygphhesbe7.cloudfront.net",
	"eu-central-1":   "https://d1mta8qj7i28i2.cloudfront.net",
	"eu-west-1":      "https://d3teyb21fexa9r.cloudfront.net",
	"eu-west-2":      "https://d1742qcu2c1ncx.cloudfront.net",
	"sa-east-1":      "https://d3c9jyj3w509b0.cloudfront.net",
}

func CreateConfiguration() string {
	logger := logger.CreateDefaultLogger()
	yourpath, yourname := ConfigurePath(logger)
	return yourpath + "/" + yourname
}

func findFile(path string) {
	logger := logger.CreateDefaultLogger()
	logger.Always("File will be in " + path)
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		showRegions(logger)
		con := createConfig()
		configuration.PrepareYaml(con, path, logger)
	}
}

func ConfigurePath(logger logger.Logger) (string, string) {
	logger.Always("Configure file could be in \n  " + makeUserPath(logger) + "\n  /etc/perun")
	yourpath := ""
	yourname := ""
	logger.GetInput("Your path ", &yourpath)
	logger.GetInput("Filename ", &yourname)
	findFile(yourpath + "/" + yourname)

	return yourpath, yourname
}

func makeUserPath(logger logger.Logger) (path string) {
	usr, err := user.Current()
	if err != nil {
		logger.Error("Error-path")
	}
	path = usr.HomeDir
	path = path + "/.config/perun"
	return
}

func showRegions(logger logger.Logger) {
	regions := makeArrayRegions()
	logger.Always("Select region")
	for i := 0; i < len(regions); i++ {
		pom := strconv.Itoa(i)
		logger.Always("Number " + pom + " region " + regions[i])
	}
}

func setRegions(logger logger.Logger) (region string) {
	var numberRegion int
	logger.GetInput("Choose region", &numberRegion)
	regions := makeArrayRegions()
	if numberRegion >= 0 && numberRegion < 14 {
		region = regions[numberRegion]
		logger.Always("Your region is: " + region)

	} else {
		logger.Error("Invalid region")

	}
	return
}

func setProfile(logger logger.Logger) (profile string) {
	logger.GetInput("Input name of profile", &profile)
	if profile != "" {
		logger.Always("Your profile is: " + profile)

	} else {
		logger.Error("Invalid profile")

	}
	return
}

func createConfig() configuration.Configuration {
	logger := logger.CreateDefaultLogger()
	myregion := setRegions(logger)
	myprofile := setProfile(logger)
	myResourceSpecificationURL := resourceSpecificationURL

	myconfig := configuration.Configuration{
		myprofile,
		myregion,
		myResourceSpecificationURL,
		false,
		3600,
		"INFO"}

	return myconfig
}

func makeArrayRegions() [14]string {
	var regionsN [14]string
	regionsN[0] = "us-east-1"
	regionsN[1] = "us-east-2"
	regionsN[2] = "us-west-1"
	regionsN[3] = "us-west-2"
	regionsN[4] = "ca-central-1"
	regionsN[5] = "ca-central-1"
	regionsN[6] = "eu-west-1"
	regionsN[7] = "eu-west-2"
	regionsN[8] = "ap-northeast-1"
	regionsN[9] = "ap-northeast-2"
	regionsN[10] = "ap-southeast-1"
	regionsN[11] = "ap-southeast-2"
	regionsN[12] = "ap-south-1"
	regionsN[13] = "sa-east-1"

	return regionsN

}
