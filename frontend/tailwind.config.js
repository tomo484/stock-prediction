/** @type {import('tailwindcss').Config} */
module.exports = {
  darkMode: 'class',
  content: [
    './src/pages/**/*.{js,ts,jsx,tsx,mdx}',
    './src/components/**/*.{js,ts,jsx,tsx,mdx}',
    './src/app/**/*.{js,ts,jsx,tsx,mdx}',
  ],
  theme: {
    extend: {
      colors: {
        primary: '#135bec',
        'background-light': '#f6f6f8',
        'background-dark': '#101622',
        'surface-light': '#ffffff',
        'surface-dark': '#1A202C',
        positive: '#48BB78',
        negative: '#E53E3E',
        'text-primary-light': '#1A202C',
        'text-primary-dark': '#E2E8F0',
        'text-secondary-light': '#4A5568',
        'text-secondary-dark': '#A0AEC0',
      },
      fontFamily: {
        display: ['Manrope', 'Noto Sans JP', 'sans-serif'],
      },
      borderRadius: {
        DEFAULT: '0.5rem',
        lg: '0.75rem',
        xl: '1rem',
        full: '9999px',
      },
    },
  },
  plugins: [],
};


