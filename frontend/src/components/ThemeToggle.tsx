import { IconMoonStars, IconSunFilled } from "@tabler/icons-react";
import { useState } from "react";

export default function ThemeToggle() {
  const [darkMode, setDarkMode] = useState(true);

  const toggleTheme = () => {
    setDarkMode(!darkMode);
    document.body.classList.toggle("dark");
  };

  return (
    <button
      onClick={toggleTheme}
      className="w-12 h-12 rounded-full flex items-center justify-center active:scale-90 transition-transform duration-150"
    >
      <div className="transition-all duration-200">
        {darkMode ? <IconSunFilled /> : <IconMoonStars />}
      </div>
    </button>
  );
}
