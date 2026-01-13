package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

type UpdateBuild struct {
	Core struct {
		BuildDate string `json:"buildDate"`
		Name      string `json:"name"`
		Sha1      string `json:"sha1"`
		Sha254    string `json:"sha256"`
		Size      int    `json:"size"`
		Url       string `json:"url"`
		Version   string `json:"version"`
	} `json:"core"`
}

var UpBuild UpdateBuild

var jenkinsCmd = &cobra.Command{
	Use:   "jenkins",
	Short: "Обновить jenkins до последней версии",
	Run: func(cmd *cobra.Command, args []string) {

		var JenkinsVersion string
		for {
			cmdGetVer := fmt.Sprintf("java -jar %s -s %s version", str.path, str.GetHttpConnection())
			StrCm := `java -jar ./jenkins-cli.jar -s http://jenkins:8080 -auth ***:*** version`
			version := exec.Command("/bin/sh", "-c", cmdGetVer)
			ver, err := version.Output()
			if err != nil {
				log.Printf("Jenkins is not available %v: %v", StrCm, err)
			} else {
				log.Printf("%v: %v", StrCm, string(ver))
				JenkinsVersion = string(ver)
				break
			}
		}

		jsonDefault, err := os.Open(str.pathJsn)
		if err != nil {
			log.Printf("Error opening file: %v", err)
		} else {
			bte, _ := io.ReadAll(jsonDefault)
			json.Unmarshal(bte, &UpBuild)
		}
		defer jsonDefault.Close()

		var newRel bool = UpBuild.Core.Version > JenkinsVersion

		if newRel {

			log.Printf("Update jenkins: %s -> %s ...", strings.TrimSpace(JenkinsVersion), UpBuild.Core.Version)

			if _, err := os.Stat(str.pathWar); errors.Is(err, os.ErrNotExist) {
				log.Printf("Absent: %v", str.pathWar)
			} else {
				target := fmt.Sprintf("%v.old_%v", str.pathWar, time.Now().Format("02012006"))
				mv := exec.Command("sudo", "mv", str.pathWar, target)
				if err := mv.Run(); err != nil {
					log.Printf("%v: %v", mv, err)
				} else {
					log.Println(mv)
				}
			}

			url := fmt.Sprintf("https://updates.jenkins.io/download/war/%s/jenkins.war", UpBuild.Core.Version)
			getJenkins := exec.Command("sudo", "curl", "-Lo", str.pathWar, url)
			getJenkins.Stdout = os.Stdout
			getJenkins.Stderr = os.Stderr
			if err := getJenkins.Run(); err != nil {
				log.Printf("%v: %v", getJenkins, err)
			} else {
				log.Printf("jenkins.war: Ok %s (%dM)", str.pathWar, UpBuild.Core.Size/1024/1024)
			}

			restart := exec.Command("sudo", "systemctl", "restart", "jenkins")
			restart.Stdout = os.Stdout
			restart.Stderr = os.Stderr
			if err = restart.Run(); err != nil {
				log.Println("Error: ", err)
			} else {
				log.Println("Jenkins is restarted ...")
			}
		}

	},
}

func init() {}
