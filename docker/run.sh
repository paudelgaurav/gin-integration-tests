while true; do
    echo "[run.sh] Starting debugging..."

    # Start Delve in headless mode, running the main package (go run . equivalent)
    dlv debug --headless --log --listen=:2345 --api-version=2 --output ./__debug_bin --accept-multiclient --continue -- . &
    
    PID=$!
    
    echo "[run.sh] Watching for changes..."

    # Monitor file changes, log each change for debugging
    inotifywait -m -e modify -e move -e create -e delete -e attrib --exclude '(__debug_bin|\.git|data)' -r . |
    while read -r directory events filename; do
        echo "[run.sh] Change detected in $directory$filename - Event: $events"
        echo "[run.sh] Stopping process id: $PID"
        kill -15 $PID  
        pkill -f __debug_bin  
        break
    done
done
