package gui

import (
	"fmt"
	"net"
	"sort"
	"strings"
)

// CategorizeInterface classifica o tipo de interface a partir do seu nome e endereço IP.
func CategorizeInterface(name string, ip net.IP) (InterfaceType, string) {
	if ip != nil && ip.IsLoopback() {
		return TypeLoopback, "Loopback"
	}

	lower := strings.ToLower(name)

	// VPNs
	if strings.HasPrefix(lower, "tun") ||
		strings.HasPrefix(lower, "tap") ||
		strings.HasPrefix(lower, "tailscale") ||
		strings.HasPrefix(lower, "wg") ||
		strings.HasPrefix(lower, "wireguard") ||
		strings.HasPrefix(lower, "nordlynx") ||
		strings.HasPrefix(lower, "proton") ||
		strings.HasPrefix(lower, "utun") ||
		strings.HasPrefix(lower, "ppp") ||
		strings.Contains(lower, "vpn") {
		return TypeVPN, "VPN"
	}

	// Docker e Bridges Virtuais
	if strings.HasPrefix(lower, "docker") ||
		strings.HasPrefix(lower, "br-") ||
		strings.HasPrefix(lower, "veth") ||
		strings.HasPrefix(lower, "cni") ||
		strings.HasPrefix(lower, "flannel") ||
		strings.HasPrefix(lower, "virbr") ||
		strings.HasPrefix(lower, "vboxnet") ||
		strings.HasPrefix(lower, "vmnet") ||
		strings.HasPrefix(lower, "hyperv") ||
		strings.HasPrefix(lower, "wsl") ||
		strings.HasPrefix(lower, "bridge") {
		return TypeDocker, "Docker / Bridge"
	}

	// Wi-Fi
	if strings.HasPrefix(lower, "wl") ||
		strings.HasPrefix(lower, "wlan") ||
		strings.HasPrefix(lower, "wifi") ||
		strings.HasPrefix(lower, "airport") ||
		strings.Contains(lower, "wireless") ||
		strings.Contains(lower, "wi-fi") {
		return TypeWiFi, "Wi-Fi"
	}

	// Ethernet
	if strings.HasPrefix(lower, "eth") ||
		strings.HasPrefix(lower, "en") ||
		strings.HasPrefix(lower, "eno") ||
		strings.HasPrefix(lower, "ens") ||
		strings.HasPrefix(lower, "enp") ||
		strings.Contains(lower, "ethernet") {
		return TypeEthernet, "Ethernet"
	}

	return TypeOther, "Rede Local"
}

// BuildAccessURL cria a URL de acesso formatada baseada no protocolo, IP/Host e TLS.
func BuildAccessURL(hostOrIP string, port int, protocol string, isTLS bool) string {
	scheme := "http"
	switch strings.ToLower(protocol) {
	case "ftp":
		scheme = "ftp"
		if isTLS {
			scheme = "ftps"
		}
	case "sftp":
		scheme = "sftp"
	default: // web
		if isTLS {
			scheme = "https"
		}
	}

	// Formata IPv6 com colchetes
	if ip := net.ParseIP(hostOrIP); ip != nil && ip.To4() == nil {
		return fmt.Sprintf("%s://[%s]:%d", scheme, hostOrIP, port)
	}

	return fmt.Sprintf("%s://%s:%d", scheme, hostOrIP, port)
}

// BuildAccessURLs retorna o acesso local e a lista de acessos LAN detectados.
func BuildAccessURLs(host string, port int, protocol string, isTLS bool) (localURL string, lanURLs []string) {
	if host == "0.0.0.0" || host == "" || host == "::" {
		localURL = BuildAccessURL("127.0.0.1", port, protocol, isTLS)
		seen := make(map[string]bool)
		seen[localURL] = true

		ifaces := DetectNetworkInterfaces(port, protocol, isTLS)
		for _, iface := range ifaces {
			if iface.IsLoopback {
				continue
			}
			if !seen[iface.URL] {
				seen[iface.URL] = true
				lanURLs = append(lanURLs, iface.URL)
			}
		}
	} else {
		localURL = BuildAccessURL(host, port, protocol, isTLS)
	}
	return localURL, lanURLs
}

