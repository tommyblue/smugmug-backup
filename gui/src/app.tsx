import { createRoot } from "react-dom/client"
import Button from "./components/Button"

const App = () => {
	return (
		<div className="container mx-auto columns-2 font-display">
			<h1 className="text-3xl font-bold underline">Hello world!</h1>
			<Button text="Click me!" onClick={() => alert("Hello world!")} />
		</div>
	)
}

const root = createRoot(document.getElementById("root"))
root.render(<App />)
