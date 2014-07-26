package main


import (
    "github.com/fromkeith/awsgo"
    "github.com/fromkeith/awsgo/cloudwatch"
    "log"
    "time"
)

type GetLogsRequest struct {
    Token               string
    LogGroupName        string
    LogStreamName       string
    Start               int64
    End                 int64
}

type GetLogsResponse struct {
    NextToken           string
    Events              []cloudwatch.LogEventsResponse
}


func (serv DashboardService) GetLogs(search GetLogsRequest) (resp GetLogsResponse) {
    get := cloudwatch.NewGetLogEventsRequest()
    get.LogGroupName = search.LogGroupName
    get.LogStreamName = search.LogStreamName
    get.NextToken = search.Token
    get.SetTimeRange(time.Unix(search.Start / 1000, 0), time.Unix(search.End / 1000, 0))
    get.Key, _ = awsgo.GetSecurityKeys()
    getResp, err := get.Request()
    if err != nil {
        log.Printf("Error getting logs %#v\n", err)
        serv.ResponseBuilder().SetResponseCode(500)
        return
    }
    resp.NextToken = getResp.NextForwardToken
    resp.Events = getResp.Events
    return
}