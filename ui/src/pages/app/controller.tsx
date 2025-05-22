import { useContext, useEffect, useRef, useState } from "react";
import Card from "./components/card";
import toast from 'react-hot-toast';
import { addVideoController, getController, updateQueueController } from "./watch2gether";
import { Header } from "./components/header";
import { useParams } from "react-router-dom";
import { GuildContext, PlayerContext } from "./components/providers";
import { Loading } from "./components/loading";
import { AddVideoCtrl } from "./components/AddVideoComponent";
import { UserPlayerBtn } from "./components/UserPlayer";
import { useMutation, useQuery } from "@tanstack/react-query";
import { Event, Media, PlayerMeta, State } from "@/types";
import { VideoPlayer } from "./components/player/VideoPlayer";
import { AudioPlayer } from "./components/player/AudioPlayer";


export const AppController = () => {

    const params = useParams()
    const { id } = params

    const [state, setState] = useState<State>()
    const [players, setPlayers] = useState<PlayerMeta[]>()
    const { showVideo, setProgress } = useContext(PlayerContext)
    const {user} = useContext(GuildContext)
    

    const { isPending, data } = useQuery({
        queryKey: ['controller', id],
        queryFn: () => getController(id!),
    })

    const addVideo = useMutation({
        mutationFn: (url: string) => {
            if (state) {
                const placeholder: Media = {
                    id: `placeholder-${Date.now()}`,
                    url: url,
                    title: "Loading video...",
                    user: user.username,
                    time: { duration: 0, progress: 0 },
                    loading: true
                }
                setState(prev => ({
                    ...prev!,
                    queue: [...prev!.queue, placeholder]
                }))
            }
            return addVideoController(id!, url)
        },
    })

    const updateQueue = useMutation({
        mutationFn: (queue: Media[]) => {
            setState(prev => ({...prev!, queue}))
            return updateQueueController(id!, queue)
        },
    })

    const connection = useRef<WebSocket>()

    useEffect(() => {
        if (data === null) {
            return
        }
        let notificationURL = ((window.location.protocol === "https:") ? "wss://" : "ws://") + window.location.host + "/api/channel/" + id + "/ws";
        const socket = new WebSocket(notificationURL)
        socket.addEventListener("open", (event) => {
            socket.send("Connection established")
        })
        socket.addEventListener("message", async (event) => {
            let evt:Event = JSON.parse(event.data)
            if (evt.action.type === "SEEK") {
                console.log("SEEK", evt)
                setProgress(evt.state.current!.time.progress)
                toast.success(`${evt.action.user} has synced the track to their position`)
            }
            setState(evt.state)
            setPlayers(evt.players)
        })
        connection.current = socket
        return () => socket.close()
    }, [data])




    return isPending || !state ?
        <Loading />
        :
        <div className="flex flex-col w-full h-full">
            <AddVideoCtrl onAddVideo={addVideo.mutate} state={state} />
            <div className='bg-violet-800 w-full overflow-auto mb-20'>
                {state.current && <Header state={state} />}
                <div className='w-full shadow-body px-4 md:px-10 text-white min-h-screen'>
                    <Card queue={state.queue} updateQueue={updateQueue.mutate} />
                </div>
            </div>
            {(showVideo && connection.current) ? <VideoPlayer state={state} connection={connection.current}/> : <AudioPlayer state={state}/>}

            {players&& <UserPlayerBtn players={players}/>}
        </div>
}
