# smtp-email-forward

This is a simple smtp server that can be hosted in a docker container and be used to forward the email to another address. This project is part of the blog post [Simple SMTP server to forward your emails](http://hodo.dev/posts/post-20-smtp-server/).

## Features

- receive email multiple domains
- save emails to GCP bucket
- forward the emails to another address using [mailgun](https://www.mailgun.com)

## Build and run

```bash
# clone the respository
git clone https://github.com/gabihodoroaga/smtp-email-forward.git
# build
cd cmd
go build
```

In order to run the project you need to create also a GCP bucket 

```bash
GCP_BUCKET=smtpd-email-forward-123
gsutil mb -c nearline gs://$GCP_BUCKET
```

and a service account

```bash
# create the service account
gcloud iam service-accounts create smtpd-email-forward-123 \
    --description="A service account to write to the bucker" \
    --display-name="smtpd-email-forward-123"
# grab the email
GCP_SERVICE_ACCOUNT=$(gcloud iam service-accounts list --format="value(email)" --filter="displayName=smtpd-email-forward-123")
# give the service account the permission to create objects
gsutil iam ch serviceAccount:$GCP_SERVICE_ACCOUNT:objectCreator gs://$GCP_BUCKET
# create and save the service account key
mkdir certs
gcloud iam service-accounts keys create certs/gcpServiceAccount.json \
  --iam-account $GCP_SERVICE_ACCOUNT
```

and the service certificates

```bash
openssl req -x509 -newkey rsa:4096 -keyout ./certs/server.key -out ./certs/server.crt -days 365 -nodes
```

In order to test he mailgun integration you need to create a new mailgun account and to setup ```MAILGUN_DOMAIN``` and ```MAILGUN_API_KEY``` environment variables

## Authors

Gabriel Hodoroaga [hodo.dev](https://hodo.dev)

## Referenced projects

- [github.com/mhale/smtpd](https://github.com/mhale/smtpd)
- [github.com/ilyakaznacheev/cleanenv](https://github.com/ilyakaznacheev/cleanenv)
- [go.uber.org/zap](https://github.com/uber-go/zap)

## TODO

- add support for SPF validation
- add support for DKIM validation 
- add support for email list
