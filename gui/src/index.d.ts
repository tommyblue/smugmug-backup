export {}

declare global {
	interface Window {
		electronAPI: {
			openFile: () => Promise<string[]>
			readFile: <T>(path: string) => Promise<T>
		}
		health: {
			check: () => Promise<void>
		}
	}
}
