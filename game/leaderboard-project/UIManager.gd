extends Node

@onready var websocket_manager = get_node("/root/Node3D/WebSocketManager") 
# Called when the node enters the scene tree for the first time.
func _ready() -> void:
	$ConnectButton.pressed.connect(_on_connect_button_pressed)
	$SendButton.pressed.connect(_on_send_button_pressed)


func _on_connect_button_pressed():
	print("THIS WAS CALLEd")
	websocket_manager.connect_to_ws()

func _on_send_button_pressed():
	websocket_manager.send_message("Hello, Server!")

# Called every frame. 'delta' is the elapsed time since the previous frame.
func _process(delta: float) -> void:
	pass
