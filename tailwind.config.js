/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./internal/views/**/*.templ",
    "./internal/views/**/*.go",
    "./templates/**/*.html",
  ],
  theme: {
    extend: {
      colors: {
        terminal: {
          bg: '#0a0a0a',
          text: '#00ff00',
          border: '#00ff0044',
        }
      },
      fontFamily: {
        mono: ['JetBrains Mono', 'Courier New', 'monospace'],
      }
    },
  },
  plugins: [],
}
