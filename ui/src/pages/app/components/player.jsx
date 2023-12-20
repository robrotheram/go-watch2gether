import React, { useContext, useState } from "react";
import { loopVideoController, pauseVideoController, playVideoController, shuffleVideoController, skipVideoController } from "../watch2gether";
import { PlayerContext, VolumeContext } from "./providers";

const Switch = () => {
    const { showVideo, setShowVideo } = useContext(PlayerContext)
    return (
        <div className="md:w-48 w-24 justify-end flex">
            {showVideo ?
                <button onClick={() => setShowVideo(false)} data-tooltip-target="tooltip-shuffle" type="button" className="p-2.5 group rounded-full mr-1 focus:outline-none focus:ring-4  focus:ring-gray-600 hover:bg-gray-600">
                    <svg xmlns="http://www.w3.org/2000/svg" className="w-5 h-5" width={24} height={24} viewBox="0 0 24 24" strokeWidth="1.5" stroke="white" fill="none" strokeLinecap="round" strokeLinejoin="round">
                        <path stroke="none" d="M0 0h24v24H0z" fill="none"></path>
                        <path d="M4 4m0 2a2 2 0 0 1 2 -2h12a2 2 0 0 1 2 2v12a2 2 0 0 1 -2 2h-12a2 2 0 0 1 -2 -2z"></path>
                        <path d="M8 4l0 16"></path>
                        <path d="M16 4l0 16"></path>
                        <path d="M4 8l4 0"></path>
                        <path d="M4 16l4 0"></path>
                        <path d="M4 12l16 0"></path>
                        <path d="M16 8l4 0"></path>
                        <path d="M16 16l4 0"></path>
                    </svg>
                    <span className="sr-only">Shuffle video</span>
                </button>
                :
                <button onClick={() => setShowVideo(true)} data-tooltip-target="tooltip-shuffle" type="button" className="p-2.5 group rounded-full mr-1 focus:outline-none focus:ring-4 focus:ring-gray-600 hover:bg-gray-600">
                    <svg xmlns="http://www.w3.org/2000/svg" className="w-5 h-5" width={24} height={24} viewBox="0 0 24 24" strokeWidth="1.5" stroke="white" fill="none" strokeLinecap="round" strokeLinejoin="round">
                        <path stroke="none" d="M0 0h24v24H0z" fill="none"></path>
                        <path d="M3 17a3 3 0 1 0 6 0a3 3 0 0 0 -6 0"></path>
                        <path d="M13 17a3 3 0 1 0 6 0a3 3 0 0 0 -6 0"></path>
                        <path d="M9 17v-13h10v13"></path>
                        <path d="M9 8h10"></path>
                    </svg>
                    <span className="sr-only">Shuffle video</span>
                </button>}
        </div>
    )
}

