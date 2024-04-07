import { useContext, useRef, useState } from "react";
import { useOnClickOutside } from "./nav";
import { GuildContext, PlayerContext } from "./providers";
import { formatTime } from "../watch2gether";
import { toast } from "react-hot-toast";

export const SyncIcon = () => {
    return <svg xmlns="http://www.w3.org/2000/svg" className="w-5 h-5" width={24} height={24} viewBox="0 0 24 24" strokeWidth="3.5" stroke="white" fill="none" strokeLinecap="round" strokeLinejoin="round">
        <path stroke="none" d="M0 0h24v24H0z" fill="none" />
        <path d="M3 7h1.948c1.913 0 3.705 .933 4.802 2.5a5.861 5.861 0 0 0 4.802 2.5h6.448" />
        <path d="M3 17h1.95a5.854 5.854 0 0 0 4.798 -2.5a5.854 5.854 0 0 1 4.798 -2.5h5.454" />
        <path d="M18 15l3 -3l-3 -3" />
    </svg>
}

const UserPlayerItem = ({player}) => {
    const {setProgress} = useContext(PlayerContext)
    const truncate = (input) => {
        const size = 12;
        if (input.length > size) {
            return input.substring(0, size) + '... ';
        }
        return input;
    };
    const handleSync = () => {
        setProgress(player.progress.progress)
        toast.success(`You have synced to ${player.user} position`)
    }
    return <li className="flex align-middle justify-between items-center">
        <span>{truncate(player.user)}: {formatTime(player.progress.progress)}</span>
        <button onClick={()=>{handleSync()}} className="rounded-full p-1 border-purple-700 icon-shadow active:bg-purple-700">
            <SyncIcon />
        </button>
    </li>
}
export const UserPlayer = ({players}) => {
    const p = players.sort(function(a, b) {
        return a.id.localeCompare(b.id);
    });
    return <div className="z-50 absolute bottom-24 left-2 right-4 md:right-auto w-60 text-white rounded-lg shadow-lg  border-purple-950" style={{ background: "rgba(0,0,0,0.8)" }} >
        <div className="max-h-60 py-2 rounded-lg overflow-y-auto hide-scrollbar shadow-lg text-violet-200 flex justify-center p-3 text-lg font-medium bg-violet-900 p">
            Currently Watching
        </div>
        <ul className="p-4 flex flex-col gap-1 max-h-48 overflow-y-auto hide-scrollbar">
            {p.map(function(player) {
                return <UserPlayerItem key={player.id} player={player}/>
            })}
        </ul>
    </div>
}

export const UserPlayerBtn = () => {
    const [visable, setVisible] = useState(false)
    const {players} = useContext(GuildContext)
    const ref = useRef();
    useOnClickOutside(ref, () => setVisible(false));

    return <div ref={ref}>
        <button onClick={()=>{setVisible(!visable)}} className="z-50 hidden  absolute w-12 h-12 bottom-4  left-4 rounded-full bg-violet-700 text-white sm:flex justify-center items-center icon-shadow active:bg-purple-700">
            <svg xmlns="http://www.w3.org/2000/svg" width={24} height={24} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth={2} strokeLinecap="round" strokeLinejoin="round" className="icon icon-tabler icons-tabler-outline icon-tabler-eyeglass-2">
                <path stroke="none" d="M0 0h24v24H0z" fill="none" />
                <path d="M8 4h-2l-3 10v2.5" />
                <path d="M16 4h2l3 10v2.5" />
                <path d="M10 16l4 0" />
                <path d="M17.5 16.5m-3.5 0a3.5 3.5 0 1 0 7 0a3.5 3.5 0 1 0 -7 0" />
                <path d="M6.5 16.5m-3.5 0a3.5 3.5 0 1 0 7 0a3.5 3.5 0 1 0 -7 0" />
            </svg>
            <div className="absolute bottom-0 border-1 border-red-50 right-0 bg-violet-900 w-4 h-4 flex justify-center items-center rounded-full">{players.length}</div>
        </button>
        {visable && <UserPlayer players={players} />}
    </div>
}