require "socket"
u2 = UDPSocket.new
u2.connect("docker", 1700)
u2.send "uuuu", 0