const VolumeControl = React.memo(() => {
    const { volume, setVolume } = useContext(VolumeContext)

    const handleChange = (event) => {
        setVolume(event.target.value);
        // Code to update the media player's volume based on the new value
    };

    const MuteBtn = ({ onClick }) => {
        return <svg xmlns="http://www.w3.org/2000/svg" className="w-5 h-5 text-gray-300  group-hover:text-white mr-1"
            fill="currentColor"
            viewBox="0 0 24 24"
            strokeWidth="2"
            stroke="currentColor"
            strokeLinecap="round"
            strokeLinejoin="round"
            onClick={onClick}
        >
            <path stroke="none" d="M0 0h24v24H0z" fill="none"></path>
            <path d="M15 8a5 5 0 0 1 1.912 4.934m-1.377 2.602a5 5 0 0 1 -.535 .464"></path>
            <path d="M17.7 5a9 9 0 0 1 2.362 11.086m-1.676 2.299a9 9 0 0 1 -.686 .615"></path>
            <path d="M9.069 5.054l.431 -.554a.8 .8 0 0 1 1.5 .5v2m0 4v8a.8 .8 0 0 1 -1.5 .5l-3.5 -4.5h-2a1 1 0 0 1 -1 -1v-4a1 1 0 0 1 1 -1h2l1.294 -1.664"></path>
            <path d="M3 3l18 18"></path>
        </svg>
    }
    const MaxVolBtn = ({ onClick }) => {
        return <svg xmlns="http://www.w3.org/2000/svg"
            className="w-5 h-5 text-gray-300  group-hover:text-white"
            fill="currentColor"
            viewBox="0 0 24 24"
            strokeWidth="2"
            stroke="currentColor"
            strokeLinecap="round"
            strokeLinejoin="round"
            onClick={onClick}
        >
            <path stroke="none" d="M0 0h24v24H0z" fill="none"></path>
            <path d="M15 8a5 5 0 0 1 0 8"></path>
            <path d="M17.7 5a9 9 0 0 1 0 14"></path>
            <path d="M6 15h-2a1 1 0 0 1 -1 -1v-4a1 1 0 0 1 1 -1h2l3.5 -4.5a.8 .8 0 0 1 1.5 .5v14a.8 .8 0 0 1 -1.5 .5l-3.5 -4.5"></path>
        </svg>
    }

    return (
        <>
            <div className="items-center md:flex hidden md:w-48 w-24 gap-1">
                <MuteBtn onClick={() => setVolume(0)} />
                <input
                    type="range"
                    min="0"
                    max="100"
                    step="1"
                    value={volume}
                    onChange={handleChange}
                    className="w-32 h-2 rounded-lg appearance-none bg-gradient-to-r from-violet-200  to-violet-700 accent-purple-500 "
                />
                <MaxVolBtn onClick={() => setVolume(100)} />
            </div>
            <div className="flex sm:hidden w-28">
                {volume === 0 ? <MuteBtn onClick={() => setVolume(100)} /> : <MaxVolBtn onClick={() => setVolume(0)} />}
            </div>
        </>
    );
});


