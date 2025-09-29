Real-Time Chat Application
This is a simple real-time chat application built with Go (backend) and HTML/CSS/JavaScript (frontend). It uses WebSockets to enable multiple clients to connect, send messages, and see real-time updates. Features include user names, message timestamps, and a sidebar showing the number of online users and their names.
Features

Real-Time Messaging: Messages are broadcast to all connected clients instantly via WebSockets.
User Names: Each user enters a username upon joining, displayed with their messages.
Timestamps: Messages include a timestamp (HH:MM:SS) for when they were sent.
Online Users: A sidebar displays the number of online users and their usernames, updating dynamically as users join or leave.
Responsive UI: Clean, modern interface with a blue-themed chat window and user-friendly design.

Prerequisites

Go: Version 1.16 or higher (for module support).
Web Browser: Any modern browser (e.g., Chrome, Firefox) for the frontend.
Git (optional): For cloning or managing the repository.

Installation

Clone the Repository (if hosted):
git clone <repository-url>
cd chatapp

Alternatively, create a directory and save the provided main.go and index.html files.

Initialize Go Module:
go mod init chatapp


Install Dependencies:Install the Gorilla WebSocket package:
go get github.com/gorilla/websocket



Project Structure

main.go: The Go backend, handling WebSocket connections and message broadcasting.
index.html: The frontend, containing HTML, CSS, and JavaScript for the chat interface.

Running the Application

Start the Server:
go run main.go

The server runs on http://localhost:8080 by default. To change the port, use the -addr flag:
go run main.go -addr=":9090"


Access the Chat:

Open http://localhost:8080 in multiple browser tabs or windows.
Enter a unique username in each tab to join the chat.
Send messages using the input field or by pressing Enter.



Usage

Joining: On page load, a prompt asks for a username. Enter a name and click "Join Chat" or press Enter.
Chatting: Type messages in the input field and send them via the "Send" button or Enter key.
UI Elements:
Chat Header: Shows "Group Chat" and your current username (e.g., "You are: Alice").
Messages: Displayed with format [username]: [message] [timestamp].
Sidebar: Lists online users and their count, updating in real-time.


Disconnecting: Closing the browser tab or refreshing the page removes the user from the online list.

Example

User Alice opens http://localhost:8080, enters "Alice", and joins.
User Bob opens another tab, enters "Bob", and joins.
Alice sends "Hello, Bob!" at 17:10:23.
Both see: Alice: Hello, Bob! 17:10:23.
Sidebar shows "Online Users (2): Alice, Bob".

Notes

The application is kept simple for clarity and does not include message persistence or advanced authentication.
The backend uses the Gorilla WebSocket library for reliable WebSocket communication.
The frontend is styled with CSS for a clean, responsive design, using a blue and white color scheme.

Troubleshooting

Server Errors: Ensure the Gorilla WebSocket package is installed (go mod tidy).
Connection Issues: Check if http://localhost:8080 is accessible and no other process is using the port.
UI Issues: Verify index.html is in the same directory as main.go.

License
This project is unlicensed and provided as-is for educational purposes. Feel free to modify and distribute as needed.