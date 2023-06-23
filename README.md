## Outline
Creates a REST API written in Go using only the standard library. This API handles a single endpoint which takes a number as an input and compares it to the current price of Bitcoin as compared to the mean price from 2 different exchanges. This comparison is done in parallel using goroutines and channels and returns the difference.

Example: User inputs $20,000 but the price of bitcoin is $22,000 on exchange A and $22,500 on exchange B, the mean between these two exchanges is $22,250 so the response is $2,250.

Server handles possible errors from each exchange (i.e. rate limiting), and also timeout-s the request if it lasts more than 5 seconds. Mean estimation functionality is verified using unit tests.
