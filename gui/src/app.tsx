import { NextUIProvider } from "@nextui-org/react"
import React from "react"
import { createRoot } from "react-dom/client"
import Button from "./components/nextui/Button"
import { Config, defaultConfig } from "./config"

type Health = {
	check: () => Promise<void>
}
const App = () => {
	const [filePath, setFilePath] = React.useState<string>("")
	const [config, setConfig] = React.useState<Config | null>(null)

	React.useEffect(() => {
		if (config === null) {
			console.log("setting default config")
			setConfig(defaultConfig)
			localStorage.setItem("config", JSON.stringify(defaultConfig))
		}
	}, [config])

	React.useEffect(() => {
		const cfg = localStorage.getItem("config")
		if (cfg) {
			console.log("config from local storage:", JSON.parse(cfg))
			setConfig(JSON.parse(cfg))
		}

		setInterval(() => {
			window.health.check()
		}, 5000)
	}, [])

	React.useEffect(() => {
		if (filePath) {
			const readFile = async () => {
				const content = await window.electronAPI.readFile<Partial<Config>>(filePath)
				const newConfig = { ...config, ...content }
				setConfig(newConfig)
			}
			readFile()
		}
	}, [filePath])

	const selectFile = async () => {
		const filePaths = await window.electronAPI.openFile()
		setFilePath(filePaths[0])
	}

	return (
		<div className="container mx-auto columns-2 font-display">
			<h1 className="text-3xl font-bold underline">Hello world!</h1>
			<Button text="Click me!" onClick={selectFile} />
			{filePath && <p>{filePath}</p>}
		</div>
	)
}

const root = createRoot(document.getElementById("root"))
root.render(
	<React.StrictMode>
		<NextUIProvider>
			<main className="dark text-foreground bg-background">
				<App />
			</main>
		</NextUIProvider>
	</React.StrictMode>
)
