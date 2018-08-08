package checkingrequiredfiles

import (
	"github.com/Appliscale/perun/configuration"
	"github.com/Appliscale/perun/configurator"
	"github.com/Appliscale/perun/context"
	"github.com/Appliscale/perun/myuser"
	"github.com/go-ini/ini"
	"io"
	"net/http"
	"os"
	"strings"
)

func isMainYAMLexists(ctx *context.Context) (bool, error) {
	homePath, pathError := myuser.GetUserHomeDir()
	if pathError != nil {
		ctx.Logger.Error(pathError.Error())
		return false, pathError
	}
	_, mainError := os.Open(homePath + "/.config/perun/main.yaml")
	if mainError != nil {
		_, mainError = os.Open(homePath + "/etc/perun/main.yaml")
		if mainError != nil {
			return false, pathError
		}
		return true, pathError
	}
	return true, pathError
}

//CheckingRequiredFiles looks for required and default files and if they don't exist
func CheckingRequiredFiles(ctx *context.Context) {
	isMainYAMLexists, mainError := isMainYAMLexists(ctx)
	if mainError != nil {
		ctx.Logger.Error(mainError.Error())
	}

	isAWSconfigExists := isAWSconfigExists(ctx)

	homePath, pathError := myuser.GetUserHomeDir()
	if pathError != nil {
		ctx.Logger.Error(pathError.Error())
	}

	if !isMainYAMLexists {
		var answer string
		ctx.Logger.GetInput("I see it's your first perun's run. I'll create default configuration file or you can do it by yourself. \nDefault? Y/N", &answer)
		if answer == "Y" || answer == "y" {
			ctx.Logger.Always("I'm creating default configuration file. After that run perun again.")
			if isAWSconfigExists {
				region := getRegionFromAWSconfig(ctx, homePath+"/.aws/config")
				ctx.Logger.Always("Creating main.yaml based on .aws/config.")
				createConfigurationFile("default", homePath, region, ctx)
			} else {
				createConfigurationFile("default", homePath, "us-east-1", ctx)
			}
			var answer string
			ctx.Logger.GetInput("Do you want to create .aws/credentials for this profile? Y/N", &answer)
			if answer == "Y" || answer == "y" {
				configurator.GetKeysFromUser(ctx, "default")
			}
		} else if answer == "N" || answer == "n" {
			ctx.Logger.Always("Here you can create your own configuration file. After that run perun again.")
			configurator.FileName(ctx)
		} else {
			ctx.Logger.Error("Invalid input")
		}
	}
	isCredentialsExists := isCredentialsExists(ctx)
	if !isCredentialsExists {
		configurator.GetKeysFromUser(ctx, ctx.Config.DefaultProfile)
		ctx.Logger.Always("Creating .aws/credentials file.")

		err := context.UpdateSessionToken(ctx.Config.DefaultProfile, ctx.Config.DefaultRegion, ctx.Config.DefaultDurationForMFA, ctx)
		if err != nil {
			ctx.Logger.Error(err.Error())
			os.Exit(1)
		}
	}
	downloadError := downloadDefaultFile()
	if downloadError != nil {
		ctx.Logger.Error(downloadError.Error())
	}
}

func createConfigurationFile(profile string, path string, region string, ctx *context.Context) {
	myConfig := configuration.Configuration{
		DefaultProfile:                 profile,
		DefaultRegion:                  region,
		SpecificationURL:               configurator.ResourceSpecificationURL,
		DefaultDecisionForMFA:          true,
		DefaultDurationForMFA:          3600,
		DefaultVerbosity:               "INFO",
		DefaultTemporaryFilesDirectory: path,
	}
	configuration.SaveToFile(myConfig, path+"/.config/perun/"+"main.yaml", ctx.Logger)
}

func isAWSconfigExists(ctx *context.Context) bool {
	homePath, pathError := myuser.GetUserHomeDir()
	if pathError != nil {
		ctx.Logger.Error(pathError.Error())
	}
	_, credentialsError := os.Open(homePath + "/.aws/config")
	if credentialsError != nil {
		return false
	}
	return true

}

func getRegionFromAWSconfig(ctx *context.Context, path string) string {
	var region string
	configuration, loadCredentialsError := ini.Load(path)
	if loadCredentialsError != nil {
		ctx.Logger.Error(loadCredentialsError.Error())
	}
	section, sectionError := configuration.GetSection("default")
	if sectionError != nil {
		ctx.Logger.Error(sectionError.Error())
	}
	region = section.Key("region").Value()

	return region
}

func downloadDefaultFile() error {
	urls := make(map[string]string)
	urls["blocked.json"] = "https://s3.amazonaws.com/perun-default-file/blocked.json"
	urls["unblocked.json"] = "https://s3.amazonaws.com/perun-default-file/unblocked.json"
	urls["style.yaml"] = "https://s3.amazonaws.com/perun-default-file/style.yaml"
	urls["specification_inconsistency.yaml"] = "https://s3.amazonaws.com/perun-default-file/specification_inconsistency.yaml"

	for file, url := range urls {
		homePath, _ := myuser.GetUserHomeDir()
		homePath += "/.config/perun/"

		if strings.Contains(file, "blocked") {
			homePath += "stack-policies/"
		}

		_, err := os.Stat(homePath)
		if os.IsNotExist(err) {
			os.Mkdir(homePath, 0755)
		}

		_, openError := os.Open(homePath + file)
		if openError != nil {
			out, creatingFileError := os.Create(homePath + file)

			if creatingFileError != nil {
				return creatingFileError
			}
			defer out.Close()

			resp, httpGetError := http.Get(url)
			if httpGetError != nil {
				return httpGetError
			}
			defer resp.Body.Close()

			_, copyError := io.Copy(out, resp.Body)
			if copyError != nil {
				return copyError
			}
		}
	}
	return nil
}

func isCredentialsExists(ctx *context.Context) bool {
	homePath, pathError := myuser.GetUserHomeDir()
	if pathError != nil {
		ctx.Logger.Error(pathError.Error())
	}
	_, credentialsError := os.Open(homePath + "/.aws/credentials")
	if credentialsError != nil {
		return false
	}
	return true
}
