import React, { useState, createContext, useMemo, useCallback } from 'react';
import { Loading } from './loading';
import { getGuilds, getSettings, getUser } from '../watch2gether';

import { useQuery } from '@tanstack/react-query';


type GuildContextType = {
  guilds: any[];
  user: any;
  settings: any;
  getGuild: (id: string) => any;
};

export const GuildContext = createContext<GuildContextType>({
  guilds: [],
  user: {},
  settings: {},
  getGuild: (id: string) => null,
});

export function GuildProvider({ children }: Readonly<{ children: React.ReactNode }>) {
  const { data: guildsData, isLoading: guildsLoading } = useQuery({
    queryKey: ['guilds'],
    queryFn: getGuilds,
    staleTime: 5 * 60 * 1000, // 5 minutes
  });

  const { data: userData, isLoading: userLoading } = useQuery({
    queryKey: ['user'],
    queryFn: getUser,
    staleTime: 5 * 60 * 1000,
  });

  const { data: settingsData, isLoading: settingsLoading } = useQuery({
    queryKey: ['settings'],
    queryFn: getSettings,
    staleTime: 5 * 60 * 1000,
  });

  const loading = guildsLoading || userLoading || settingsLoading;
  const guilds = guildsData ?? [];
  const user = userData;
  const settings = settingsData;

  const getGuild = useCallback((id: string) => {
    return guilds.find(g => g.id === id);
  }, [guilds]);

  const value = useMemo(() => ({ guilds, user, settings, getGuild }), [guilds, user, settings]);

  return (
    <GuildContext.Provider value={value}>
      {loading ? (
        <div className='text-white wrap-login min-h-screen w-full flex justify-center items-center'>
          <Loading />
        </div>
      ) : (
        children
      )}
    </GuildContext.Provider>
  );
}

type PlayerContextType = {
  showVideo: boolean;
  setShowVideo: React.Dispatch<React.SetStateAction<boolean>>;
  progress: number;
  setProgress: React.Dispatch<React.SetStateAction<number>>;
  volume: number;
  setVolume: React.Dispatch<React.SetStateAction<number>>;
};

export const PlayerContext = createContext<PlayerContextType>({
  showVideo: false,
  volume: 0,
  progress: 0,
  setShowVideo: () => {},
  setProgress: () => {},
  setVolume: () => {},
});

export function PlayerProvider({ children }: Readonly<{ children: React.ReactNode }>) {
  const [showVideo, setShowVideo] = useState(false);
  const [progress, setProgress] = useState(0);
  const [volume, setVolume] = useState(0);

  const value = useMemo(() => ({ showVideo, setShowVideo, progress, setProgress, volume, setVolume }),[showVideo, progress, volume]);

  return (
    <PlayerContext.Provider value={value}>
      {children}
    </PlayerContext.Provider>
  );
}

export function Provider({ children }: Readonly<{ children: React.ReactNode }>) {
  return (
    <GuildProvider>
      <PlayerProvider>
        {children}
      </PlayerProvider>
    </GuildProvider>
  );
}