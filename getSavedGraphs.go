package main

import (
    "encoding/json"
    "github.com/fromkeith/awsgo"
    "github.com/fromkeith/awsgo/dynamo"
    "log"
    "time"
)


type GetSavedGraphResponse struct {
    Graphs                  []SaveGraphRequest
}

func (serv DashboardService) GetSavedGraphs() (resp GetSavedGraphResponse) {
    scan := dynamo.NewScanRequest()
    scan.TableName = DYNAMO_GRAPHS_TABLE
    scan.Host.Region = DYNAMO_REGION
    scan.Key, _ = awsgo.GetSecurityKeys()
    scanResp, err := scan.Request()
    if err != nil {
        log.Println("Error scanning table! %#v", err)
        serv.ResponseBuilder().SetResponseCode(500)
        return
    }

    resp.Graphs = make([]SaveGraphRequest, 0, 100)

    copySavedGraphResponse(scanResp, &resp)

    if len(scanResp.LastEvaluatedKey) > 0 {
        lastKey := scanResp.LastEvaluatedKey["GraphId"]
        for backOff := 0; len(scanResp.LastEvaluatedKey) != 0; {
            if scanResp.LastEvaluatedKey["GraphId"] == lastKey {
                backOff ++
            }
            lastKey = scanResp.LastEvaluatedKey["GraphId"]
            if backOff > 5 {
                log.Println("Exceed backoff!")
                return
            }
            time.Sleep(time.Duration(backOff * 100))
            retryRequest := dynamo.NewScanRequest()
            retryRequest.ExclusiveStartKey = map[string]interface{}{
                "GraphId" : lastKey,
            }
            retryRequest.TableName = DYNAMO_GRAPHS_TABLE
            retryRequest.Host.Region = DYNAMO_REGION
            retryRequest.Key, _ = awsgo.GetSecurityKeys()
            scanResp, err = retryRequest.Request()
            if err != nil {
                continue
            }
            if scanResp.LastEvaluatedKey["GraphId"] != lastKey {
                copySavedGraphResponse(scanResp, &resp)
            }
        }
    }
    return
}


func copySavedGraphResponse(scanResp * dynamo.ScanResponse, resp * GetSavedGraphResponse) {
    for i := range scanResp.Items {
        var g SaveGraphRequest
        g.Id = scanResp.Items[i]["GraphId"].(string)

        if err := json.Unmarshal([]byte(scanResp.Items[i]["Metrics"].(string)), &g.Metrics); err != nil {
            log.Println("Errror unmarhsalling metrics: %v %#v", err, err)
            continue
        }
        g.Time.Start = int64(scanResp.Items[i]["Time.Start"].(float64))
        g.Time.End = int64(scanResp.Items[i]["Time.End"].(float64))
        g.Time.Period = int(scanResp.Items[i]["Time.Period"].(float64))
        g.Name = scanResp.Items[i]["Name"].(string)
        resp.Graphs = append(resp.Graphs, g)
    }
}