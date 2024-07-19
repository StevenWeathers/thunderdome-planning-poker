module.exports = {
  darkMode: 'class',
  content: [
    './src/**/*.{svelte,ts,js}',
    './public/index.html'
  ],
  theme: {
    extend: {
      colors: {
        'yellow-thunder': '#ffdd57',
      },
      fontFamily: {
        rajdhani: ['Rajdhani', 'Arial Narrow', 'sans-serif'],
      }
    }
  }
}