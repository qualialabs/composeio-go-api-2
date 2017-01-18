package composeio

import (
  "crypto/tls"
  "crypto/x509"
  "io/ioutil"

  "fmt"
  "net"
  "os"
  "log"
  "strings"

  "gopkg.in/mgo.v2"

)

// ##############################
// Client is the object that handles talking to compose
type Client struct {
  SslPemPath string
  AdminMongodbURL string
}

// NewClient returns a new composeio.Client which can be used to access the API
// methods. The expected argument is the composeio token.
func NewClient(admin_mongodb_url string, ssl_pem_path string) *Client {
  return &Client{
    AdminMongodbURL:     admin_mongodb_url,
    SslPemPath: ssl_pem_path,
  }
}

// ################################

type Mongodb struct {
  Account string      `json:"account,omitempty"`
  Deployment  string      `json:"deployment,omitempty"`
  Name    string   `json:"name,omitempty"`
}

type User struct{
  Username string `json:"username,omitempty"`
  Password string `json:"password,omitempty"`
  ReadOnly bool `json:"readOnly,omitempty"`
}


// #############################

func (client *Client) CreateMongodbUser(mongodb *Mongodb, user *User) error {

  roots := x509.NewCertPool()
  pem_path := client.SslPemPath
  if ca, err := ioutil.ReadFile(pem_path); err == nil { 
    roots.AppendCertsFromPEM(ca)
  }

  tlsConfig := &tls.Config{}
  tlsConfig.RootCAs = roots

  //connect URL:
  // "mongodb://<username>:<password>@<hostname>:<port>,<hostname>:<port>/<db-name>
  admin_mongodb_url := client.AdminMongodbURL
  admin_mongodb_url = strings.TrimSuffix(admin_mongodb_url, "?ssl=true")


  dialInfo, err := mgo.ParseURL(admin_mongodb_url)
  if err != nil {
    fmt.Println("Failed to parse URI: ", err)
    os.Exit(1)
  }

  dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
    conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
    return conn, err
  }

  session, err := mgo.DialWithInfo(dialInfo)
  if err != nil {
    fmt.Println("Failed to connect: ", err)
    os.Exit(1)
  }


  db := session.DB(mongodb.Name)
  err = db.AddUser(user.Username, user.Password, user.ReadOnly)  //should use UpSert, but this is much easier
  if err != nil {
    log.Fatal(err)
    return nil
  } else {
    log.Println("[DEBUG] created" + user.Username )
    return nil
  }

  // if err != nil {
  //   return err
  // } else {
  //   body, _ := ioutil.ReadAll(resp.Body)
  //   log.Println("[DEBUG] Get response from composeio response: ", string(body))
  //   return nil 
  // }

}


func (client *Client) UpdateMongodbUser(mongodb *Mongodb, user *User) error {

  client.DeleteMongodbUser(mongodb, user)
  client.CreateMongodbUser(mongodb, user)
  return nil
}

func (client *Client) DeleteMongodbUser(mongodb *Mongodb, user *User) error {


  roots := x509.NewCertPool()
  pem_path := client.SslPemPath
  if ca, err := ioutil.ReadFile(pem_path); err == nil { 
    roots.AppendCertsFromPEM(ca)
  }

  tlsConfig := &tls.Config{}
  tlsConfig.RootCAs = roots

  admin_mongodb_url := client.AdminMongodbURL
  admin_mongodb_url = strings.TrimSuffix(admin_mongodb_url, "?ssl=true")

  dialInfo, err := mgo.ParseURL(admin_mongodb_url)
  if err != nil {
    fmt.Println("Failed to parse URI: ", err)
    os.Exit(1)
  }

  dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
    conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
    return conn, err
  }
  session, err := mgo.DialWithInfo(dialInfo)
  if err != nil {
    fmt.Println("Failed to connect: ", err)
    os.Exit(1)
  }


  db := session.DB(mongodb.Name)
  err = db.RemoveUser(user.Username)
  if err != nil {
    log.Fatal(err)
    return nil
  } else {
    log.Println("[DEBUG] created" + user.Username )
    return nil
  }

}



