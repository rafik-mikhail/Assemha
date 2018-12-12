# Base container image
FROM golang:1.8-alpine AS builder

# Using Alpine's apk tool, install git which
# is used by Go to download packages
RUN apk --no-cache -U add git

# Install package manager
RUN go get -u github.com/kardianos/govendor

# Copy app files into container
WORKDIR /go/src/app
COPY . .

# Install dependencies
RUN govendor sync

# Build the app
RUN govendor build -o /go/src/app/splitwise

#------------------------------------------------#
FROM golang:1.8-alpine AS builderr

#Certificates to send mail
RUN apk update && apk add --no-cache ca-certificates
RUN update-ca-certificates
#------------------------------------------------#

# Smallest container image
FROM alpine

# Copy built executable from base image to here
COPY --from=builder /go/src/app/splitwise /
COPY --from=builderr /etc/ssl/certs /etc/ssl/certs

# Run the app
CMD [ "/splitwise" ]