// DetectNetworkInterfaces descobre e categoriza todos os adaptadores de rede ativos no sistema operacional.
func DetectNetworkInterfaces(port int, protocol string, isTLS bool) []NetworkInterface {
	var result []NetworkInterface
	seenIPs := make(map[string]bool)

	// Adiciona sempre Localhost Loopback como primeira interface
	loopbackURL := BuildAccessURL("127.0.0.1", port, protocol, isTLS)
	result = append(result, NetworkInterface{
		Name:          "lo (Localhost)",
		IP:            "127.0.0.1",
		Type:          TypeLoopback,
		TypeLabel:     "Loopback",
		IsRecommended: true,
		IsLoopback:    true,
		URL:           loopbackURL,
	})
	seenIPs["127.0.0.1"] = true

	ifaces, err := net.Interfaces()
	if err != nil {
		return result
	}

	var physicalIfaces []NetworkInterface
	var virtualIfaces []NetworkInterface

	for _, iface := range ifaces {
		// Ignora interfaces desativadas ou loopback
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if ip == nil || ip.IsLoopback() {
				continue
			}

			ipStr := ip.String()
			if seenIPs[ipStr] {
				continue
			}
			seenIPs[ipStr] = true

			ifaceType, label := CategorizeInterface(iface.Name, ip)
			url := BuildAccessURL(ipStr, port, protocol, isTLS)

			entry := NetworkInterface{
				Name:          iface.Name,
				IP:            ipStr,
				Type:          ifaceType,
				TypeLabel:     label,
				IsRecommended: false,
				IsLoopback:    false,
				URL:           url,
			}

			// Prioriza interfaces físicas (Wi-Fi e Ethernet) sobre virtuais/bridges/docker
			if ifaceType == TypeWiFi || ifaceType == TypeEthernet {
				physicalIfaces = append(physicalIfaces, entry)
			} else {
				virtualIfaces = append(virtualIfaces, entry)
			}
		}
	}

	// Ordena interfaces físicas IPv4 no topo
	sort.Slice(physicalIfaces, func(i, j int) bool {
		ipI := net.ParseIP(physicalIfaces[i].IP)
		ipJ := net.ParseIP(physicalIfaces[j].IP)
		if ipI != nil && ipJ != nil {
			if ipI.To4() != nil && ipJ.To4() == nil {
				return true
			}
			if ipI.To4() == nil && ipJ.To4() != nil {
				return false
			}
		}
		return physicalIfaces[i].Type == TypeWiFi && physicalIfaces[j].Type != TypeWiFi
	})

	// Marca a primeira interface física primária como Recomendada
	if len(physicalIfaces) > 0 {
		physicalIfaces[0].IsRecommended = true
	}

	result = append(result, physicalIfaces...)
	result = append(result, virtualIfaces...)

	return result
}

// FormatShareMessage gera uma mensagem completa e formatada com os links de acesso pronta para envio.
func FormatShareMessage(cfg ServerConfig, ifaces []NetworkInterface) string {
	protoName := "Web (HTTP)"
	if cfg.UseTLS {
		protoName = "Web Segura (HTTPS)"
	}
	if cfg.Protocol == ProtocolFTP {
		protoName = "FTP"
		if cfg.UseTLS {
			protoName = "FTPS (TLS)"
		}
	} else if cfg.Protocol == ProtocolSFTP {
		protoName = "SFTP (SSH)"
	}

	var sb strings.Builder
	sb.WriteString("📁 *File Server - Compartilhamento de Arquivos*\n")
	sb.WriteString(fmt.Sprintf("📂 *Diretório:* `%s`\n", cfg.TargetDir))
	sb.WriteString(fmt.Sprintf("🔒 *Protocolo:* %s\n", protoName))
	sb.WriteString("🌐 *Links de Acesso:*\n")

	for _, iface := range ifaces {
		recBadge := ""
		if iface.IsRecommended {
			recBadge = " ⭐ (Recomendado)"
		}
		sb.WriteString(fmt.Sprintf("  • *%s* [%s]%s: %s\n", iface.Name, iface.TypeLabel, recBadge, iface.URL))
	}

	if (cfg.Protocol == ProtocolFTP || cfg.Protocol == ProtocolSFTP) && cfg.User != "" {
		sb.WriteString("\n🔑 *Credenciais de Acesso:*\n")
		sb.WriteString(fmt.Sprintf("  • Usuário: `%s`\n", cfg.User))
		if cfg.Pass != "" {
			sb.WriteString(fmt.Sprintf("  • Senha: `%s`\n", cfg.Pass))
		}
		if cfg.ReadOnly {
			sb.WriteString("  • Modo: Somente Leitura (Download)\n")
		}
	}

	return sb.String()
}
