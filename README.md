# GO-Challenge

## Description

Http web service handling simple math operations with implemented caching of repeating operands and operations.

Math operations supported:

- add
- subtract
- multipply
- divide

## Requirements

- go version go1.18

## How to install and run

Run the web server by executing this command in the terminal:

```
$ go run main.go
```

You can now send the requests to your localhost:8090.

## Hot to use

Server handles only GET requests.
Required url and query parameters are:

**localhost:8090/{operation}?x={integer}&y={integer}**

For example:

http://localhost:8090/divide?x=4&y=2
http://localhost:8090/add?x=5&y=2

Big integers and floats are not supported.

## Deployment options

Determined by the context of the app usage in regards to availability, reliability, cost, internal or public usage.

Example in AWS:

Web app is used by the public users, we created a domain math.com and registered it in the Route53.
EC2 instance is created, if that same instance will be reused for other apps or services than we create an image for our app and run it in container.

Elastic IP is associated with this instance and mapped to a DNS record in Route53 under math.com hosted zone.
Basic networking is required with VPC, Public Subnet and Internet Gateway configuration.

If we expect high load and have a need for high availability we can create an autoscaling group with load balancer in front.

Other options are EKS (Kuberneetes managed service), ECS (AWS EKS equivalent with options to deploy with EC2 or Fargate)

With some modifications we can create Serverless deployment in AWS Lambda with API Gateway.
