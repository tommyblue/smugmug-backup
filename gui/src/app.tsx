import React from "react"
import { createRoot } from "react-dom/client"
import Button from "./components/Button"

const App = () => {
	React.useEffect(() => {
		setInterval(() => {
			window.health.check()
		}, 5000)
	}, [])

	const selectFile = async () => {
		const filePaths = await window.electronAPI.openFile()
		console.log(filePaths)
	}

	return (
		<div className="container mx-auto columns-2 font-display">
			<h1 className="text-3xl font-bold underline">Hello world!</h1>
			<Button text="Click me!" onClick={selectFile} />
		</div>
	)
}

const root = createRoot(document.getElementById("root"))
root.render(<App />)
