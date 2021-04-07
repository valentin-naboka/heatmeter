package report

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html"
	"golang.org/x/net/html/charset"
)

const sessionKey string = "CABSESSID"

type Submitter struct {
	client     http.Client
	host       string
	dumpFolder string
}

func NewSubmitter(host string, dumpFolder string) *Submitter {
	return &Submitter{client: http.Client{Timeout: time.Duration(time.Minute)},
		host:       host,
		dumpFolder: dumpFolder}
}

func (s Submitter) postForm(path string, data url.Values, session string) *http.Response {

	request, err := http.NewRequest("POST", s.host+path, strings.NewReader(data.Encode()))
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	request.AddCookie(&http.Cookie{Name: sessionKey, Value: session})

	resp, err := s.client.Do(request)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	return resp
}

func (s Submitter) getSessionValue() string {
	req, err := http.NewRequest("GET", s.host, nil)
	if err != nil {
		log.Println(err)
		return ""
	}

	resp, err := s.client.Do(req)
	if err != nil {
		log.Println(err)
		return ""
	}
	defer resp.Body.Close()

	for _, c := range resp.Cookies() {
		if c == nil {
			panic("nil Cookie found")
		}
		if c.Name == sessionKey {
			return c.Value
		}
	}
	log.Printf("%s cookie not found", sessionKey)
	return ""
}

func (s Submitter) getPhoneAndFormURL(session, userCabinetURL string) (string, string, bool) {
	req, err := http.NewRequest("GET", s.host+userCabinetURL, nil)
	req.AddCookie(&http.Cookie{Name: sessionKey, Value: session})
	if err != nil {
		log.Fatal(err)
		return "", "", false
	}

	resp, err := s.client.Do(req)
	if err != nil {
		log.Fatal(err)
		return "", "", false
	}
	defer resp.Body.Close()
	utf8Body := s.convertToUTF8(resp)
	if utf8Body == nil {
		return "", "", false
	}

	if resp.StatusCode != 200 {
		log.Printf("status code %d", resp.StatusCode)
		s.dumpHTML(bytes.NewBuffer(utf8Body))
		return "", "", false
	}
	doc := s.getHTMLRoot(bytes.NewBuffer(utf8Body))

	if doc == nil {
		return "", "", false
	}

	submitURL, ok := findAttribute("action", userCabinetURL, doc)
	if !ok {
		s.dumpHTML(bytes.NewBuffer(utf8Body))
		return "", "", false
	}

	phone, ok := findPhone(doc)
	if !ok {
		s.dumpHTML(bytes.NewBuffer(utf8Body))
		return "", "", false
	}

	return submitURL, phone, true
}

func (s Submitter) login(session, email, pass string) (string, string, bool) {

	data := url.Values{}
	data.Set("email", email)
	data.Set("passw", pass)
	data.Set("email-signin", "")

	redirect := false
	const cabinetURL string = "/cabinet"
	s.client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		redirect = req.URL.Path == cabinetURL
		return nil
	}
	resp := s.postForm("/cabinet/login", data, session)
	defer resp.Body.Close()
	utf8Body := s.convertToUTF8(resp)
	if utf8Body == nil {
		return "", "", false
	}

	s.client.CheckRedirect = nil

	if resp.StatusCode != 200 {
		log.Printf("status code %d", resp.StatusCode)
		s.dumpHTML(bytes.NewBuffer(utf8Body))
		return "", "", false
	}

	if !redirect {
		log.Printf("redirect to %s expected\n", cabinetURL)
		s.dumpHTML(bytes.NewBuffer(utf8Body))
		return "", "", false
	}

	doc := s.getHTMLRoot(bytes.NewBuffer(utf8Body))
	if doc == nil {
		return "", "", false
	}

	userCabinetURL, ok := findAttribute("href", "/cabinet/heat/", doc)
	if !ok {
		return "", "", false
	}

	return s.getPhoneAndFormURL(session, userCabinetURL)
}

