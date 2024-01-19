import { useContext } from "react"
import { AudioPlayer } from "./AudioPlayer"
import { VideoPlayer } from "./VideoPlayer"
import { PlayerContext } from "../providers"

export const Player = ({state, connection}) => {
    const { showVideo} = useContext(PlayerContext)
    return showVideo ? <VideoPlayer state={state} connection={connection}/> : <AudioPlayer state={state}/>
}