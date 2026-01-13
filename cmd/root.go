package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

type Param struct {
	user     string
	password string
	url      string
	api      string
	path     string
	pathWar  string
	pathJsn  string
	all      bool
	jenkins  bool
	plugin   bool
}

var str = &Param{}

type Update interface {
	UpdateJenkins()
	UpdatePlugin()
	UpdateAll()
}

type UpDependency interface {
	GetHttpConnection() string
	UploadJar() error
}

func (p Param) UpdateJenkins() {
	jenkinsCmd.SetArgs([]string{})
	jenkinsCmd.Execute()
}

func (p Param) UpdatePlugin() {
	pluginCmd.SetArgs([]string{})
	pluginCmd.Execute()
}

func (p Param) UpdateAll() {
	p.UpdateJenkins()
	p.UpdatePlugin()
}

func (p Param) GetHttpConnection() string {

	p.api = fmt.Sprintf("%s -auth %s:%s", p.url, p.user, p.password)
	return p.api
}

func (p *Param) UploadJar() error {
	if _, err := os.Stat(p.path); errors.Is(err, os.ErrNotExist) {
		url := fmt.Sprintf("%s/jnlpJars/jenkins-cli.jar", p.url)
		log.Println("Downloading Jenkins CLI Jar: ", url)
		cmd := exec.Command("curl", "-O", url)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Errorf("ошибка при загрузке jenkins-cli.jar: %s", err)
			os.Exit(1)
		} else {
			p.path = "./jenkins-cli.jar"
			log.Println("jenkins-cli.jar: Ok", p.path)
		}
	}
	return nil
}

func updateSystem(up Update, all bool, plugin bool, jenkins bool) {
	if all {
		up.UpdateAll()
	} else if plugin {
		up.UpdatePlugin()
	} else if jenkins {
		up.UpdateJenkins()
	}
}

func dependency(d UpDependency) { // TODO: realisation with gorutine + context, or reverse run command
	d.GetHttpConnection()
	if err := d.UploadJar(); err != nil {
		log.Println(err)
		os.Exit(1)
	}

	_, err := exec.LookPath("java")
	if err != nil {
		fmt.Errorf("Не найден 'java' исполняемый файл: %v", err)
	}
	cmdJava := "java --version | head -n 1"
	java := exec.Command("/bin/sh", "-c", cmdJava)
	out, _ := java.Output()
	log.Println("java: OK,", strings.TrimSpace(string(out)))
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "update-jenkins",
	Short: "\n\tОбновление Jenkins.",
	Run:   func(cmd *cobra.Command, args []string) {},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// fmt.Println(str)

		if !str.plugin && !str.jenkins && !str.all {
			cmd.Help()
			os.Exit(0)
		}

		dependency(str)
		updateSystem(str, str.all, str.plugin, str.jenkins)
	},
}

func Execute() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolP("help", "h", false, "update [command] --help")
	rootCmd.PersistentFlags().BoolVarP(&str.all, "all-update", "", false, "update --all-update обновить полностью")
	rootCmd.PersistentFlags().BoolVarP(&str.jenkins, "jenkins", "", false, "update --jenkins обновить версию jenkins")
	rootCmd.PersistentFlags().BoolVarP(&str.plugin, "plugin", "", false, "update --plugin обновить плагины")
	rootCmd.PersistentFlags().StringVarP(&str.user, "usr", "u", "", "имя пользователя")
	rootCmd.PersistentFlags().StringVarP(&str.password, "pwd", "p", "", "пароль пользователя")
	rootCmd.PersistentFlags().StringVar(&str.url, "url", "http://jenkins.build2:8080", "адрес хоста")
	rootCmd.PersistentFlags().StringVar(&str.path, "path-cli", "/var/lib/jenkins/jenkins-cli.jar", "путь к файлу jenkins-cli.jar")
	rootCmd.PersistentFlags().StringVar(&str.pathWar, "path-war", "/usr/share/java/jenkins.war", "путь к файлу jenkins.war")
	rootCmd.PersistentFlags().StringVar(&str.pathJsn, "path-json", "/var/lib/jenkins/updates/default.json", "путь к файлу default.json")

}
