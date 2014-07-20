package main

import (
	// "encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"github.com/bitly/go-simplejson"
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
				<input type="text" name="str" placeholder="Type Tags..." id="str">
				<input type="submit" value=".. and see the images!">
			</form>
		</body>
	</html>`

var upperTemplate = template.Must(template.New("showimage").Parse(upperTemplateHTML))

func showimage(w http.ResponseWriter, r *http.Request) {
	tag := r.FormValue("str")
	safeTag := url.QueryEscape(tag)
	apiKey := "e7ef66cea848474a3e1fe3de117f4670"
	fullUrl := fmt.Sprintf("https://api.flickr.com/services/rest/?method=flickr.photos.search&api_key=%s&tags=%s&format=json&nojsoncallback=1&extras=url_z", apiKey, safeTag)

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
	fmt.Printf("resp.Body: %+v\n", resp.Body)
	defer resp.Body.Close()

	body, dataReadErr := ioutil.ReadAll(resp.Body)
	if dataReadErr != nil {
		log.Fatal("ReadAll: ", dataReadErr)
		return
	}
	// fmt.Printf("ReadAll: %+v\n", body)

	// res := make(map[string]map[string][]map[string]interface{}, 0)
	// var res map[string]interface{}
	// err = json.Unmarshal(body, &res)

	js, err := simplejson.NewJson(body)
    if err != nil {
        log.Fatal("Failed to read json: ", err)
    }

 //    m := res.(map[string]interface{})
	// photos := res["photos"].(map[string]interface{})
	// photo := photos["photo"].([]map[string]interface{})


 //    foomap := m["foo"]
 //    v := foomap.(map[string]interface{})

	// fmt.Printf("%+v\n", &res)

    photo := js.Get("photos").Get("photo").GetIndex(0)
    owner, err := photo.Get("owner").String()
    if err != nil {
        log.Fatal("Failed to get owner: ", err)
    }
    id, err := photo.Get("id").String()
    if err != nil {
        log.Fatal("Failed to get id: ", err)
    }
    photoUrl, err := photo.Get("url_z").String()
    if err != nil {
        log.Fatal("Failed to get photoUrl: ", err)
    }
    fmt.Printf("owner, id, photoUrl: %v, %v, %v\n", owner, id, photoUrl)

	tempErr := upperTemplate.Execute(w, photoUrl)
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

	


