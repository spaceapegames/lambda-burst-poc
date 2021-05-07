## Lambda Burst

Repo containing accompaniment code for this [Space Ape Games blog article]().

### Overview

Please see the above linked article for a full explanation.

Essentially this code builds the following:

- A container image that is deployable to both Lambda and Fargate.
- An ECS Cluster with a Fargate service that uses the image.
- A Lambda function that also uses the image.
- An ALB with a separate Target Group for both Lambda and Fargate.

The application being deployed is a Golang-based web application which calculates all of the prime numbers between 0
and `MAX_PRIME`.

It runs on Fargate, and has a rate-limiter which, when triggered, redirects client requests to the Lambda backend.

### Building an Image for both Lambda and Fargate

There are some quirks to consider when building an image that is compatible for both formats.

First, the code needs to be able to respond to both HTTP requests (Fargate) and Lambda API Gateway Events.

We use the [apex-gateway](https://github.com/apex/gateway) library here but there is an [official AWS offering](https://github.com/awslabs/aws-lambda-go-api-proxy).

We then switch between the two with this statement:

```go
func (s *Server) Serve() {
	if s.lambdaMode {
		log.Println("running in Lambda mode")
		gateway.ListenAndServe(":8000", s.Router)
	} else {
		log.Printf("running in http mode on port %d", s.port)
		http.ListenAndServe(fmt.Sprintf(":%d", s.port), s.Router)
	}
}
```

For the image, we:

- Use the official [AWS Lambda image](https://hub.docker.com/r/amazon/aws-lambda-go) such that Lambda is supported out-of-the-box.
- Copy the pre-built binary to `/var/task`, which is where the entrypoint of that image expects it to be.

This gives us an image which will work for Lambda. The image itself contains a `lambda-entrypoint.sh` that understands how
to invoke the binary in the Lambda context.

For the Fargate tasks, we simply override the entrypoint to point directly at the binary, bypassing the Lambda gubbins.