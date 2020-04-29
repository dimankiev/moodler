package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"
)

type about struct {
	name    string
	version string
	author  string
}

// If you make modification of this program - just add after dimankiev "x your_username"
// For example: author: "dimankiev x modderUsername"
// ONLY FOR EDUCATIONAL PURPOSES | COMMERCIAL USAGE IS FORBIDDEN
var aboutProgram = about{name: "[Moodler Toolkit] Online Bot", version: "1.0.0", author: "dimankiev"}

// Configuration represents the storage for Moodle configuration
type Configuration struct {
	Domain   string `json:"moodleDomain"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func readConfiguration() Configuration {
	data, err := ioutil.ReadFile("../../../../settings.json") // Open json file with credentials
	if err != nil {
		fmt.Print(err)
	}

	var obj Configuration // Declare var for cred. storage

	// Unmarshall json
	err = json.Unmarshal(data, &obj)
	if err != nil {
		fmt.Println("error:", err)
	}
	return obj
}

// Jar represents the storage for Moodle cookies
type Jar struct {
	lk      sync.Mutex
	cookies map[string][]*http.Cookie
}

// NewJar creates a Jar
func NewJar() *Jar {
	jar := new(Jar)
	jar.cookies = make(map[string][]*http.Cookie)
	return jar
}

// SetCookies handles the receipt of the cookies in a reply for the
// given URL.  It may or may not choose to save the cookies, depending
// on the jar's policy and implementation.
func (jar *Jar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	jar.lk.Lock()
	jar.cookies[u.Host] = cookies
	jar.lk.Unlock()
}

// Cookies returns the cookies to send in a request for the given URL.
// It is up to the implementation to honor the standard cookie use
// restrictions such as in RFC 6265.
func (jar *Jar) Cookies(u *url.URL) []*http.Cookie {
	return jar.cookies[u.Host]
}

// ParsedResponse represents clear and easy-to-use object which is accurately parsed request
type ParsedResponse struct {
	URL          string
	StatusCode   int
	ResponseBody string
	Error        error
}

func parseResp(err error, resp *http.Response, response ParsedResponse) ParsedResponse {
	if err != nil {
		response.Error = err
	} else {
		response.StatusCode = resp.StatusCode
		defer resp.Body.Close()
		body, errBodyRead := ioutil.ReadAll(resp.Body)
		if errBodyRead != nil {
			response.Error = errBodyRead
		} else {
			response.ResponseBody = string(body)
		}
	}
	return response
}

func doGetRequest(client http.Client, url string) ParsedResponse {
	resp, err := client.Get(url)
	response := ParsedResponse{URL: url}
	return parseResp(err, resp, response)
}

func doPostFormRequest(client http.Client, url string, formData url.Values) ParsedResponse {
	resp, err := client.PostForm(url, formData)
	response := ParsedResponse{URL: url}
	return parseResp(err, resp, response)
}

func getToken(client http.Client, domain string) string {
	tokenLookup := doGetRequest(client, fmt.Sprintf("https://%s/login/index.php", domain))
	tokenRegexp := regexp.MustCompile(`<input type="hidden" name="logintoken" value="([a-zA-Z0-9]+)">`)
	if tokenLookup.Error == nil && tokenRegexp.MatchString(tokenLookup.ResponseBody) {
		return tokenRegexp.FindStringSubmatch(tokenLookup.ResponseBody)[1]
	}
	return "none"
}

func authorize(client http.Client, token string, config Configuration) bool {
	doPostFormRequest(client, "https://"+config.Domain+"/login/index.php", url.Values{"logintoken": {token}, "username": {config.Username}, "password": {config.Password}})
	request := doGetRequest(client, "https://"+config.Domain+"/")
	loginVerify := regexp.MustCompile(`<div class="logininfo">([A-Za-zА-Яа-яА-ЩЬЮЯҐЄІЇа-щьюяґєії\s]|['\x60’ʼ])+<a href="https://` + strings.Replace(config.Domain, ".", "\\.", -1) + `/user/profile\.php\?id=[0-9]+" title="([A-Za-zА-Яа-яА-ЩЬЮЯҐЄІЇа-щьюяґєії\s]|['\x60’ʼ])+">`)
	if token != "none" {
		if request.Error != nil {
			return false
		}
		if request.StatusCode == 200 && loginVerify.MatchString(request.ResponseBody) {
			return true
		}
	}
	return false
}

func becomeOnline(client http.Client, domain string) {
	fmt.Printf("[%v] Request was sent! Success: %v\n", time.Now(), doGetRequest(client, "https://"+domain+"/").StatusCode == 200)
}

func main() {
	fmt.Printf("%v v.%v (by %v)\n", aboutProgram.name, aboutProgram.version, aboutProgram.author)
	fmt.Println("Reading configuration...")
	config := readConfiguration()
	jar := NewJar()
	client := http.Client{Jar: jar}
	fmt.Print("Authorizing... ")
	token := getToken(client, config.Domain)
	login := authorize(client, token, config)
	if login == true {
		fmt.Print("Success\n")
		ticker := time.NewTicker(5 * time.Second)
		quit := make(chan struct{})
		for {
			select {
			case <-ticker.C:
				becomeOnline(client, config.Domain)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	} else {
		fmt.Print("Failed\n")
		fmt.Println("Please check your settings and Internet connection...")
	}
}
