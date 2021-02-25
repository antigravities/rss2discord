package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/mmcdole/gofeed"
)

type outgoingWebhook struct {
	Content string `json:"content"`
}

func main() {
	var file *os.File

	feed := flag.String("feed", "", "The feed to fetch from. This can be any feed type https://github.com/mmcdole/gofeed supports.")
	hook := flag.String("hook", "", "The Discord Webhook to send to (i.e. https://discord.com/api/webhooks/...)")
	fn := flag.String("data", "", "Location of a database file to write. If provided, rss2discord will \"remember\" the last item it sent and not post again until there's a new item at the top of the feed. This file can be shared across feeds.")
	dry := flag.Bool("dry", false, "Execute a dry-run: fetch the feed and write data files if applicable, but don't post anywhere.")

	flag.Parse()

	if *feed == "" || *hook == "" {
		log.Fatalf("Any of -feed or -hook are missing. Specify them and try again. For help, " + os.Args[0] + " -help.")
	}

	fp := gofeed.NewParser()
	fd, err := fp.ParseURL(*feed)
	if err != nil {
		log.Fatalf("Error fetching feed: %v", err)
	}

	if len(fd.Items) == 0 {
		log.Printf("Feed had no items")
		return
	}

	data := make(map[string]string)

	if *fn != "" {
		existed := true
		if _, err := os.Stat(*fn); os.IsNotExist(err) {
			existed = false
		}

		file, err = os.OpenFile(*fn, os.O_CREATE|os.O_RDWR, 0666)
		if err != nil {
			log.Printf("Error opening data file for reading: %v", err)
			log.Fatal("To continue without a data file, omit the option and run again.")
		}
		defer file.Close()

		if existed {
			err = json.NewDecoder(file).Decode(&data)
			if err != nil {
				log.Printf("Error reading data file: %v", err)
				log.Fatal("To continue without a data file, omit the option and run again.")
			}
		}
	}

	hadItem := false
	item, ok := data[*feed]
	if ok {
		if item == fd.Items[0].Title {
			hadItem = true
		}
	}

	if !hadItem && !*dry {
		link := ""
		if fd.Items[0].Link != "" {
			link = " - <" + fd.Items[0].Link + ">"
		}

		bx, err := json.Marshal(&outgoingWebhook{
			Content: fd.Items[0].Title + link,
		})
		if err != nil {
			log.Fatalf("Error marshaling outgoing Webhook: %v", err)
		}

		res, err := http.Post(*hook+"?wait=true", "application/json", bytes.NewReader(bx))
		if err != nil {
			log.Fatalf("Error posting to Webhook: %v", err)
		}
		defer res.Body.Close()

		if res.StatusCode != 200 {
			msg, _ := io.ReadAll(res.Body)

			log.Fatalf("Error posting to Webhook: expected status code 200, got %d: %s", res.StatusCode, msg)
		}
	}

	data[*feed] = fd.Items[0].Title

	bx, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("Error marshaling data file: %v", err)
	}

	file.Truncate(0)
	file.Seek(0, 0)
	file.Write(bx)
	// defer file.Close()
}
