# Smol

This is a small web application that takes small payloads. I built it as a small, and hopefully easy challenge.

It is written in Go because I was learning go while I put it together!

## Usage
```
git clone https://www.github.com/jhicks-r7/smol
cd smol
go run .
```

git clone the repo.
`git clone https://www.github.com/jhicks-r7/smol`

cd to the directory
`cd smol`

run
`go run .`



`go run . -address {LISTENER ADDRESS} -port {LISTENER PORT}`

Optionally, specify a listener address and port. The defaults are 127.0.0.1 and 8000
```
  -address string
        IPv4 Address to Listen On (default "127.0.0.1")
  -port string
        TCP Port to Listen On (default "8000")
```
Once smol is up and running, access the application via your web browser of choice and have fun!
