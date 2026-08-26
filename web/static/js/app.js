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
