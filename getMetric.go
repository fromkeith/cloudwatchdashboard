package main

import (
    "github.com/fromkeith/awsgo"
    "github.com/fromkeith/awsgo/cloudwatch"
    "log"
    "time"
    "strings"
)




type resultDataPoints struct {
    Time                int64
    Value               float64
}
type MetricResults struct {
    Datapoints          []resultDataPoints
    Unit                string
    Statistic           string
}

type Dimensions struct {
    Name                string
    Value               string
}
type MetricSearch struct {
    Namespace           string
    MetricName          string
    StartTime           int64
    EndTime             int64
    Period              int
    Statistic           string
    Unit                string
    Dimensions          []Dimensions
}



func (serv DashboardService) GetMetrics(m MetricSearch) MetricResults {
    cw := cloudwatch.NewGetMetricStatisticsRequest()
    cw.Host.Region = serv.Context.Request().Header.Get("Region")
    cw.StartTime = time.Unix(m.StartTime / 1000, 0)
    cw.EndTime = time.Unix(m.EndTime / 1000, 0)
    cw.Dimensions = make([]cloudwatch.MetricDimensions, len(m.Dimensions))
    for i := range m.Dimensions {
        cw.Dimensions[i].Name = m.Dimensions[i].Name
        cw.Dimensions[i].Value = m.Dimensions[i].Value
    }
    cw.MetricName = m.MetricName
    cw.Namespace = m.Namespace
    if m.Period < 60 {
        m.Period = 60
    }
    cw.Period = m.Period
    cw.Statistics = []string{m.Statistic}
    cw.Key, _ = awsgo.GetSecurityKeys()
    resp, err := cw.Request()
    if err != nil {
        serv.ResponseBuilder().SetResponseCode(500)
        log.Printf("cloudwatch error: %#v\n", err)
        return MetricResults{}
    }
    result := MetricResults{}
    result.Datapoints = make([]resultDataPoints, len(resp.GetMetricStatisticsResult.Datapoints))

    if len(resp.GetMetricStatisticsResult.Datapoints) <= 0 {
        return result
    }
    result.Unit = resp.GetMetricStatisticsResult.Datapoints[0].Unit
    result.Statistic = m.Statistic

    var getter func (dp cloudwatch.GetMetricResultDatapoint) float64
    switch strings.ToLower(m.Statistic) {
        case "average":
            getter = getAverage
        case "sum":
            getter = getSum
        case "samplecount":
            getter = getSampleCount
        case "maximum":
            getter = getMaximum
        case "minimum":
            getter = getMinimum
        default:
            serv.ResponseBuilder().SetResponseCode(500)
            log.Printf("unknown metric: %s\n", m.Statistic)
            return result
    }

    for i := range resp.GetMetricStatisticsResult.Datapoints {
        result.Datapoints[i].Time = resp.GetMetricStatisticsResult.Datapoints[i].Timestamp.Unix()
        result.Datapoints[i].Value = getter(resp.GetMetricStatisticsResult.Datapoints[i])
    }

    return result
}

func getAverage(dp cloudwatch.GetMetricResultDatapoint) float64 {
    return dp.Average
}

func getSum(dp cloudwatch.GetMetricResultDatapoint) float64 {
    return dp.Sum
}
func getSampleCount(dp cloudwatch.GetMetricResultDatapoint) float64 {
    return dp.SampleCount
}
func getMaximum(dp cloudwatch.GetMetricResultDatapoint) float64 {
    return dp.Maximum
}
func getMinimum(dp cloudwatch.GetMetricResultDatapoint) float64 {
    return dp.Minimum
}
