import React, { useContext, useEffect } from 'react';
import ReactPlayer from 'react-player'
import { getRoomId, playVideoController, skipVideoController } from '../watch2gether';
import { VolumeContext } from './providers';

export const VideoPlayer = ({state, connection}) => {
    const playerRef = React.useRef(null);
    const { volume } = useContext(VolumeContext)
    const onEnded = () => {
        if (!state.loop){
            skipVideoController();
        }
    }

    const onStart = () => {
        playVideoController();
    } 


    const handleProgress = (video_state) => {
        let s = Object.assign({}, state)
        s.current.time = {
            duration: s.current.time.duration,
            progress: Math.floor(video_state.playedSeconds)*1000000000
        }

        const evt =  {
            id: getRoomId(),
            action:{
                type: "UPDATE_DURATION"
            },
            state: {
                Current: s.current
            }
        }
        console.log("SENDING", JSON.stringify(evt))
        connection.send(JSON.stringify(evt))
    };

   
    return(
        <div className='w-full flex justify-center' style={{"maxHeight": "650px", height:"100%"}}>
            <ReactPlayer
                ref={playerRef}
                url={state.current.url}
                width='100%'
                height='100%'
                muted={volume === 0}
                volume={volume / 100}
                onStart={onStart}
                onEnded={onEnded}
                onProgress={handleProgress}
                playing={state.status === "PLAY" }
                loop={state.loop}                
            />
        </div>

    );
}
