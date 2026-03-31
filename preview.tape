# Bubly CLI App Demo
# This tape demonstrates the main features of the Bubly CLI app

Output preview.gif

Require build/bubly

Set FontSize 16
Set Width 1000
Set Height 800
Set Padding 50
Set Margin 20
Set MarginFill "#222222"
Set BorderRadius 10
Set WindowBar "Colorful"
Set WindowBarSize 40
Set TypingSpeed 100ms

# Start the app
Type "./build/bubly"
Enter

# Wait for the app to start
Sleep 4s

# Navigate to "Download Youtube video" option and select it
Enter

# Wait for URL input prompt
Sleep 2s

# Enter YouTube URL
Type "https://www.youtube.com/watch?v=jNQXAC9IVRw"
Enter

# Wait for formats to load
Sleep 8s

# Select first format
Enter

# Wait for download to complete
Sleep 10s

# Go back to main menu
Backspace

# Navigate to "Download Youtube audio" option and select it
Down
Enter

# Wait for URL input prompt
Sleep 2s

# Enter YouTube URL
Type "https://www.youtube.com/watch?v=jNQXAC9IVRw"
Enter

# Wait for formats to load
Sleep 8s

# Select first format
Enter

# Wait for download to complete
Sleep 10s

# Go back to main menu
Backspace

# Navigate to "Download Youtube subtitles" option and select it
Down
Down
Enter

# Wait for URL input prompt
Sleep 2s

# Enter YouTube URL
Type "https://www.youtube.com/watch?v=jNQXAC9IVRw"
Enter

# Wait for languages to load
Sleep 8s

# Select first language
Enter

# Wait for download to complete
Sleep 8s

# Quit the application
Escape

Sleep 2s
