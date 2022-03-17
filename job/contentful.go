package job

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/S-Ryouta/sample-blog/db"
	"github.com/S-Ryouta/sample-blog/models"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

// NOTE: JSON-to-GOでcontentfulのレスポンス構造体を作成
//       ref: https://mholt.github.io/json-to-go/

type EntitiesResponse struct {
	Sys struct {
		Type string `json:"type"`
	} `json:"sys"`
	Total int `json:"total"`
	Skip  int `json:"skip"`
	Limit int `json:"limit"`
	Items []struct {
		Metadata struct {
			Tags []interface{} `json:"tags"`
		} `json:"metadata"`
		Sys struct {
			Space struct {
				Sys struct {
					Type     string `json:"type"`
					LinkType string `json:"linkType"`
					ID       string `json:"id"`
				} `json:"sys"`
			} `json:"space"`
			ID          string    `json:"id"`
			Type        string    `json:"type"`
			CreatedAt   time.Time `json:"createdAt"`
			UpdatedAt   time.Time `json:"updatedAt"`
			Environment struct {
				Sys struct {
					ID       string `json:"id"`
					Type     string `json:"type"`
					LinkType string `json:"linkType"`
				} `json:"sys"`
			} `json:"environment"`
			Revision    int `json:"revision"`
			ContentType struct {
				Sys struct {
					Type     string `json:"type"`
					LinkType string `json:"linkType"`
					ID       string `json:"id"`
				} `json:"sys"`
			} `json:"contentType"`
			Locale string `json:"locale"`
		} `json:"sys"`
		Fields struct {
			Title string `json:"title"`
			Body  struct {
				Data struct {
				} `json:"data"`
				Content []struct {
					Data struct {
					} `json:"data"`
					Content []struct {
						Data struct {
						} `json:"data"`
						Marks    []interface{} `json:"marks"`
						Value    string        `json:"value"`
						NodeType string        `json:"nodeType"`
					} `json:"content"`
					NodeType string `json:"nodeType"`
				} `json:"content"`
				NodeType string `json:"nodeType"`
			} `json:"body"`
		} `json:"fields"`
	} `json:"items"`
}

func FetchContentful(ctx context.Context) {
	counter := 0
	waitTime := 30 * time.Second
	ticker := time.NewTicker(waitTime)
	defer ticker.Stop()
	child, childCancel := context.WithCancel(ctx)
	defer childCancel()

	// NOTE: infinite loop
	for {
		select {
		case t := <-ticker.C:
			counter++
			requestID := counter
			log.Println("[DEBUG] START taskNo=", requestID, "t=", t)

			errCh := make(chan error, 1)
			go func() { // 登録したタスクをブロックせずに実行
				errCh <- getRequest()
			}()

			go func() {
				// error channelにリクエストの結果が返ってくるのを待つ
				select {
				case err := <-errCh:
					if err != nil {
						// Deamonの強制終了
						log.Println("[ERROR] ", err)
					}
					log.Println("[DEBUG] END requestNo=", requestID)
				}
			}()
		case <-child.Done():
			return
		}
	}
}

func getRequest() error {
	spaceId := os.Getenv("CONTENTFUL_SPACE_ID")
	accessToken := os.Getenv("CONTENTFUL_ACCESS_TOKEN")
	url := "https://cdn.contentful.com/spaces/" + spaceId + "/entries?access_token=" + accessToken

	req, err := http.NewRequest(http.MethodGet, url, nil)
	client := new(http.Client)
	response, _ := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)
	jsonBytes := ([]byte)(string(body))
	data := new(EntitiesResponse)

	if err := json.Unmarshal(jsonBytes, data); err != nil {
		fmt.Println("JSON Unmarshal error:", err)
		return nil
	}

	db := db.Connect()
	mysqlClient, e := db.DB()

	if e != nil {
		fmt.Println("failed to connect to DB:", err)
		return nil
	}
	defer mysqlClient.Close()

	for _, item := range data.Items {
		body, _ := json.Marshal(item.Fields.Body)
		uuid := item.Sys.ID
		title := item.Fields.Title
		description := ""
		bodyString := string(body)

		entry := models.Entry{ID: uuid, Title: title, Description: description, Body: bodyString}
		models.AddOrUpdateEntry(db, entry)
	}

	return nil
}
