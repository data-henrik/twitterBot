# TwitterBot for IBM Cloud Code Engine
Some fun and experiments with Golang, Twitter and IBM Cloud Code Engine. Tweet at schedule with the message composed of latest IBM Cloud blog entries.


## (Rough) Instructions

1. Set up Code Engine (CE) project
2. create registry in CE
3. configure .env with credentials / secrets
4. create secret from file in CE
5. build the container image, either in CE or using the Container Registry
6. create the CE app from the image and pass the configured secrets / credentials
7. set up the CE ping subscription and pass the secret key
   
