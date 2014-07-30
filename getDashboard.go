package main


import (
    "github.com/fromkeith/awsgo"
    "github.com/fromkeith/awsgo/dynamo"
    "log"
)


type GetDashboardResponse struct {
    dashboardInfo
    Graphs                  []string
}

func (serv DashboardService) GetDashboard(id string) (resp GetDashboardResponse) {
    if id == "" {
        serv.ResponseBuilder().SetResponseCode(400)
        return
    }

    getDashboard := dynamo.NewGetItemRequest()
    getDashboard.Search["DashboardId"] = id
    getDashboard.TableName = DYNAMO_DASHBOARDS_TABLE
    getDashboard.Host.Region = DYNAMO_REGION
    getDashboard.Key, _ = awsgo.GetSecurityKeys()
    getDashboardResp, err := getDashboard.Request()
    if err != nil {
        serv.ResponseBuilder().SetResponseCode(500)
        log.Printf("Error getting dashboard: %v", err)
        return
    }
    if len(getDashboardResp.Item) == 0 {
        serv.ResponseBuilder().SetResponseCode(404)
        return
    }

    resp.Name = getDashboardResp.Item["Name"].(string)
    resp.DashboardId = getDashboardResp.Item["DashboardId"].(string)
    resp.Graphs = make([]string, 0, 20)


    getGraphs := dynamo.NewQueryRequest()
    getGraphs.AddKeyCondition("DashboardId", []interface{}{id}, dynamo.ComparisonOperator_EQ)
    getGraphs.Select = dynamo.Select_ALL_ATTRIBUTES
    getGraphs.Host.Region = DYNAMO_REGION
    getGraphs.TableName = DYNAMO_DASHBOARD_GRAPHS_TABLE
    getGraphs.Key, _ = awsgo.GetSecurityKeys()
    getGraphsResp, err := getGraphs.Request()
    if err != nil {
        serv.ResponseBuilder().SetResponseCode(500)
        log.Printf("Error getting graphs: %v", err)
        return
    }
    for i := range getGraphsResp.Items {
        resp.Graphs = append(resp.Graphs, getGraphsResp.Items[i]["GraphId"].(string))
    }
    return
}