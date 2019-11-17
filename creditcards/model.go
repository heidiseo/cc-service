package creditcards

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
