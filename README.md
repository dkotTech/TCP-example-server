# TCP Test Work

This project demonstrates a TCP-based client-server application written in Go. The server provides a Proof-of-Work (PoW) challenge to the client, and upon successful completion, sends a random quote. The client can also send echo messages to the server, which are returned after a delay.

## How It Works

### Server
The server listens for incoming TCP connections and handles them concurrently. For each connection:
1. **Challenge Generation**: The server generates a random PoW challenge and sends it to the client.
2. **Challenge Verification**: The server verifies the client's solution to the PoW challenge.
3. **Quote Response**: Upon successful verification, the server sends a random quote to the client.
4. **Echo Handling**: The server listens for echo messages from the client and responds after a delay.

### Client
The client connects to the server and interacts as follows:
1. **Challenge Solving**: The client receives the PoW challenge, solves it, and sends the solution back to the server.
2. **Quote Display**: The client displays the random quote received from the server.
3. **Echo Messages**: The client allows the user to input messages, which are sent to the server and echoed back.

## Running the Application

### Server
To start the server, run the following command in one terminal:
```bash
go run cmd/server/main.go
```

### Client
```bash
go run cmd/client/main.go
```

# Why SHA-256 is Used as the Proof-of-Work Method

SHA-256 is chosen for the PoW mechanism because it is secure, efficient, and widely supported. It ensures the integrity of the challenge by producing unique hashes that are resistant to tampering. The difficulty of the challenge can be easily adjusted by requiring a specific number of leading zeros in the hash, making it scalable and predictable.

Additionally, SHA-256 is computationally efficient, allowing the client to solve challenges quickly while ensuring fairness. Its widespread use in critical systems like cryptocurrencies demonstrates its reliability and trustworthiness. The server can also verify solutions easily, ensuring smooth communication between the client and server.

By using SHA-256, this project leverages a proven cryptographic method to implement a secure and efficient PoW mechanism.

# Build Dockers

```bash
docker build -f Dockerfile.server -t server:v0.0.1 .
docker run -p 8080:8080 server:v0.0.1
```

```bash
docker build -f Dockerfile.client -t client:v0.0.1 .
docker run -it client:v0.0.1 -address=172.17.0.1:8080
```
