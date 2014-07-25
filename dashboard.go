package main

import (
    "encoding/json"
    "github.com/fromkeith/gorest"
    "io/ioutil"
    "log"
    "net/http"
    "os"
)

const (
    DYNAMO_REGION = "us-west-2"
    DYNAMO_GRAPHS_TABLE = "dashboard.graphs"
)

type DashboardService struct {
    gorest.RestService          `root:"/r" consumes:"application/json" produces:"application/json"`

    getMetrics          gorest.EndPoint     `method:"GET" path:"/metric?{search:MetricSearch}" output:"MetricResults"`
    listMetrics         gorest.EndPoint     `method:"GET" path:"/metric/list/?{token:string}" output:"ListMetricsResponse"`
    saveGraph           gorest.EndPoint     `method:"POST" path:"/graph/?{id:string}" postdata:"SaveGraphRequest" output:"SaveGraphResponse"`
    getSavedGraphs      gorest.EndPoint     `method:"GET" path:"/graphs" output:"GetSavedGraphResponse"`
}


type config struct {
    Files           string
}

func main() {
    conf := loadConfig()

    gorest.RegisterService(new(DashboardService))
    http.Handle("/r/", gorest.Handle())
    http.Handle("/", http.FileServer(http.Dir(conf.Files)))
    http.ListenAndServe(":8080", nil)
}

func loadConfig() config {
    f, err := os.Open("dashboard.json")
    if err != nil {
        log.Fatalf("Error loading config! %v", err)
        return config{}
    }
    defer f.Close()
    if b, err := ioutil.ReadAll(f); err != nil {
        log.Fatalf("Error loading config! %v", err)
        return config{}
    } else {
        var conf config
        if err := json.Unmarshal(b, &conf); err != nil {
            log.Fatalf("Error loading config! %v", err)
            return config{}
        }
        return conf
    }
}
