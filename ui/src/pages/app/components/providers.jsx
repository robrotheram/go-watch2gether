import React, { useState, createContext, useEffect } from 'react';
import { Loading } from './loading';
import { getGuilds, getSettings, getUser } from '../watch2gether';

export const GuildContext = createContext();
const GuildProvider = ({ children }) => {
  const [guilds, setGuilds] = useState([]);
  const [loading, setLoading] = useState(true)
  const [user, setUser] = useState({});
  const [settings, setSettings] = useState({})
  const [players, setPlayers] = useState([])

  useEffect(() => {
    async function get() {
      const g = await getGuilds();
      setUser(await getUser())
      setSettings(await getSettings())
      if (g != null) {
        setGuilds(g);
      }
      setLoading(false)
    }
    if (guilds.length == 0) { get() };
  }, [])

  const getGuild = (id) => {
    return guilds.filter(g => g.id === id)[0]
  }

  return (
    <GuildContext.Provider value={{ guilds, user, settings, players, setPlayers, getGuild }}>
      {loading ? <div className='text-white wrap-login min-h-screen w-full flex justify-center items-center'><Loading /></div>
        : children}
    </GuildContext.Provider>
  );
};

export const PlayerContext = createContext();
const PlayerProvider = ({ children }) => {
  const [showVideo, setShowVideo] = useState(false);
  const [progress, setProgress] = useState(0);
  const [volume, setVolume] = useState(0);
  return (
    <PlayerContext.Provider value={{ showVideo, setShowVideo, progress, setProgress, volume, setVolume }}>
      {children}
    </PlayerContext.Provider>
  );
};




export const Provider = ({ children }) => {
  return (
    <>
      <GuildProvider>
        <PlayerProvider>
          {children}
        </PlayerProvider>
      </GuildProvider>
    </>
  );
};