package main


import (
    "github.com/fromkeith/awsgo"
    "github.com/fromkeith/awsgo/cloudwatch"
    "log"
)


type ListMetricsResponse struct {
    Metrics             []cloudwatch.MetricInfo
    NextToken           string
}

func (serv DashboardService) ListMetrics(token string) (myResp ListMetricsResponse) {
    region := serv.Context.Request().Header.Get("Region")
    list := cloudwatch.NewListMetricsRequest()
    list.NextToken = token
    list.Host.Region = region
    list.Key, _ = awsgo.GetSecurityKeys()
    resp, err := list.Request()
    if err != nil {
        serv.ResponseBuilder().SetResponseCode(500)
        log.Printf("[Err] ListMetrics - %v\n", err)
        return
    }
    myResp.Metrics = resp.ListMetricsResult.Metrics
    myResp.NextToken = resp.ListMetricsResult.NextToken
    return myResp
}


