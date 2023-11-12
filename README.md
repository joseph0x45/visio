# VISIO
## Introduction
Visio is an opensource cloud based service that provides algorithms for face detection and recognition. It is still in active development and is not feature complete yet but is stable enough to be
used in small to medium scale projects.

## How it work
I did not reinvent the wheel, Visio is currently based on go-face which is a go library that provides bindings for dlib. This can change anytime soon though. The end goal is to have a system as 
performant as possible, a rewrite using a library that would not call C code from Go can happen at anytime, but for now I found go-face to be quicker to work with to get an MVP even though somewhat
limited, for example on the type of images that it can process (only jpeg format is supported currently)

## Usage
To get started with Visio, first go to the [website](https://getvisio.cloud) and authenticate using your github account. Then you will be redirected to the console. Once there, go to the `keys` tab
and create an API key. You can now start calling Visio from your services using the API keys. For node users, a npm package should be on the way. Feel free to contribute to it :)

Here are the available features and how you can use them:

Note that all requests should be authenticated using the "Authorization" header which value should be the generated API key

* Detecting a face 
POST https://api.getvisio.cloud

This request takes in a file field in the multipart form body and returns the descriptor of a face if detected on the image. If no face is detected or more than one face is detected, a 400 error is
returned

* Creating a face
POST https://api.getvisio.cloud/v1/faces

This request detect a face in the image, stores it's descriptor in the database and returns an id corresponding to that face. THe user can use that id to reference the created face in other requests

* Comparing two faces with their id
This request takes in 2 ids corresponding to faces that should have been created on Visio and compare them, a verdict is returned

POST https://api.getvisio.cloud/v1/faces/compare


* Comparing two faces with an id and an image
This request takes in an id and an image, the id corresponding to a face created on Visio. The face on the image is detected and compared to the face corresponding to the id sent

POST https://api.getvisio.cloud/v1/faces/compare-mixt
