import { useState } from "react";
import { clearVideoController, shuffleVideoController } from "../watch2gether";
import { State } from "@/types";

type AddVideoCtrlProps = {
    onAddVideo: (url: string) => void
    state?: State
}

export const AddVideoCtrl = ({ onAddVideo, state }:AddVideoCtrlProps) => {
    const [video, setVideo] = useState("");
    const addVideo = async () => {
        if (video.length == 0) {
            return
        }
        onAddVideo(video)
        setVideo("")
    }
    const handleKeyPress = (e: React.KeyboardEvent<HTMLInputElement>) => {
        if (e.key === 'Enter') {
            addVideo()
        }
    }
    return (
        <div className="w-full bg-violet-950 shadow-2xl z-20 flex items-center">
            <div className="inset-y-0 left-0 flex items-center pl-3 pointer-events-none">
                <svg fill="none" className="w-5 h-5 text-gray-100" stroke="currentColor" strokeWidth="1.5" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg" aria-hidden="true">
                    <path strokeLinecap="round" strokeLinejoin="round" d="M3.375 19.5h17.25m-17.25 0a1.125 1.125 0 01-1.125-1.125M3.375 19.5h1.5C5.496 19.5 6 18.996 6 18.375m-3.75 0V5.625m0 12.75v-1.5c0-.621.504-1.125 1.125-1.125m18.375 2.625V5.625m0 12.75c0 .621-.504 1.125-1.125 1.125m1.125-1.125v-1.5c0-.621-.504-1.125-1.125-1.125m0 3.75h-1.5A1.125 1.125 0 0118 18.375M20.625 4.5H3.375m17.25 0c.621 0 1.125.504 1.125 1.125M20.625 4.5h-1.5C18.504 4.5 18 5.004 18 5.625m3.75 0v1.5c0 .621-.504 1.125-1.125 1.125M3.375 4.5c-.621 0-1.125.504-1.125 1.125M3.375 4.5h1.5C5.496 4.5 6 5.004 6 5.625m-3.75 0v1.5c0 .621.504 1.125 1.125 1.125m0 0h1.5m-1.5 0c-.621 0-1.125.504-1.125 1.125v1.5c0 .621.504 1.125 1.125 1.125m1.5-3.75C5.496 8.25 6 7.746 6 7.125v-1.5M4.875 8.25C5.496 8.25 6 8.754 6 9.375v1.5m0-5.25v5.25m0-5.25C6 5.004 6.504 4.5 7.125 4.5h9.75c.621 0 1.125.504 1.125 1.125m1.125 2.625h1.5m-1.5 0A1.125 1.125 0 0118 7.125v-1.5m1.125 2.625c-.621 0-1.125.504-1.125 1.125v1.5m2.625-2.625c.621 0 1.125.504 1.125 1.125v1.5c0 .621-.504 1.125-1.125 1.125M18 5.625v5.25M7.125 12h9.75m-9.75 0A1.125 1.125 0 016 10.875M7.125 12C6.504 12 6 12.504 6 13.125m0-2.25C6 11.496 5.496 12 4.875 12M18 10.875c0 .621-.504 1.125-1.125 1.125M18 10.875c0 .621.504 1.125 1.125 1.125m-2.25 0c.621 0 1.125.504 1.125 1.125m-12 5.25v-5.25m0 5.25c0 .621.504 1.125 1.125 1.125h9.75c.621 0 1.125-.504 1.125-1.125m-12 0v-1.5c0-.621-.504-1.125-1.125-1.125M18 18.375v-5.25m0 5.25v-1.5c0-.621.504-1.125 1.125-1.125M18 13.125v1.5c0 .621.504 1.125 1.125 1.125M18 13.125c0-.621.504-1.125 1.125-1.125M6 13.125v1.5c0 .621-.504 1.125-1.125 1.125M6 13.125C6 12.504 5.496 12 4.875 12m-1.5 0h1.5m-1.5 0c-.621 0-1.125.504-1.125 1.125v1.5c0 .621.504 1.125 1.125 1.125M19.125 12h1.5m0 0c.621 0 1.125.504 1.125 1.125v1.5c0 .621-.504 1.125-1.125 1.125m-17.25 0h1.5m14.25 0h1.5"></path>
                </svg>
            </div>
            <input onKeyDown={handleKeyPress} type="text" value={video} onChange={(e) => setVideo(e.target.value)} className="block w-full p-4 pl-10 text-xl bg-transparent text-white focus:ring-0 focus:outline-none " placeholder="Add New video" required />
            <button onClick={() => addVideo()} className="text-white whitespace-nowrap py-2 px-4 mr-2 bg-purple-700 hover:bg-purple-800 focus:ring-4 focus:outline-none focus:ring-purple-300 font-medium rounded-lg">Add Video</button>
            {state&&<div className="inline-flex rounded-md shadow-sm mr-2">
                <button type="button" onClick={() => shuffleVideoController(state.id)} className="inline-flex items-center px-4 py-2 text-sm font-medium text-gray-900 border border-gray-200 rounded-l-lg hover:bg-gray-100 hover:text-purple-700 focus:z-10 focus:ring-2 focus:ring-purple-700 focus:text-purple-700 bg-zinc-100 ">
                    <svg xmlns="http://www.w3.org/2000/svg" className="icon icon-tabler icon-tabler-square-letter-x" width={24} height={24} viewBox="0 0 24 24" strokeWidth={2} stroke="currentColor" fill="none" strokeLinecap="round" strokeLinejoin="round">
                        <path stroke="none" d="M0 0h24v24H0z" fill="none" />
                        <path d="M18 4l3 3l-3 3" />
                        <path d="M18 20l3 -3l-3 -3" />
                        <path d="M3 7h3a5 5 0 0 1 5 5a5 5 0 0 0 5 5h5" />
                        <path d="M21 7h-5a4.978 4.978 0 0 0 -3 1m-4 8a4.984 4.984 0 0 1 -3 1h-3" />
                    </svg>
                </button>
                <button type="button" onClick={() => clearVideoController(state.id)} className="font-bold inline-flex items-center px-4 py-2 text-gray-900 bg-white border border-gray-200 rounded-r-md hover:bg-gray-100 hover:text-purple-700 focus:z-10 focus:ring-2 focus:ring-purple-700 focus:text-purple-700">
                    <svg xmlns="http://www.w3.org/2000/svg" className="icon icon-tabler icon-tabler-square-letter-x" width={24} height={24} viewBox="0 0 24 24" strokeWidth={2} stroke="currentColor" fill="none" strokeLinecap="round" strokeLinejoin="round">
                        <path stroke="none" d="M0 0h24v24H0z" fill="none" />
                        <path d="M3 3m0 2a2 2 0 0 1 2 -2h14a2 2 0 0 1 2 2v14a2 2 0 0 1 -2 2h-14a2 2 0 0 1 -2 -2z" />
                        <path d="M10 8l4 8" />
                        <path d="M10 16l4 -8" />
                    </svg>
                </button>
            </div>}
        </div>
    )
}