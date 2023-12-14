import { createContext, useContext, useEffect, useMemo, useRef, useState } from "react";
import Card from "./Cards"
import { Header, VideoHeader } from "./Header"
import Player from "./Player"
import toast from 'react-hot-toast';
import { addVideoController, getChannelPlaylists, getController, getSocket, updateQueueController } from "../api/watch2gether";
import { useNavigate } from "react-router-dom";
import { PlaylistBtn } from "../pages/app/playlists/playtlist";
import useWebSocket, { ReadyState } from 'react-use-websocket';
import { PlayerContext, Provider, SocketContext } from "./Provider";
import { VideoPlayer } from "./Video";


export const AddVideoCtrl = ({ onAddVideo }) => {
    const [video, setVideo] = useState("");
    const addVideo = async () => {
        if (video.length == 0) {
            return
        }
        onAddVideo(video)
        setVideo("")
    }
    const handleKeyPress = (e) => {
        if (e.key == 'Enter') {
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
        </div>

    )
}

const Controller = () => {
    const {state} = useContext(SocketContext)
    const { showVideo } = useContext(PlayerContext)
    const [loading, setLoading] = useState(true)
    const [playlists, setPlaylists] = useState([])
    const updatePlaylists = async () => { setPlaylists(await getChannelPlaylists()) }

    useEffect(() => {
        updatePlaylists()
        setLoading(false)
    }, []);

    if (loading) {
        return (
            <div className="min-h-screen w-full flex justify-center items-center bg-black">
                <div className="bg-violet-800 w-48 h-48  absolute animate-ping rounded-full delay-10s shadow-2xl"></div>
                <div className="bg-violet-700 w-32 h-32  absolute animate-ping rounded-full delay-5s shadow-xl"></div>
                <div className="bg-violet-600 w-16 h-16  absolute animate-ping rounded-full shadow-xl"></div>
            </div>
        )
    }

    const addVideo = async (video) => {
        try {
            await addVideoController(video)
            toast.success("Video is being added to the queue please wait");
        } catch (e) {
            toast.error("Unable to add video: invalid video url");
        }
    }

    const updateQueue = async (queue) => {
        try {
            await updateQueueController(queue)
            toast.success("Queue updated")
        } catch (e) {
            console.info(e)
            toast.error("Sorry there was an issue updating the queue")
        }
    }
    return (
            <div className="flex flex-col w-full h-full">
                <AddVideoCtrl onAddVideo={addVideo} />
                <div className='bg-violet-800 w-full' style={{ "overflow": "auto" }}>
                    {state.Current.id && (showVideo ? <VideoHeader state={state} />:<Header state={state} />)}
                    <div className='w-full shadow-body px-4 md:px-10 text-white min-h-screen'>
                            <Card queue={state.Queue} updateQueue={updateQueue} />
                    </div>
                </div>
                <Player state={state} />
                <div className="absolute bottom-9 right-0">
                    <PlaylistBtn playlists={playlists} />
                </div>
            </div>

    )
}
export const AppController = () => <Provider><Controller/></Provider>