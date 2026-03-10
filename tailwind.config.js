/** @type {import('tailwindcss').Config} */
module.exports = {
	content: [
		"./internal/**/*.templ",
		"./**/*.go",
	],
	theme: {
		extend: {
			// <div class="bg-background text-foreground border-red hover:border-red-dim">
			colors: {
				background: '#282828',
				foreground: '#ebdbb2',
				primary: '#fabd2f',
				'primary-dim': '#d79921',
				secondary: '#8ec07c',
				'secondary-dim': '#6b8e4a',
				muted: '#928374',
				'muted-dim': '#7c6f64',
				red: '#cc241d',
				'red-dim': '#9d0006',
				green: '#b8bb26',
				'green-dim': '#79740e',
				yellow: '#fabd2f',
				'yellow-dim': '#b57614',
				blue: '#458588',
				'blue-dim': '#076678',
				purple: '#b16286',
				'purple-dim': '#8e407d',
				aqua: '#689d6a',
				'aqua-dim': '#427b58',
				gray: '#928374',
				'gray-dim': '#7c6f64',
			},
		},
	},
	plugins: [],
	darkMode: 'media',
}
