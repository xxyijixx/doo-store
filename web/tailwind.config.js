/** @type {import('tailwindcss').Config} */
module.exports = {
    darkMode: ["class"],
    content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}", // 确保包含所有的文件
],
theme: {
		extend: {
			borderRadius: {
				lg: 'var(--radius)',
				md: 'calc(var(--radius) - 2px)',
				sm: 'calc(var(--radius) - 4px)'
			},
			colors: {
				'theme-color': '#84c56a',
				'theme-success': '#e8f5e2',
				'theme-warning':'#fdf2f2',
				'uninsbg-warn': '#FEE2E2',
				background: 'hsl(var(--background))',
				foreground: 'hsl(var(--foreground))',
				card: {
					DEFAULT: 'hsl(var(--card))',
					foreground: 'hsl(var(--card-foreground))'
				},
				popover: {
					DEFAULT: 'hsl(var(--popover))',
					foreground: 'hsl(var(--popover-foreground))'
				},
				primary: {
					DEFAULT: 'hsl(var(--primary))',
					foreground: 'hsl(var(--primary-foreground))'
				},
				secondary: {
					DEFAULT: 'hsl(var(--secondary))',
					foreground: 'hsl(var(--secondary-foreground))'
				},
				muted: {
					DEFAULT: 'hsl(var(--muted))',
					foreground: 'hsl(var(--muted-foreground))'
				},
				accent: {
					DEFAULT: 'hsl(var(--accent))',
					foreground: 'hsl(var(--accent-foreground))'
				},
				destructive: {
					DEFAULT: 'hsl(var(--destructive))',
					foreground: 'hsl(var(--destructive-foreground))'
				},
				border: 'hsl(var(--border))',
				input: 'hsl(var(--input))',
				ring: 'hsl(var(--ring))',
				chart: {
					'1': 'hsl(var(--chart-1))',
					'2': 'hsl(var(--chart-2))',
					'3': 'hsl(var(--chart-3))',
					'4': 'hsl(var(--chart-4))',
					'5': 'hsl(var(--chart-5))'
				}
			},
			
			transitionProperty: {
				'transform': 'transform',
			},
			transitionDuration: {
				'300': '300ms',
			},
			transitionTimingFunction: {
				'in-out': 'cubic-bezier(0.4, 0, 0.2, 1)',
			},
			rotate: {
				'360': '360deg'
			},
			screens: {
				'sm': '642px',    // 自定义小屏幕断点
				'mdd': '768px',    // 自定义中等屏幕断点
				'md': '890px',    // 自定义中等屏幕断点
				'lg': '1200px',   // 自定义大屏幕断点
				'xl': '1480px',    // 超大屏幕
				'2xl': '1920px',   // 超超大屏幕
			},
			transform: {
				'rotate-90': 'rotate(90deg)',
			},
		}
	},
	plugins: [
		require("tailwindcss-animate"),
		require('@tailwindcss/line-clamp')
	],
}
