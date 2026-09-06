# We need to determine what the phone and laptop can already communicate over.
Can the laptop reach the phone itself through the phone's hotspot?

--> ipconfig

Windows IP Configuration

Ethernet adapter Ethernet:

   Media State . . . . . . . . . . . : Media disconnected
   Connection-specific DNS Suffix  . : 

Ethernet adapter Ethernet 3:

   Connection-specific DNS Suffix  . : 
   Link-local IPv6 Address . . . . . : fe80::2c57:498:6d2b:4b2b%43
   IPv4 Address. . . . . . . . . . . : 10.145.201.101
   Subnet Mask . . . . . . . . . . . : 255.255.255.0
   Default Gateway . . . . . . . . . : 10.145.201.7

Wireless LAN adapter Local Area Connection* 1:

   Media State . . . . . . . . . . . : Media disconnected
   Connection-specific DNS Suffix  . : 

Wireless LAN adapter Local Area Connection* 2:

   Media State . . . . . . . . . . . : Media disconnected
   Connection-specific DNS Suffix  . : 

Wireless LAN adapter Wi-Fi:

   Connection-specific DNS Suffix  . : 
   IPv4 Address. . . . . . . . . . . : 10.235.157.157
   Subnet Mask . . . . . . . . . . . : 255.255.255.0
   Default Gateway . . . . . . . . . : 10.235.157.208

# Testing wifi getway

--> ping 10.235.157.208

Pinging 10.235.157.208 with 32 bytes of data:
Reply from 10.235.157.208: bytes=32 time=5ms TTL=64
Reply from 10.235.157.208: bytes=32 time=2ms TTL=64
Reply from 10.235.157.208: bytes=32 time=4ms TTL=64
Reply from 10.235.157.208: bytes=32 time=5ms TTL=64

Ping statistics for 10.235.157.208:
    Packets: Sent = 4, Received = 4, Lost = 0 (0% loss),
Approximate round trip times in milli-seconds:
    Minimum = 2ms, Maximum = 5ms, Average = 4ms

==> Conclusion: laptop can communicate with 10.235.157.208 over wi-fi

# Identifying the phone

--> arp -a 
looks for the IP address assigned to the phone itself

Interface: 10.235.157.157 --- 0xa
  Internet Address      Physical Address      Type
  10.235.157.208        5e-d1-97-60-20-46     dynamic   
  10.235.157.255        ff-ff-ff-ff-ff-ff     static    
  224.0.0.22            01-00-5e-00-00-16     static    
  224.0.0.251           01-00-5e-00-00-fb     static    
  224.0.0.252           01-00-5e-00-00-fc     static    
  239.255.255.250       01-00-5e-7f-ff-fa     static    
  255.255.255.255       ff-ff-ff-ff-ff-ff     static    

Interface: 10.145.201.101 --- 0x2b
  Internet Address      Physical Address      Type
  10.145.201.7          56-39-50-3c-05-f0     dynamic   
  10.145.201.255        ff-ff-ff-ff-ff-ff     static    
  224.0.0.22            01-00-5e-00-00-16     static    
  224.0.0.251           01-00-5e-00-00-fb     static    
  224.0.0.252           01-00-5e-00-00-fc     static    
  239.255.255.250       01-00-5e-7f-ff-fa     static    
  255.255.255.255       ff-ff-ff-ff-ff-ff     static 

==> Conclusion: laptop can ping 10.235.157.208, ARP does not prove that 10.235.157.208 is the phone itself

# 