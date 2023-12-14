import React, { useContext, useEffect, useMemo } from 'react';
import ReactPlayer from 'react-player'
import { SocketContext, VolumeContext } from './Provider';

const AppVideo = ({id, video}) => {
    const { volume } = useContext(VolumeContext)
    const {state, sendMessage, setState } = useContext(SocketContext)
    const playerRef = React.useRef(null);
    useEffect(() => {
        if (state.State === "PLAYING" && playerRef !== null) {
            console.log(playerRef) //.play()
        }
    }, [state.State])

    const onProgress = ({ playedSeconds }) => {
        state.Proccessing = playedSeconds * 1000000000
        const evnt = {
            "Source": "WEB",
            "Action": "UPDATE",
            "Message": "",
            "State": state
        }
        setState(state)
        sendMessage(JSON.stringify(evnt))
    }

    const onEnded = () => {
        const evnt = {
            "Source": "WEB",
            "Action": "FINISHED",
            "Message": "",
            "State": state
        }
        sendMessage(JSON.stringify(evnt))
    }

    const onStart = () => {
        state.Proccessing = 0
        const evnt = {
            "Source": "WEB",
            "Action": "UPDATE",
            "Message": "",
            "State": state
        }
     
        sendMessage(JSON.stringify(evnt))
    }
    console.log("REREBDER")
    return(
        <div className='w-full flex justify-center' style={{"maxHeight": "650px"}}>
            <ReactPlayer
                ref={playerRef}
                url={`/api/channel/${state.id}/stream?id=${id}`}
                width='auto'
                height='100%'
                muted={volume === 0}
                volume={volume / 100}
                onProgress={onProgress}
                onStart={onStart}
                onEnded={onEnded}
                playing={state.Current.Action === "PLAY" }
                loop={state.Loop}
                controls
                config={{
                    file: {
                        forceVideo: video,
                        forceAudio: !video,
                        forceHLS: true,
                        forceDLS: true
                    }
                }}
            />
        </div>

    );
}

export const VideoPlayer = React.memo(AppVideo)