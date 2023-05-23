## store-monitoring
step 1 :-
clone this repo under $GOPATH/src/github.com/

step 2 :- 
download all the CVs and put this inside data folder like below
     #### Files structure 
     
     ```
        data/
        
          ├── business_hours.csv
          
          ├── stores.csv
          
          └── timezones.csv
        ```

step 3:- start docker daemon because this service is running mysql container 

step 4:- make run-local
 it will run docker compose up and also the server , data will be loaded in database so it might take few minutes depending on system

Following are the apis:- 
API - 1 :-  it triggers the report , you can only trigger 10 concurrent reports at the same time. No caching has been used , so it might take same time after each trigger
```
curl --location 'http://0.0.0.0:8080/trigger_report'
```

API -2 :- it will give the csv generated from above api , if the report generation is not completed it will give "Running" in the response otherwise will give csv data in response.

```
curl --location 'http://0.0.0.0:8080/get_report?report_id={report-id}'
```


# Few other details 

1) if you restart the server all tables will be dropped and all the data from the table will be removed .
2) each report is pushed in "/data" folder
