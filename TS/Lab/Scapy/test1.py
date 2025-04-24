from scapy.all import IP, ICMP, send

target_ip = "192.168.164.12"  # Wpisz tutaj adres IP celu, który chcesz pingować
spoofed_ip = "192.168.164.103" # Wpisz tutaj adres IP, który ma być widoczny jako źródło
#print(getmacbyip("192.168.7.12"))  # pobiera MAC i zapisuje do tablicy ARP


packet = IP(src=spoofed_ip, dst=target_ip)/ICMP()/"Hej!" # Dodajemy ładunek do pakietu ICMP

print("Przygotowany pakiet:")
packet.show()
#48 40 238 37
# Scapy samo zajmie się routingiem i warstwą 2 (np. Ethernet).
print(f"\nWysyłanie pinga z {spoofed_ip} do {target_ip}...")
i = 0
while True:
    send(packet, verbose=0) # verbose=0 wyłącza domyślny komunikat "Sent 1 packets."
    print(f"Pakiet wysłany. {i}")
    i = i + 1
    #sudo -E ../../../../ts/bin/python3 test1.py 