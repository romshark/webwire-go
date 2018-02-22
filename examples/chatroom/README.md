# Example: Chat Room

This example demonstrates the advanced use of the request-response topology, signals in both directions,
authentication and sessions. Clients connect to the server to listen for broadcasted messages
and publish their own messages in the public chat room.
Messages are sent anonymously by default, though clients can authenticate themselves
signing in to one of the predefined user accounts to make the server associate their messages with their names.
