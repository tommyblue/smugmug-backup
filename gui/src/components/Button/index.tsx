interface IButtonProps {
	text: string
	onClick: React.MouseEventHandler<HTMLButtonElement>
}

const Button = (props: IButtonProps) => {
	return (
		<button
			className="flex w-full justify-center rounded-md bg-indigo-600 px-3 py-1.5 text-sm font-semibold leading-6 text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600"
			onClick={props.onClick}
		>
			{props.text}
		</button>
	)
}

export default Button
