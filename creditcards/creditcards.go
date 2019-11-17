package creditcards

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
)

//handler receives the user info, passes it to CSCard and ScoredCard APIs, format and sort the responses
func handler(w http.ResponseWriter, r *http.Request) {
	var newUserInfo UserInfo
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, "please enter user info")
	}
	err = json.Unmarshal(reqBody, &newUserInfo)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, "please enter the body in JSON")
	}
	var creditcards []CreditCard
	csCardsResults := newUserInfo.getCSCards()
	for _, csCardsResult := range csCardsResults {
		creditcards = append(creditcards, csCardsResult)
	}
	scoredCardsResults := newUserInfo.getScoredCards()
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

//getCSCards sends a post request to CSCard API endpoint and formats the response
func (userInfo *UserInfo) getCSCards() []CreditCard {
	csCardEndPoint := os.Getenv("CSCARDS_ENDPOINT")

	var jsonStr = []byte(fmt.Sprintf(`{
		"fullName": "%s %s",
		"dateOfBirth": "%s",
		"creditScore": %d
	}`, userInfo.FirstName, userInfo.LastName, userInfo.DOB, userInfo.CreditScore))

	req, err := http.NewRequest("POST", csCardEndPoint, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var csCardResult []CSCardResponse

	var creditCardResults []CreditCard

	err = json.Unmarshal(body, &csCardResult)
	if err != nil {
		log.Fatal(err)
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
	return creditCardResults

}

//getScoredCards sends a post request to ScoredCard API endpoint and formats the response
func (userInfo *UserInfo) getScoredCards() []CreditCard {
	scoredCardEndPoint := os.Getenv("SCOREDCARDS_ENDPOINT")

	var jsonStr = []byte(fmt.Sprintf(`{
		"first-name": "%s",
		"last-name": "%s",
		"date-of-birth": "%s",
		"score": %d,
		"employment-status": "%s",
		"salary": %d
	}`, userInfo.FirstName, userInfo.LastName, userInfo.DOB, userInfo.CreditScore, userInfo.EmpStatus, userInfo.Salary))
	req, err := http.NewRequest("POST", scoredCardEndPoint, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var scoredCardResult []ScoredCardResponse

	var creditCardResults []CreditCard

	err = json.Unmarshal(body, &scoredCardResult)
	if err != nil {
		log.Fatal(err)
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
	return creditCardResults
}
