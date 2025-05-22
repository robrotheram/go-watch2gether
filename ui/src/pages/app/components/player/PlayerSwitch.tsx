import React, { useContext } from "react";
import { PlayerContext } from "../providers";

export const PlayerSwitch = () => {
    const { showVideo, setShowVideo} = useContext(PlayerContext)
    return (showVideo ?
        <button onClick={() => setShowVideo(false)} data-tooltip-target="tooltip-shuffle" type="button" className="p-2.5 group rounded-full mr-1 focus:outline-none focus:ring-4  focus:ring-gray-600 hover:bg-gray-600">
            <svg  xmlns="http://www.w3.org/2000/svg"  width="24"  height="24"  viewBox="0 0 24 24"  fill="currentColor"  className="text-white">
                <path stroke="none" d="M0 0h24v24H0z" fill="none"/>
                <path d="M14.983 3l.123 .006c2.014 .214 3.527 .672 4.966 1.673a1 1 0 0 1 .371 .488c1.876 5.315 2.373 9.987 1.451 12.28c-1.003 2.005 -2.606 3.553 -4.394 3.553c-.732 0 -1.693 -.968 -2.328 -2.045a21.512 21.512 0 0 0 2.103 -.493a1 1 0 1 0 -.55 -1.924c-3.32 .95 -6.13 .95 -9.45 0a1 1 0 0 0 -.55 1.924c.717 .204 1.416 .37 2.103 .494c-.635 1.075 -1.596 2.044 -2.328 2.044c-1.788 0 -3.391 -1.548 -4.428 -3.629c-.888 -2.217 -.39 -6.89 1.485 -12.204a1 1 0 0 1 .371 -.488c1.439 -1.001 2.952 -1.459 4.966 -1.673a1 1 0 0 1 .935 .435l.063 .107l.651 1.285l.137 -.016a12.97 12.97 0 0 1 2.643 0l.134 .016l.65 -1.284a1 1 0 0 1 .754 -.54l.122 -.009zm-5.983 7a2 2 0 0 0 -1.977 1.697l-.018 .154l-.005 .149l.005 .15a2 2 0 1 0 1.995 -2.15zm6 0a2 2 0 0 0 -1.977 1.697l-.018 .154l-.005 .149l.005 .15a2 2 0 1 0 1.995 -2.15z" />
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