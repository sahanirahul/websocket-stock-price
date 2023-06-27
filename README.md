## Go Api Server README

This Go project consists of 3 component:
- API server (serving 2 *GET API*)
    1. `"http://localhost:19093/underlying-prices"`

        Returns list of underlying instruments with latest prices

        Responses are:

        StatusCode: 200
        ```css
        {
            "payload": [
                {
                    "symbol": "SHELL",
                    "underlying": "SHELL",
                    "token": 6873616,
                    "instrument_type": "EQ",
                    "expiry": "",
                    "strike": 0,
                    "price": 22193.39033446994
                }
            ],
            "success": true
        }
        ```

        StatusCode: 422 (UnprocessableEntity)
        ```css
        {
            "error": "error message",
            "success": false
        }
        ```
    2. `"http://localhost:19093/underlying-prices/:symbol"`

        Returns list of derivative instrument with latest prices for the requested underlying symbol

        Responses are:

        StatusCode: 200
        ```css
        {
            "payload": [
                {
                    "symbol": "BANKNIFTY20Mar32000.00CE",
                    "underlying": "BANKNIFTY",
                    "token": 5886662,
                    "instrument_type": "CE",
                    "expiry": "2020-03-26",
                    "strike": 32000,
                    "price": 238.5717775283092
                }
            ],
            "success": true
        }
        ```

        StatusCode: 422 (UnprocessableEntity)
        ```css
        {
            "error": "error message",
            "success": false
        }
        ```
        StatusCode: 400 (Bad Request if sumbol not present in URI)
        ```css
        {
            "error": "invalid_request",
            "success": false
        }
        ```
    
- Cron Jobs (2 Jobs) : Runs in a predefined interval of time
    1. Equity-Cron

        `every 15 min`
        - Fetch latest list of Underlyings from the broker
        - URI used: `https://prototype.sbulltech.com/api/underlyings`
        - Creates entry in database and subscribes the new underlyings to webscoket for price updates
        - Deletes the entry of old underlyings not present in the new list and unsubscribes from websocket

    2. Derivative-Cron

        `every 1 min`
        - Fetch latest list of derivatives for each underlyings from the broker 
        - URI used : `https://prototype.sbulltech.com/api/derivatives/:token`
        - Creates entry in database and subscribes the new derivatives to webscoket for price updates
        - Deletes the entry of old derivatives not present in the new list and unsubscribes from websocket


- Websocket : Read/Write from Websocket server
    - WebsocketURL used: `wss://prototype.sbulltech.com/api/ws`
    - Reads message and updates price
    - Subscribe for price update


#### Database Used: REDIS
- Data Model

    Each Instrument is stored with token as key
    ```json
        "5886662":{
            "symbol": "BANKNIFTY20Mar32000.00CE",
            "underlying": "BANKNIFTY",
            "token": 5886662,
            "instrument_type": "CE",
            "expiry": "2020-03-26",
            "strike": 32000,
            "price": 238.5717775283092
        }
    ```
    List of latest underlyings is stored as a Set of token against the key "EQ:ALLUNDERYINGS"
    ```json
        "EQ:ALLUNDERLYINGS":[5886662,5886664,5886667]
    ```

    List of latest derivatives for each underlyings is stored as a Set of token against the key "DERIVATIVES:{TOKEN}"
    ```json
        "DERIVATIVES:5886662":[8137866,813786,8857866]
    ```

    Symbol to Token mapping is stored
    ```json
        "SHELL":6873616
    ```

#### how to run

- As Docker Container
    - Make use you have docker and docker-compose installed on your system
    - cd /path/to/project/dir
    - RUN `docker-compose up` to START the services
    - RUN  `docker-compose down` to STOP the services

- On local machine
    - Make sure redis is install and running
    - cd /path/to/project/dir
    - RUN `go run .`
    - *REDIS DB connection details should be added in config/config.local.json file*


#### worker pool
- the program contains a generic worker pool implementation in go for parallel processing of request. At the same time the workers can be configured as per requirement. This is better than simply spawning go routines, if we were to spawn multiple go routine for a single request, the system can crash when there is a outburst of request. Also spawning a new go routine is costlier than a worker go routine that processes task by fetching from a job queue (channel)

- Here we have 2 worker pool:
    1. Websocket worker pool (20 Workers specified by WEBSOCKET_WORKER_POOL_SIZE env variable in docker-compose.yaml)
        - Used by Websocket class to update prices when received from webscoket server
    2. Service worker pool (10 Workers specified by SERVICE_WORKER_POOL_SIZE env variable in docker-compose.yaml)
        - Used by Service class to update derivative list
