#!/bin/sh
for pid_dir in /proc/[1-9]*/exe; do
    case "$(readlink "$pid_dir" 2>/dev/null || true)" in
        */olcrtc) exit 0 ;;
    esac
done
exit 1
