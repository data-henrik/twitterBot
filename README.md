# TwitterBot for IBM Cloud Code Engine
Some fun and experiments with Golang, Twitter and IBM Cloud Code Engine. Tweet at schedule with the message composed of latest IBM Cloud blog entries.


## (Rough) Instructions

1. Set up Code Engine (CE) project
2. create registry in CE
3. configure .env with credentials / secrets
4. create secret from file in CE
5. build the container image, either in CE or using the Container Registry
6. create the CE app from the image and pass the configured secrets / credentials
7. set up the CE ping subscription and pass the secret key, e.g., 
   ```
   ibmcloud ce sub ping create -n tweety --destination twitterbot --path /tweet
       --schedule '07 4,8,13,17 * * *' --data 'SECRET_KEY=SET_YOUR_SECRET' --ct 'application/x-www-form-urlencoded'
   ```
   or
   ```
   ibmcloud ce sub ping create -n tweety --destination twitterbot --path /tweet --data
    '{"secret_key":"SET_YOUR_SECRET","tweet_string2":"Written in #Golang by @data_henrik and running on #IBMCloud #CodeEngine"}'
   --schedule '07 9,17 * * *'
   ```

### Local testing
1. configure .env
2. `go run twitterBot.go`
3. `curl -X POST localhost:8080/tweet -H "Content-Type: application/json" --data '{"secret_key":"SET_YOUR_SECRET", "feed":"https://blog.4loeser.net/feeds/posts/default","tweet_string1":"@data_henrik recently wrote about %s. The #blog post is available at %s. "}`