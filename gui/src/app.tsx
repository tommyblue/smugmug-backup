import { NextUIProvider } from "@nextui-org/react"
import React from "react"
import { createRoot } from "react-dom/client"
import Button from "./components/Button"

interface Config {
	[key: string]: string
}

const App = () => {
	const [filePath, setFilePath] = React.useState("")
	const [config, setConfig] = React.useState({})

	React.useEffect(() => {
		setInterval(() => {
			window.health.check()
		}, 5000)
	}, [])

	React.useEffect(() => {
		if (filePath) {
			const readFile = async () => {
				const content: Config = await window.electronAPI.readFile(filePath)
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
	<NextUIProvider>
		<App />
	</NextUIProvider>
)
