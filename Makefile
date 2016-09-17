all: slackin

slackin: slackin.go main.go html.go
	CGO_ENABLED=0 go build slackin.go main.go html.go
