import { useEffect, useRef, useState } from 'react'
import './App.css'
import Sidebar from '../../components/Sidebar'
import { getGuilds, getRoomId, getUser } from '../../api/watch2gether'
import { Outlet, useParams, useSearchParams } from 'react-router-dom'
import { Nav } from '../../components/Nav'
import { Notifications } from '../../components/Notifications'
import useWebSocket from 'react-use-websocket';


function App() {
 
  
  const [user, setUser] = useState({});

  const updateUser = async() => {
    setUser(await getUser())
  }
  useEffect(() => {updateUser()}, [])


  return (
    <>
      <Nav user={user}/>
      <main className='w-full mt-8  fixed top-8 bottom-0 flex'>
        <Notifications/>
        <Outlet />
      </main>
    </>
  )
}



export default App
