import { Link, useLocation } from "react-router-dom";
import { getRoomId } from "../api/watch2gether";
import { useEffect, useState } from "react";

export const GuildIcon = ({ guild, active }) => {
    let className = "ml-1.5 mr-2.5 my-2  h-12 w-12 relative inline-flex items-center justify-center  overflow-hidden rounded-full bg-gray-100  dark:bg-gray-600 border-purple-700 hover:border-2 icon-shadow"
    className += active ? " border-2 icon-active" : " border-0";

    if (guild.icon !== undefined && guild.icon !== '') {
        return (<img className={className} alt={guild.name} src={`https://cdn.discordapp.com/icons/${guild.id}/${guild.icon}.png`} />);
    }
    return (
        <div className={className}>
            <span className="font-medium text-l text-gray-600 dark:text-gray-300">{guild.name.substring(0, 2)}</span>
        </div>
    );
};


const RoomList = ({guilds}) => {
    const location = useLocation()
    const [active, setActive] = useState("")
    useEffect(() => {
      setActive(getRoomId())
    }, [location])
    return <>
        {guilds.map(guild => <Link key={guild.id} to={`/app/${guild.id}`}><GuildIcon guild={guild} key={guild.id} active={guild.id===active}/></Link>)}
    </>
}

export default RoomList