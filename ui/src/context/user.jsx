import React, { createContext, useState, useEffect } from "react";
import { useLocation, useNavigate } from "react-router";
import { getUser } from "../api/watch2gether";
// create context
const UserContext = createContext();
const defaultUser =  {
  id: '',
  room: '',
  auth: false,
  username: '',
  icon: '',
  guilds: [],
  video_id: '',
  isHost: false,
  playing: false,
}

const UserContextProvider = ({ children }) => {
  const nav = useNavigate();
  const loc = useLocation();
  const [user, setUser] = useState(defaultUser);
  const [loading, setLoading] = useState(false);
  const redirect = loc.pathname.startsWith("/")? loc.pathname : "/app";
  
  useEffect(() => {
    const fetchUser = async() => {
        setLoading(true);
        try{
          let user = await getUser()
          setUser(user)
          nav(""+loc.pathname )
        }catch(e){
          nav("/")
        }
        setLoading(false)
    };
    fetchUser();
  }, []);

  return (
    // the Provider gives access to the context to its children
    <UserContext.Provider value={[user, loading]}>
      {children}
    </UserContext.Provider>
  );
};

export { UserContext, UserContextProvider };