# Bubly CLI App Demo
#
# Showcases the interface: the main menu, all three modes, URL entry, and the
# format/language pickers. It deliberately stops BEFORE starting a download so
# the recording never hangs on yt-dlp's throttle sleeps or fails when the host
# (e.g. CI) can't reach YouTube.
#
# The app's palette is dark text on a light surface, so a light VHS theme is
# required or the menus/URL text render invisibly dark-on-dark.

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
Set TypingSpeed 80ms

# Light theme so the app's dark text and gray menu items stay visible.
Set Theme { "name": "Bubly Light", "background": "#ffffff", "foreground": "#1f2937", "selection": "#d1d5db", "cursor": "#7d56f4", "black": "#1f2937", "red": "#b91c1c", "green": "#16a34a", "yellow": "#f97316", "blue": "#2563eb", "purple": "#7d56f4", "cyan": "#0891b2", "white": "#e5e7eb", "brightBlack": "#374151", "brightRed": "#dc2626", "brightGreen": "#22c55e", "brightYellow": "#fb923c", "brightBlue": "#3b82f6", "brightPurple": "#8b5cf6", "brightCyan": "#06b6d4", "brightWhite": "#f9fafb" }

# --- Launch ------------------------------------------------------------------
Type "./build/bubly"
Enter
Sleep 3s

# Browse the main menu to show all three modes.
Down
Sleep 800ms
Down
Sleep 800ms
Up
Up
Sleep 1s

# --- Flow 1: Video formats ---------------------------------------------------
Enter
Sleep 1500ms
Type "https://www.youtube.com/watch?v=jNQXAC9IVRw"
Enter

# Wait for yt-dlp to list formats, then browse the picker.
Sleep 6s
Down
Sleep 700ms
Down
Sleep 700ms
Up
Sleep 1s

# Back to the main menu.
Backspace
Sleep 1500ms

# --- Flow 2: Audio formats ---------------------------------------------------
Down
Enter
Sleep 1500ms
Type "https://www.youtube.com/watch?v=jNQXAC9IVRw"
Enter
Sleep 6s
Down
Sleep 700ms
Up
Sleep 1s

Backspace
Sleep 1500ms

# --- Flow 3: Subtitle languages ----------------------------------------------
Down
Down
Enter
Sleep 1500ms
Type "https://www.youtube.com/watch?v=jNQXAC9IVRw"
Enter
Sleep 6s
Down
Sleep 700ms
Up
Sleep 1s

# --- Quit --------------------------------------------------------------------
Escape
Sleep 2s
