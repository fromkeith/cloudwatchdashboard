package main

import (
    "github.com/fromkeith/awsgo"
    "github.com/fromkeith/awsgo/dynamo"
    "log"
    "time"
)

type dashboardInfo struct {
    Name                string
    DashboardId         string
}

type GetDashboardsResponse struct {
    Dashboards          []dashboardInfo
}

func (serv DashboardService) GetDashboards() (resp GetDashboardsResponse) {
    scan := dynamo.NewScanRequest()
    scan.TableName = DYNAMO_DASHBOARDS_TABLE
    scan.Host.Region = DYNAMO_REGION
    scan.Key, _ = awsgo.GetSecurityKeys()
    scanResp, err := scan.Request()
    if err != nil {
        log.Println("Error scanning table! %#v", err)
        serv.ResponseBuilder().SetResponseCode(500)
        return
    }
    resp.Dashboards = make([]dashboardInfo, 0, 100)

    copyDashboardsResponse(scanResp, &resp)

    if len(scanResp.LastEvaluatedKey) > 0 {
        lastKey := scanResp.LastEvaluatedKey["DashboardId"]
        for backOff := 0; len(scanResp.LastEvaluatedKey) != 0; {
            if scanResp.LastEvaluatedKey["DashboardId"] == lastKey {
                backOff ++
            }
            lastKey = scanResp.LastEvaluatedKey["DashboardId"]
            if backOff > 5 {
                log.Println("Exceed backoff!")
                return
            }
            time.Sleep(time.Duration(backOff * 100))
            retryRequest := dynamo.NewScanRequest()
            retryRequest.ExclusiveStartKey = map[string]interface{}{
                "DashboardId" : lastKey,
            }
            retryRequest.TableName = DYNAMO_DASHBOARDS_TABLE
            retryRequest.Host.Region = DYNAMO_REGION
            retryRequest.Key, _ = awsgo.GetSecurityKeys()
            scanResp, err = retryRequest.Request()
            if err != nil {
                continue
            }
            if scanResp.LastEvaluatedKey["DashboardId"] != lastKey {
                copyDashboardsResponse(scanResp, &resp)
            }
        }
    }
    return
}


func copyDashboardsResponse(scanResp * dynamo.ScanResponse, resp * GetDashboardsResponse) {
    for i := range scanResp.Items {
        var dash dashboardInfo
        dash.Name = scanResp.Items[i]["Name"].(string)
        dash.DashboardId = scanResp.Items[i]["DashboardId"].(string)
        resp.Dashboards = append(resp.Dashboards, dash)
    }
}