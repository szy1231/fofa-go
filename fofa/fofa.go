// Package fofa implements some fofa-api utility functions.
package fofa

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/buger/jsonparser"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// Fofa a fofa client can be used to make queries
type Fofa struct {
	Email []byte
	Key   []byte
	*http.Client
}

// User struct for fofa user
type User struct {
	Email      string `json:"email"`
	UserName   string `json:"username"`
	Fcoin      int    `json:"fcoin"`
	Vip        bool   `json:"isvip"`
	VipLevel   int    `json:"vip_level"`
	IsVerified bool   `json:"is_verified"`
	Avatar     string `json:"avatar"`
	Err        string `json:"errmsg,omitempty"`
}

// Data fofa result
type Data struct {
	Mode    string     `json:"mode"`
	Error   bool       `json:"error"`
	Query   string     `json:"query"`
	Page    int        `json:"page"`
	Size    int        `json:"size"`
	Results [][]string `json:"results"`
}

// SingleData single results data
type SingleData struct {
	Mode    string   `json:"mode"`
	Error   bool     `json:"error"`
	Query   string   `json:"query"`
	Page    int      `json:"page"`
	Size    int      `json:"size"`
	Results []string `json:"results"`
}

const (
	defaultAPIUrl = "https://fofa.so/api/v1/search/all?"
)

// NewFofaClient create a fofa client
func NewFofaClient(email, key []byte) *Fofa {

	transCfg := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	return &Fofa{
		Email: email,
		Key:   key,
		Client: &http.Client{
			Transport: transCfg, // disable tls verify
		},
	}
}

// Get overwrite http.Get
func (ff *Fofa) Get(u string) ([]byte, error) {

	body, err := ff.Client.Get(u)
	if err != nil {
		return nil, err
	}
	defer body.Body.Close()
	content, err := ioutil.ReadAll(body.Body)
	if err != nil {
		return nil, err
	}
	return content, nil
}

// QueryAsJSON make a fofa query and return json data as result
// echo 'domain="nosec.org"' | base64 - | xargs -I{}
// curl "https://fofa.so/api/v1/search/all?email=${FOFA_EMAIL}&key=${FOFA_KEY}&qbase64={}"
// Default: host title ip domain port country city
func (ff *Fofa) QueryAsJSON(page uint, size uint, full string, args ...[]byte) ([]byte, error) {
	var (
		query  = []byte(nil)
		fields = []byte("domain,host,ip,port,title,country,city")
		q      = []byte(nil)
	)
	switch {
	case len(args) == 1 || (len(args) == 2 && args[1] == nil):
		query = args[0]
	case len(args) == 2:
		query = args[0]
		fields = args[1]
	}

	q = []byte(base64.StdEncoding.EncodeToString(query))
	q = bytes.Join([][]byte{[]byte(defaultAPIUrl),
		[]byte("email="), ff.Email,
		[]byte("&key="), ff.Key,
		[]byte("&qbase64="), q,
		[]byte("&fields="), fields,
		[]byte("&page="), []byte(strconv.Itoa(int(page))),
		[]byte("&size="), []byte(strconv.Itoa(int(size))),
		[]byte("&full="), []byte(full),
	}, []byte(""))
	content, err := ff.Get(string(q))
	if err != nil {
		return nil, err
	}
	errmsg, err := jsonparser.GetString(content, "errmsg")
	if err == nil {
		err = errors.New(errmsg)
	} else {
		err = nil
	}
	return content, err
}

// QueryAsObject make a fofa query and
// return object data as result
// echo 'domain="nosec.org"' | base64 - | xargs -I{}
// curl "https://fofa.so/api/v1/search/all?email=${FOFA_EMAIL}&key=${FOFA_KEY}&qbase64={}"
func (ff *Fofa) QueryAsObject(page uint, size uint, full string, args ...[]byte) (data Data, err error) {

	var content []byte

	content, err = ff.QueryAsJSON(page, size, full, args...)
	if err != nil {
		return Data{}, err
	}

	errmsg, err := jsonparser.GetString(content, "errmsg")
	// err equals to nil on error
	if err == nil {
		return Data{}, errors.New(errmsg)
	}

	if len(args) == 2 && len(strings.Split(string(args[1]), ",")) == 1 && args[1] != nil {
		return ff.toData(content)
	}
	err = json.Unmarshal(content, &data)
	return
}

//toData tmpData struct to data struct
func (ff *Fofa) toData(content []byte) (data Data, err error) {
	var tmpData SingleData
	err = json.Unmarshal(content, &tmpData)
	if err != nil {
		return Data{}, err
	}
	arr := make([][]string, len(tmpData.Results))
	for k, v := range tmpData.Results {
		arr[k] = []string{v}
	}

	data = Data{
		Mode:    tmpData.Mode,
		Error:   tmpData.Error,
		Query:   tmpData.Query,
		Page:    tmpData.Page,
		Size:    tmpData.Size,
		Results: arr,
	}

	return
}

// UserInfo get user information
func (ff *Fofa) UserInfo() (user *User, err error) {
	user = new(User)
	queryStr := strings.Join([]string{"https://fofa.so/api/v1/info/my?email=", string(ff.Email), "&key=", string(ff.Key)}, "")
	content, err := ff.Get(queryStr)

	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(content, user); err != nil {
		return nil, err
	}

	if len(user.Err) != 0 {
		return nil, errors.New(user.Err)
	}

	return user, nil
}

//String user to string
func (u *User) String() string {
	data, err := json.Marshal(u)
	if err != nil {
		log.Fatalf("json marshal failed. err: %s\n", err)
		return ""
	}
	return string(data)
}

//String data to string
func (r *Data) String() string {
	data, err := json.Marshal(r)
	if err != nil {
		log.Fatalf("json marshal failed. err: %s\n", err)
		return ""
	}
	return string(data)
}