func findAttributeInternal(key, valuePattern string, n *html.Node) (string, bool) {
	for _, a := range n.Attr {
		if a.Key == key && strings.Contains(a.Val, valuePattern) {
			return a.Val, true
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if value, ok := findAttributeInternal(key, valuePattern, c); ok {
			return value, ok
		}
	}
	return "", false
}

func findAttribute(key, valuePattern string, n *html.Node) (string, bool) {
	value, ok := findAttributeInternal(key, valuePattern, n)
	if !ok {
		log.Printf("unable to find atrribute with key: %s, value pattern: %s", key, valuePattern)
	}
	return value, ok

}

func findPhoneInternal(n *html.Node) (string, bool) {
	phoneIDFound := false
	for _, a := range n.Attr {
		if a.Key == "id" && a.Val == "phone" {
			phoneIDFound = true
			continue
		}

		if phoneIDFound && a.Key == "value" {
			return a.Val, true
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if value, ok := findPhoneInternal(c); ok {
			return value, ok
		}
	}
	return "", false
}

func findPhone(n *html.Node) (string, bool) {
	phone, ok := findPhoneInternal(n)
	if !ok {
		log.Println("unable to find atrribute with id phone")
	}
	return phone, ok

}

func (s Submitter) submitData(formURL,
	session,
	energy,
	volume,
	power,
	volumeFlow,
	flowTemp,
	ReturnTemp,
	operatingTime,
	errorTime,
	phone,
	email string) bool {

	data := url.Values{}

	data.Set("LastDtPok", time.Now().Format("02.01.2006 15:04"))
	data.Set("LastValuePok", energy)
	data.Set("LastRashPok", volume)
	data.Set("HeatRash", volumeFlow)
	data.Set("HeatPower", power)
	data.Set("TempPod", flowTemp)
	data.Set("TempO", ReturnTemp)
	data.Set("LastTimePok", operatingTime)
	data.Set("ErrorTime", errorTime)
	data.Set("ErrorKod", "")
	data.Set("Comment", "")
	data.Set("phone", phone)
	data.Set("email", email)
	data.Set("send", "")

	resp := s.postForm(formURL, data, session)
	defer resp.Body.Close()
	utf8Body := s.convertToUTF8(resp)
	if utf8Body == nil {
		return false
	}

	if resp.StatusCode != 200 {
		log.Printf("status code %d", resp.StatusCode)
		s.dumpHTML(bytes.NewBuffer(utf8Body))
		return false
	}

	doc := s.getHTMLRoot(bytes.NewBuffer(utf8Body))
	if doc == nil {
		return false
	}

	if _, ok := findAttribute("class", "alert alert-warning", doc); ok {
		log.Println("errors ocurred due to measurement submition")
		s.dumpHTML(bytes.NewBuffer(utf8Body))
		return false
	}

	return true
}

func (s Submitter) Execute(email, pass string,
	energy,
	volume,
	power,
	volumeFlow,
	flowTemp,
	ReturnTemp,
	operatingTime,
	errorTime string) bool {

	session := s.getSessionValue()
	if session == "" {
		return false
	}

	submitURL, phone, ok := s.login(session, email, pass)
	if !ok {
		return false
	}

	return s.submitData(submitURL,
		session,
		energy,
		volume,
		power,
		volumeFlow,
		flowTemp,
		ReturnTemp,
		operatingTime,
		errorTime,
		phone,
		email)
}

func (s *Submitter) convertToUTF8(r *http.Response) []byte {
	utf8, err := charset.NewReader(r.Body, r.Header.Get("Content-Type"))
	if err != nil {
		log.Println("encoding error:", err)
		return nil
	}

	result, err := ioutil.ReadAll(utf8)
	if err != nil {
		log.Println("unable to convert to utf8:", err)
		return nil
	}

	return result
}

func (s *Submitter) getHTMLRoot(r io.Reader) *html.Node {

	root, err := html.Parse(r)
	if err != nil {
		s.dumpHTML(r)
		log.Println("unable to parse html:", err)
	}
	return root
}

func (s *Submitter) dumpHTML(r io.Reader) {
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		log.Println("io error:", err)
		return
	}

	err = ioutil.WriteFile(s.dumpFolder+time.Now().Format("02_01_2006_15-04_05.000")+".html", buf, 0600)
	if err != nil {
		log.Println("io error:", err)
	}
}
