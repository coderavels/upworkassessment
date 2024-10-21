# Virtual Bookstore

This repository contains the frontend and the backend logic for the virtual bookstore. It fetches the books from a third party client `assess` and presents it organised in shelves of fixed width.

## Setup Instructions
```
cd code
make setup
make start
```
NOTE: This would start the service at port 8022 by default. To change the port, update the PORT env in the .env file

## To run tests
```
make test
```

## Future Items
1. Safegaurd api routes with authentication
2. Cleanup the UI for better UX
3. Introduce data caching at server end for better response time.

