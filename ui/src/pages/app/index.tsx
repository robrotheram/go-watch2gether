import { Outlet} from 'react-router-dom'
import { Nav } from './components/nav'
import { Notifications } from './components/notifications'
import { Provider } from './components/providers'

export default function () { 
  return (
      <Provider>
        <Nav/>
        <main className='w-full mt-8  fixed top-8 bottom-0 flex bg-gradient-to-b from-violet-900  to-black'>
          <Notifications/>
          <Outlet />
        </main>
      </Provider>
  )
}