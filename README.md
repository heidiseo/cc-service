# ClearScore Backend Test

This API has a single endpoint (POST) that consumes user financial details and returns recommended credit cards based on the credit score.

## Running and API endpoint

* To run the project both locally(PORT=5000) and deployed version run `go build -o bin/go-getting-started -v .` and `heroku local web`.
* To call the API using localhost, endpoint is `http://localhost:5000/creditcard`.
* As Golang requires Golang Environments to set up locally in order to run the file, I have deployed the API to Prod. To call the Prod API using Postman/Insomnia, endpoint is `https://heidi-cs-cc-service.herokuapp.com/v1/creditcard`.

## Built With:
* `go` version go1.13

## Environments
PORT(localhost:5000), CSCARDS_ENDPOINT(CSCards API endpoint) and SCOREDCARDS_ENDPOINT(ScoredCards API endpoint).

## Functions(main.go)
    `handler` function
        * receives the user financial details from the body of the post request 
        * passes the information to `getCSCards` and `getScoredCards` APIs
        * receives the formated credit cards result in CreditCard struct.
        * combines the results from both APIs
        * sorts the results by card score
    
    `getCSCards` function
        * sends a post request to CSCards API with the information received from the body of the creditcard post request
        * stores the result in CreditCard struct
        * calculates the card score based on the eligibility and the APR received
        * returns all the credit card results

    `getScoredCards` function
        * sends a post request to ScoredCards API with the information received from the body of the creditcard post request
        * stores the result in CreditCard struct
        * combines attributes and introductory-offers
        * calculates the card score based on the approval-rating and APR received
        * returns all the credit card results


## Tests(main_test.go)

To run the tests do `go test`.