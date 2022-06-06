package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/spf13/cobra"
)

var query []string
var header []string
var postBody string
var putBody string

var (
	rootCmd = &cobra.Command{
		Use:   "httpcli",
		Short: "",
		Long:  ``,
		Args: func(cmd *cobra.Command, args []string) error {
			return cobra.OnlyValidArgs(cmd, args)
		},
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println("Please enter url")
			} else {
				get(args[0])
			}
		},
	}

	getCmd = &cobra.Command{
		Use:   "get",
		Short: "Get information",
		Long: `Get information from url by using command
		get <URL>`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println("Please Enter url")
			} else {
				get(args[0])
				if len(query) > 0 {
					Query(args[0])
				}
			}
		},
	}

	postCmd = &cobra.Command{
		Use:   "post",
		Short: "Post information",
		Long: `Post information by using command
		post <URL> --json '{ "key": "value" }'`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println("Please Enter url")
			} else {
				url := args[0]
				json_body := []byte(postBody)

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
				if len(query) > 0 {
					Query(args[0])
				}

			}
		},
	}

	deleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete information",
		Long: `Delete information in the url by using command
		delete <URL>`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println("Please Enter url")
			} else {
				url := args[0]
				req, err := http.NewRequest("DELETE", url, nil)
				if err != nil {
					log.Fatalln(err)
				}

				client := &http.Client{}
				resp, err := client.Do(req)
				if err != nil {
					log.Fatalln(err)
				}
				fmt.Println("Status: ", http.StatusText(resp.StatusCode))
				defer resp.Body.Close()

				if len(query) > 0 {
					Query(args[0])
				}

			}

		},
	}

	putCmd = &cobra.Command{
		Use:   "put",
		Short: "Put updated information",
		Long: `Put an updated information to the url by using command
		put <url> --json "{ 'key': 'value' }"`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println("Please Enter url")
			} else {
				url := args[0]
				json_body := []byte(putBody)

				req, err := http.NewRequest("PUT", url, bytes.NewBuffer(json_body))
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

				if len(query) > 0 {
					Query(args[0])
				}

			}
		},
	}
)

func Header(r *http.Request) {
	for i := 0; i < len(header); i++ {
		res1 := strings.Split(header[i], "=")
		r.Header.Add(res1[0], res1[1])
	}
}

func Query(str string) {
	u, _ := url.Parse(str)

	q, _ := url.ParseQuery(u.RawQuery)
	for i := 0; i < len(query); i++ {
		res1 := strings.Split(query[i], "=")
		q.Add(res1[0], res1[1])
	}

	u.RawQuery = q.Encode()

	fmt.Println(u)

}

func get(url string) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln(err)
	}
	if len(header) > 0 {
		Header(req)
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
	rootCmd.PersistentFlags().StringSliceVar(&query, "query", []string{}, "Query")
	rootCmd.PersistentFlags().StringSliceVar(&header, "header", []string{}, "Header")

	postCmd.Flags().StringVar(&postBody, "json", "", "Enter json body")
	rootCmd.AddCommand(postCmd)

	rootCmd.AddCommand(deleteCmd)

	putCmd.Flags().StringVarP(&putBody, "json", "", "", "Enter json boy")
	rootCmd.AddCommand(putCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}

}
