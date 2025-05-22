import { useContext, useEffect, useRef, useState } from "react"
import logo from '../../../assets/logo.svg'
import { Link, useParams } from "react-router-dom";
import { GuildIcon } from "./roomList";
import { GuildContext } from "./providers";
import { DiscordGuild } from "@/types";

// Hook

export function useOnClickOutside(ref:React.RefObject<HTMLDivElement>, handler:(value: React.SetStateAction<boolean>) => void) {
  useEffect(
    () => {
      const listener = (event: MouseEvent | TouchEvent) => {
        if (!ref.current || ref.current.contains(event.target as Node)) {
          return;
        }
        handler(false);
      };
      document.addEventListener("mousedown", listener);
      document.addEventListener("touchstart", listener);
      return () => {
        document.removeEventListener("mousedown", listener);
        document.removeEventListener("touchstart", listener);
      };
    },
    [ref, handler]
  );
}



export const NavLogo = ({ guild }:{guild:DiscordGuild}) => {
  if (guild !== null && guild !== undefined) {
    return <GuildIcon guild={guild} />
  }
  return <img src={logo} className='h-16 w-16 bg-purple-700' alt="logo" />
}


export const Nav = () => {
  const params = useParams()
  const [isModalOpen, setIsModalOpen] = useState(false);
  const { user, settings, getGuild } = useContext(GuildContext)

  const guild = getGuild(params.id!)
  const ref = useRef<HTMLDivElement>(null);
  useOnClickOutside(ref, () => setIsModalOpen(false));
  return (
    <header className='flex shadow-lg fixed w-full top-0 justify-between bg-zinc-900 z-10'>
      <Link to={`/app/`} ><NavLogo guild={guild} /></Link>
      <h1 className='text-3xl font-dosis font-bold'>Watch2Gether</h1>
      <div ref={ref} className="flex items-center md:order-2 mr-4">
        <button onClick={() => setIsModalOpen(!isModalOpen)} type="button" className="float-right flex mr-4 text-sm rounded-full md:mr-0 focus:ring-4 focus:ring-gray-300 dark:focus:ring-gray-600" id="user-menu-button" aria-expanded="false" data-dropdown-toggle="user-dropdown" data-dropdown-placement="bottom">
          <img className="w-8 h-8 rounded-full" src={user.avatar_icon} alt="user profile" />
        </button>
        {isModalOpen && <div className="z-50 fixed top-10 right-2 my-4 text-base list-none divide-y divide-gray-100 rounded-lg shadow bg-zinc-800" id="user-dropdown">
          <ul className="px-4 py-1" aria-labelledby="user-menu-button">
            <li>
              <a href={`https://discord.com/oauth2/authorize?client_id=${settings.bot}&scope=bot`} className="block text-center px-8 py-3 text-md hover:rounded-lg  text-gray-700 hover:bg-gray-100 dark:hover:bg-gray-600 dark:text-gray-200 dark:hover:text-white">
                invite bot
              </a>
            </li>
            <li>
              <a href="/auth/logout" className="block text-center px-8 py-3 text-md hover:rounded-lg  text-gray-700 hover:bg-gray-100 dark:hover:bg-gray-600 dark:text-gray-200 dark:hover:text-white">
                Sign out
              </a>
            </li>
            <hr style={{ margin: "0.5rem 0" }} />
            <li className="text-center py-2">
              <span>Release: {import.meta.env.VITE_APP_VERSION}</span>
            </li>
          </ul>
        </div>}
      </div>
    </header>
  )
}