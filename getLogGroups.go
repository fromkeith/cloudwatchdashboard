package main


import (
    "github.com/fromkeith/awsgo"
    "github.com/fromkeith/awsgo/cloudwatch"
    "log"
)

type GetLogGroupsResponse struct {
    GroupNames          []string
    Token               string
}

func (serv DashboardService) GetLogGroups(token string) (resp GetLogGroupsResponse) {
    group := cloudwatch.NewDescribeLogGroupsRequest()
    group.NextToken = token
    group.Key, _ = awsgo.GetSecurityKeys()
    groupResp, err := group.Request()
    if err != nil {
        log.Println("Error getting log groups %#v", err)
        serv.ResponseBuilder().SetResponseCode(500)
        return
    }

    resp.GroupNames = make([]string, len(groupResp.LogGroups))

    for i := range groupResp.LogGroups {
        resp.GroupNames[i] = groupResp.LogGroups[i].LogGroupName
    }
    resp.Token = groupResp.NextToken
    return
}