# CaseForBiTaksi

1. A **Driver Location API** that uses location data stored in a MongoDB collection,
2. A **Matching API Service** that finds the nearest driver with the rider location using the **Driver Location API**.

## How can you test it locally?
1. Clone the repository
```
git clone https://github.com/mervemor/CaseForBiTaksi.git
```
2. You need to run both APIs **Driver Location API** and **Matching API Service**
```
/DriverLocationAPI/cmd/server/main.go
/MatchingAPI/cmd/server/main.go
```
3. You can use these curls. Paste it to Postman or terminal.

For POST /find-driver and authenticated is true
```
curl --location 'http://localhost:8081/find-driver' \
--header 'Content-Type: application/json' \
--header 'Apikey: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRoZW50aWNhdGVkIjp0cnVlfQ.E-rHE1oDGHS9L7caiSnTlkEXaHvCvXCfAgz1xaaClbs' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRoZW50aWNhdGVkIjp0cnVlfQ.E-rHE1oDGHS9L7caiSnTlkEXaHvCvXCfAgz1xaaClbs' \
--data '{
    "type": "Point",
    "coordinates": [40.2144912, 29.0390297],
    "radius": 5000
}'
```
if authenticated is false
```
curl --location 'http://localhost:8081/find-driver' \
--header 'Content-Type: application/json' \
--header 'Apikey: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRoZW50aWNhdGVkIjpmYWxzZX0.RC8NOw_GTtW3C4HNtv441FNPmscHXX6kfIq_DiIAdtg' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRoZW50aWNhdGVkIjpmYWxzZX0.RC8NOw_GTtW3C4HNtv441FNPmscHXX6kfIq_DiIAdtg' \
--data '{
    "type": "Point",
    "coordinates": [40.2144912, 29.0390297],
    "radius": 5000
}'
```

For POST /upsert-driver. If there is no data, it inserts, if data exists, it updates. It can be made for more than one driver at the same time.
```
curl --location 'http://localhost:8080/upsert-driver' \
--header 'Content-Type: application/json' \
--data '
[
    {
        "id": "65fd797d77cdafd494963987",
        "location": {
            "type": "Point",
            "coordinates": [40.94289771, 29.0390297]
        }
    }

]'
```

## Documentation

After running the APIs, you can access the Swagger pages.
- DriverLocationAPI http://localhost:8080/swagger/index.html
- MatchingAPI http://localhost:8081/swagger/index.html
