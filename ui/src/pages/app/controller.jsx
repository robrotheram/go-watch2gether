import { useContext, useEffect, useRef, useState } from "react";
import Card from "./components/card";
import toast from 'react-hot-toast';
import { addVideoController, getChannelPlaylists, getSocket, updateQueueController, getController, createController } from "./watch2gether";
import { Header } from "./components/header";
import { useNavigate } from "react-router-dom";
import { PlaylistBtn } from "./playlist";
import { PlayerContext } from "./components/providers";
import { Loading } from "./components/loading";
import { useHotkeys } from 'react-hotkeys-hook'
import { Player } from "./components/player";
import { AddVideoCtrl } from "./components/header/Controls";


export const AppController = () => {
    const navigate = useNavigate();
    const [loading, setLoading] = useState(true)
    const [playlists, setPlaylists] = useState([])
    const [notificationURL, setNotificationURL] = useState(null)
    const [debug, setDebug] = useState(false)
    const { showVideo, setProgress } = useContext(PlayerContext)
    
    const [state, setState] = useState({
        id: "",
        status: "STOPPED",
        queue: [],
        current: {}
    })


    useHotkeys('ctrl+shift+b', () => setDebug(!debug), [debug])
    


    const updatePlaylists = async () => {
        try {
            setPlaylists(await getChannelPlaylists())
        } catch (error) { }
    }

    const getState = async () => {
        setLoading(true)
        try {
            setState(await getController())
            setNotificationURL(getSocket())
        } catch (error) {
            try {
                setState(await createController())
                setNotificationURL(getSocket())
            } catch (error) {
                console.log("ERROR", error)
                navigate("/app")
            }
        }
        setLoading(false)
    }

    

    const connection = useRef(null)

    useEffect(() => {
        if (notificationURL === null) {
            return
        }
        const socket = new WebSocket(notificationURL)
        socket.addEventListener("open", (event) => {
            socket.send("Connection established")
        })
        socket.addEventListener("message", (event) => {
            let evt = JSON.parse(event.data)
            if (evt.action.type === "SEEK"){
                console.log("SEEK", evt)
                setProgress(evt.state.current.time.progress)
                toast.success(`${evt.action.user} has synced the track to their position`)
            }
            setState(evt.state)
        })
        connection.current = socket
        return () => socket.close()
    }, [notificationURL])


    useEffect(() => {
        updatePlaylists()
        getState()
    }, []);

    const addVideo = async (video) => {
        try {
            await addVideoController(video)
        } catch (e) {
            console.log("ADD VIDOE ERROR", e)
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

    return loading ?
        <Loading />
        :
        <div className="flex flex-col w-full h-full">
            <AddVideoCtrl onAddVideo={addVideo} controls/>
            <div className='bg-violet-800 w-full overflow-auto mb-20'>
                {state.current.id && <Header state={state} />}
                <div className='w-full shadow-body px-4 md:px-10 text-white min-h-screen'>
                    <Card queue={state.queue} updateQueue={updateQueue} />
                </div>
            </div>
            <Player state={state} connection={connection.current}/>
            <div className="absolute md:bottom-9 bottom-12 md:right-0 -right-2">
                <PlaylistBtn playlists={playlists} />
            </div>

            {debug && <div className="w-full h-1/3 z-50 p-4">
                <pre  className="bg-stone-900 text-white rounded-lg p-8 h-full overflow-auto">
                    <code>
                        {JSON.stringify(state, null, 2)}
                    </code>
                </pre>
            </div>}
        </div>
}
