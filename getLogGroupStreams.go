package main


import (
    "github.com/fromkeith/awsgo"
    "github.com/fromkeith/awsgo/cloudwatch"
    "log"
)

type GetLogGroupStreamsResponse struct {
    StreamNames         []string
    Token               string
}

func (serv DashboardService) GetLogStreams(logGroup, token string) (resp GetLogGroupStreamsResponse) {
    streams := cloudwatch.NewDescribeLogStreamsRequest()
    streams.NextToken = token
    streams.LogGroupName = logGroup
    streams.Key, _ = awsgo.GetSecurityKeys()
    streamResp, err := streams.Request()
    if err != nil {
        log.Println("Error getting log streams %#v", err)
        serv.ResponseBuilder().SetResponseCode(500)
        return
    }

    resp.StreamNames = make([]string, len(streamResp.LogStreams))

    for i := range streamResp.LogStreams {
        resp.StreamNames[i] = streamResp.LogStreams[i].LogStreamName
    }
    resp.Token = streamResp.NextToken
    return
}