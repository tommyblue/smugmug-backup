type Auth = {
	api_key: string
	api_secret: string
	user_token: string
	user_secret: string
}

type Store = {
	destination: string
	file_names: string
	use_metadata_times: boolean
	force_metadata_times: boolean
	write_csv: boolean
	force_video_download: boolean
	concurrent_albums: number
	concurrent_downloads: number
}

export type Config = {
	auth: Auth
	store: Store
}

export const defaultConfig: Config = {
	auth: {
		api_key: "",
		api_secret: "",
		user_token: "",
		user_secret: "",
	},
	store: {
		destination: "",
		file_names: "",
		use_metadata_times: true,
		force_metadata_times: true,
		write_csv: true,
		force_video_download: true,
		concurrent_albums: 5,
		concurrent_downloads: 10,
	},
}
