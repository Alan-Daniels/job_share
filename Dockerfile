ARG BUILD_FROM
FROM $BUILD_FROM

# Install requirements
#RUN \
#	apk add --no-cache \
#	go

# Set destination for COPY
WORKDIR /app

# Download Go modules
#COPY go.mod go.sum ./
#RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
#COPY *.go ./

# Build
#RUN CGO_ENABLED=0 GOOS=linux go build -buildvcs=false -o /job_share

COPY job_share /
CMD [ "/job_share" ]
