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
    DYNAMO_DASHBOARDS_TABLE = "dashboard.dashboard"
    DYNAMO_DASHBOARD_GRAPHS_TABLE = "dashboard.dashboard.graphs"
    DYNAMO_LOGIN_TABLE = "dashboard.login"
)

type DashboardService struct {
    gorest.RestService          `root:"/r" consumes:"application/json" produces:"application/json" realm:"dashauth"`

    getMetrics          gorest.EndPoint     `method:"GET" path:"/metric?{search:MetricSearch}" output:"MetricResults" role:"authed"`
    listMetrics         gorest.EndPoint     `method:"GET" path:"/metric/list/?{token:string}" output:"ListMetricsResponse" role:"authed"`
    saveGraph           gorest.EndPoint     `method:"POST" path:"/graph/?{id:string}" postdata:"SaveGraphRequest" output:"SaveGraphResponse" role:"authed"`
    getSavedGraphs      gorest.EndPoint     `method:"GET" path:"/graphs" output:"GetSavedGraphResponse" role:"authed"`

    saveDashboard       gorest.EndPoint     `method:"POST" path:"/dashboard/{id:string}" postdata:"SaveDashboardRequest" output:"SaveDashboardResponse" role:"authed"`
    createDashboard     gorest.EndPoint     `method:"PUT" path:"/dashboard" postdata:"PutDashboardRequest" output:"PutDashboardResponse" role:"authed"`
    getDashboards       gorest.EndPoint     `method:"GET" path:"/dashboards" output:"GetDashboardsResponse" role:"authed"`
    getDashboard        gorest.EndPoint     `method:"GET" path:"/dashboard/{id:string}" output:"GetDashboardResponse" role:"authed"`

    getLogGroups        gorest.EndPoint     `method:"GET" path:"/loggroups?{token:string}" output:"GetLogGroupsResponse" role:"authed"`
    getLogStreams       gorest.EndPoint     `method:"GET" path:"/loggroup/{name:string}/streams?{token:string}" output:"GetLogGroupStreamsResponse" role:"authed"`
    getLogs             gorest.EndPoint     `method:"GET" path:"/logs?{search:GetLogsRequest}" output:"GetLogsResponse" role:"authed"`

    login               gorest.EndPoint     `method:"POST" path:"/login" postdata:"LoginRequest" output:"LoginResponse" role:"open"`
}


type config struct {
    Files           string
}

func main() {
    conf := loadConfig()

    gorest.RegisterRealmAuthorizer("dashauth", SimpleAuther)
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
