from struct import pack, unpack
import socket
sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
whoami = 'william'.encode('utf-8')
payload = pack("BBB{}s".format(len(whoami)), 1,1, len(whoami), whoami)
sock.sendto(payload, 0, ('127.0.0.1', 20111))
