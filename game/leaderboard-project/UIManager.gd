extends Node

@onready var websocket_manager = get_node("/root/Node3D/WebSocketManager") 
# Called when the node enters the scene tree for the first time.
func _ready() -> void:
	$ConnectButton.pressed.connect(_on_connect_button_pressed)
	$SendButton.pressed.connect(_on_send_button_pressed)
	$Leaderboard.pressed.connect(_on_leaderboard_button_pressed)

func _on_connect_button_pressed():
	print("THIS WAS CALLEd")
	websocket_manager.connect_to_ws()

func _on_send_button_pressed():
	var message_dict = {"command": "IncrementKillCount"}

	var json_string = JSON.stringify(message_dict)
	websocket_manager.send_message(json_string)

func _on_leaderboard_button_pressed():
	var message_dict = {"command": "GetLeaderboard"}
	var json_string = JSON.stringify(message_dict)

	websocket_manager.send_message(json_string)

# Called every frame. 'delta' is the elapsed time since the previous frame.
func _process(delta: float) -> void:
	pass
