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
	"time"

	"github.com/spf13/cobra"
)

var (
	query, header []string

	rootCmd = &cobra.Command{
		Use:   "httpcli",
		Short: "",
		Long:  ``,
		Args: func(cmd *cobra.Command, args []string) error {
			return cobra.OnlyValidArgs(cmd, args)
		},
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				log.Fatalln("Please enter url")
			}
			get(args[0])
		},
	}
)

func getCmd() *cobra.Command {
	var get = &cobra.Command{
		Use:   "get",
		Short: "Get information",
		Long: `Get information from url by using command
		get <URL>`,
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			get(args[0])
			if len(query) > 0 {
				addQuery(args[0])
			}
		},
	}

	return get
}

func postCmd() *cobra.Command {
	var postBody string
	var postCmd = &cobra.Command{
		Use:   "post",
		Short: "Post information",
		Long: `Post information by using command
		post <URL> --json '{ "key": "value" }'`,
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			url := args[0]
			json_body := []byte(postBody)
			req, err := http.NewRequest("POST", url, bytes.NewBuffer(json_body))
			if err != nil {
				log.Fatalln(err)
			}
			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{
				Timeout: 10 * time.Second,
			}
			resp, err := client.Do(req)
			if err != nil {
				log.Fatalln(err)
			}
			defer resp.Body.Close()
			Body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatalln(err)
			}

			res, err := prettyString(string(Body))
			if err != nil {
				log.Fatalln(err)
			}

			fmt.Println(res)
			if len(query) > 0 {
				addQuery(args[0])
			}
		},
	}

	postCmd.Flags().StringVar(&postBody, "json", "", "Enter json body")
	return postCmd
}

func deleteCmd() *cobra.Command {
	var deleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete information",
		Long: `Delete information in the url by using command
		delete <URL>`,
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			url := args[0]
			req, err := http.NewRequest("DELETE", url, nil)
			if err != nil {
				log.Fatalln(err)
			}

			client := &http.Client{
				Timeout: 10 * time.Second,
			}
			resp, err := client.Do(req)
			if err != nil {
				log.Fatalln(err)
			}
			fmt.Println("Status: ", http.StatusText(resp.StatusCode))
			defer resp.Body.Close()

			if len(query) > 0 {
				addQuery(args[0])
			}
		},
	}

	return deleteCmd
}

func putCmd() *cobra.Command {
	var putBody string
	var putCmd = &cobra.Command{
		Use:   "put",
		Short: "Put updated information",
		Long: `Put an updated information to the url by using command
		put <url> --json "{ 'key': 'value' }"`,
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			url := args[0]
			json_body := []byte(putBody)

			req, err := http.NewRequest("PUT", url, bytes.NewBuffer(json_body))
			if err != nil {
				log.Fatalln(err)
			}
			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{
				Timeout: 10 * time.Second,
			}
			resp, err := client.Do(req)
			if err != nil {
				log.Fatalln(err)
			}
			defer resp.Body.Close()
			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatalln(err)
			}

			res, err := prettyString(string(b))
			if err != nil {
				log.Fatalln(err)
			}

			fmt.Println(res)

			if len(query) > 0 {
				addQuery(args[0])
			}
		},
	}

	putCmd.Flags().StringVarP(&putBody, "json", "", "", "Enter json boy")
	return putCmd
}

func addHeader(r *http.Request) {
	for i := 0; i < len(header); i++ {
		header := strings.Split(header[i], "=")
		r.Header.Add(header[0], header[1])
	}
}

func addQuery(str string) {
	u, err := url.Parse(str)
	if err != nil {
		log.Fatalln(err)
	}

	q, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		log.Fatalln(err)
	}
	for i := 0; i < len(query); i++ {
		header := strings.Split(query[i], "=")
		q.Add(header[0], header[1])
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
		addHeader(req)
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

	res, err := prettyString(string(b))
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(res)
}

func prettyString(str string) (string, error) {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, []byte(str), "", "   "); err != nil {
		return "", err
	}
	return prettyJSON.String(), nil
}

func init() {
	rootCmd.PersistentFlags().StringSliceVar(&query, "query", []string{}, "Query")
	rootCmd.PersistentFlags().StringSliceVar(&header, "header", []string{}, "Header")

	rootCmd.AddCommand(getCmd())
	rootCmd.AddCommand(postCmd())
	rootCmd.AddCommand(deleteCmd())
	rootCmd.AddCommand(putCmd())
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}

}
