
import { useEffect, useState } from 'react';
import { getGuilds } from './api/watch2gether';
import logo from './assets/logo.svg'
import user from './data/users.json'


const GuildIcon = ({ guild }) => {
    if (guild.icon !== undefined && guild.icon !== '') {
        return (<img className='ml-1.5 mr-2.5 my-2 h-12 w-12 rounded-full mb-3 border-purple-700 border-0 hover:border-2 icon-shadow shadow-white' alt={guild.name} src={`https://cdn.discordapp.com/icons/${guild.id}/${guild.icon}.png`} />);
    }
    return (
        <div className="ml-1.5 mr-2.5 my-2  h-12 w-12 relative inline-flex items-center justify-center  overflow-hidden rounded-full bg-gray-100  dark:bg-gray-600 border-0 border-purple-700 hover:border-2 icon-shadow">
            <span className="font-medium text-l text-gray-600 dark:text-gray-300">{guild.name.substring(0, 2)}</span>
        </div>
    );
};


const RoomList = ({guilds}) => {
    return <>
        {guilds.map(guild => <a key={guild.id} href={`/app/${guild.id}`}><GuildIcon guild={guild} key={guild.id}/></a>)}
    </>
}

export default RoomList