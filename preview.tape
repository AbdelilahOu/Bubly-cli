# Bubly CLI App Demo
# Demonstrates the video, audio, and subtitle download flows.
#
# The app's palette uses dark text on a light surface, so a light VHS theme
# is required or the menus/URL text render invisibly dark-on-dark.

Output preview.gif

Require build/bubly

# --- Recording settings ------------------------------------------------------
Set FontSize 16
Set Width 1100
Set Height 760
Set Padding 40
Set Margin 20
Set MarginFill "#1f2430"
Set BorderRadius 10
Set WindowBar "Colorful"
Set WindowBarSize 40
Set TypingSpeed 90ms

# Light theme so the app's dark text and gray menu items stay visible.
Set Theme { "name": "Bubly Light", "background": "#ffffff", "foreground": "#1f2937", "selection": "#d1d5db", "cursor": "#7d56f4", "black": "#1f2937", "red": "#b91c1c", "green": "#16a34a", "yellow": "#f97316", "blue": "#2563eb", "purple": "#7d56f4", "cyan": "#0891b2", "white": "#e5e7eb", "brightBlack": "#374151", "brightRed": "#dc2626", "brightGreen": "#22c55e", "brightYellow": "#fb923c", "brightBlue": "#3b82f6", "brightPurple": "#8b5cf6", "brightCyan": "#06b6d4", "brightWhite": "#f9fafb" }

# --- Launch ------------------------------------------------------------------
Type "./build/bubly"
Enter
Sleep 3s

# --- Flow 1: Download video --------------------------------------------------
# "Download Youtube video" is the first option; select it.
Enter
Sleep 1500ms

# A tiny, stable, always-available video keeps format-fetch + download fast.
Type "https://www.youtube.com/watch?v=jNQXAC9IVRw"
Enter

# Wait for yt-dlp to list formats, then browse the list.
Sleep 6s
Down
Sleep 700ms
Down
Sleep 700ms
Up
Sleep 700ms

# Download the selected format (yt-dlp adds a throttle sleep before fetching).
Enter
Sleep 12s

# Back to the main menu.
Backspace
Sleep 1s

# --- Flow 2: Download audio --------------------------------------------------
Down
Enter
Sleep 1500ms

Type "https://www.youtube.com/watch?v=jNQXAC9IVRw"
Enter
Sleep 6s
Down
Sleep 700ms
Enter
Sleep 12s

Backspace
Sleep 1s

# --- Flow 3: Download subtitles ----------------------------------------------
Down
Down
Enter
Sleep 1500ms

Type "https://www.youtube.com/watch?v=jNQXAC9IVRw"
Enter
Sleep 6s
Down
Sleep 700ms
Enter
Sleep 10s

# --- Quit --------------------------------------------------------------------
Escape
Sleep 2s
