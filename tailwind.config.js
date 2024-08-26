/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./views/**/*.{html,js,go,templ}",
  ],
  theme: {
    extend: {},
  },
  plugins: [require("daisyui")],
}