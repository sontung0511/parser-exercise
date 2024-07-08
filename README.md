# Interview Exercise

## 1. Objective

Your goal is to create a simple but functional blockchain parser using Go. This parser will interact with blockchain data by monitoring specific addresses and gathering related transactions. Focus on the Ethereum blockchain (ETH) and use the JSON-RPC interface to interact with the network. You can connect to the Ethereum network through any public gateway. 

Finally, you'll build a basic HTTP API to interact with the parser.

## 2. Task

### Server

You need to implement a blockchain parser that can do the following:

1. Monitor Specified Addresses: The parser should take a list of blockchain addresses and keep an eye on them for any transactions.
2. Collect Transactions: For each monitored address, the parser should gather all related transactions starting from the subscribed date, without needing historical data.
3. HTTP API Exposure: The parser should provide an HTTP API to interact with the gathered data. This API will let users check the current block number, subscribe to an address, and get transactions for subscribed addresses.

### Others

Choose a way to store and manage the list of subscribed addresses and their transactions. For simplicity, in-memory storage works. But to show off your skills, you could use PostgreSQL (SQL) or MongoDB (NoSQL). 

If you go with a database, please set it up using Docker Compose and write appropriate setup instructions.

## 3. APIs

#### Get the Current Block Number
This endpoint should return the current block number of the blockchain, so users know the most recent block that has been processed.

#### Subscribe to an Address
This endpoint should let users subscribe to a specific blockchain address. The request body should contain the address to be monitored. If the address is successfully added to the list of monitored addresses, the API should confirm the subscription.

#### Get Transactions for a Subscribed Address
This endpoint should return the list of transactions for a given subscribed address. The address should be passed as a query parameter. If the address is not subscribed, the API should return a 400 status code indicating a bad request.

## 4. Recommendation

A few rules to follow:

- Use version control to manage your code, preferably Git, with a README.md file for introduction or explanation.
- Your application should run smoothly without any issues.

We recommend:

- Keep your code clean, easy to read, and maintainable.
- The architecture of the solution should be scalable.
- Write some tests for your code.
- Using Docker Compose is a plus if you need to set up infrastructure.

Your time is valuable, we know. Please aim to spend no more than 8 hours on this exercise in total and submit it within a week. If you have additional ideas to make the server more robust, feel free to share them in the README. We'd love to see them.

----

Kindly upload your completed code to a private repository on GitHub. Once done, add https://github.com/qphuongl as a collaborator. After you've finished, please email your contact person to let them know.