const Player = ({ state }) => {
    const formatTime = (seconds) => {
        if (seconds === undefined) {
            seconds = 0
        }
        let iso = new Date(seconds / 1000000).toISOString()
        return iso.substring(11, iso.length - 5);
    }
    const playerProgress = (current, total) => {
        let pct = current / total * 100
        return Math.min(Math.max(pct, 0), 100)
    }

    const handleShuffle = () => {
        shuffleVideoController();
    }
    const handlePlay = () => {
        playVideoController();
    }
    const handlePause = () => {
        pauseVideoController();
    }
    const handleSkip = () => {
        skipVideoController();
    }
    const handleLoop = () => {
        loopVideoController();
    }
    return (
        <div className="flex-shrink grid w-full h-20 grid-cols-1 px-8 pt-1  md:grid-cols-3 bg-zinc-900">
            <div></div>
            {state.active ?
                <div className="flex items-center w-full">
                    <div className="w-full">
                        <div className="flex items-center justify-center mx-auto mb-1">
                            <div className="flex md:w-48 w-28">
                                <Switch />
                            </div>

                            <button onClick={() => handleShuffle()} data-tooltip-target="tooltip-shuffle" type="button" className="p-2.5 group rounded-full  mr-1 focus:outline-none focus:ring-4 focus:ring-gray-600 hover:bg-gray-600">
                                <svg className="w-5 h-5 text-gray-300  group-hover:text-white" fill="currentColor" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 512 512" aria-hidden="true"><path d="M403.8 34.4c12-5 25.7-2.2 34.9 6.9l64 64c6 6 9.4 14.1 9.4 22.6s-3.4 16.6-9.4 22.6l-64 64c-9.2 9.2-22.9 11.9-34.9 6.9s-19.8-16.6-19.8-29.6V160H352c-10.1 0-19.6 4.7-25.6 12.8L284 229.3 244 176l31.2-41.6C293.3 110.2 321.8 96 352 96h32V64c0-12.9 7.8-24.6 19.8-29.6zM164 282.7L204 336l-31.2 41.6C154.7 401.8 126.2 416 96 416H32c-17.7 0-32-14.3-32-32s14.3-32 32-32H96c10.1 0 19.6-4.7 25.6-12.8L164 282.7zm274.6 188c-9.2 9.2-22.9 11.9-34.9 6.9s-19.8-16.6-19.8-29.6V416H352c-30.2 0-58.7-14.2-76.8-38.4L121.6 172.8c-6-8.1-15.5-12.8-25.6-12.8H32c-17.7 0-32-14.3-32-32s14.3-32 32-32H96c30.2 0 58.7 14.2 76.8 38.4L326.4 339.2c6 8.1 15.5 12.8 25.6 12.8h32V320c0-12.9 7.8-24.6 19.8-29.6s25.7-2.2 34.9 6.9l64 64c6 6 9.4 14.1 9.4 22.6s-3.4 16.6-9.4 22.6l-64 64z" fill="currentColor" /></svg>
                                <span className="sr-only">Shuffle video</span>
                            </button>
                            <div id="tooltip-shuffle" role="tooltip" className="absolute z-10 invisible inline-block px-3 py-2 text-sm font-medium text-white transition-opacity duration-300 rounded-lg shadow-sm opacity-0 tooltip bg-gray-700">
                                Shuffle video
                                <div className="tooltip-arrow" data-popper-arrow></div>
                            </div>

                            {state.status === "PLAY" ?
                                <button onClick={() => handlePause()} data-tooltip-target="tooltip-pause" type="button" className="inline-flex items-center justify-center p-2.5 mx-2 font-medium bg-purple-600 rounded-full hover:bg-purple-700 group focus:ring-4 focus:outline-none focus:ring-purple-800">
                                    <svg className="w-4 h-4 text-white" viewBox="0 0 10 14" fill="currentColor" xmlns="http://www.w3.org/2000/svg" aria-hidden="true">
                                        <path fillRule="evenodd" clipRule="evenodd" d="M0.625 1.375C0.625 1.02982 0.904823 0.75 1.25 0.75H2.5C2.84518 0.75 3.125 1.02982 3.125 1.375V12.625C3.125 12.9702 2.84518 13.25 2.5 13.25H1.25C1.08424 13.25 0.925268 13.1842 0.808058 13.0669C0.690848 12.9497 0.625 12.7908 0.625 12.625L0.625 1.375ZM6.875 1.375C6.875 1.02982 7.15482 0.75 7.5 0.75H8.75C8.91576 0.75 9.07473 0.815848 9.19194 0.933058C9.30915 1.05027 9.375 1.20924 9.375 1.375L9.375 12.625C9.375 12.9702 9.09518 13.25 8.75 13.25H7.5C7.15482 13.25 6.875 12.9702 6.875 12.625V1.375Z" fill="currentColor" />
                                    </svg>
                                    <span className="sr-only">Pause video</span>

                                </button>
                                : <button onClick={() => handlePlay()} data-tooltip-target="tooltip-pause" type="button" className="inline-flex items-center justify-center p-2.5 mx-2 font-medium bg-purple-600 rounded-full hover:bg-purple-700 group focus:ring-4 focus:outline-none focus:ring-purple-800">
                                    <svg className="w-4 h-4 text-white" aria-hidden="true" fill="none" stroke="currentColor" strokeWidth="2.5" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                                        <path d="M5.25 5.653c0-.856.917-1.398 1.667-.986l11.54 6.348a1.125 1.125 0 010 1.971l-11.54 6.347a1.125 1.125 0 01-1.667-.985V5.653z" strokeLinecap="round" strokeLinejoin="round"></path>
                                    </svg>
                                    <span className="sr-only">Play video</span>
                                </button>
                            }

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
                        {state.current.id && <div className="flex items-center justify-between space-x-2">
                            <span className="text-sm font-medium  text-gray-400">{formatTime(state.current.time.progress)}</span>
                            <div className="w-full  rounded-full h-1.5 bg-gray-800">
                                <div className="bg-purple-600 h-1.5 rounded-full" style={{ "width": `${playerProgress(state.current.time.progress, state.current.time.duration)}%` }}></div>
                            </div>
                            <span className="text-sm font-medium text-gray-400">{formatTime(state.current.time.duration)}</span>
                        </div>}
                    </div>
                </div>
                :
                <div className="flex items-center w-full text-center text-white">
                    <div className="w-full text-center">Player is not active. <br /> Join the bot to one of the voice channels</div>
                </div>
            }
        </div>
    )
}

export default Player