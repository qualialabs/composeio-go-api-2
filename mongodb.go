package composeio

import (
  "net/http"
  "log"
  // "encoding/json"
  // "fmt"
  "bytes"
  "time"
  "io/ioutil"
  "strings"
  "errors"

)

// ##############################
// Client is the object that handles talking to the Datadog API. This maintains
// state information for a particular application connection.
type Client struct {
  ComposeioToken string

  // URL to the Compose API to use
  URL string

  //The Http Client that is used to make requests
  HttpClient *http.Client
}

// NewClient returns a new composeio.Client which can be used to access the API
// methods. The expected argument is the composeio token.
func NewClient(composeio_token string) *Client {
  return &Client{
    ComposeioToken:     composeio_token,
    URL: "https://api.compose.io",
    HttpClient: http.DefaultClient,
  }
}

// ################################

type Mongodb struct {
  Account string      `json:"account,omitempty"`
  Deployment  string      `json:"deployment,omitempty"`
  Name    string   `json:"name,omitempty"`
}

type Collection struct{
  Name string `json:"name,omitempty"`
}

type User struct{
  Username string `json:"username,omitempty"`
  Password string `json:"password,omitempty"`
  ReadOnly bool `json:"readOnly,omitempty"`
}


// #############################


func (client *Client) CreateMongodb(mongodb *Mongodb) error {
    // could not create mongodb directly, must create a collection

    collection := &Collection{
      Name: "dummy_collection",
    }
    log.Println("[DEBUG] Creating dummpy colllection to create mongodb")
    err := client.CreateMongodbCollection(mongodb, collection)
    if err != nil {
      return err
    }

    time.Sleep(time.Millisecond * 5000)

    log.Println("[DEBUG] Deleting dummpy colllection from mongodb")
    err = client.DeleteMongodbCollection(mongodb, collection)
    if err != nil {
      return err
    } else {
      return nil      
    }

}

func (client *Client) ReadMongodb(mongodb *Mongodb) error {

  url := client.URL + "/deployments/" + mongodb.Account + "/" + mongodb.Deployment + "/mongodb/" + mongodb.Name + "/stats"
  log.Println("[DEBUG] URL:>", url)
  req, err := http.NewRequest("GET", url, nil)
  if err != nil {
    return err
  }
  req.Header.Set("Content-Type", "application/json")
  req.Header.Set("Accept-Version", "2014-06")
  req.Header.Set("Authorization", "Bearer " + client.ComposeioToken)

  resp, err := client.HttpClient.Do(req)

  if err != nil {
    return err
  } else {
    body, _ := ioutil.ReadAll(resp.Body)
    log.Println("[DEBUG] Get response from composeio response: ", string(body))
    if strings.Contains(string(body), `"dataSize":0`) {
      // err := "error"
      log.Println("[DEBUG] not found ")
      return errors.New("not found")
    } else {
      return nil       
    }
  }


}


func (client *Client) CreateMongodbCollection(mongodb *Mongodb, collection *Collection) error {

  url := client.URL + "/deployments/" + mongodb.Account + "/" + mongodb.Deployment + "/mongodb/" + mongodb.Name + "/collections"
  log.Println("[DEBUG] URL:>", url)
  collection_name := collection.Name
  jsonStr := []byte(`{"name":"`+ collection_name+`"}`)
  req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
  if err != nil {
    return err
  }
  req.Header.Set("Content-Type", "application/json")
  req.Header.Set("Accept-Version", "2014-06")
  req.Header.Set("Authorization", "Bearer " + client.ComposeioToken)

  resp, err := client.HttpClient.Do(req)

  if err != nil {
    return err
  } else {
    body, _ := ioutil.ReadAll(resp.Body)
    log.Println("[DEBUG] Get response from composeio response: ", string(body))
    return nil 
  }


}

func (client *Client) DeleteMongodbCollection(mongodb *Mongodb, collection *Collection) error {

  url := client.URL + "/deployments/" + mongodb.Account + "/" + mongodb.Deployment + "/mongodb/" + mongodb.Name + "/collections/" + collection.Name
  log.Println("[DEBUG] URL:>" + url)
  req, err := http.NewRequest("DELETE", url, nil)
  if err != nil {
    return err
  }
  req.Header.Set("Content-Type", "application/json")
  req.Header.Set("Accept-Version", "2014-06")
  req.Header.Set("Authorization", "Bearer " + client.ComposeioToken)

  
  resp, err := client.HttpClient.Do(req)

  if err != nil {
    return err
  } else {
    body, _ := ioutil.ReadAll(resp.Body)
    log.Println("[DEBUG] Get response from composeio response: ", string(body))
    return nil 
  }

}

