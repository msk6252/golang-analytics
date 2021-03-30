package main

import (
  "fmt"
  "io/ioutil"
  "golang.org/x/oauth2"
  "golang.org/x/oauth2/google"
  "github.com/joho/godotenv"
  "os"
  "log"
  analytics "google.golang.org/api/analytics/v3"
)

func main() {
  key, _ := ioutil.ReadFile("./secret.json")

  viewId := GetEnv("ViewID")

  jwtConf, err := google.JWTConfigFromJSON(
    key,
    analytics.AnalyticsReadonlyScope,
  )

  if err != nil {
    log.Fatalln(err)
  }

  httpClient := jwtConf.Client(oauth2.NoContext)

  svc, err := analytics.New(httpClient)

  if err != nil {
    log.Fatalln(err)
  }

  result, err := svc.Data.Ga.Get("ga:" + viewId, "7daysAgo", "today", "ga:pageviews").Dimensions("ga:pagePath").Filters("ga:pagePath=~^/archives/").Sort("-ga:pageviews").MaxResults(20).Do()
  for _, row := range result.Rows {
    fmt.Printf("%s,%s\n", row[0], row[1])
  }
}

func GetEnv(key string) string {
  err := godotenv.Load()

  if err != nil {
    log.Fatalln("Get env error:", err)
  }

  value := os.Getenv(key)

  if value == "" {
    log.Fatalln("Empty .env")
  }

  return value
}
