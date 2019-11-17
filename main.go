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
	FirstName   string `json:"firstname"`
	LastName    string `json:"lastname"`
	DOB         string `json:"dob"`
	CreditScore int    `json:"credit-score"`
	EmpStatus   string `json:"employment-status"`
	Salary      int    `json:"salary"`
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
		fmt.Fprintf(w, "please enter the body in JSON")
		return
	}
	var creditcards []CreditCard
	csCardsResults, err := newUserInfo.GetCSCards()
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, "unable to retrieve CSCards")
		return
	}
	for _, csCardsResult := range csCardsResults {
		creditcards = append(creditcards, csCardsResult)
	}
	scoredCardsResults, err := newUserInfo.GetScoredCards()
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, "unable to retrieve ScoredCards")
		return
	}
	for _, scoredCardResult := range scoredCardsResults {
		creditcards = append(creditcards, scoredCardResult)
	}

	sort.SliceStable(creditcards, func(i, j int) bool {
		return creditcards[j].CardScore < creditcards[i].CardScore
	})

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(creditcards)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "failed to response in JSON")
	}
}

//GetCSCards sends a post request to CSCard API endpoint and formats the response
func (userInfo *UserInfo) GetCSCards() ([]CreditCard, error) {
	// csCardEndPoint := os.Getenv("CSCARDS_ENDPOINT")
	csCardEndPoint := "https://y4xvbk1ki5.execute-api.us-west-2.amazonaws.com/CS/v1/cards"

	if csCardEndPoint == "" {
		log.Fatal("$CSCARDS_ENDPOINT must be set")
		return nil, fmt.Errorf("CSCards API endpoint is incorrect")
	}

	var jsonStr = []byte(fmt.Sprintf(`{
		"fullName": "%s %s",
		"dateOfBirth": "%s",
		"creditScore": %d
	}`, userInfo.FirstName, userInfo.LastName, userInfo.DOB, userInfo.CreditScore))

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

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var csCardResult []CSCardResponse

	var creditCardResults []CreditCard

	err = json.Unmarshal(body, &csCardResult)
	if err != nil {
		return nil, fmt.Errorf("unable to reach CSCards API due to the incorrect body")
	}
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
	return creditCardResults, nil

}

//GetScoredCards sends a post request to ScoredCard API endpoint and formats the response
func (userInfo *UserInfo) GetScoredCards() ([]CreditCard, error) {
	// scoredCardEndPoint := os.Getenv("SCOREDCARDS_ENDPOINT")
	scoredCardEndPoint := "https://m33dnjs979.execute-api.us-west-2.amazonaws.com/CS/v2/creditcards"

	if scoredCardEndPoint == "" {
		log.Fatal("$SCOREDCARDS_ENDPOINT must be set")
		return nil, fmt.Errorf("ScoredCards API endpoint is incorrect")
	}

	var jsonStr = []byte(fmt.Sprintf(`{
		"first-name": "%s",
		"last-name": "%s",
		"date-of-birth": "%s",
		"score": %d,
		"employment-status": "%s",
		"salary": %d
	}`, userInfo.FirstName, userInfo.LastName, userInfo.DOB, userInfo.CreditScore, userInfo.EmpStatus, userInfo.Salary))
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

	body, _ := ioutil.ReadAll(resp.Body)

	var scoredCardResult []ScoredCardResponse

	var creditCardResults []CreditCard

	err = json.Unmarshal(body, &scoredCardResult)
	if err != nil {
		return nil, fmt.Errorf("unable to reach ScoredCards API due to the incorrect body")
	}
	var features []string
	for _, result := range scoredCardResult {
		sc := math.Pow(1/result.Apr, 2)
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
	return creditCardResults, nil
}
