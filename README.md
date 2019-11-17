# ClearScore Backend Test

This API has a single endpoint (POST) that consumes user financial details and returns recommended credit cards based on the credit score.

## Built With:
* `go` version go1.13

## Functions
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


## API endpoint and 

## Running locally

* To run the project locally do `go main.go`.
* To call the API using localhost, end point is `http://localhost:8081/creditcard`.


## Testing and Linting

To run the tests do `go test`.