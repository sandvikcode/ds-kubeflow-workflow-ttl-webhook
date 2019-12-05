# docker build -t mutateme:local .

FROM golang:1.12-alpine AS build 
ENV CGO_ENABLED 0
ENV GOOS linux

RUN apk add git make openssl

WORKDIR /go/src/k8s-kubeflow-mutate-webhook
COPY cmd cmd
COPY pkg pkg 

RUN go get k8s.io/api/admission/v1beta1
RUN go get github.com/argoproj/argo/pkg/apis/workflow/v1alpha1
RUN go get k8s.io/apimachinery/pkg/apis/meta/v1
RUN go get github.com/stretchr/testify/assert

# To jump to intermiddate stages in the build to debug 
# docker run -it --rm 29b94b0299a8 sh
RUN go test -v ./... -cover
RUN go build -v -o mutate_workflow cmd/main.go

# This is the second part of the build. 
FROM alpine
RUN apk --no-cache add ca-certificates && mkdir -p /app
WORKDIR /app
COPY --from=build /go/src/k8s-kubeflow-mutate-webhook/mutate_workflow .
CMD ["/app/mutate_workflow"]
