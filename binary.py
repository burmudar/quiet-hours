from struct import pack, unpack
import socket
sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
whoami = 'william'.encode('utf-8')
payload = pack("BBB{}s".format(len(whoami)), 1,1, len(whoami), whoami)
print(f"Sending\n{payload}")
sock.sendto(payload, 0, ('127.0.0.1', 20111))
reply = sock.recv(512)
version, packet_type, is_quiet, wake_up, s = unpack("5B", reply[:5])
whoru = unpack(f"{s}s", reply[5:])[0].decode('utf-8')
print(f'''Received
      Version: {version}
      PackeType: {packet_type}
      Is Quiet Time: {is_quiet == 1}
      Wake up in: {wake_up} hours
      Who R U: {whoru}
      ''')
