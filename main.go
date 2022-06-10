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
				header, err := cmd.Flags().GetStringSlice("header")
				if err != nil {
					log.Fatalln(err)
				}
				get(args[0], header)
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
				header, err := cmd.Flags().GetStringSlice("header")
				if err != nil {
					log.Fatalln(err)
				}
				get(args[0], header)
				query, err := cmd.Flags().GetStringSlice("query")
				if err != nil {
					log.Fatalln(err)
				}
				if len(query) > 0 {
					addQuery(args[0], query)
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
			// Function is too long. You should break this function down to smaller functions.

			if len(args) == 0 {
				fmt.Println("Please Enter url")
			} else {
				// The else block is unneccessary.
				url := args[0]
				postBody, err := cmd.Flags().GetString("json")
				if err != nil {
					log.Fatalln(err)
				}

				json_body := []byte(postBody) // Don't use underscore_case case in Golang. Please use camelCase instead.

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

				query, err := cmd.Flags().GetStringSlice("query")
				if err != nil {
					log.Fatalln(err)
				}
				if len(query) > 0 {
					addQuery(args[0], query)
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

				client := &http.Client{
					Timeout: 10 * time.Second,
				}
				resp, err := client.Do(req)
				if err != nil {
					log.Fatalln(err)
				}
				fmt.Println("Status: ", http.StatusText(resp.StatusCode))
				defer resp.Body.Close()

				query, err := cmd.Flags().GetStringSlice("query")
				if err != nil {
					log.Fatalln(err)
				}
				if len(query) > 0 {
					addQuery(args[0], query)
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
				putBody, err := cmd.Flags().GetString("json")
				if err != nil {
					log.Fatalln(err)
				}
				jsonBody := []byte(putBody)

				req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonBody))
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

				query, err := cmd.Flags().GetStringSlice("query")
				if err != nil {
					log.Fatalln(err)
				}
				if len(query) > 0 {
					addQuery(args[0], query)
				}

			}
		},
	}
)

// Function should be named `addHeadersToRequest`
func addHeader(header []string, r *http.Request) {
	for i := 0; i < len(header); i++ {
		h := strings.Split(header[i], "=")
		r.Header.Add(h[0], h[1])
	}
}

// This function doesn't return anything. The product of this function doesn't seem
// to be used anywhere?
//
// Could u please clarify what this function does?
func addQuery(str string, query []string) {
	// Variable naming is hard to read.
	// You should name variables meaningfully instead of using a single character.
	u, err := url.Parse(str)
	if err != nil {
		log.Fatalln(err)
	}

	q, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		log.Fatalln(err)
	}
	for i := 0; i < len(query); i++ {
		qu := strings.Split(query[i], "=")
		q.Add(qu[0], qu[1])
	}

	u.RawQuery = q.Encode()

	fmt.Println(u)

}

func get(url string, header []string) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln(err)
	}
	if len(header) > 0 {
		addHeader(header, req)
	}
	req.Header.Set("Accept", "application/json")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}
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

// init function should be at the top of the file.
func init() {
	// Group the same methods together.
	// The way these statements are organized is a little hard to read.
	rootCmd.AddCommand(getCmd)
	rootCmd.PersistentFlags().StringSlice("query", []string{}, "query")
	rootCmd.PersistentFlags().StringSlice("header", []string{}, "header")

	postCmd.Flags().String("json", "", "Post body")
	rootCmd.AddCommand(postCmd)

	rootCmd.AddCommand(deleteCmd)

	putCmd.Flags().String("json", "", "Put body")
	rootCmd.AddCommand(putCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}

}
