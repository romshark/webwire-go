# Example: PubSub

This example demonstrates the use of the server-side signals.
The client connects to the server and listens for N incoming signals (6 by default)
until it disconnects, while the server constantly broadcasts the current time
to all connected clients.
