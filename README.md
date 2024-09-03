To extract kafka messages into a file from kafka server.

## UseCase: 
To extract all messages in bulk from a PROD/QA env to analyse any messages.
        Especially for DLQ messages this could be more helpful.

## Usage:

#### WITH SSL:
```bash
./kafctl -b <broker> -t <topic> -g <consumer-group> -o <consumer_out.json> -s -f <ssl-config.json>
./kafctl -b kafka-service.com:9093 -t topic_internal -g xconsumer-id-5 -o consumer_out.json -s -f config.json 
```

#### WITHOUT SSL:
```bash
./kafctl -b <broker> -t <topic> -g <consumer-group> -o <consumer_out.json>
./kafctl -b localhost:9092 -t purchases -g xconsumer-id-5 -o consumer_out.json
```
```bash
config.json:
{
    "securityProtocol": "SSL",
    "sslCaLocation": "./ssl/cacerts.pem",
    "sslCertLocation": "./ssl/public_key.pem",
    "sslKeyLocation": "./ssl/private_key.pem",
    "sslKeyPassword": "mykeypassword"
  }
  ```