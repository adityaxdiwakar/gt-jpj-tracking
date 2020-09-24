package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/PuerkitoBio/goquery"
	"github.com/bwmarrin/discordgo"
	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"
	"golang.org/x/net/context"
)

var rdb *redis.Client
var db *sql.DB
var conf tomlConfig
var ctx = context.Background()

type tomlConfig struct {
	Redis    redisCredentials
	Database postgresCredentials
	Webhook  []string
}

type redisCredentials struct {
	Address  string
	Password string
	DB       int
}

type postgresCredentials struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

func init() {
	if _, err := toml.DecodeFile("config.toml", &conf); err != nil {
		log.Fatalf("error: could not parse configuration %v\n", err)
	}

	rdb = redis.NewClient(&redis.Options{
		Addr:     conf.Redis.Address,
		Password: conf.Redis.Password,
		DB:       conf.Redis.DB,
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("error: could not make connection with redis: %v\n", err)
	}

	pSqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s "+
		"sslmode=disable", conf.Database.Host, conf.Database.Port,
		conf.Database.User, conf.Database.Password, conf.Database.DBName)

	db, err = sql.Open("postgres", pSqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
}

type JPJApiReport struct {
	Payload    []JPJPayload `json:"payload"`
	StatusCode int          `json:"status_code"`
}

type JPJPayload struct {
	ID       int    `json:"id"`
	Date     string `json:"date"`
	Reported int    `json:"reported"`
	Total    int    `json:"total"`
}

func averagePayload(slice []JPJPayload) float64 {
	sum := 0
	for _, k := range slice {
		sum += k.Reported
	}
	return float64(sum) / float64(len(slice))
}

func main() {
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).DialContext,
	}

	client := http.Client{Transport: transport}

	req, err := http.NewRequest("GET", "https://health.gatech.edu/coronavirus/health-alerts", nil)
	if err != nil {
		log.Fatal(err)
	}

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	teaser := doc.Find(".super-block__teaser").Nodes[0]
	table := goquery.NewDocumentFromNode(teaser).Children().Nodes[0]
	row := goquery.NewDocumentFromNode(table).Children().Nodes[1]
	tr := goquery.NewDocumentFromNode(row).Children().Nodes[0]
	data := goquery.NewDocumentFromNode(tr).Children().Nodes

	date := data[0].FirstChild.Data
	reported := data[1].FirstChild.Data
	aggregation := data[2].FirstChild.Data

	reported = strings.Replace(reported, "*", "", 1)

	previousDate, err := rdb.Get(ctx, "gt.cases.lastdate").Result()
	if previousDate != date {

		sqlStatement := `
        INSERT INTO cases (date, reported, total) 
        VALUES ($1, $2, $3)`

		_, err := db.Exec(sqlStatement, date, reported, aggregation)
		if err != nil {
			log.Fatal(err)
		} else {
			// only set redis value if DB insertion was successful
			rdb.Set(ctx, "gt.cases.lastdate", date, 0)
		}

		resp, err := client.Get("https://api.aditya.diwakar.io/gt-jpj/cases")
		if err != nil {
			log.Fatal(err)
		}

		defer resp.Body.Close()

		var jpjResponse JPJApiReport
		err = json.NewDecoder(resp.Body).Decode(&jpjResponse)
		if err != nil {
			log.Fatal(err)
		}

		payloadLength := len(jpjResponse.Payload)

		sevenDayMA := averagePayload(jpjResponse.Payload[payloadLength-7 : payloadLength])
		thirtyDayMA := averagePayload(jpjResponse.Payload[payloadLength-30 : payloadLength])

		webhookMessage := discordgo.WebhookParams{
			Username:  "GT Stamps Health Services",
			AvatarURL: "https://img.aditya.diwakar.io/stamps.png",
			Embeds: []*discordgo.MessageEmbed{
				{
					Title: fmt.Sprintf("[%s] GT COVID-19 Update", date),
					URL:   "https://health.gatech.edu/coronavirus/health-alerts",
					Color: 11772777,
					Footer: &discordgo.MessageEmbedFooter{
						Text: "Made with ❤️ by Aditya Diwakar",
					},
					Fields: []*discordgo.MessageEmbedField{
						{
							Name:   "Reported Today",
							Value:  reported,
							Inline: true,
						},
						{
							Name:   "Total",
							Value:  aggregation,
							Inline: true,
						},
						{
							Name:   "7/30 Day MA",
							Value:  fmt.Sprintf("%.1f/%.1f", sevenDayMA, thirtyDayMA),
							Inline: true,
						},
					},
				},
			},
		}

		jsonStr, _ := json.Marshal(webhookMessage)

		for _, wh := range conf.Webhook {
			req, err := http.NewRequest("POST", wh, bytes.NewBuffer(jsonStr))
			if err != nil {
				log.Println(err)
				continue
			}
			req.Header.Set("Content-Type", "application/json")

			res, err = client.Do(req)
			if err != nil {
				log.Println(err)
				continue
			}
		}

	}
}
