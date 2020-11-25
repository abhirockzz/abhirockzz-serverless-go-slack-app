package function

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

const slackResponseStaticText string = "*With* :heart: *from funcy*"

//Funcy implements Slash command backend
func Funcy(w http.ResponseWriter, r *http.Request) {
	fmt.Println("funcy function triggered")

	signingSecret := os.Getenv("SLACK_SIGNING_SECRET")
	apiKey := os.Getenv("GIPHY_API_KEY")

	if signingSecret == "" || apiKey == "" {
		fmt.Println("Required environment variable(s) missing")
		http.Error(w, "Failed to process request. Please contact the admin", http.StatusUnauthorized)
		return
	}

	slackTimestamp := r.Header.Get("X-Slack-Request-Timestamp")

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("unable to read slack data in message body", err)
		http.Error(w, "Failed to process request", http.StatusBadRequest)
		return
	}
	slackSigningBaseString := "v0:" + slackTimestamp + ":" + string(b)
	slackSignature := r.Header.Get("X-Slack-Signature")

	if !matchSignature(slackSignature, signingSecret, slackSigningBaseString) {
		fmt.Println("Signature did not match!")
		http.Error(w, "Function was not invoked by Slack", http.StatusForbidden)
		return
	}

	fmt.Println("Slack request verified successfully")

	//parse the application/x-www-form-urlencoded data sent by Slack
	vals, err := parse(b)
	if err != nil {
		fmt.Println("unable to parse data sent by slack", err)
		http.Error(w, "Failed to process request", http.StatusBadRequest)
		return
	}
	giphyTag := vals.Get("text")
	fmt.Println("Invoking GIPHY API for keyword", giphyTag)

	giphyResp, err := http.Get("http://api.giphy.com/v1/gifs/random?tag=" + giphyTag + "&api_key=" + apiKey)
	if err != nil {
		fmt.Println("giphy did not respond", err)
		http.Error(w, "Failed to process request", http.StatusFailedDependency)
		return
	}

	resp, err := ioutil.ReadAll(giphyResp.Body)
	if err != nil {
		fmt.Println("could not read giphy response", err)
		http.Error(w, "Failed to process request", http.StatusInternalServerError)
		return
	}

	var gr GiphyResponse
	json.Unmarshal(resp, &gr)
	title := gr.Data.Title
	url := gr.Data.Images.Downsized.URL

	fmt.Println("Got response from GIPHY. Image URL", url)

	slackResponse := SlackResponse{Text: slackResponseStaticText, Attachments: []Attachment{{Text: title, ImageURL: url}}}

	//slack needs the content-type to be set explicitly - https://api.slack.com/slash-commands#responding_immediate_response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(slackResponse)
	fmt.Println("Sent response to Slack")
}

func matchSignature(slackSignature, signingSecret, slackSigningBaseString string) bool {

	//calculate SHA256 of the slackSigningBaseString using signingSecret
	mac := hmac.New(sha256.New, []byte(signingSecret))
	mac.Write([]byte(slackSigningBaseString))

	//hex encode the SHA256
	calculatedSignature := "v0=" + hex.EncodeToString(mac.Sum(nil))

	match := hmac.Equal([]byte(slackSignature), []byte(calculatedSignature))
	return match
}

//adapted from from net/http/request.go --> func parsePostForm(r *Request) (vs url.Values, err error)
func parse(b []byte) (url.Values, error) {
	vals, e := url.ParseQuery(string(b))
	if e != nil {
		fmt.Println("unable to parse", e)
		return nil, e
	}
	return vals, nil
}
