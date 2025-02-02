extends Node

@export var websocket_url = "ws://localhost:8080/ws"

var socket = WebSocketPeer.new()
var is_connecting = false

# Called when the node enters the scene tree for the first time.
func connect_to_ws() -> void:
	print("Trying to Connect...")
	var err = socket.connect_to_url(websocket_url)

	if err != OK:
		print("Unable to connect", err)
		set_process(false)
	else:
		# Wait for the socket to connect.
		is_connecting = true

		# Send data.
		print("ready state: ",  socket.get_ready_state())
		#socket.send_text("Test packet")
		print("Sent test packet")

func send_message(message: String):
	if socket.get_ready_state() == WebSocketPeer.STATE_OPEN:
		socket.send_text(message)
		print("Message sent:", message)
	else:
		print("WebSocket is not connected")

# Called every frame. 'delta' is the elapsed time since the previous frame.
func _process(delta: float) -> void:
	if is_connecting:
		socket.poll()
		var state = socket.get_ready_state()

		if state == WebSocketPeer.STATE_OPEN:
			print("Connected to WebSocket!")
			socket.send_text("Test packet")
			print("Sent test packet")
			is_connecting = false  # Stop trying to connect
		elif state == WebSocketPeer.STATE_CLOSING:
			print("WebSocket is closing...")
		elif state == WebSocketPeer.STATE_CLOSED:
			print("WebSocket connection closed.")
			set_process(false)  # Stop processing after disconnect
