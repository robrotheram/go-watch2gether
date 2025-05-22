import React, { useContext } from "react";
import { PlayerContext } from "../providers";

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


export const VolumeControl = React.memo(() => {
    const { volume, setVolume } = useContext(PlayerContext)

    const handleChange = (event) => {
        setVolume(event.target.value);
        // Code to update the media player's volume based on the new value
    };
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
            <div className="flex sm:hidden items-center">
                {volume === 0 ? <MuteBtn onClick={() => setVolume(100)} /> : <MaxVolBtn onClick={() => setVolume(0)} />}
            </div>
        </>
    );
});