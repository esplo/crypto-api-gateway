# Gateway example

this program receives a request from user's browser, and withdraw a constant amount of XMR from coinhive.

## settings

cost to call an endpoint can be changed by modifying `cost` values in `api.yaml`

## deploy

1. create Coinhive API secret
1. cp app.sample.yaml app.yaml
1. set the environment variables in app.yaml
1. gcloud app deploy
