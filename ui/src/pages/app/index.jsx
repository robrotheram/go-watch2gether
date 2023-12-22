import { useEffect, useState } from 'react'
import { getSettings, getUser } from './watch2gether'
import { Outlet} from 'react-router-dom'
import { Nav } from './components/nav'
import { Notifications } from './components/notifications'
import { Provider } from './components/providers'

export default function () { 
  const [user, setUser] = useState({});
  const [settings, setSettings] = useState({})

  const update = async() => {
    setUser(await getUser())
    setSettings(await getSettings())
  }
  useEffect(() => {update()}, [])
  return (
    <>
      <Provider>
        <Nav user={user} bot={settings.bot}/>
        <main className='w-full mt-8  fixed top-8 bottom-0 flex bg-gradient-to-b from-violet-900  to-black'>
          <Notifications/>
          <Outlet />
        </main>
      </Provider>
    </>
  )
}