import grpc from "@grpc/grpc-js"
import protoLoader from "@grpc/proto-loader"
import { app, BrowserWindow, dialog, ipcMain } from "electron"
import fs from "fs"
import path from "path"

// Handle creating/removing shortcuts on Windows when installing/uninstalling.
if (require("electron-squirrel-startup")) {
	app.quit()
}

let mainWindow: BrowserWindow | null
let serverAddr = "localhost:8089"

let serverProc = require("child_process").spawn("./server")
serverProc.stdout.on("data", (data: any) => {
	const jsonData = JSON.parse(data)
	if (jsonData.listen) {
		console.log(`Server listening on ${jsonData.listen}`)
		serverAddr = jsonData.listen
	}
})

serverProc.stderr.on("data", (data: any) => {
	console.error(`stderr: ${data}`)
})

serverProc.on("close", (code: any) => {
	console.log(`child process exited with code ${code}`)
})
// serverProc.on("exit", (code, sig) => {
// 	// finishing
// 	console.log("serverProc exit", code, sig)
// })
// serverProc.on("error", error => {
// 	// error handling
// 	console.log("serverProc error", error)
// })

const createWindow = () => {
	// Create the browser window.
	mainWindow = new BrowserWindow({
		width: 1600,
		height: 1000,
		webPreferences: {
			preload: path.join(__dirname, "preload.js"),
		},
	})

	// and load the index.html of the app.
	if (MAIN_WINDOW_VITE_DEV_SERVER_URL) {
		mainWindow.loadURL(MAIN_WINDOW_VITE_DEV_SERVER_URL)
	} else {
		mainWindow.loadFile(path.join(__dirname, `../renderer/${MAIN_WINDOW_VITE_NAME}/index.html`))
	}

	// Open the DevTools.
	mainWindow.webContents.openDevTools()
}

// This method will be called when Electron has finished
// initialization and is ready to create browser windows.
// Some APIs can only be used after this event occurs.
app.on("ready", createWindow)

// Quit when all windows are closed, except on macOS. There, it's common
// for applications and their menu bar to stay active until the user quits
// explicitly with Cmd + Q.
app.on("window-all-closed", () => {
	if (process.platform !== "darwin") {
		app.quit()
	}
})

app.on("will-quit", () => {
	serverProc.kill()
})

app.on("activate", () => {
	// On OS X it's common to re-create a window in the app when the
	// dock icon is clicked and there are no other windows open.
	if (BrowserWindow.getAllWindows().length === 0) {
		createWindow()
	}
})

// In this file you can include the rest of your app's specific main process
// code. You can also put them in separate files and import them here.
ipcMain.handle("dialog:openFile", async event => {
	const result = await dialog.showOpenDialog(mainWindow, {
		properties: ["openFile"],
	})
	return result.filePaths
})

ipcMain.handle("dialog:readFile", async (event, filePath) => {
	return fs.readFileSync(filePath, "utf-8")
})

ipcMain.handle("health:check", async event => {
	// TODO: use better way to load proto file
	const PROTO_PATH = path.join(__dirname, "../../src/proto/health.proto")
	const packageDefinition = protoLoader.loadSync(PROTO_PATH, {
		keepCase: true,
		longs: String,
		enums: String,
		defaults: true,
		oneofs: true,
	})

	const healthProto = grpc.loadPackageDefinition(packageDefinition).grpc.health.v1
	// TODO: get the server address from the renderer
	const client = new healthProto.Health(serverAddr, grpc.credentials.createInsecure())
	client.Check({ service: "" }, (err, response) => {
		if (err) {
			console.error("Errore nel controllo dello stato:", err)
			response = "unknown"
		} else {
			console.log("Stato del server gRPC:", response.status)
			response = response.status
		}
	})
})
