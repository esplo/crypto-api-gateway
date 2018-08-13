# Paid API by mining Cryptocurrency

This repository contains three programs to illustrate a paid API architecture based on Coinhive.

[visit demo](https://mining-gateway.firebaseapp.com/)

[![demo](https://github.com/esplo/crypto-api-gateway/blob/master/cryp-paid-api.gif)]

- backend sample -> `api-sample`
    - This is just an ordinary API... No need to consider mining or payment.
    - go, GAE
- gateway sample -> `gateway`
    - A hard worker, which communicates with Coinhive, processes a payment, and calls the backend API.
    - go, GAE
- client sample -> `client-sample`
    - This program calls an endpoint on the gateway, and displays what is happening.
    - js, static site (i.e. Firebase Hosting, Github Pages, S3, and so on)
