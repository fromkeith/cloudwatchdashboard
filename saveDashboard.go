package main

import (
    "github.com/fromkeith/awsgo"
    "github.com/fromkeith/awsgo/dynamo"
    "log"
)

type SaveDashboardRequest struct {
    Name                    string
    AddedGraphs             []string
    RemovedGraphs           []string
}

type SaveDashboardResponse struct {

}

func (serv DashboardService) SaveDashboard(req SaveDashboardRequest, id string) (resp SaveDashboardResponse) {
    if id == "" {
        serv.ResponseBuilder().SetResponseCode(400)
        return
    }
    if len(req.Name) > 0 {
        update := dynamo.NewUpdateItemRequest()
        update.TableName = DYNAMO_DASHBOARDS_TABLE
        update.UpdateKey["DashboardId"] = id
        update.Update["Name"] = dynamo.AttributeUpdates{dynamo.AttributeUpdate_Action_Put, req.Name}
        update.Host.Region = DYNAMO_REGION
        update.Key, _ = awsgo.GetSecurityKeys()
        _, err := update.Request()
        if err != nil {
            serv.ResponseBuilder().SetResponseCode(500)
            log.Printf("Failed to update dashboard: %v", err)
            return
        }
    }
    if len(req.AddedGraphs) == 0 && len(req.RemovedGraphs) == 0 {
        return
    }
    editGraphs := dynamo.NewBatchWriteItemRequest()
    for i := range req.RemovedGraphs {
        editGraphs.AddDeleteRequest(DYNAMO_DASHBOARD_GRAPHS_TABLE,
            map[string]interface{}{
                "DashboardId": id,
                "GraphId": req.RemovedGraphs[i],
            },
        )
    }
    for i := range req.AddedGraphs {
        editGraphs.AddPutRequest(DYNAMO_DASHBOARD_GRAPHS_TABLE,
            map[string]interface{}{
                "DashboardId": id,
                "GraphId": req.AddedGraphs[i],
            },
        )
    }
    editGraphs.Host.Region = DYNAMO_REGION
    editGraphs.Key, _ = awsgo.GetSecurityKeys()
    _, err := editGraphs.RequestSplit()
    if err != nil {
        serv.ResponseBuilder().SetResponseCode(500)
        log.Printf("Failed to update dashboard's graphs: %v", err)
        return
    }
    return
}

