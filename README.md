# Building

```bash
git clone https://github.com/kunalsinghdadhwal/redilink
cd redilink
docker-compose up -d
```

# Test 

#### 1. Creating a new link
```bash
curl -X POST -H "Content-type: application/json" -d '{"url" : "www.google.com"}' localhost:3000/api/v1

# for pretty print
curl -X POST -H "Content-type: application/json" -d '{"url" : "www.google.com"}' localhost:3000/api/v1 | jq
```

#### 2. Visit the link sent in response 


#### Thanks
