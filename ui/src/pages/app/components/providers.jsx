import React, { useState, createContext } from 'react';

export const VolumeContext = createContext();
const VolumeProvider = ({ children }) => {
  const [volume, setVolume] = useState(0);
  return (
    <VolumeContext.Provider value={{ volume, setVolume }}>
      {children}
    </VolumeContext.Provider>
  );
};

export const PlayerContext = createContext();
const PlayerProvider = ({ children }) => {
  const [showVideo, setShowVideo] = useState(false);
  const [progress, setProgress] = useState(0);
  return (
    <PlayerContext.Provider value={{ showVideo, setShowVideo, progress, setProgress}}>
      {children}
    </PlayerContext.Provider>
  );
};




export const Provider = ({ children }) => {
  return (
    <>
      <VolumeProvider>
          <PlayerProvider>
          {children}
          </PlayerProvider>
      </VolumeProvider>
    </>
  );
};