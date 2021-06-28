package cmd

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func RunMizuFetch(fetch *MizuFetchOptions) {
	resp, err := http.Get(fmt.Sprintf("http://localhost:%v/api/har?from=%v&to=%v", fetch.MizuPort, fetch.FromTimestamp, fetch.ToTimestamp))
	if err != nil {
		log.Fatal(err)
	}

	defer func() { _ = resp.Body.Close() }()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	zipReader, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		log.Fatal(err)
	}
	_ = Unzip(zipReader, fetch.Directory)

}

// func FilterRequests(body []byte) {
// 	enforcePolicy, _ := decodeEnforcePolicy()
// 	var HARFetched shared.HARFetched
// 	err := json.Unmarshal(body, &HARFetched)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	for key, entry := range HARFetched.Log.Entries {
// 		response := entry.Response
// 		for _, value := range enforcePolicy.Rules {
// 			if value.Type == "json" {
// 				var bodyJsonMap map[string]interface{}
// 				err := json.Unmarshal(entry.Content.Text, &bodyJsonMap)
// 				if err != nil {
// 					var bodyJsonMap []interface{}
// 					err := json.Unmarshal(entry.Content.Text, &bodyJsonMap)
// 					if err != nil {
// 						fmt.Println(err)
// 					}
// 				} else {
// 					result := map[string]bool{}
// 					if bodyJsonMap[value.Key] != value.Value {

// 					}
// 				}
// 			} else {

// 			}
// 		}
// 	}
// }

// func decodeEnforcePolicy() (shared.RulesPolicy, error) {
// 	content, err := ioutil.ReadFile("/app/enforce-policy/enforce-policy.yaml")
// 	enforcePolicy := shared.RulesPolicy{}
// 	if err != nil {
// 		return enforcePolicy, err
// 	}
// 	err = yaml.Unmarshal([]byte(content), &enforcePolicy)
// 	if err != nil {
// 		return enforcePolicy, err
// 	}
// 	invalidIndex := enforcePolicy.ValidateRulesPolicy()
// 	if len(invalidIndex) != 0 {
// 		for i := range invalidIndex {
// 			fmt.Println("only json and header types are supported on rule")
// 			enforcePolicy.RemoveNotValidPolicy(invalidIndex[i])
// 		}
// 	}
// 	return enforcePolicy, nil
// }

func Unzip(reader *zip.Reader, dest string) error {
	dest, _ = filepath.Abs(dest)
	_ = os.MkdirAll(dest, os.ModePerm)

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		path := filepath.Join(dest, f.Name)

		// Check for ZipSlip (Directory traversal)
		if !strings.HasPrefix(path, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", path)
		}

		if f.FileInfo().IsDir() {
			_ = os.MkdirAll(path, f.Mode())
		} else {
			_ = os.MkdirAll(filepath.Dir(path), f.Mode())
			fmt.Print("writing HAR file [ ", path, " ] .. ")
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
				fmt.Println(" done")
			}()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, f := range reader.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return err
		}
	}

	return nil
}
