import React, { useState, createContext, useEffect } from 'react';
import { useWebSocket } from 'react-use-websocket/dist/lib/use-websocket';
import { getController, getSocket } from '../api/watch2gether';
import { toast } from 'react-hot-toast';

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
  const test = (check) => {console.log("HEELO"), setShowVideo(check)}
  return (
    <PlayerContext.Provider value={{ showVideo, test}}>
      {children}
    </PlayerContext.Provider>
  );
};

export const SocketContext = createContext();
const SocketProvider = ({ children }) => {
  const [socketUrl] = useState(getSocket());
  const { sendMessage, lastMessage, readyState } = useWebSocket(socketUrl);
  const [state, setState] = useState({   
    Current: {
      duration: 0,
      proccessing: 0
    }, 
    Queue: []
  })

  useEffect(() => {
    if (lastMessage !== null) {
      let data = JSON.parse(lastMessage.data)
      if (data.Action === "MESSAGE"){
        toast.error(data.Message)
      }else{
        console.log("WS",data)
        toast.success(data.Action)
        setState({...state, Current:data})
      }
    }
  }, [lastMessage]);

  useEffect(() => {
    const fetchData = async () => {
      let _state = await getController();
      setState(_state)
    }
    fetchData().catch(console.error);;
  }, []);

  return (
    <SocketContext.Provider value={{ state, sendMessage, setState }}>
      {children}
    </SocketContext.Provider>
  );
};

export const Provider = ({ children }) => {
  return (
    <>
      <VolumeProvider>
        <SocketProvider>
          <PlayerProvider>
          {children}
          </PlayerProvider>
        </SocketProvider>
      </VolumeProvider>
    </>
  );
};
