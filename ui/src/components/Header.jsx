import { useRef, useState } from "react";
import { useOnClickOutside } from "./Nav";
import { loadFromPlaylist } from "../api/watch2gether";
import { toast } from "react-hot-toast";
import { Link } from "react-router-dom";
import ReactPlayer from 'react-player'
import { useContext } from "react";
import { UserContext } from "../context/user";



const Header = ({ state }) => {
  return (
    <div className="flex shadow-head w-full flex-col ">
      <div className="flex md:items-end w-full shadow-head p-0 md:pt-8 md:pb-8 md:px-16 flex-col md:flex-row ">
        {state.player_type === "MUSIC"  && <img className="mr-0 md:rounded-xl shadow-xl w-full md:h-48 md:w-48 md:ml-8 object-cover mt-0 hidden sm:block" src={state.Current.thumbnail} />}
        <div className="flex grow flex-col sm:ml-8 justify-start  md:justify-center p-4 md:p-0 md:pb-2">
          {state.player_type === "VIDEO" && <div style={{ maxWidth: "100%", height: "calc( 100vh - 27rem )", aspectRatio: "16 / 9", marginBottom: "2rem" }}>
            <ReactPlayer url='https://www.youtube.com/watch?v=ysz5S6PUM-U' width='100%' height='100%'
              config={{
                youtube: {
                  playerVars: { showinfo: 0, autohide: 0 }
                }
              }} />
          </div>}

          <h4 className="mt-0 mb-2 uppercase text-white tracking-widest text-xs">Now Playing</h4>
          <h1 className="mt-0 text-white text-3xl md:text-5xl">{state.Current.title}</h1>



        </div>
      </div>
      





    </div>
  )
}

export default Header