// See the Electron documentation for details on how to use preload scripts:
// https://www.electronjs.org/docs/latest/tutorial/process-model#preload-scripts
const { contextBridge, ipcRenderer } = require("electron")

contextBridge.exposeInMainWorld("electronAPI", {
	openFile: () => ipcRenderer.invoke("dialog:openFile"),
	readFile: async (path: string) => ipcRenderer.invoke("dialog:readFile", path),
})

contextBridge.exposeInMainWorld("health", {
	check: () => ipcRenderer.invoke("health:check"),
})
