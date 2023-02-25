import axios from "axios";
import React, { createContext, useState, useEffect } from "react";
import { useLocation, useNavigate } from "react-router";
import { BASE_URL } from "./config";

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
  const redirect = loc.pathname.startsWith("/app/room/")? loc.pathname : "/app";
  
  useEffect(() => {
    const fetchUser = () => {
        setLoading(true);
        axios.get(`${BASE_URL}/auth/user`).then((res) => {
          console.log("DATA", res)
            setUser({
              ...res.data.user, guilds: res.data.guilds
            })
            setLoading(false);
            nav(redirect)
          }).catch(err => {
            setLoading(false)
          })  
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