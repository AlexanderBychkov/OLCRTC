#!/bin/bash
# Update Olcbox subscription files with current room IDs.
# Run via cron every minute.

SUB_DIR="/opt/olcrtc/subs"

get_env() {
    local container="$1" varname="$2"
    docker inspect "$container" \
        --format '{{range .Config.Env}}{{println .}}{{end}}' 2>/dev/null \
        | grep "^${varname}=" | cut -d= -f2-
}

update_sub() {
    local n="$1" container="$2" vol="$3"
    local room key client transport fps batch

    room=$(cat "/var/lib/docker/volumes/${vol}/_data/room_id" 2>/dev/null)
    [ -z "$room" ] && return

    key=$(get_env "$container" OLCRTC_KEY)
    client=$(get_env "$container" OLCRTC_CLIENT_ID)
    transport=$(get_env "$container" OLCRTC_TRANSPORT)
    fps=$(get_env "$container" OLCRTC_VP8_FPS)
    batch=$(get_env "$container" OLCRTC_VP8_BATCH)

    [ -z "$key" ] || [ -z "$client" ] && return

    local uri="olcrtc://wbstream?${transport}<vp8-fps=${fps}&vp8-batch=${batch}>@${room}#${key}%${client}\$server"
    local out="${SUB_DIR}/sub${n}.txt"

    # Only write if changed
    if [ ! -f "$out" ] || [ "$(cat "$out")" != "$uri" ]; then
        echo "$uri" > "$out"
        logger -t olcrtc-subs "sub${n}.txt updated: room=${room}"
    fi
}

update_sub 1 olcrtc-wb-1 olcrtc-state-wb-1
update_sub 2 olcrtc-wb-2 olcrtc-state-wb-2
