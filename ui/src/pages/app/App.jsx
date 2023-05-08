import { useEffect, useRef, useState } from 'react'
import './App.css'
import Sidebar from '../../components/Sidebar'
import { getGuilds, getRoomId, getUser } from '../../api/watch2gether'
import { Outlet, useParams, useSearchParams } from 'react-router-dom'
import { Nav } from '../../components/Nav'
import { Notifications } from '../../components/Notifications'



function App() {
 
  const [guilds, setGuilds] = useState([]);
  const [user, setUser] = useState({});

  useEffect(() => {
    async function get() {
      setUser(await getUser())
      const g = await getGuilds();
      if (g != null) {
        setGuilds(g);
      }
    }
    if (guilds.length == 0) { get() };
  }, [])

  return (
    <>
      <Nav user={user} guilds={guilds} />
     
      <main className='w-full mt-8  fixed top-8 bottom-0 flex'>
        <Notifications/>
        <Sidebar guilds={guilds} />
        <Outlet />
      </main>
    </>
  )
}



export default App
