/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        brand: {
          DEFAULT: '#9400FF',
          dark: '#6F00C2',
          light: '#E6D1FF',
        },
        background: '#ffffff',
        foreground: '#111111',
        card: '#ffffff',
        'card-foreground': '#111111',
        popover: '#ffffff',
        'popover-foreground': '#111111',
        primary: '#9400FF',
        'primary-foreground': '#ffffff',
        secondary: '#f5f5f5',
        'secondary-foreground': '#111111',
        muted: '#f3f4f6',
        'muted-foreground': '#6b7280',
        accent: '#f3f4f6',
        'accent-foreground': '#111111',
        destructive: '#6F00C2',
        'destructive-foreground': '#ffffff',
        border: '#e5e7eb',
        input: '#e5e7eb',
        ring: '#9400FF',
      },
      fontFamily: {
        sans: ['Inter', 'system-ui', 'sans-serif'],
        body: ['Inter', 'system-ui', 'sans-serif'],
        heading: ['"Instrument Serif"', 'serif'],
      }
    },
  },
  plugins: [
    require('@tailwindcss/typography'),
  ],
}
