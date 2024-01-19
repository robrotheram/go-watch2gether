import React from "react";
import { loopVideoController, pauseVideoController, playVideoController, skipVideoController } from "../../watch2gether";
import { PlayBtn } from "./PlayBtn";
import { PlayerSwitch } from "./PlayerSwitch";
import { VolumeControl } from "./VolumeControl";

export const AudioPlayer = ({ state }) => {
    const formatTime = (seconds) => {
        if (seconds === undefined) {
            seconds = 0
        }
        let iso = new Date(seconds / 1000000).toISOString()
        return iso.substring(11, iso.length - 5);
    }

    const progressPercentage = (current, total) => {
        if (total === -1) {
            return 100
        }
        let pct = current / total * 100
        return Math.min(Math.max(pct, 0), 100)
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
        <div className="w-full h-20 grid-cols-1 px-8 pt-1  bg-zinc-900 absolute bottom-0 flex justify-center flex-col items-center">
            <div className="flex items-center mb-1 lg:w-auto w-3/5">
                <div className="flex w-64 justify-end gap-1">
                    <PlayerSwitch />
                </div>
                <div className="flex w-24 justify-center">
                    <PlayBtn play={handlePlay} pause={handlePause} status={state.status} />
                </div>
                <div className="flex w-64 justify-end gap-1">
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
            {state.current.id &&
                <div className="flex items-center justify-between space-x-2 w-4/5">
                    <span className="text-sm font-medium  text-gray-400 w-16">{formatTime(state.current.time.progress)}</span>
                    <div className="w-full  rounded-full h-1.5 bg-gray-800 ">
                        <div className="bg-purple-600 h-1.5 rounded-full " style={{ "width": `${progressPercentage(state.current.time.progress, state.current.time.duration)}%` }}></div>
                    </div>
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
    )
}
