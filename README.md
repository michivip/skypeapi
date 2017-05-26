# Golang Skype API #

### Warning ###
**This API is still under heavy development. Bugs can occur - please report them in the issue section.**
**Therefore please do not use this API for professional use.**
### Introduction ###
This API is based on the 
[Bot Framework by Microsoft](https://docs.microsoft.com/de-de/bot-framework/rest-api/bot-framework-rest-connector-api-reference).
 For further details on how to register a bot see the official [Tutorial](https://docs.microsoft.com/en-us/bot-framework/rest-api/bot-framework-rest-connector-quickstart).
### Features ###
* requesting an authentication token
* creating a valid HTTPS endpoint and parsing activity objects you can work with
### Requirements ###
* files to setup SSL endpoint (both of them have to be valid CA certificates)
    * certificate file (e.g. *fullchain.pem*)
    * private key file (e.g. *privkey.pem*)
* some Golang knowledge
### License ###
The source code is licensed under the MIT license. For further details see the LICENSE file.
### Examples ###
You can find some examples inside the examples folder. **Warning:** These are just examples. Error handling and other 
good practices are not used.
