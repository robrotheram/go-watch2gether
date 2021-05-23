import React, { useEffect, useRef, useState } from 'react';
import ReactPlayer from 'react-player'
import { Empty, Card } from "antd"

import { useDispatch, useSelector } from "react-redux";
import { play, pause, updateSeek, handleFinish } from '../../../store/video/video.actions'

const Player = () => {
    const dispatch = useDispatch();

    const inputRange = useRef(null);
    const playing = useSelector(state => state.video.playing);
    const url = useSelector(state => state.video.url);
    const seek = useSelector(state => state.video.seek_to_user);
    const [seeking, setSeeking] = useState(false);
    
    useEffect(() => {
        setSeeking(true)
        try {
            console.log("SeekTO", seek )
            inputRange.current.seekTo(parseFloat(seek.progress_seconds))
        }
        catch{}
        setSeeking(false)
    }, [seek])
    
    useEffect(() => {
        setSeeking(false)
    }, [url])

    const handleProgress = state => {
        if (!seeking){
            dispatch(updateSeek(state.played, state.playedSeconds))
        }
    }
    const handleEnded = state => {
        console.log("UPDATE_SEEK")
        setSeeking(true)
        dispatch(handleFinish())
    }


    return (
        <Card className="videoPlayer" style={{ "height": "100%", "width": "100%" }}>
            {url !== "" ?
                <ReactPlayer
                    ref={inputRange}
                    width="100%"
                    height="100%"

                    playing={playing}
                    url={url}
                    controls={true}
                    onPause={() => { dispatch(pause()) }}
                    onPlay={() => { dispatch(play()) }}
                    onProgress={handleProgress}
                    onEnded={handleEnded}
                />
                : <Empty style={{ "height": "600px", "width": "100%", "paddingTop": "180px" }} />}
        </Card>
    );
};

export default Player;