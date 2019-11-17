package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"sort"

	"github.com/gorilla/mux"
)

//UserInfo is the information received as a body of the post request to /creditcard
type UserInfo struct {
	FirstName   string `json:"firstname" binding:"required"`
	LastName    string `json:"lastname" binding:"required"`
	DOB         string `json:"dob" binding:"required"`
	CreditScore int    `json:"credit-score" binding:"required"`
	EmpStatus   string `json:"employment-status" binding:"required"`
	Salary      int    `json:"salary" binding:"required"`
}

//CreditCard is the response of /creditcard endpoint if successful
type CreditCard struct {
	Provider  string   `json:"provider"`
	Name      string   `json:"name"`
	ApplyURL  string   `json:"apply-url"`
	Apr       float64  `json:"apr"`
	Features  []string `json:"features"`
	CardScore float64  `json:"card-score"`
}

//CSCardResponse is the response of /cards endpoint if successful
type CSCardResponse struct {
	CardName    string   `json:"cardName,omitempty"`
	URL         string   `json:"url,omitempty"`
	Apr         float64  `json:"apr,omitempty"`
	Eligibility float64  `json:"eligibility,omitempty"`
	Features    []string `json:"features,omitempty"`
}

//ScoredCardResponse is the response of /creditcards endpoint if successful
type ScoredCardResponse struct {
	Card           string   `json:"card,omitempty"`
	ApplyURL       string   `json:"apply-url,omitempty"`
	Apr            float64  `json:"annual-percentage-rate,omitempty"`
	ApprovalRating float64  `json:"approval-rating,omitempty"`
	Attributes     []string `json:"attributes,omitempty"`
	IntroOffers    []string `json:"introductory-offers,omitempty"`
}

//CreditCards contains all credit cards from CSCards and ScoredCards
type CreditCards []CreditCard

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	r := mux.NewRouter()
	r.HandleFunc("/v1/creditcard", Handler).Methods(http.MethodPost)
	err := http.ListenAndServe(":"+port, r)

	if err != nil {
		log.Fatal("error occurred")
	}
}

//Handler receives the user info, passes it to CSCard and ScoredCard APIs, format and sort the responses
func Handler(w http.ResponseWriter, r *http.Request) {
	var newUserInfo UserInfo
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, "please enter user info")
		return
	}
	err = json.Unmarshal(reqBody, &newUserInfo)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, "please enter the body in right JSON format")
		return
	}

	//creates an empty result array
	var creditcards []CreditCard

	//gets credit cards information from CSCards
	csCardsResults, err := newUserInfo.GetCSCards()
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, "unable to retrieve CSCards")
		return
	}

	//appends the result array with credit cards received from CSCards
	for _, csCardsResult := range csCardsResults {
		creditcards = append(creditcards, csCardsResult)
	}

	//gets credit cards information from ScoredCards
	scoredCardsResults, err := newUserInfo.GetScoredCards()
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, "unable to retrieve ScoredCards")
		return
	}
	//appends the result array with credit cards received from ScoredCards
	for _, scoredCardResult := range scoredCardsResults {
		creditcards = append(creditcards, scoredCardResult)
	}

	//sorts the result by card score
	sort.SliceStable(creditcards, func(i, j int) bool {
		return creditcards[j].CardScore < creditcards[i].CardScore
	})

	//converts the result to json for the response
	err = json.NewEncoder(w).Encode(creditcards)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "failed to response in JSON")
	}
	//responds with http response status code to 200 if successful
	w.WriteHeader(http.StatusOK)
}

//GetCSCards sends a post request to CSCard API endpoint and formats the response
func (userInfo *UserInfo) GetCSCards() ([]CreditCard, error) {
	//csCardEndPoint := os.Getenv("CSCARDS_ENDPOINT")
	//above was not working with test environment, put the link here as well as in .env
	csCardEndPoint := "https://y4xvbk1ki5.execute-api.us-west-2.amazonaws.com/CS/v1/cards"

	//makes a body for the POST request with user information received
	var jsonStr = []byte(fmt.Sprintf(`{
		"fullName": "%s %s",
		"dateOfBirth": "%s",
		"creditScore": %d
	}`, userInfo.FirstName, userInfo.LastName, userInfo.DOB, userInfo.CreditScore))

	//makes a POST request with the body from above
	req, err := http.NewRequest("POST", csCardEndPoint, bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, fmt.Errorf("unable to make a post request due to the incorrect body")
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	//retrieves the response body from the POST request
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var csCardResult []CSCardResponse

	var creditCardResults []CreditCard

	//converts json response to CSCardResponse structure
	err = json.Unmarshal(body, &csCardResult)
	if err != nil {
		return nil, fmt.Errorf("unable to reach CSCards API due to the incorrect body")
	}

	//iterates elements of CSCardResponse, convert it to CreditCard struct and appending it to the result array
	for _, result := range csCardResult {
		sc := math.Pow(1/result.Apr, 2)
		creditCard := CreditCard{
			Provider:  "CSCards",
			Name:      result.CardName,
			ApplyURL:  result.URL,
			Apr:       result.Apr,
			Features:  result.Features,
			CardScore: math.Floor((result.Eligibility*sc*10)*1000) / 1000,
		}

		creditCardResults = append(creditCardResults, creditCard)
	}
	//returns the result array of all credit cards received
	return creditCardResults, nil

}

//GetScoredCards sends a post request to ScoredCard API endpoint and formats the response
func (userInfo *UserInfo) GetScoredCards() ([]CreditCard, error) {
	//scoredCardEndPoint := os.Getenv("SCOREDCARDS_ENDPOINT")
	//above was not working with test environment, put the link here as well as in .env
	scoredCardEndPoint := "https://m33dnjs979.execute-api.us-west-2.amazonaws.com/CS/v2/creditcards"

	//makes a body for the POST request with user information received
	var jsonStr = []byte(fmt.Sprintf(`{
		"first-name": "%s",
		"last-name": "%s",
		"date-of-birth": "%s",
		"score": %d,
		"employment-status": "%s",
		"salary": %d
	}`, userInfo.FirstName, userInfo.LastName, userInfo.DOB, userInfo.CreditScore, userInfo.EmpStatus, userInfo.Salary))

	//makes a POST request with the body from above
	req, err := http.NewRequest("POST", scoredCardEndPoint, bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, fmt.Errorf("unable to make a post request due to the incorrect body")
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	//retrieves the response body from the POST request
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var scoredCardResult []ScoredCardResponse

	var creditCardResults []CreditCard

	//converts json response to ScoredCardResponse structure
	err = json.Unmarshal(body, &scoredCardResult)
	if err != nil {
		return nil, fmt.Errorf("unable to reach ScoredCards API due to the incorrect body")
	}
	var features []string

	//iterates elements of ScoredCardResponse, convert it to CreditCard struct and appending it to the result array
	for _, result := range scoredCardResult {
		sc := math.Pow(1/result.Apr, 2)
		//combines attributes and introductory offers into one feature array
		for _, attr := range result.Attributes {
			features = append(features, attr)
		}
		for _, introOffer := range result.IntroOffers {
			features = append(features, introOffer)
		}

		creditCard := CreditCard{
			Provider:  "ScoredCards",
			Name:      result.Card,
			ApplyURL:  result.ApplyURL,
			Apr:       result.Apr,
			Features:  features,
			CardScore: math.Floor((result.ApprovalRating*100*sc)*1000) / 1000,
		}
		creditCardResults = append(creditCardResults, creditCard)
	}
	//returns the result array of all credit cards received
	return creditCardResults, nil
}
