# firebase-user-api-go


export GOOGLE_APPLICATION_CREDENTIALS="/filepath/serviceAccountKey.json"

export FIREBASE_WEB_API_KEY="[YOUR FIREBASE WEB API KEY]"

go run firebaseUserManager.go


## How to Use

```go
curl -X POST \
  http://localhost:8080/createuser \
  -H 'cache-control: no-cache' \
  -d '{
  "Email": "user@example.com",
  "Password":"secret",
  "DisplayName":"User Fullname"
}'

curl -X POST \
  http://localhost:8080/signinuser \
  -H 'cache-control: no-cache' \
  -d '{
  "email":"user@example.com",
  "password":"secret"
}'

curl -X POST \
  http://localhost:8080/updateuser \
  -H 'cache-control: no-cache' \
  -d '{
  "Email": "user@example.com",
  "Password":"newPassword",
  "DisplayName":"Different name"
}'
```
