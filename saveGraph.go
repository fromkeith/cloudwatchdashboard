package main


import (
    "code.google.com/p/go-uuid/uuid"
    "github.com/fromkeith/awsgo"
    "github.com/fromkeith/awsgo/dynamo"
    "encoding/json"
    "log"
)

type metricDef struct {
    Metric              string
    Dimension           string
    Statistic           string
}

type timeDef struct {
    Start               int64
    End                 int64
    Period              int
    Relative            string
}

type SaveGraphRequest struct {
    Metrics             []metricDef
    Time                timeDef
    Id                  string
    Name                string
}

type SaveGraphResponse struct {
    Id                  string
}

func (serv DashboardService) SaveGraph(req SaveGraphRequest, id string) (resp SaveGraphResponse) {

    if id == "" {
        id = uuid.New()
    }
    resp.Id = id

    put := dynamo.NewPutItemRequest()
    put.TableName = DYNAMO_GRAPHS_TABLE
    put.Item["GraphId"] = id
    put.Item["Name"] = req.Name
    // probably want to change this to support times relative to now
    put.Item["Time.Start"] = req.Time.Start
    put.Item["Time.End"] = req.Time.End
    put.Item["Time.Period"] = req.Time.Period
    if req.Time.Relative != "" {
        put.Item["Time.Relative"] = req.Time.Relative
    }
    // max limit is 64KB.. so for now this will be fine
    b, _ := json.Marshal(req.Metrics)
    put.Item["Metrics"] = string(b)
    put.Host.Region = DYNAMO_REGION
    put.Key, _ = awsgo.GetSecurityKeys()
    _, err := put.Request()
    if err != nil {
        log.Printf("Error saving graph: %#v", err)
        serv.ResponseBuilder().SetResponseCode(500)
        return
    }
    return
}