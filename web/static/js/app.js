// File Server Frontend Application Scripts

function fileExplorer() {
    return {
        searchQuery: '',
        showUploadModal: false,
        isUploading: false,
        dragOver: false,
        isDraggingOnWindow: false,
        selectedFiles: [],
        dragCounter: 0,

        matchesSearch(name) {
            if (!this.searchQuery) return true;
            return name.toLowerCase().includes(this.searchQuery.toLowerCase());
        },

        formatBytes(bytes) {
            if (bytes === 0) return '0 B';
            const k = 1024;
            const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
            const i = Math.floor(Math.log(bytes) / Math.log(k));
            return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i];
        },

        onWindowDragOver(e) {
            if (e.dataTransfer && e.dataTransfer.types && Array.from(e.dataTransfer.types).includes('Files')) {
                this.isDraggingOnWindow = true;
            }
        },

        onWindowDragLeave(e) {
            this.isDraggingOnWindow = false;
        },

        onWindowDrop(e) {
            this.isDraggingOnWindow = false;
            if (e.dataTransfer && e.dataTransfer.files && e.dataTransfer.files.length > 0) {
                this.selectedFiles = Array.from(e.dataTransfer.files);
                this.showUploadModal = true;
                this.autoSubmitDropFiles(e.dataTransfer.files);
            }
        },

        handleDrop(e) {
            this.dragOver = false;
            if (e.dataTransfer && e.dataTransfer.files && e.dataTransfer.files.length > 0) {
                this.selectedFiles = Array.from(e.dataTransfer.files);
            }
        },

        autoSubmitDropFiles(files) {
            const input = this.$refs.fileInput;
            if (input) {
                const dt = new DataTransfer();
                for (let i = 0; i < files.length; i++) {
                    dt.items.add(files[i]);
                }
                input.files = dt.files;
            }
        }
    };
}

