# ccfg Demo — VHS Recording Script
# Usage: vhs docs/demo.tape
# Requires: https://github.com/charmbracelet/vhs

Output docs/demo.gif

# Terminal settings
Set Shell "bash"
Set FontSize 16
Set Width 1400
Set Height 800
Set Padding 20
Set TypingSpeed 80ms
Set Theme "Catppuccin Mocha"

# Start ccfg
Type "ccfg"
Sleep 500ms
Enter
Sleep 2s

# Navigate the tree — expand User scope
Type "j"
Sleep 400ms
Type "j"
Sleep 400ms
Type "j"
Sleep 400ms
Enter
Sleep 1s

# Select a config file
Type "j"
Sleep 400ms
Enter
Sleep 2s

# Switch to right panel to scroll content
Type "l"
Sleep 1s
Type "j"
Sleep 300ms
Type "j"
Sleep 300ms
Type "j"
Sleep 300ms
Type "j"
Sleep 300ms
Type "j"
Sleep 1s

# Switch back to left panel
Type "h"
Sleep 500ms

# Toggle merged view
Type "m"
Sleep 3s

# Turn off merged view
Type "m"
Sleep 1s

# Search for a setting
Type "/"
Sleep 500ms
Type "mcp"
Sleep 1s
Enter
Sleep 2s

# Exit search
Escape
Sleep 1s

# Switch to ranking tabs
Type "1"
Sleep 2s
Type "2"
Sleep 2s
Type "3"
Sleep 2s

# Toggle scope
Type "s"
Sleep 2s

# Back to tools tab
Type "1"
Sleep 1s

# Quit
Type "q"
Sleep 500ms
