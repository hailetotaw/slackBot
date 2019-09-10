# using image 1.12
FROM golang:1.12

LABEL maintainer="Hailegebreal Mamo <hmamo@mum.com>"

# setting working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

#copy the google translate token json
COPY allinone-1126-9140291eec07.json ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o main .

# set the time zone for google translate same us your pc timezone
ENV TZ=America/New_York
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# Expose port 5000 to the outside world
EXPOSE 5000
# Command to run the executable
CMD ["./main"]
