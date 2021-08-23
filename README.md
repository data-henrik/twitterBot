# TwitterBot for IBM Cloud Code Engine
Some fun and experiments with Golang, Twitter and IBM Cloud Code Engine. Tweet at schedule with the message composed of latest IBM Cloud blog entries.


## (Rough) Instructions

1. Set up a [Code Engine (CE) project](https://cloud.ibm.com/docs/codeengine?topic=codeengine-manage-project).
2. [Add access to a container registry](https://cloud.ibm.com/docs/codeengine?topic=codeengine-add-registry) in Code Engine.
3. Configure a local file **.env** with credentials / secrets. Check [.env.template](.env.template) for a template.
4. [Create a secret](https://cloud.ibm.com/docs/codeengine?topic=codeengine-configmap-secret#secret-create) from file.
5. Build the container image, either in CE or using the Container Registry
6. [Create the CE app](https://cloud.ibm.com/docs/codeengine?topic=codeengine-cli#cli-application-create) from the image and pass the configured secrets / credentials.
7. Set up the [CE cron subscription](https://cloud.ibm.com/docs/codeengine?topic=codeengine-subscribe-cron-tutorial) and pass the secret key, e.g., 
   ```
   ibmcloud ce sub cron create -n tweety --destination twitterbot --path /tweet
       --schedule '07 4,8,13,17 * * *' --data 'SECRET_KEY=SET_YOUR_SECRET' --ct 'application/x-www-form-urlencoded'
   ```
   or
   ```
   ibmcloud ce sub cron create -n tweety --destination twitterbot --path /tweet --data
    '{"secret_key":"SET_YOUR_SECRET","tweet_string2":"Written in #Golang by @data_henrik and running on #IBMCloud #CodeEngine"}' 
    --content-type 'application/json' --schedule '07 9,17 * * *'
   ```

### Local testing
1. configure .env
2. `go run twitterBot.go`
3. `curl -X POST localhost:8080/tweet -H "Content-Type: application/json" --data '{"secret_key":"SET_YOUR_SECRET", "feed":"https://blog.4loeser.net/feeds/posts/default","tweet_string1":"@data_henrik recently wrote about %s. The #blog post is available at %s. "}`
