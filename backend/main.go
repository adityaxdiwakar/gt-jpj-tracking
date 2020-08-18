package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	_ "github.com/lib/pq"
)

var db *sql.DB
var conf tomlConfig

type tomlConfig struct {
	Database postgresCredentials
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

	pSqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s "+
		"sslmode=disable", conf.Database.Host, conf.Database.Port,
		conf.Database.User, conf.Database.Password, conf.Database.DBName)

	var err error
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
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/gt-jpj", homePage)
	r.Get("/gt-jpj/cases", getAllCases)
	r.Get("/gt-jpj/testing", getAllSurveys)

	http.ListenAndServe(":3000", r)
}

func homePage(w http.ResponseWriter, r *http.Request) {
	data := StringResponse{
		Payload: "Welcome to the Georgia Institute of Technology JPJ Tracking API - Not affiliated with Georgia Tech",
		Code:    200,
	}

	json.NewEncoder(w).Encode(data)
}

type CasesRow struct {
	ID       int    `json:"id"`
	Date     string `json:"date"`
	Reported int    `json:"reported"`
	Total    int    `json:"total"`
}

type CaseResponse struct {
	Payload []CasesRow `json:"payload"`
	Code    int        `json:"status_code"`
}

type StringResponse struct {
	Payload string `json:"payload"`
	Code    int    `json:"status_code"`
}

func getAllCases(w http.ResponseWriter, r *http.Request) {
	statement := `SELECT * FROM cases`
	rows, err := db.Query(statement)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(StringResponse{
			Code:    500,
			Payload: "Internal Server Error",
		})
		return
	}

	caseData := make([]CasesRow, 0)
	defer rows.Close()

	for rows.Next() {
		day := CasesRow{}

		if err := rows.Scan(&day.ID, &day.Date, &day.Reported, &day.Total); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(StringResponse{
				Code:    500,
				Payload: "Internal Server Error",
			})
			return
		}

		day.Date = strings.TrimSpace(day.Date)

		caseData = append(caseData, day)
	}

	uniqueCaseData := make([]CasesRow, 0)
	uniqueMap := make(map[string]bool)

	for i, day := range caseData {
		if uniqueMap[day.Date] == false {
			uniqueMap[day.Date] = true
			caseData[i].ID = len(uniqueMap)
			uniqueCaseData = append(uniqueCaseData, caseData[i])
		}
	}

	data := CaseResponse{
		Payload: uniqueCaseData,
		Code:    200,
	}

	json.NewEncoder(w).Encode(data)

}

type SurveysRow struct {
	ID           int    `json:"id"`
	Date         string `json:"date"`
	Positive     int    `json:"reported"`
	Administered int    `json:"total"`
}

type SurveyResponse struct {
	Payload []SurveysRow `json:"payload"`
	Code    int          `json:"status_code"`
}

func getAllSurveys(w http.ResponseWriter, r *http.Request) {
	statement := `SELECT * FROM surveys`
	rows, err := db.Query(statement)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(StringResponse{
			Code:    500,
			Payload: "Internal Server Error",
		})
		return
	}

	surveyData := make([]SurveysRow, 0)
	defer rows.Close()

	for rows.Next() {
		day := SurveysRow{}

		if err := rows.Scan(&day.ID, &day.Date, &day.Positive, &day.Administered); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(StringResponse{
				Code:    500,
				Payload: "Internal Server Error",
			})
			return
		}

		day.Date = strings.TrimSpace(day.Date)

		surveyData = append(surveyData, day)
	}

	uniqueSurveyData := make([]SurveysRow, 0)
	uniqueMap := make(map[string]bool)

	for i, day := range surveyData {
		if uniqueMap[day.Date] == false {
			uniqueMap[day.Date] = true
			surveyData[i].ID = len(uniqueMap)
			uniqueSurveyData = append(uniqueSurveyData, surveyData[i])
		}
	}

	log.Println(uniqueSurveyData)

	data := SurveyResponse{
		Payload: uniqueSurveyData,
		Code:    200,
	}

	json.NewEncoder(w).Encode(data)

}
