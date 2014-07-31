package main

import (
    "github.com/fromkeith/gorest"
    "encoding/base64"
    "crypto/rand"
    "github.com/fromkeith/awsgo"
    "github.com/fromkeith/awsgo/dynamo"
    "code.google.com/p/go.crypto/bcrypt"
    "net/http"

    "sync"
    "time"
)

var (
    // this is jsut a single instane, in memory session data.
    // not scalable... but simple enough for now
    authMapLock         sync.Mutex
    authMap             = make(map[string]time.Time)
)

type LoginRequest struct {
    Username        string
    Password        string
}
type LoginResponse struct {
    Token           string
}

type SimpleAuthData struct {
    Token               string
}

func (sess SimpleAuthData) SessionId() string {
    return sess.Token
}



// @return inRealm, inRole, sess
func SimpleAuther(xsrftoken string, role string, request *http.Request) (bool, bool, gorest.SessionData) {
    if role == "open" {
        return true, true, SimpleAuthData{"Open"}
    }
    if role != "authed" {
        return false, false, nil
    }
    authMapLock.Lock()
    defer authMapLock.Unlock()
    if exp, ok := authMap[xsrftoken]; ok {
        if exp.After(time.Now()) {
            return true, true, SimpleAuthData{xsrftoken}
        }
        delete(authMap, xsrftoken)
    }
    return false, false, nil
}

func (serv DashboardService) Login(req LoginRequest) (resp LoginResponse) {
    if req.Username == "" {
        serv.ResponseBuilder().SetResponseCode(400)
        return
    }
    if req.Password == "" {
        serv.ResponseBuilder().SetResponseCode(400)
        return
    }
    get := dynamo.NewGetItemRequest()
    get.Search["Username"] = req.Username
    get.Host.Region = DYNAMO_REGION
    get.TableName = DYNAMO_LOGIN_TABLE
    get.Key, _ = awsgo.GetSecurityKeys()
    getResp, err := get.Request()
    if err != nil {
        serv.ResponseBuilder().SetResponseCode(500)
        return
    }
    if len(getResp.Item) == 0 {
        serv.ResponseBuilder().SetResponseCode(400)
        return
    }

    salt := getResp.Item["Salt"].(string)
    properPassword64 := getResp.Item["Password"].(string)
    properPasswordBytes, _ := base64.StdEncoding.DecodeString(properPassword64)

    if err = bcrypt.CompareHashAndPassword(properPasswordBytes, []byte(req.Password + salt)); err != nil {
        serv.ResponseBuilder().SetResponseCode(400)
        return
    }
    randToken := make([]byte, 64)
    rand.Read(randToken)
    resp.Token = base64.StdEncoding.EncodeToString(randToken)
    authMapLock.Lock()
    defer authMapLock.Unlock()
    authMap[resp.Token] = time.Now().Add(2 * time.Hour)
    return
}