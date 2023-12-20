/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    './components/**/*.templ',
  ],
  theme: {
    extend: {
      colors: {
        dark: '#13141B',
        yellow: '#D8CF00',
        light: '#F8F8F2',
        blue: '#47DEFF',
        gray: '#6272A4',
        pink: '#FF79C6',
        green: '#2AFF98',
      },
    },
  },
  plugins: [],
};

