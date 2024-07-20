const { nextui } = require("@nextui-org/react")

/** @type {import('tailwindcss').Config} */
// const defaultTheme = require("tailwindcss/defaultTheme")
module.exports = {
	content: ["./index.html", "./src/**/*.{js,ts,jsx,tsx}", "./node_modules/@nextui-org/theme/dist/**/*.{js,ts,jsx,tsx}"],
	theme: {
		extend: {},
	},
	darkMode: "class",
	plugins: [nextui()],
}
