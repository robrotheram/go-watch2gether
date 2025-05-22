import React, { useContext, useEffect, useState } from "react";
import { formatTime, loopVideoController, pauseVideoController, playVideoController, seekVideoController, skipVideoController } from "../../watch2gether";
import { FullScreen, useFullScreenHandle } from "react-full-screen";
import ReactPlayer from "react-player";
import { toast } from "react-hot-toast";
import { FullScreenBtn } from "./FullScreenBtn";
import { PlayerSwitch } from "./PlayerSwitch";
import { PlayBtn } from "./PlayBtn";
import { VolumeControl } from "./VolumeControl";
import { PlayerContext } from "../providers";
import { useHotkeys } from "react-hotkeys-hook";
import { State } from "@/types";


const microseconds = 1000000000;
type PlayerProps = {
    state: State
    connection: WebSocket
}
export const VideoPlayer = ({ state, connection }:PlayerProps) => {
    const playerRef = React.useRef<any>(null);
    const [playerProgress, setPlayerProgress] = useState(0)
    const [updating, setUpdating] = useState(false)
    const { progress, volume } = useContext(PlayerContext)
    const handle = useFullScreenHandle();

    useHotkeys('f', () => handle.active ? handle.exit() : handle.enter() , [handle])
    useHotkeys('space', () => state.status === "PLAY"? handlePause(): handlePlay(), [state.status])

    useEffect(() => {
        if (!updating && (Math.abs((playerProgress - progress) / microseconds)) > 2) {
            setPlayerProgress(progress)
            playerRef.current.seekTo(progress / microseconds)
        }
        console.log("PRGORSS UPDSATED")
    }, [progress])


    useEffect(() => {
        console.log("LOADING")
    }, [playerRef])

    const handleOnSeek = (evt:any) => {
        setUpdating(false)
        playerRef.current.seekTo(evt.target.value / microseconds);
    }

    const updateSeek = () => {
        toast.success("Syncing everyone to your position")
        seekVideoController(state.id, playerProgress)
    }

    const handleProgessChange = (evt:any) => {
        setPlayerProgress(evt.target.value)
    }

    const handlePlay = () => {
        playVideoController(state.id);
    }
    const handlePause = () => {
        pauseVideoController(state.id);
    }
    const handleSkip = () => {
        skipVideoController(state.id);
    }
    const handleLoop = () => {
        loopVideoController(state.id);
    }
    const onEnded = () => {
        if (!state.loop) {
            skipVideoController(state.id);
        }
    }

    const handleProgress = (video_state:any) => {
        let s = {...state}
        let seconds = Math.floor(video_state.playedSeconds) * microseconds;
        s.current!.time = {
            duration: s.current!.time.duration,
            progress: seconds
        }

        const evt = {
            id: state.id,
            action: {
                type: "UPDATE_DURATION"
            },
            state: {
                Current: s.current
            }
        }
        connection.send(JSON.stringify(evt))
        setPlayerProgress(seconds)
    };

    const getMediaUrl = () => {
        if (!state.current) {
            return ""
        }

       switch(state.current.type){
            case "YOUTUBE_LIVE": return state.current.url
            case "YOUTUBE": return state.current.url
            default: return state.current.audio_url
        }
    }  

    return (
        <FullScreen handle={handle}>
            <div className={`w-full flex justify-center bg-black fixed ${handle.active ? "top-0 bottom-0" : "top-16 bottom-20"} z-20`}>
                <ReactPlayer
                    ref={playerRef}
                    url={getMediaUrl()}
                    width='100%'
                    height='100%'
                    muted={volume === 0}
                    volume={volume / 100}
                    onStart={handlePlay}
                    onEnded={onEnded}
                    onProgress={handleProgress}
                    onPlay={handlePlay}
                    playing={state.status === "PLAY"}
                    loop={state.loop}
                    controls={false}
                />
            </div>
            <div className="w-full h-20 grid-cols-1 px-8 pt-1  bg-zinc-900 absolute bottom-0 flex justify-center flex-col items-center">
                <div className="flex items-center mb-1 lg:w-auto">
                    <div className="flex md:w-64  justify-end gap-1">
                        <button onClick={() => updateSeek()} data-tooltip-target="tooltip-shuffle" type="button" className="p-2.5 group rounded-full  mr-1 focus:outline-none focus:ring-4 focus:ring-gray-600 hover:bg-gray-600">
                            <svg xmlns="http://www.w3.org/2000/svg" className="w-5 h-5" width={24} height={24} viewBox="0 0 24 24" strokeWidth="3.5" stroke="white" fill="none" strokeLinecap="round" strokeLinejoin="round">
                                <path stroke="none" d="M0 0h24v24H0z" fill="none" />
                                <path d="M3 7h1.948c1.913 0 3.705 .933 4.802 2.5a5.861 5.861 0 0 0 4.802 2.5h6.448" />
                                <path d="M3 17h1.95a5.854 5.854 0 0 0 4.798 -2.5a5.854 5.854 0 0 1 4.798 -2.5h5.454" />
                                <path d="M18 15l3 -3l-3 -3" />
                            </svg>
                        </button>
                        <FullScreenBtn status={handle.active} show={() => {handle.enter();}} hide={handle.exit} />
                        <PlayerSwitch/>
                    </div>
                    <div className="flex md:w-24 sm:w-12 justify-center">
                        <PlayBtn play={handlePlay} pause={handlePause} status={state.status} />
                    </div>
                    <div className="flex md:w-64 sm:w-12 justify-end gap-1">
                        <button onClick={() => handleSkip()} data-tooltip-target="tooltip-next" type="button" className="p-2.5 group rounded-full mr-1 focus:outline-none focus:ring-4 focus:ring-gray-200  hover:bg-gray-600">
                            <svg className="w-5 h-5 text-gray-300  group-hover:text-white" fill="currentColor" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 320 512" aria-hidden="true"><path d="M52.5 440.6c-9.5 7.9-22.8 9.7-34.1 4.4S0 428.4 0 416V96C0 83.6 7.2 72.3 18.4 67s24.5-3.6 34.1 4.4l192 160L256 241V96c0-17.7 14.3-32 32-32s32 14.3 32 32V416c0 17.7-14.3 32-32 32s-32-14.3-32-32V271l-11.5 9.6-192 160z" fill="currentColor" /></svg>
                            <span className="sr-only">Next video</span>
                        </button>
                        <div id="tooltip-next" role="tooltip" className="absolute z-10 invisible inline-block px-3 py-2 text-sm font-medium text-white transition-opacity duration-300  rounded-lg shadow-sm opacity-0 tooltip bg-gray-700">
                            Next video
                            <div className="tooltip-arrow" data-popper-arrow></div>
                        </div>
                        <button onClick={() => handleLoop()} data-tooltip-target="tooltip-restart" type="button" className="relative p-2.5 group rounded-full mr-1 focus:outline-none focus:ring-4 focus:ring-gray-600 hover:bg-gray-600">
                            <svg className="w-5 h-5 text-gray-300 group-hover:text-white" fill="currentColor" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 512 512" aria-hidden="true"><path d="M0 224c0 17.7 14.3 32 32 32s32-14.3 32-32c0-53 43-96 96-96H320v32c0 12.9 7.8 24.6 19.8 29.6s25.7 2.2 34.9-6.9l64-64c12.5-12.5 12.5-32.8 0-45.3l-64-64c-9.2-9.2-22.9-11.9-34.9-6.9S320 19.1 320 32V64H160C71.6 64 0 135.6 0 224zm512 64c0-17.7-14.3-32-32-32s-32 14.3-32 32c0 53-43 96-96 96H192V352c0-12.9-7.8-24.6-19.8-29.6s-25.7-2.2-34.9 6.9l-64 64c-12.5 12.5-12.5 32.8 0 45.3l64 64c9.2 9.2 22.9 11.9 34.9 6.9s19.8-16.6 19.8-29.6V448H352c88.4 0 160-71.6 160-160z" fill="currentColor" /></svg>
                            <span className="sr-only">Loop video</span>
                            {state.loop && <div className="absolute bottom-1 right-1 w-3 h-3 text-xs font-bold text-white bg-red-500 border-2 border-white rounded-full" />}
                        </button>
                        <VolumeControl />
                    </div>
                </div>
                {state.current &&
                    <div className="flex items-center justify-between space-x-2 w-4/5">
                        <span className="text-sm font-medium  text-gray-400 w-16">{formatTime(playerProgress)}</span>
                        <input
                            type="range"
                            min={0}
                            max={state.current.time.duration}
                            step={1}
                            value={playerProgress}

                            onMouseDown={() => {
                                setUpdating(true)
                                console.log("mouse_Down")
                            }}
                            onMouseUp={handleOnSeek}
                            onTouchStart={() => { setUpdating(true) }}
                            onTouchEnd={() => { setUpdating(false) }}
                            onChange={handleProgessChange}
                            className="w-full h-2 rounded-lg appearance-none cursor-pointer bg-gray-700 accent-purple-500"

                        /> 
                        {state.current.time.duration > -1 ?
                            <span className="text-sm font-medium text-gray-400">{formatTime(state.current.time.duration)}</span>
                            :
                            <div className="text-sm font-medium text-gray-400 inline-flex items-center gap-1">
                                <span className="w-2 h-2 animate-ping  bg-red-800 rounded-full block" /> live
                            </div>
                        }
                    </div>
                }
            </div>
        </FullScreen>
    )
}