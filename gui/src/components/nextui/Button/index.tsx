import { Button as NextButton } from "@nextui-org/react"

interface IButtonProps {
	text: string
	onClick: React.MouseEventHandler<HTMLButtonElement>
}

const Button = (props: IButtonProps) => {
	return (
		<NextButton size="lg" color="primary" radius="md" onClick={props.onClick}>
			{props.text}
		</NextButton>
	)
}

export default Button
