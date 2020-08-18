package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
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

func main() {
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).DialContext,
	}

	client := http.Client{Transport: transport}

	req, err := http.NewRequest("GET", "https://health.gatech.edu/surveillance-testing-program-results", nil)
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

	// header scraping to find the date
	header := goquery.NewDocumentFromNode(table).Children().Nodes[0]
	headerTr := goquery.NewDocumentFromNode(header).Children().Nodes[0]
	headerTh := goquery.NewDocumentFromNode(headerTr).Children().Nodes[0]
	headerP := goquery.NewDocumentFromNode(headerTh).Children().Nodes[0]
	date := headerP.FirstChild.LastChild.FirstChild.Data

	row := goquery.NewDocumentFromNode(table).Children().Nodes[1]
	tr := goquery.NewDocumentFromNode(row).Children().Nodes
	tdA := goquery.NewDocumentFromNode(tr[0]).Children().Nodes[1]
	tdB := goquery.NewDocumentFromNode(tr[1]).Children().Nodes[1]

	positive := tdA.FirstChild.NextSibling.FirstChild.Data
	total := tdB.FirstChild.NextSibling.FirstChild.Data

	positive = strings.Replace(positive, ",", "", -1)
	total = strings.Replace(total, ",", "", -1)

	positiveInt, err := strconv.Atoi(positive)
	if err != nil {
		log.Fatal(err)
	}
	totalInt, err := strconv.Atoi(total)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(positiveInt, totalInt)

	previousDate, _ := rdb.Get(ctx, "gt.survey.lastdate").Result()
	if previousDate != date {
		rdb.Set(ctx, "gt.survey.lastdate", date, 0)

		sqlStatement := `
        INSERT INTO survey (date, positive, administered) 
        VALUES ($1, $2, $3)`

		_, err := db.Exec(sqlStatement, date, positiveInt, totalInt)
		if err != nil {
			log.Fatal(err)
		}

		webhookMessage := discordgo.WebhookParams{
			Username:  "GT Stamps Health Services",
			AvatarURL: "https://img.aditya.diwakar.io/stamps.png",
			Embeds: []*discordgo.MessageEmbed{
				{
					Title: fmt.Sprintf("[%s] Surveillance Testing Program Results ", date),
					URL:   "https://health.gatech.edu/coronavirus/health-alerts",
					Color: 11772777,
					Footer: &discordgo.MessageEmbedFooter{
						Text: "Made with ❤️ by Aditya Diwakar",
					},
					Fields: []*discordgo.MessageEmbedField{
						{
							Name:   "Tested Positive (All Time)",
							Value:  positive,
							Inline: true,
						},
						{
							Name:   "Tests Administered",
							Value:  total,
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

			_, err = client.Do(req)
			if err != nil {
				log.Println(err)
				continue
			}
		}

	}
}
