import { useEffect, useRef, useState } from 'react'
import { getGuilds, getRoomId, getUser } from './watch2gether'
import { Outlet, useParams, useSearchParams } from 'react-router-dom'
import { Nav } from './components/nav'
import { Notifications } from './components/notifications'

export default function () { 
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