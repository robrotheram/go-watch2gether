import React, { useContext } from "react";
import { PlayerContext } from "../providers";

export const PlayerSwitch = () => {
    const { showVideo, setShowVideo} = useContext(PlayerContext)
    return (showVideo ?
        <button onClick={() => setShowVideo(false)} data-tooltip-target="tooltip-shuffle" type="button" className="p-2.5 group rounded-full mr-1 focus:outline-none focus:ring-4  focus:ring-gray-600 hover:bg-gray-600">
            <svg xmlns="http://www.w3.org/2000/svg" className="w-5 h-5" width={24} height={24} viewBox="0 0 24 24" strokeWidth="1.5" stroke="white" fill="none" strokeLinecap="round" strokeLinejoin="round">
                <path stroke="none" d="M0 0h24v24H0z" fill="none"></path>
                <path d="M3 17a3 3 0 1 0 6 0a3 3 0 0 0 -6 0"></path>
                <path d="M13 17a3 3 0 1 0 6 0a3 3 0 0 0 -6 0"></path>
                <path d="M9 17v-13h10v13"></path>
                <path d="M9 8h10"></path>
            </svg>
            <span className="sr-only">Shuffle video</span>
        </button>
        :
        <button onClick={() => setShowVideo(true)} data-tooltip-target="tooltip-shuffle" type="button" className="p-2.5 group rounded-full mr-1 focus:outline-none focus:ring-4 focus:ring-gray-600 hover:bg-gray-600">
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
    )
}