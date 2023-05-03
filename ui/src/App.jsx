import { useEffect, useRef, useState } from 'react'
import logo from './assets/logo.svg'
import './App.css'
import Player from './Player'
import Card from './Cards'
import Header from './Header'
import Sidebar from './Sidebar'
import RoomList from './RoomList'
import Controller from './Controller'
import { getGuilds, getUser } from './api/watch2gether'


// Hook
function useOnClickOutside(ref, handler) {
  useEffect(
    () => {
      const listener = (event) => {
        // Do nothing if clicking ref's element or descendent elements
        if (!ref.current || ref.current.contains(event.target)) {
          return;
        }
        handler(event);
      };
      document.addEventListener("mousedown", listener);
      document.addEventListener("touchstart", listener);
      return () => {
        document.removeEventListener("mousedown", listener);
        document.removeEventListener("touchstart", listener);
      };
    },
    // Add ref and handler to effect dependencies
    // It's worth noting that because passed in handler is a new ...
    // ... function on every render that will cause this effect ...
    // ... callback/cleanup to run every render. It's not a big deal ...
    // ... but to optimize you can wrap handler in useCallback before ...
    // ... passing it into this hook.
    [ref, handler]
  );
}


function App() {
  const [showMenu, setShowMenu] = useState(false)

  const ref = useRef();
  const [isModalOpen, setModalOpen] = useState(false);
  useOnClickOutside(ref, () => setModalOpen(false));
  
  const [guilds, setGuilds] = useState([]);
  const [user, setUser] = useState({});
  
    useEffect(() => {
      async function get() {
            setUser(await getUser())
            const g = await getGuilds();
            if(g != null){
              setGuilds(g);
            }
        }
        if (guilds.length == 0 ){get()};
    }, [])

  return (
    <>
      <header className='flex shadow-lg fixed z-20 w-full top-0 justify-between bg-zinc-900'>
        <img src={logo} className='h-16 w-16 bg-purple-700' alt="logo" onClick={() => setShowMenu(!showMenu)} />
        <h1 className='text-3xl font-dosis font-bold'>Watch2Gether</h1>
        <div ref={ref} className="flex items-center md:order-2 mr-4">
          <button onClick={()=>setModalOpen(!isModalOpen)} type="button" className="float-right flex mr-4 text-sm rounded-full md:mr-0 focus:ring-4 focus:ring-gray-300 dark:focus:ring-gray-600" id="user-menu-button" aria-expanded="false" data-dropdown-toggle="user-dropdown" data-dropdown-placement="bottom">
            <img className="w-8 h-8 rounded-full" src={user.avatar_icon} alt="user photo" />
          </button>
          {isModalOpen && <div   className="z-50 fixed top-10 right-2 my-4 text-base list-none divide-y divide-gray-100 rounded-lg shadow bg-zinc-800" id="user-dropdown">
           <ul className="" aria-labelledby="user-menu-button">
              {/* <li className='p-4 text-center'>{user.username}</li> */}
              <li>
                <a href="/auth/logout" className="block text-center px-8 py-4 text-md hover:rounded-lg  text-gray-700 hover:bg-gray-100 dark:hover:bg-gray-600 dark:text-gray-200 dark:hover:text-white">
                  Sign out
                </a>
              </li>
            </ul>
          </div>}
        </div>
      </header>
      {showMenu && <section className="w-16 bottom-24 top-16 fixed md:hidden bg-black z-30 text-white overflow-y-auto hide-scrollbar shadow-xl shadow-left" >
        <RoomList guilds={guilds}/>
      </section>
      }
      <main className='w-full mt-8  fixed top-8 bottom-24 flex'>
        <Sidebar guilds={guilds}/>
        <Controller/>
      </main>
    </>
  )
}

export default App
