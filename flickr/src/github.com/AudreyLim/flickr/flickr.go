package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/showimage", showimage)
	fmt.Println("listening...")
	err := http.ListenAndServe(GetPort(), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func GetPort() string {
	var port = os.Getenv("PORT")
	if port == "" {
		port = "4747"
		fmt.Println("INFO: No PORT environment variable detected, defaulting to " + port)
	}
	return ":" + port
}

func handler (w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, rootForm)
}

const rootForm = 
	`<!DOCTYPE html>
	<html>
		<head>
		<meta charset="utf-8">
			<title>Flickr photos</title>
		</head>
		<body>
			<h1>Flickr photos</h1>
			<p>Find photos by tags!</p>
			<form action="/showimage" method="post" accept-charset="utf-8">
				<input type="text" name="str" value="Type Tags..." id="str">
				<input type="submit" value=".. and see the images!">
			</form>
		</body>
	</html>`

var upperTemplate = template.Must(template.New("showimage").Parse(upperTemplateHTML))

func showimage(w http.ResponseWriter, r *http.Request) {
	tag := r.FormValue("str")
	safeTag := url.QueryEscape(tag)
	apiKey := "e7ef66cea848474a3e1fe3de117f4670"
	fullUrl := fmt.Sprintf("https://api.flickr.com/services/rest/?method=flickr.photos.search&api_key=%stags=%s", apiKey, safeTag)

	client := &http.Client{}

	req, err := http.NewRequest("GET", fullUrl, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
		return
	}

	resp, requestErr := client.Do(req)
	if requestErr != nil {
		log.Fatal("Do: ", requestErr)
		return
	}

	defer resp.Body.Close()



	body, dataReadErr := ioutil.ReadAll(resp.Body)
	if dataReadErr != nil {
		log.Fatal("ReadAll: ", dataReadErr)
		return
	}

	res := make(map[string]map[string]map[string]interface{}, 0)

	json.Unmarshal(body, &res)

	owner, _ := res["photos"]["photo"]["owner"]
	id, _ := res["photos"]["photo"]["id"]

	queryUrl := fmt.Sprintf("http://flickr.com/photos/%s/%s", owner, id)

	

	tempErr := upperTemplate.Execute(w, queryUrl)
	if tempErr != nil {
		http.Error(w, tempErr.Error(), http.StatusInternalServerError)
	}
}

const upperTemplateHTML =
`<!DOCTYPE html>
<html>
	<head>
		<meta charset="utf-8">
		<title>Display images</title>
	</head>
	<body>
		<h3>Images</h3>
		<img src="{{html .}}" alt="Image" />
	</body>
</html>`

	


