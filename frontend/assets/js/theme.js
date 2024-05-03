//accepts dark or light
function setTheme(theme) {
    document.documentElement.setAttribute('data-bs-theme', theme);
    localStorage.setItem('theme', theme);
}

function getCurrentTheme() {
    return document.documentElement.getAttribute('data-bs-theme');
}

function setThemeFromLocalStorage() {
    const theme = localStorage.getItem('theme');
    if (theme) {
        setTheme(theme);
        return true;
    }
}

function setThemeFromMediaQuery() {
    if (window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches) {
        setTheme('dark');
    } else {
        setTheme('light');
    }
}

function toggleTheme() {
    const theme = getCurrentTheme();
    theme === 'dark' ? setTheme('light') : setTheme('dark');
}

if (!setThemeFromLocalStorage()) {
   setThemeFromMediaQuery()
}