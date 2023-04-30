FROM ubuntu:latest 

# Set the working directory
WORKDIR /app

# Copy the source code to the container
COPY . /app

# Install any necessary dependencies
RUN apt-get update && apt-get install -y make protobuf-compiler protoc-gen-go golang-go ca-certificates

RUN update-ca-certificates

# Compile the program using Make
RUN make build

# Expose port 2000
EXPOSE 2000

# Start the compiled program using Make
CMD make run-proxy
