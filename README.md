# The Beer Game

This is an implementation of [The Beer Distribution Game](https://en.wikipedia.org/wiki/Beer_distribution_game).

## Setup

The app requires [Go](https://golang.org/) and [NodeJS + NPM](https://nodejs.org/). Once you have both installed you also need to install the client dependancies:

```
cd client
npm install
```

## Development

To run the server:
```
cd server
go run beergame
```

To run the client:
```
cd client
npm run dev
```

## Deployment

To run in Docker:
```
cd client
npm run build
cd ../server
docker build -t beergame .
docker run --rm -p 80:80 beergame
```
