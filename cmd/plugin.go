package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var pluginCmd = &cobra.Command{
	Use:   "plugin",
	Short: "Обновить плагины",
	Run: func(cmd *cobra.Command, args []string) {
		cmdStr := fmt.Sprintf("java -jar %s -s %s list-plugins | grep ')$' | awk '{print $1}' | tr '\n' ' '", str.path, str.GetHttpConnection())
		list := exec.Command("/bin/sh", "-c", cmdStr)
		listPlugins, err := list.Output()
		if err != nil {
			fmt.Println("Error command::", err)
		} else {
			strCm := `java -jar ./jenkins-cli.jar -s http://jenkins:8080 -auth ***:*** | grep ')$' | awk '{print $1}' | tr '\n' ' '`
			log.Println("Start:", strCm)
		}

		cmdUpdateStr := fmt.Sprintf("java -jar %s -s %s install-plugin %s -restart", str.path, str.GetHttpConnection(), strings.TrimSpace(string(listPlugins)))
		cmdUpdate := exec.Command("/bin/sh", "-c", cmdUpdateStr)
		cmdUpdate.Stdout = os.Stdout
		cmdUpdate.Stderr = os.Stderr
		if err = cmdUpdate.Run(); err != nil {
			log.Println("Error command::", err)
			if strings.TrimSpace(string(listPlugins)) == "" {
				log.Println("No plugins for updates.")
			}
		} else {
			strCm := `java -jar ./jenkins-cli.jar -s http://jenkins:8080 -auth ***:*** install-plugin ` + strings.TrimSpace(string(listPlugins)) + ` -restart`
			log.Println("Start:", strCm)
		}
	},
}

func init() {}
