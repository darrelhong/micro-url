# MicroURL

Yet another URL shortener, but with a twist

Login with GitHub to get you very own database

Live demo: See above 

All short URLs will be stored in your tenant database

## Project Structure

`main.go`

Initialise database connection and server

`routes.go`

Add handlers to server

`handlers`

Contains server handler logic for each path

`store`

Logic for interacting with database
