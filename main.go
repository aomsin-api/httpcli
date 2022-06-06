package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/spf13/cobra"
)

var body string

var (
	rootCmd = &cobra.Command{
		Use:   "httpcli",
		Short: "An example cobra program",
		Long:  `This is a simple example`,
	}

	getCmd = &cobra.Command{
		Use:   "get",
		Short: "",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println("Please Enter url")
			} else {
				printData(args[0])
			}
		},
	}

	postCmd = &cobra.Command{
		Use:   "post",
		Short: "",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println("Please Enter url")
			} else {
				url := args[0]
				json_body := []byte(body)
				fmt.Println(url, bytes.NewBuffer(json_body))
				req, err := http.NewRequest("POST", url, bytes.NewBuffer(json_body))
				if err != nil {
					log.Fatalln(err)
				}
				req.Header.Set("Content-Type", "application/json")

				client := &http.Client{}
				resp, err := client.Do(req)
				if err != nil {
					panic(err)
				}
				defer resp.Body.Close()
				Body, _ := ioutil.ReadAll(resp.Body)

				res, err := PrettyString(string(Body))
				if err != nil {
					log.Fatalln(err)
				}

				fmt.Println(res)

			}
		},
	}
)

func printData(url string) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Set("Accept", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	res, err := PrettyString(string(b))
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(res)
}

func PrettyString(str string) (string, error) {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, []byte(str), "", "   "); err != nil {
		return "", err
	}
	return prettyJSON.String(), nil
}

func init() {
	rootCmd.AddCommand(getCmd)

	postCmd.Flags().StringVarP(&body, "json", "j", "", "Enter json body")
	rootCmd.AddCommand(postCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}

}