func (client *Client) CreateMongodbUser(mongodb *Mongodb, user *User) error {

  url := client.URL + "/deployments/" + mongodb.Account + "/" + mongodb.Deployment + "/mongodb/" + mongodb.Name + "/users"
  log.Println("[DEBUG] URL:>", url)
  username := user.Username
  password := user.Password
  var readOnly = "false"
  if user.ReadOnly {
    readOnly = "true"
  } 
  jsonStr := []byte(`{"username":"`+ username +`", "password":"`+ password + `", "readOnly":`+ readOnly + `}`)
  req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
  if err != nil {
    return err
  }
  req.Header.Set("Content-Type", "application/json")
  req.Header.Set("Accept-Version", "2014-06")
  req.Header.Set("Authorization", "Bearer " + client.ComposeioToken)

  
  resp, err := client.HttpClient.Do(req)

  if err != nil {
    return err
  } else {
    body, _ := ioutil.ReadAll(resp.Body)
    log.Println("[DEBUG] Get response from composeio response: ", string(body))
    return nil 
  }
}

// ############ the api not work #############
// func (client *Client) ReadMongodbUser(mongodb *Mongodb, user *User) error {

//   url := client.URL + "/deployments/" + mongodb.Account + "/" + mongodb.Deployment + "/mongodb/" + mongodb.Name + "/users"
//   log.Println("[DEBUG] URL:>", url)
//   req, err := http.NewRequest("GET", url, nil)
//   if err != nil {
//     return err
//   }
//   req.Header.Set("Content-Type", "application/json")
//   req.Header.Set("Accept-Version", "2014-06")
//   req.Header.Set("Authorization", "Bearer " + client.ComposeioToken)

//   resp, err := client.HttpClient.Do(req)

//   if err != nil {
//     return err
//   } else {
//     body, _ := ioutil.ReadAll(resp.Body)
//     log.Println("[DEBUG] Get response from composeio response: ", string(body))
//     if !strings.Contains(string(body), user.Username) {
//       // err := "error"
//       log.Println("[DEBUG] "+  user.Username + " not found ")
//       return errors.New(user.Username + " not found")
//     } else {
//       return nil       
//     }
//   }

// }

func (client *Client) UpdateMongodbUser(mongodb *Mongodb, user *User) error {

  url := client.URL + "/deployments/" + mongodb.Account + "/" + mongodb.Deployment + "/mongodb/" + mongodb.Name + "/users/" + user.Username
  log.Println("[DEBUG] URL:>", url)

  password := user.Password
  var readOnly = "false"
  if user.ReadOnly {
    readOnly = "true"
  } 
  jsonStr := []byte(`{"password":"`+ password + `", "readOnly":`+ readOnly + `}`)
  req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonStr))
  if err != nil {
    return err
  }
  req.Header.Set("Content-Type", "application/json")
  req.Header.Set("Accept-Version", "2014-06")
  req.Header.Set("Authorization", "Bearer " + client.ComposeioToken)

  
  resp, err := client.HttpClient.Do(req)

  if err != nil {
    return err
  } else {
    body, _ := ioutil.ReadAll(resp.Body)
    log.Println("[DEBUG] Get response from composeio response: ", string(body))
    return nil 
  }
}

func (client *Client) DeleteMongodbUser(mongodb *Mongodb, user *User) error {

  url := client.URL + "/deployments/" + mongodb.Account + "/" + mongodb.Deployment + "/mongodb/" + mongodb.Name + "/users/" + user.Username
  log.Println("[DEBUG] URL:>", url)
  req, err := http.NewRequest("DELETE", url, nil )
  if err != nil {
    return err
  }
  req.Header.Set("Content-Type", "application/json")
  req.Header.Set("Accept-Version", "2014-06")
  req.Header.Set("Authorization", "Bearer " + client.ComposeioToken)

  
  resp, err := client.HttpClient.Do(req)

  if err != nil {
    return err
  } else {
    body, _ := ioutil.ReadAll(resp.Body)
    log.Println("[DEBUG] Get response from composeio response: ", string(body))
    return nil 
  }
}



