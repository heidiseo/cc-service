package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/onsi/gomega"
)

//TestHandler makes a mock http request and tests if the response is correct
func TestHandler(t *testing.T) {
	//test tool
	g := gomega.NewGomegaWithT(t)

	//makes a mock request body for POST request for creditcard
	reqBody := []byte(`{
		"firstname": "John",
		"lastname": "Smith",
		"dob": "1991/04/18",
		"credit-score": 500,
		"employment-status": "FULL_TIME",
		"salary": 30000
	}`)
	//converts the format	[]byte to bytes.Reader
	body := bytes.NewReader(reqBody)
	//makes a mock http request wit the mock body
	req, err := http.NewRequest("POST", "/Handler", body)
	if err != nil {
		t.Fatal(err)
	}

	//creates a ResponseRecorder to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Handler)

	//directly passes in the Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	//checks the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	//checks the response body is what we expect.
	respBody, _ := ioutil.ReadAll(rr.Body)
	var creditCardResults []CreditCard
	//convers json to CreditCard struct
	err = json.Unmarshal(respBody, &creditCardResults)
	//expected response in CreditCard struct
	test := []CreditCard{
		{
			Provider:  "ScoredCards",
			Name:      "ScoredCard Builder",
			ApplyURL:  "http://www.example.com/apply",
			Apr:       19.4,
			Features:  []string{"Supports ApplePay", "Interest free purchases for 1 month"},
			CardScore: 0.212,
		},
		{
			Provider:  "CSCards",
			Name:      "SuperSaver Card",
			ApplyURL:  "http://www.example.com/apply",
			Apr:       21.4,
			Features:  nil,
			CardScore: 0.137,
		},
		{
			Provider:  "CSCards",
			Name:      "SuperSpender Card",
			ApplyURL:  "http://www.example.com/apply",
			Apr:       19.2,
			Features:  []string{"Interest free purchases for 6 months"},
			CardScore: 0.135,
		},
	}
	//compares the response received and the expected response above
	g.Expect(creditCardResults).To(gomega.Equal(test))
}

//TestGetCSCards tests if it receives the correct information from CSCards API
func TestGetCSCards(t *testing.T) {
	//test tool
	g := gomega.NewGomegaWithT(t)
	//makes test cases with struct in order to iterate instead of repeating
	tests := []struct {
		Message string
		UInfo   UserInfo
		Error   string
	}{
		{Message: "should fail as the body of the request is missing user's date of birth",
			UInfo: UserInfo{
				FirstName:   "John",
				LastName:    "Smith",
				CreditScore: 500,
				EmpStatus:   "FULL_TIME",
				Salary:      30000,
			},
			Error: "unable to reach CSCards API due to the incorrect body",
		},
		{Message: "should not fail as the body of the request is correct and the request was successfully made",
			UInfo: UserInfo{
				FirstName:   "John",
				LastName:    "Smith",
				DOB:         "1991/04/18",
				CreditScore: 500,
				EmpStatus:   "FULL_TIME",
				Salary:      30000,
			},
			Error: "",
		},
	}
	//makes the expected response from CSCards
	csCards := []CreditCard{
		{
			Provider:  "CSCards",
			Name:      "SuperSaver Card",
			ApplyURL:  "http://www.example.com/apply",
			Apr:       21.4,
			CardScore: 0.137,
		},
		{
			Provider:  "CSCards",
			Name:      "SuperSpender Card",
			ApplyURL:  "http://www.example.com/apply",
			Apr:       19.2,
			Features:  []string{"Interest free purchases for 6 months"},
			CardScore: 0.135,
		},
	}

	//iterates the tests, checks error codes an compares the response body with the expected response above
	for _, test := range tests {
		t.Run(test.Message, func(t *testing.T) {
			creditcards, err := test.UInfo.GetCSCards()
			if err != nil {
				g.Expect(err).To(gomega.HaveOccurred())
				g.Expect(err.Error()).To(gomega.Equal(test.Error))
			} else {
				g.Expect(err).NotTo(gomega.HaveOccurred())
				g.Expect(creditcards).To(gomega.Equal(csCards))
			}
		})
	}
}

func TestGetScoredCards(t *testing.T) {
	//test tool
	g := gomega.NewGomegaWithT(t)
	//makes test cases with struct in order to iterate instead of repeating
	tests := []struct {
		Message string
		UInfo   UserInfo
		Error   string
	}{
		{Message: "should fail as the body of the request is user's employment status",
			UInfo: UserInfo{
				FirstName:   "John",
				LastName:    "Smith",
				DOB:         "1991/04/18",
				CreditScore: 500,
				Salary:      30000,
			},
			Error: "unable to reach ScoredCards API due to the incorrect body",
		},
		{Message: "should not fail as the body of the request is correct and the request was successfully made",
			UInfo: UserInfo{
				FirstName:   "John",
				LastName:    "Smith",
				DOB:         "1991/04/18",
				CreditScore: 500,
				EmpStatus:   "FULL_TIME",
				Salary:      30000,
			},
			Error: "",
		},
	}
	//makes the expected response from ScoredCards
	csCards := []CreditCard{
		{
			Provider:  "ScoredCards",
			Name:      "ScoredCard Builder",
			ApplyURL:  "http://www.example.com/apply",
			Apr:       19.4,
			Features:  []string{"Supports ApplePay", "Interest free purchases for 1 month"},
			CardScore: 0.212,
		},
	}
	//iterates the tests, checks error codes an compares the response body with the expected response above
	for _, test := range tests {
		t.Run(test.Message, func(t *testing.T) {
			creditcards, err := test.UInfo.GetScoredCards()
			if err != nil {
				g.Expect(err).To(gomega.HaveOccurred())
				g.Expect(err.Error()).To(gomega.Equal(test.Error))
			} else {
				g.Expect(err).NotTo(gomega.HaveOccurred())
				g.Expect(creditcards).To(gomega.Equal(csCards))
			}
		})
	}
}
