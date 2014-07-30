package main


import (
    "code.google.com/p/go-uuid/uuid"
    "github.com/fromkeith/awsgo"
    "github.com/fromkeith/awsgo/dynamo"
    "log"
)

type PutDashboardRequest struct {
    Name                string
}
type PutDashboardResponse struct {
    DashboardId         string
}

func (serv DashboardService) CreateDashboard(req PutDashboardRequest) (resp PutDashboardResponse) {
    if req.Name == "" {
        serv.ResponseBuilder().SetResponseCode(400)
        return
    }
    resp.DashboardId = uuid.New()

    put := dynamo.NewPutItemRequest()
    put.TableName = DYNAMO_DASHBOARDS_TABLE
    put.Item["DashboardId"] = resp.DashboardId
    put.Item["Name"] = req.Name
    put.Host.Region = DYNAMO_REGION
    put.Key, _ = awsgo.GetSecurityKeys()
    _, err := put.Request()
    if err != nil {
        serv.ResponseBuilder().SetResponseCode(500)
        log.Printf("Failed to put dashboard: %v", err)
        return
    }
    return
}