// GUI Launcher Component
function guiLauncher(initialDir) {
    return {
        protocol: 'web',
        targetDir: initialDir || '',
        host: '0.0.0.0',
        port: 8080,
        useTLS: false,
        tlsMode: 'self-signed',
        tlsCert: '',
        tlsKey: '',
        user: 'admin',
        pass: '',
        readOnly: false,
        passivePorts: '',
        authKey: '',
        hostKey: '',

        isRunning: false,
        isLoading: false,
        statusMessage: '',
        errorMessage: '',
        searchQuery: '',

        interfaces: [],
        logs: [],
        copiedId: null,
        showQRModal: false,
        qrURL: '',
        eventSource: null,

        init() {
            this.fetchStatus();
            this.fetchInterfaces();
            this.initLogStream();
        },

        onProtocolChange(newProto) {
            this.protocol = newProto;
            if (!this.isRunning) {
                if (newProto === 'ftp') {
                    this.port = 2121;
                } else if (newProto === 'sftp') {
                    this.port = 2222;
                } else {
                    this.port = 8080;
                }
                this.fetchInterfaces();
            }
        },

        generatePassword() {
            const chars = 'ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnpqrstuvwxyz23456789!@#$%&*';
            let result = '';
            const array = new Uint8Array(12);
            window.crypto.getRandomValues(array);
            for (let i = 0; i < 12; i++) {
                result += chars[array[i] % chars.length];
            }
            this.pass = result;
        },

        async pickFolder() {
            try {
                const res = await fetch('/api/picker/folder', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ currentDir: this.targetDir })
                });
                const data = await res.json();
                if (data.success && data.path) {
                    this.targetDir = data.path;
                }
            } catch (err) {
                console.error('Erro ao abrir diálogo de pastas:', err);
            }
        },

        async fetchStatus() {
            try {
                const res = await fetch('/api/status');
                const data = await res.json();
                this.isRunning = data.isRunning;
                if (data.isRunning) {
                    this.protocol = data.protocol;
                    this.host = data.host;
                    this.port = data.port;
                    this.targetDir = data.targetDir;
                    this.useTLS = data.useTLS;
                    this.interfaces = data.interfaces || [];
                }
            } catch (err) {
                console.error('Erro ao buscar status:', err);
            }
        },

        async fetchInterfaces() {
            try {
                const res = await fetch(`/api/interfaces?port=${this.port}&protocol=${this.protocol}&tls=${this.useTLS}`);
                const data = await res.json();
                if (Array.isArray(data)) {
                    this.interfaces = data;
                }
            } catch (err) {
                console.error('Erro ao buscar interfaces:', err);
            }
        },

        get filteredInterfaces() {
            if (!this.searchQuery.trim()) {
                return this.interfaces;
            }
            const q = this.searchQuery.toLowerCase();
            return this.interfaces.filter(iface => 
                (iface.name && iface.name.toLowerCase().includes(q)) ||
                (iface.ip && iface.ip.toLowerCase().includes(q)) ||
                (iface.typeLabel && iface.typeLabel.toLowerCase().includes(q)) ||
                (iface.url && iface.url.toLowerCase().includes(q))
            );
        },

        async toggleServer() {
            this.isLoading = true;
            this.errorMessage = '';

            try {
                if (this.isRunning) {
                    const res = await fetch('/api/server/stop', { method: 'POST' });
                    const data = await res.json();
                    if (res.ok) {
                        this.isRunning = false;
                        this.fetchInterfaces();
                    } else {
                        this.errorMessage = data.error || 'Falha ao parar servidor';
                    }
                } else {
                    const payload = {
                        protocol: this.protocol,
                        host: this.host,
                        port: parseInt(this.port, 10),
                        targetDir: this.targetDir,
                        useTLS: this.useTLS,
                        tlsCert: this.tlsMode === 'custom' ? this.tlsCert : '',
                        tlsKey: this.tlsMode === 'custom' ? this.tlsKey : '',
                        user: this.user,
                        pass: this.pass,
                        readOnly: this.readOnly,
                        passivePorts: this.passivePorts,
                        authKey: this.authKey,
                        hostKey: this.hostKey
                    };

                    const res = await fetch('/api/server/start', {
                        method: 'POST',
                        headers: { 'Content-Type': 'application/json' },
                        body: JSON.stringify(payload)
                    });
                    const data = await res.json();
                    if (res.ok) {
                        this.isRunning = true;
                        this.interfaces = data.interfaces || [];
                    } else {
                        this.errorMessage = data.error || 'Falha ao iniciar servidor';
                    }
                }
            } catch (err) {
                this.errorMessage = err.message || 'Erro de comunicação';
            } finally {
                this.isLoading = false;
            }
        },

        async openBrowser() {
            try {
                await fetch('/api/app/open-browser', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ url: this.interfaces[0]?.url || `http://127.0.0.1:${this.port}` })
                });
            } catch (err) {
                console.error('Erro ao abrir navegador:', err);
            }
        },

        async copyText(text, id) {
            try {
                await navigator.clipboard.writeText(text);
                this.copiedId = id;
                setTimeout(() => {
                    if (this.copiedId === id) this.copiedId = null;
                }, 2000);
            } catch (err) {
                console.error('Erro ao copiar:', err);
            }
        },

        async copyAllIPs() {
            const list = this.interfaces.map(i => `${i.name} (${i.typeLabel}): ${i.url}`).join('\n');
            await this.copyText(list, 'all-ips');
        },

        async shareLink() {
            try {
                const res = await fetch('/api/share-message');
                const data = await res.json();
                if (data.message) {
                    await this.copyText(data.message, 'share-msg');
                }
            } catch (err) {
                console.error('Erro ao gerar mensagem de compartilhamento:', err);
            }
        },

        openQR(url) {
            this.qrURL = url;
            this.showQRModal = true;
            this.$nextTick(() => {
                this.renderQRCode(url);
            });
        },

        renderQRCode(text) {
            const container = document.getElementById('qrcode-canvas');
            if (!container) return;
            container.innerHTML = '';
            
            // Gerador visual de QR Code SVG usando serviço rápido ou representação vetorial
            const size = 200;
            const encoded = encodeURIComponent(text);
            const qrImg = document.createElement('img');
            qrImg.src = `https://api.qrserver.com/v1/create-qr-code/?size=${size}x${size}&data=${encoded}&margin=2&format=svg`;
            qrImg.alt = 'QR Code';
            qrImg.className = 'w-48 h-48 rounded-lg shadow-md bg-white p-2 border border-slate-700';
            qrImg.onerror = () => {
                // Fallback caso offline
                container.innerHTML = `<div class="p-4 text-center text-xs text-slate-400 border border-slate-700 rounded-lg">Acesso direto:<br><span class="font-mono text-indigo-400 select-all">${text}</span></div>`;
            };
            container.appendChild(qrImg);
        },

        initLogStream() {
            if (window.EventSource) {
                this.eventSource = new EventSource('/api/logs/stream');
                this.eventSource.onmessage = (e) => {
                    try {
                        const parsed = JSON.parse(e.data);
                        if (parsed.message) {
                            this.logs.push(parsed.message);
                            if (this.logs.length > 500) {
                                this.logs.shift();
                            }
                            this.$nextTick(() => {
                                const terminal = this.$refs.logConsole;
                                if (terminal) {
                                    terminal.scrollTop = terminal.scrollHeight;
                                }
                            });
                        }
                    } catch (err) {
                        console.error('Erro ao processar log SSE:', err);
                    }
                };
            }
        },

        clearLogs() {
            this.logs = [];
        },

        async closeApp() {
            try {
                await fetch('/api/app/close', { method: 'POST' });
            } catch (err) {}
            window.close();
        }
    };
}
