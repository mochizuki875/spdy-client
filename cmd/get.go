package cmd

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/spf13/cobra"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

var getUrl string

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "GET Request to URL.",
	Long: `Example:
	GET request to kubelet /exec API to run a command in a container.
	spdy-client get --url="https://<Kubernetes Node IP Address>:10250/exec/default/nginx/nginx?command=/bin/bash&input=1&output=1&tty=1"`,
	Run: getRun,
}

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.Flags().StringVar(&getUrl, "url", "", "Request URL.")
}

func getRun(cmd *cobra.Command, args []string) {
	method := "POST"

	u, err := url.Parse(getUrl)
	if err != nil {
		panic(err)
	}

	server := u.Host
	apiPath := u.Path
	queryCommands := u.RawQuery
	serverFullAddress := fmt.Sprintf("https://%s", server)

	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
	}

	config := &rest.Config{
		Host:    serverFullAddress,
		APIPath: apiPath,
		TLSClientConfig: rest.TLSClientConfig{
			Insecure: true,
		},
		Transport:   tr,
		BearerToken: "",
	}

	urlObject := &url.URL{
		Scheme:   "https",
		Opaque:   "",
		User:     nil,
		Host:     server,
		Path:     apiPath,
		RawPath:  "",
		RawQuery: queryCommands,
	}

	exec, err := remotecommand.NewSPDYExecutor(config, method, urlObject)
	if err != nil {
		fmt.Println(err)
	}

	err = exec.StreamWithContext(context.TODO(), remotecommand.StreamOptions{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Tty:    true,
	})

	if err != nil {
		fmt.Println(err)
	}

}
